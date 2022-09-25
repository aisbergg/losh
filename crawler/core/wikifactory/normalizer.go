// Copyright 2022 André Lehmann
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package wikifactory

import (
	"context"
	"fmt"
	"html"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"losh/crawler/core/wikifactory/wfclient"
	"losh/internal/core/product/models"
	"losh/internal/infra/dgraph/dgclient"
	"losh/internal/lib/fileformats"

	"github.com/abadojack/whatlanggo"
	"github.com/aisbergg/go-errors/pkg/errors"
	"github.com/aisbergg/go-pathlib/pkg/pathlib"
	"github.com/microcosm-cc/bluemonday"
)

var excludeFiles = []string{
	"ACKNOWLEDGMENTS",
	"AUTHORS",
	"CHANGELOG",
	"CODE_OF_CONDUCT",
	"CODEOWNERS",
	"CONTRIBUTING",
	"CONTRIBUTORS",
	"FUNDING",
	"ISSUE_TEMPLATE",
	"LICENSE",
	"PULL_REQUEST_TEMPLATE",
	"README",
	"SECURITY",
	"SUPPORT",
	"USERGUIDE",
	"USERMANUAL",
}
var licenseMapping = map[string]string{
	"CC-BY-4.0":    "CC-BY-4.0",
	"CC0-1.0":      "CC0-1.0",
	"MIT":          "MIT",
	"BSD-2-Clause": "BSD-2-Clause",
	"CC-BY-SA-4.0": "CC-BY-SA-4.0",
	"GPL-3.0":      "GPL-3.0-only",
	"OHL":          "TAPR-OHL-1.0",
	"CERN OHL":     "CERN-OHL-1.2",
}
var host = &models.Host{
	Domain: p("wikifactory.com"),
	Name:   p("Wikifactory"),
}
var activeTimeThreshold = 2 * 365 * 24 * time.Hour // 2 years

type rawData struct {
	timestamp   time.Time
	projectInfo *wfclient.ProjectFullFragment
	// group information, because I can't seem to get the group information from the project info
	groupInfo *wfclient.GetGroup_Initiative_Result
}

// NormalizeProduct creates a normalized product from the given Wikifactory
// project information.
func (c *WikifactoryCrawler) NormalizeProduct(ctx context.Context, timestamp time.Time, wfPrjInfo *wfclient.ProjectFullFragment) (*models.Product, error) {
	product := &models.Product{}

	// releases
	licensor, err := c.normLicensor(ctx, wfPrjInfo, timestamp)
	if err != nil {
		return nil, err
	}
	releases := c.normReleases(wfPrjInfo, licensor, timestamp)
	latestRelease := releases[0] // latest release is always the first

	// crawler info
	product.DiscoveredAt = latestRelease.DiscoveredAt
	product.LastIndexedAt = latestRelease.LastIndexedAt
	product.DataSource = latestRelease.DataSource

	// product info
	productURL := &models.ProductURL{
		Domain: "wikifactory.com",
		Owner:  *wfPrjInfo.ParentSlug,
		Repo:   *wfPrjInfo.Slug,
	}
	// Xid format: domain.tld/owner/repo/file-path
	product.Xid = asXid(productURL.Domain, productURL.Owner, productURL.Repo, "")
	product.Name = latestRelease.Name
	product.Description = latestRelease.Description
	product.DocumentationLanguage = latestRelease.DocumentationLanguage
	product.Version = latestRelease.Version
	product.License = latestRelease.License
	product.Licensor = licensor
	product.Website = product.DataSource.URL
	product.State = c.normState(wfPrjInfo)
	product.LastUpdatedAt = wfPrjInfo.LastUpdated
	product.Release = latestRelease
	product.Releases = releases
	for _, release := range releases {
		release.Product = product
	}

	// product.RenamedTo = "XXX"   // TODO
	// product.RenamedFrom = "XXX" // TODO
	// product.ForkOf = "XXX"      // TODO
	product.Forks = []*models.Product{} // TODO
	product.ForkCount = wfPrjInfo.ForkCount
	product.StarCount = wfPrjInfo.StarCount
	product.Tags = normTags(wfPrjInfo.Tags)
	// product.Category = "XXX" // TODO

	return product, nil
}

func (c *WikifactoryCrawler) normState(prjInfo *wfclient.ProjectFullFragment) *dgclient.ProductState {
	state := dgclient.ProductStateUndetermined
	// states other active/inactive are not reported by Wikifactory
	if time.Now().Sub(*prjInfo.LastUpdated) < activeTimeThreshold {
		state = dgclient.ProductStateActive
	} else {
		state = dgclient.ProductStateInactive
	}
	return &state
}

// normReleases returns the release of the product.
func (c *WikifactoryCrawler) normReleases(prjInfo *wfclient.ProjectFullFragment, owner models.UserOrGroup, timestamp time.Time) []*models.Component {
	wfContribs := prjInfo.Contributions.Edges
	releases := make([]*models.Component, 0, len(wfContribs))

	// image (same for every release)
	var image *models.File

	for i, edge := range wfContribs {
		wfContrib := edge.Node
		version := wfContrib.Version
		crawlerMeta := &models.CrawlerMetaImpl{
			DiscoveredAt:  &timestamp,
			LastIndexedAt: &timestamp,
			DataSource:    c.normRepository(owner, *version, prjInfo),
		}
		files := c.normFiles(wfContrib, crawlerMeta)

		release := &models.Component{}
		release.DiscoveredAt = crawlerMeta.DiscoveredAt
		release.LastIndexedAt = crawlerMeta.LastIndexedAt
		release.DataSource = crawlerMeta.DataSource

		// Xid format: domain.tld/owner/repo/ref/file-path/component-name
		release.Xid = asXid(*release.DataSource.Host.Domain, *release.DataSource.Owner.GetName(), *release.DataSource.Name, *release.DataSource.Reference, "", *prjInfo.Name)
		release.Name = sp(prjInfo.Name)
		release.Description = c.normDescription(prjInfo)
		release.Version = wfContrib.Version
		release.CreatedAt = wfContrib.DateCreated
		release.Releases = make([]*models.Component, 0, len(wfContribs))
		if i == 0 {
			release.IsLatest = p(true)
		}
		release.Repository = crawlerMeta.DataSource
		release.License = c.normLicense(prjInfo.License)
		release.Licensor = owner
		release.DocumentationLanguage = c.normDocumentationLanguage(*release.Description)
		release.TechnologyReadinessLevel = p(dgclient.TechnologyReadinessLevelUndetermined)
		release.DocumentationReadinessLevel = p(dgclient.DocumentationReadinessLevelUndetermined)
		// release.Attestation = "XXX" // TODO
		// release.Publication = "XXX" // TODO
		// release.CompliesWith = "XXX" // TODO
		// release.CpcPatentClass = "XXX" // TODO
		// release.Tsdc = "XXX" // TODO

		// TODO: mandatory information need to be added to subcomponents
		// subCmps := c.getSubComponents(files, prjInfo)
		// release.Components = subCmps

		release.Software = []*models.Software{}
		if image == nil {
			image = c.getImage(prjInfo.Image, crawlerMeta)
		}
		release.Image = image
		release.Readme = c.normInfoFile([]string{"README"}, files)
		release.ContributionGuide = c.normInfoFile([]string{"CONTRIBUTING"}, files)
		release.Bom = c.normInfoFile([]string{"BOM", "BILLOFMATERIALS"}, files)
		release.ManufacturingInstructions = c.normInfoFile([]string{"MANUFACTURINGINSTRUCTIONS", "MANUFACTURING"}, files)
		release.UserManual = c.normInfoFile([]string{"USERGUIDE", "USERMANUAL"}, files)
		// release.Product = "XXX" // TODO
		// release.UsedIn = "XXX" // TODO
		// release.Source = "XXX" // TODO
		// release.Export = "XXX" // TODO
		// release.Auxiliary = "XXX" // TODO
		// release.Organization = "XXX" // TODO
		// release.Mass = "XXX" // TODO
		// release.OuterDimensions = "XXX" // TODO
		// release.Material = "XXX" // TODO
		// release.ManufacturingProcess = "XXX" // TODO
		// release.ProductionMetadata = "XXX" // TODO

		releases = append(releases, release)
	}

	// link releases to each other
	for i := 0; i < len(releases); i++ {
		release := releases[i]
		for j := 0; j < len(releases); j++ {
			release.Releases = append(release.Releases, releases[j])
		}
	}
	// for _, release := range releases {
	// 	release.Releases = append(release.Releases, release)
	// }

	return releases
}

// normRepository returns the source of the component.
func (c *WikifactoryCrawler) normRepository(owner models.UserOrGroup, ref string, prjInfo *wfclient.ProjectFullFragment) *models.Repository {
	productURL := &models.ProductURL{
		Domain: "wikifactory.com",
		Owner:  *prjInfo.ParentSlug,
		Repo:   *prjInfo.Slug,
		Ref:    ref,
	}
	repoURL := productURL.RepositoryURL()
	permaURL := productURL.PermaURL()
	// Xid format: domain.tld/owner/repo/ref/file-path
	xid := strings.Join([]string{productURL.Domain, productURL.Owner, productURL.Repo, productURL.Ref, "-"}, "/")
	return &models.Repository{
		Xid:       p(xid),            // TODO
		URL:       p(repoURL),        // TODO
		PermaURL:  p(permaURL),       // TODO
		Host:      host,              // TODO
		Owner:     owner,             // TODO
		Name:      prjInfo.Slug,      // TODO
		Reference: stringOrNil(&ref), // TODO
	}
}

// func (c *WikifactoryCrawler) normRepository(owner models.UserOrGroup, repo, ref *string, prjInfo *wfclient.ProjectFullFragment) *models.Repository {
// 	productURL := &models.ProductURL{
// 		Domain: "wikifactory.com",
// 		Owner:  *prjInfo.ParentSlug,
// 		Repo:   *prjInfo.Slug,
// 		Ref:    ref,
// 	}
// 	repoURL := productURL.RepositoryURL()
// 	permaURL := productURL.PermaURL()
// 	// Xid format: domain.tld/owner/repo/ref/file-path
// 	// xid := strings.Join([]string{productURL.Domain, productURL.Owner, productURL.Repo, productURL.Ref, "-"}, "/")

// 	xid := asXid(*host.Domain, owner, *productURL.Repo, *productURL.Ref, "")
// 	return &models.Repository{
// 		Xid:       xid,
// 		URL:       p(repoURL),
// 		PermaURL:  p(permaURL),
// 		Host:      host,
// 		Owner:     owner,
// 		Name:      prjInfo.Slug,
// 		Reference: stringOrNil(ref),
// 	}
// }

// normDescription returns the normDescription of the product.
func (c *WikifactoryCrawler) normDescription(prjInfo *wfclient.ProjectFullFragment) *string {
	htmlDescription := sp(prjInfo.Description)
	if htmlDescription == nil || *htmlDescription == "" {
		return nil
	}

	// sanitize description by removing all HTML tags
	return p(stripHTMLTags(*htmlDescription))
}

// normDocumentationLanguage returns the documentation language of the product.
func (c *WikifactoryCrawler) normDocumentationLanguage(documentation string) *string {
	defaultLang := "en"
	if documentation == "" {
		return &defaultLang
	}
	lang := whatlanggo.DetectLang(documentation).Iso6391()
	if lang == "" {
		return &defaultLang
	}
	return &lang
}

// normFiles returns the files in the Wikifactory project.
func (c *WikifactoryCrawler) normFiles(contrib *wfclient.ContributionFragment, crawlerMeta *models.CrawlerMetaImpl) []*models.File {
	files := make([]*models.File, 0, 10)
	for _, wfFileMeta := range contrib.Files {
		wfFile := wfFileMeta.File
		if wfFile == nil { // skip directories
			continue
		}
		wfFile.Path = &wfFile.Filename
		dirName := wfFileMeta.Dirname
		if dirName != nil {
			path := fmt.Sprintf("%s/%s", *dirName, wfFile.Filename)
			wfFile.Path = &path // reuse field, it doesn't contain usable data
		}
		file := c.normFile(wfFile, crawlerMeta)
		files = append(files, file)
	}
	return files
}

// normFile returns a file from the Wikifactory project.
func (c *WikifactoryCrawler) normFile(wfFile *wfclient.FileFragment, crawlerMeta *models.CrawlerMetaImpl) *models.File {
	if wfFile == nil {
		return nil
	}
	file := &models.File{}
	file.Path = wfFile.Path
	file.Name = p(filepath.Base(*file.Path))
	file.MimeType = stringOrNil(&wfFile.MimeType)
	file.URL = wfFile.URL
	file.CreatedAt = wfFile.DateCreated
	file.DiscoveredAt = crawlerMeta.DiscoveredAt
	file.LastIndexedAt = crawlerMeta.LastIndexedAt
	file.DataSource = crawlerMeta.DataSource

	owner := ""
	if crawlerMeta.DataSource.Owner != nil {
		owner = *crawlerMeta.DataSource.Owner.GetName()
	}
	repo := ""
	if crawlerMeta.DataSource.Name != nil {
		repo = *crawlerMeta.DataSource.Name
	}
	ref := ""
	if crawlerMeta.DataSource.Reference != nil {
		ref = *crawlerMeta.DataSource.Reference
	}
	file.Xid = asXid(*crawlerMeta.DataSource.Host.Domain, owner, repo, ref, *file.Path)

	return file
}

// normInfoFile returns the info file of the product.
func (c *WikifactoryCrawler) getImage(wfFile *wfclient.FileFragment, crawlerMeta *models.CrawlerMetaImpl) *models.File {
	image := c.normFile(wfFile, crawlerMeta)
	if image == nil {
		return nil
	}
	image.Xid = asXid(*crawlerMeta.DataSource.Host.Domain, *crawlerMeta.DataSource.Owner.GetName(), "", "", *image.Path)
	return image
}

// normInfoFile returns the info file of the product.
func (c *WikifactoryCrawler) normInfoFile(names []string, files []*models.File) *models.File {
	for _, file := range files {
		// only consider files in root dir
		parts := strings.Split(strings.TrimLeft(*file.Path, "/"), "/")
		if len(parts) > 1 {
			continue
		}
		for _, name := range names {
			filename := *file.Name
			if pos := strings.LastIndexByte(filename, '.'); pos != -1 {
				filename = filename[:pos]
			}
			filename = strings.TrimSpace(filename)
			filename = strings.Replace(filename, " ", "", -1)
			filename = strings.Replace(filename, "-", "", -1)
			filename = strings.Replace(filename, "_", "", -1)
			filename = strings.ToUpper(filename)
			if filename == name {
				return file
			}
		}
	}
	return nil
}

// translateLicense translates the license IDs used by Wikifactory into SPDX
// license IDs.
func (*WikifactoryCrawler) translateLicense(wfLcsStr string) string {
	wfLcsStr = strings.TrimSpace(wfLcsStr)
	if lcsStr, ok := licenseMapping[strings.ToUpper(wfLcsStr)]; ok {
		return lcsStr
	}
	return wfLcsStr
}

// normLicense returns the license of the product.
func (c *WikifactoryCrawler) normLicense(wfLicense *wfclient.ProjectFullFragment_License) *models.License {
	if wfLicense == nil {
		return nil
	}
	lcsStr := c.translateLicense(*wfLicense.Abreviation)
	return c.productService.GetCachedLicenseByIDOrName(lcsStr)
}

// normLicensor returns the owner of the product.
func (c *WikifactoryCrawler) normLicensor(ctx context.Context, prjInfo *wfclient.ProjectFullFragment, timestamp time.Time) (models.UserOrGroup, error) {
	switch *prjInfo.ParentContent.Type {
	case "initiative": // is group
		// need more information
		getGroup, err := c.wfClient.GetGroup(ctx, *prjInfo.ParentSlug)
		if err != nil || getGroup.Initiative.Result == nil {
			return nil, errors.Wrapf(err, "failed to get information of group '%s'", *prjInfo.ParentContent.Slug)
		}
		groupInfo := getGroup.Initiative.Result

		xid := asXid(*host.Domain, *groupInfo.Slug)
		url := "https://" + *xid
		group := &models.Group{
			Xid:         xid,
			Host:        host,
			Name:        groupInfo.Slug,
			FullName:    groupInfo.Title,
			Description: stringOrNil(sp(groupInfo.Description)),
			Members:     []models.UserOrGroup{}, //TODO
			URL:         &url,
		}

		if groupInfo.Avatar != nil {
			crawlerMeta := &models.CrawlerMetaImpl{
				DiscoveredAt:  &timestamp,
				LastIndexedAt: &timestamp,
				DataSource:    c.normRepository(group, "latest", prjInfo),
			}
			avatar := c.normFile(groupInfo.Avatar, crawlerMeta)
			if avatar != nil {
				// Xid format: domain.tld/owner/repo/ref/file-path
				avatar.Xid = asXid(*host.Domain, *group.Name, "", "", *avatar.Path)
				group.Avatar = avatar
			}
		}

		return group, nil

	default: // is user
		xid := asXid(*host.Domain, *prjInfo.Creator.Profile.Username)
		url := "https://" + *xid
		user := &models.User{
			Xid:         xid,
			Host:        host,
			Name:        prjInfo.Creator.Profile.Username,
			FullName:    stringOrNil(sp(prjInfo.Creator.Profile.FullName)),
			Email:       stringOrNil(sp(prjInfo.Creator.Profile.Email)),
			Description: stringOrNil(sp(prjInfo.Creator.Profile.Description)),
			Locale:      stringOrNil(prjInfo.Creator.Profile.Locale),
			URL:         &url,
		}

		if prjInfo.Creator.Profile.Avatar != nil {
			crawlerMeta := &models.CrawlerMetaImpl{
				DiscoveredAt:  &timestamp,
				LastIndexedAt: &timestamp,
				DataSource:    c.normRepository(user, "latest", prjInfo),
			}
			avatar := c.normFile(prjInfo.Creator.Profile.Avatar, crawlerMeta)
			if avatar != nil {
				// Xid format: domain.tld/owner/repo/ref/file-path
				avatar.Xid = asXid(*host.Domain, *user.Name, "", "", *avatar.Path)
				user.Avatar = avatar
			}
		}

		return user, nil
	}
}

// getSubComponents returns the sub components of the release.
// XXX: need a better way of identifying components
func (c *WikifactoryCrawler) getSubComponents(files []*models.File, prjInfo *wfclient.ProjectFullFragment) []*models.Component {
	// filter out readme and other files
	filtered := make([]*models.File, 0, len(files))
	for _, file := range files {
		filename := *file.Name
		if pos := strings.LastIndexByte(filename, '.'); pos != -1 {
			filename = filename[:pos]
		}
		filename = strings.TrimSpace(filename)
		filename = strings.Replace(filename, " ", "", -1)
		filename = strings.Replace(filename, "-", "", -1)
		filename = strings.Replace(filename, "_", "", -1)
		filename = strings.ToUpper(filename)
		for _, excl := range excludeFiles {
			if filename == excl {
				continue
			}
		}
		filtered = append(filtered, file)
	}

	type FileWrap struct {
		file *models.File
		path pathlib.PurePath
	}

	fileWraps := make([]FileWrap, 0, len(filtered))
	for _, file := range filtered {
		fileWraps = append(fileWraps, FileWrap{file, pathlib.NewPurePosixPath(*file.Path)})
	}

	// put files in buckets
	buckets := make(map[string][]FileWrap, len(filtered))
	for _, fileWrap := range fileWraps {
		ps, _ := fileWrap.path.WithSuffix("")
		normalizedName := strings.ToLower(ps.String())
		bucket, ok := buckets[normalizedName]
		if !ok {
			bucket = make([]FileWrap, 0, 1)
		}
		buckets[normalizedName] = append(bucket, fileWrap)
	}

	// figure out what files are the sources, the exports and the images
	cmps := make([]*models.Component, 0, len(buckets))
	for _, bucket := range buckets {
		cmp := &models.Component{}
		for _, fileWrap := range bucket {
			ext := fileWrap.path.Suffix()

			//  get sources and exports by extension
			if fileformats.IsCADFile(ext) {
				isSource := fileformats.IsSourceFile(ext)
				if isSource {
					if cmp.Source == nil {
						cmp.Source = fileWrap.file
					} else {
						cmp.Export = append(cmp.Export, fileWrap.file)
					}
				} else {
					cmp.Export = append(cmp.Export, fileWrap.file)
				}
				continue
			} else if fileformats.IsPCBFile(ext) {
				isSource := fileformats.IsSourceFile(ext)
				if isSource {
					if cmp.Source == nil {
						cmp.Source = fileWrap.file
					} else {
						cmp.Export = append(cmp.Export, fileWrap.file)
					}
				} else {
					cmp.Export = append(cmp.Export, fileWrap.file)
				}
				continue
			}

			// get first image by extension
			if fileformats.IsImageFile(ext) {
				if cmp.Image == nil {
					cmp.Image = fileWrap.file
				}
				continue
			}
		}

		// if no source file was identified, but exports, then use the exports instead
		if cmp.Source == nil && len(cmp.Export) > 0 {
			cmp.Source = cmp.Export[0]
			cmp.Export = cmp.Export[1:]
		}

		// # only add, if a source file was identified
		if cmp.Source != nil {
			cmp.Name = cmp.Source.Name
			cmps = append(cmps, cmp)
		}
	}

	if len(cmps) == 0 {
		return nil
	}
	return cmps
}

// TODO: make this better, e.g. https://github.com/abhinav/goldmark-hashtag
// https://meta.stackexchange.com/a/298981
var validTagPattern = regexp.MustCompile(`^[a-z0-9_.+-]+$`)

// parseTags parses the given string into a slice of tags.
func normTags(wfTags []*wfclient.ProjectFullFragment_Tags) []*models.Tag {
	if wfTags == nil || len(wfTags) == 0 {
		return nil
	}
	tags := make([]*models.Tag, 0, len(wfTags))
	for _, tag := range wfTags {
		name := strings.ToLower(strings.TrimSpace(*tag.Name))
		if len(name) == 0 ||
			len(name) > 128 ||
			!validTagPattern.MatchString(name) {
			continue
		}
		tags = append(tags, &models.Tag{
			Name: &name,
		})
	}
	if len(tags) == 0 {
		return nil
	}
	return tags
}

// stringOrNil returns the string pointer if it is non nil and contains a string
// else returns nil.
func stringOrNil(s *string) *string {
	if s == nil || *s == "" {
		return nil
	}
	return s
}

func stringOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// stripHTMLTags removes HTML tags from the given string.
func stripHTMLTags(h string) string {
	p := bluemonday.StrictPolicy()
	h = p.Sanitize(h)
	h = html.UnescapeString(h)
	return h
}

func sp(s *string) *string {
	if s == nil {
		return nil
	}
	*s = strings.TrimSpace(*s)
	return s
}

func s(s string) string {
	return strings.TrimSpace(s)
}

func p[T any](v T) *T {
	return &v
}

func asXid(args ...string) *string {
	if args == nil || len(args) == 0 {
		return nil
	}
	var b strings.Builder
	for i := 0; i < len(args); i++ {
		if i > 0 {
			b.WriteString("/")
		}
		if args[i] == "" {
			b.WriteString("-")
		} else {
			b.WriteString(url.PathEscape(args[i]))
		}
	}
	xid := b.String()
	return &xid
}
