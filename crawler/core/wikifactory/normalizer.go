package wikifactory

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"losh/crawler/core/wikifactory/wfclient"
	"losh/internal/core/product/models"
	"losh/internal/lib/fileformats"
	"losh/internal/license"

	"github.com/abadojack/whatlanggo"
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

type WikifactoryNormalizer struct {
	licenseCache *license.LicenseCache
}

// NewWikifactoryNormalizer creates a new WikifactoryNormalizer.
func NewWikifactoryNormalizer(licenseCache *license.LicenseCache) *WikifactoryNormalizer {
	return &WikifactoryNormalizer{
		licenseCache: licenseCache,
	}
}

// NormalizeProduct creates a normalized product from the given Wikifactory
// project information.
func (wn *WikifactoryNormalizer) NormalizeProduct(timestamp time.Time, wfProjectInfo *wfclient.GetFullProjectBySlug) *models.Product {
	product := &models.Product{}

	// releases
	owner := wn.getOwner(wfProjectInfo, timestamp)
	releases := wn.getReleases(wfProjectInfo, owner, timestamp)
	latestRelease := releases[0] // latest release is always the first

	// crawler info
	product.DiscoveredAt = latestRelease.DiscoveredAt
	product.LastIndexedAt = latestRelease.LastIndexedAt
	product.DataSource = latestRelease.DataSource

	// product info
	productURL := &models.ProductURL{
		Domain: "wikifactory.com",
		Owner:  *wfProjectInfo.Project.Result.ParentSlug,
		Repo:   *wfProjectInfo.Project.Result.Slug,
	}
	product.Xid = fmt.Sprintf("%s/%s/%s", productURL.Domain, productURL.Owner, productURL.Repo)
	product.Name = latestRelease.Name
	product.Owner = owner
	product.Description = latestRelease.Description
	product.Website = &product.DataSource.URL
	product.Version = latestRelease.Version
	product.Release = latestRelease
	product.Releases = releases
	// product.RenamedTo = "XXX"   // TODO
	// product.RenamedFrom = "XXX" // TODO
	// product.ForkOf = "XXX"      // TODO
	product.Forks = []*models.Product{} // TODO
	// product.Tags = "XXX"     // TODO
	// product.Category = "XXX" // TODO

	return product

	// Crawler Start:
	// - Upload Host
	// - [check diagram]

	// Add
	// - Get Product Entry from DB
	// - Exists -> Skip
	// - Update Product Entry
	// - Add Product to DB
}

// getReleases returns the release of the product.
func (wn *WikifactoryNormalizer) getReleases(wfProjectInfo *wfclient.GetFullProjectBySlug, owner models.UserOrGroup, timestamp time.Time) []*models.Component {
	wfContribs := wfProjectInfo.Project.Result.Contributions.Edges
	releases := make([]*models.Component, 0, len(wfContribs))

	for i, edge := range wfContribs {
		wfContrib := edge.Node
		version := *wfContrib.Version
		crawlerMeta := &models.CrawlerMetaImpl{
			DiscoveredAt:  timestamp,
			LastIndexedAt: timestamp,
			DataSource:    wn.getRepository(owner, version, wfProjectInfo),
		}
		files := wn.getFiles(wfContrib, crawlerMeta)

		release := &models.Component{}
		release.DiscoveredAt = crawlerMeta.DiscoveredAt
		release.LastIndexedAt = crawlerMeta.LastIndexedAt
		release.DataSource = crawlerMeta.DataSource

		release.Name = *wfProjectInfo.Project.Result.Name
		release.Description = wn.getDescription(wfProjectInfo)
		release.Owner = owner
		release.Version = *wfContrib.Version
		release.CreatedAt = *wfProjectInfo.Project.Result.DateCreated
		release.Releases = []*models.Component{} // TODO
		if i == 0 {
			release.IsLatest = true
		}
		release.Repository = crawlerMeta.DataSource
		release.License = wn.getLicense(wfProjectInfo.Project.Result.License.Abreviation)
		release.Licensor = owner
		release.DocumentationLanguage = wn.getDocumentationLanguage(release.Description)
		// release.TechnologyReadinessLevel = "XXX" // TODO
		// release.DocumentationReadinessLevel = "XXX" // TODO
		// release.Attestation = "XXX" // TODO
		// release.Publication = "XXX" // TODO
		// release.CompliesWith = "XXX" // TODO
		// release.CpcPatentClass = "XXX" // TODO
		// release.Tsdc = "XXX" // TODO
		// subComponents := wn.getSubComponents(files, wfProjectInfo)
		release.Components = []*models.Component{} // TODO
		release.Software = []*models.Software{}
		release.Image = wn.getFile(wfProjectInfo.Project.Result.Image, crawlerMeta)
		release.Readme = wn.getInfoFile([]string{"README"}, files)
		release.ContributionGuide = wn.getInfoFile([]string{"CONTRIBUTING"}, files)
		release.Bom = wn.getInfoFile([]string{"BOM", "BILLOFMATERIALS"}, files)
		release.ManufacturingInstructions = wn.getInfoFile([]string{"MANUFACTURINGINSTRUCTIONS", "MANUFACTURING"}, files)
		release.UserManual = wn.getInfoFile([]string{"USERGUIDE", "USERMANUAL"}, files)
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

	return releases
}

// getRepository returns the source of the component.
func (wn *WikifactoryNormalizer) getRepository(owner models.UserOrGroup, version string, wfProjectInfo *wfclient.GetFullProjectBySlug) models.Repository {
	productURL := &models.ProductURL{
		Domain: "wikifactory.com",
		Owner:  *wfProjectInfo.Project.Result.ParentSlug,
		Repo:   *wfProjectInfo.Project.Result.Slug,
		Tag:    version,
	}
	repoURL := productURL.RepositoryURL()
	permaURL := productURL.PermaURL()
	xid := strings.Join([]string{productURL.Domain, productURL.Owner, productURL.Repo, productURL.Tag}, "/")
	return models.Repository{
		Xid:      xid,                               // TODO
		URL:      repoURL,                           // TODO
		PermaURL: permaURL,                          // TODO
		Host:     wn.getHost(),                      // TODO
		Owner:    owner,                             // TODO
		Name:     wfProjectInfo.Project.Result.Slug, // TODO
		Tag:      stringOrNil(&version),             // TODO
	}
}

// getDescription returns the getDescription of the product.
func (wn *WikifactoryNormalizer) getDescription(wfProjectInfo *wfclient.GetFullProjectBySlug) string {
	htmlDescription := wfProjectInfo.Project.Result.Description
	if htmlDescription == nil || *htmlDescription == "" {
		return ""
	}

	// sanitize description by removing all HTML tags
	return stripHTMLTags(*htmlDescription)
}

// getDocumentationLanguage returns the documentation language of the product.
func (wn *WikifactoryNormalizer) getDocumentationLanguage(documentation string) string {
	defaultLang := "en"
	if documentation == "" {
		return defaultLang
	}
	lang := whatlanggo.DetectLang(documentation).Iso6391()
	if lang == "" {
		return defaultLang
	}
	return lang
}

// getFiles returns the files in the Wikifactory project.
func (wn *WikifactoryNormalizer) getFiles(contrib *wfclient.ContributionFragment, crawlerMeta *models.CrawlerMetaImpl) []*models.File {
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
		file := wn.getFile(wfFile, crawlerMeta)
		files = append(files, file)
	}
	return files
}

// getFile returns a file from the Wikifactory project.
func (wn *WikifactoryNormalizer) getFile(wfFile *wfclient.FileFragment, crawlerMeta *models.CrawlerMetaImpl) *models.File {
	if wfFile == nil {
		return nil
	}
	file := &models.File{}
	file.Path = *wfFile.Path
	file.Name = filepath.Base(file.Path)
	file.MimeType = stringOrNil(&wfFile.MimeType)
	file.URL = *wfFile.Permalink
	file.CreatedAt = wfFile.DateCreated
	file.DiscoveredAt = crawlerMeta.DiscoveredAt
	file.LastIndexedAt = crawlerMeta.LastIndexedAt
	file.DataSource = crawlerMeta.DataSource
	return file
}

// getHost returns the host of the product.
func (wn *WikifactoryNormalizer) getHost() models.Host {
	return models.Host{
		Domain: "wikifactory.com",
		Name:   "Wikifactory",
	}
}

// getInfoFile returns the info file of the product.
func (wn *WikifactoryNormalizer) getInfoFile(names []string, files []*models.File) *models.File {
	for _, file := range files {
		// only consider files in root dir
		parts := strings.Split(strings.TrimLeft(file.Path, "/"), "/")
		if len(parts) > 1 {
			continue
		}
		for _, name := range names {
			filename := file.Name
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
func (*WikifactoryNormalizer) translateLicense(wfLcsStr string) string {
	wfLcsStr = strings.TrimSpace(wfLcsStr)
	if lcsStr, ok := licenseMapping[strings.ToUpper(wfLcsStr)]; ok {
		return lcsStr
	}
	return wfLcsStr
}

// getLicense returns the license of the product.
func (wn *WikifactoryNormalizer) getLicense(wfLicenseString *string) models.License {
	if wfLicenseString == nil {
		return models.License{}
	}
	*wfLicenseString = wn.translateLicense(*wfLicenseString)
	license := wn.licenseCache.GetByIDOrName(*wfLicenseString)
	if license == nil {
		return models.License{}
	}
	return *license
}

// getOwner returns the owner of the product.
func (wn *WikifactoryNormalizer) getOwner(wfProjectInfo *wfclient.GetFullProjectBySlug, timestamp time.Time) (owner models.UserOrGroup) {
	creator := wfProjectInfo.Project.Result.Creator
	platform := wn.getHost()
	xid := platform.Domain + "/" + *creator.Profile.Username
	url := "https://" + xid
	var avatar *models.File
	if wfProjectInfo.Project.Result.Creator.Profile.Avatar != nil {
		crawlerMeta := &models.CrawlerMetaImpl{
			DiscoveredAt:  timestamp,
			LastIndexedAt: timestamp,
			DataSource:    models.Repository{},
		}
		avatar = wn.getFile(wfProjectInfo.Project.Result.Creator.Profile.Avatar, crawlerMeta)
	}

	// is group
	if *wfProjectInfo.Project.Result.ParentContent.Type == "initiative" {
		owner = &models.Group{
			Xid:     xid,
			Host:    platform,
			Name:    *creator.Profile.Username,
			Email:   stringOrNil(creator.Profile.Email),
			Members: []models.UserOrGroup{},
			Avatar:  avatar,
			URL:     &url, // TODO
			// MemberOf: "", // TODO
			Products: []*models.Product{}, // TODO
		}

	} else {
		// is user
		owner = models.User{
			Xid:      xid,
			Host:     platform,
			Name:     *creator.Profile.Username,
			FullName: stringOrNil(creator.Profile.FullName),
			Email:    stringOrNil(creator.Profile.Email),
			Locale:   stringOrNil(creator.Profile.Locale),
			Avatar:   avatar,
			URL:      &url,
			// MemberOf: "", // TODO
			Products: []*models.Product{}, // TODO
		}
	}

	if avatar != nil {
		avatar.DataSource = wn.getRepository(owner, "latest", wfProjectInfo)
	}

	return
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

// getSubComponents returns the sub components of the release.
func (wn *WikifactoryNormalizer) getSubComponents(files []*models.File, wfProjectInfo *wfclient.GetFullProjectBySlug) []*models.Component {
	// filter out readme and other files
	filtered := make([]*models.File, 0, len(files))
	for _, file := range files {
		filename := file.Name
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
		fileWraps = append(fileWraps, FileWrap{file, pathlib.NewPurePosixPath(file.Path)})
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
	parts := make([]*models.Component, 0, len(buckets))
	for _, bucket := range buckets {
		part := &models.Component{}
		for _, fileWrap := range bucket {
			ext := fileWrap.path.Suffix()

			//  get sources and exports by extension
			if fileformats.IsCADFile(ext) {
				isSource := fileformats.IsSourceFile(ext)
				if isSource {
					if part.Source == nil {
						part.Source = fileWrap.file
					} else {
						part.Export = append(part.Export, fileWrap.file)
					}
				} else {
					part.Export = append(part.Export, fileWrap.file)
				}
				continue
			} else if fileformats.IsPCBFile(ext) {
				isSource := fileformats.IsSourceFile(ext)
				if isSource {
					if part.Source == nil {
						part.Source = fileWrap.file
					} else {
						part.Export = append(part.Export, fileWrap.file)
					}
				} else {
					part.Export = append(part.Export, fileWrap.file)
				}
				continue
			}

			// get first image by extension
			if fileformats.IsImageFile(ext) {
				if part.Image == nil {
					part.Image = fileWrap.file
				}
				continue
			}
		}

		// if no source file was identified, but exports, then use the exports instead
		if part.Source == nil && len(part.Export) > 0 {
			part.Source = part.Export[0]
			part.Export = part.Export[1:]
		}

		// # only add, if a source file was identified
		if part.Source != nil {
			part.Name = part.Source.Name
			parts = append(parts, part)
		}
	}

	return parts
}

// stripHTMLTags removes HTML tags from the given string.
func stripHTMLTags(html string) string {
	p := bluemonday.StrictPolicy()
	return p.Sanitize(html)
}
