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

// Package spdxorg provides a license provider that downloads licenses from the
// SPDX.org license list.
package spdxorg

import (
	"context"
	"embed"
	"encoding/json"
	"net/http"
	"strings"

	"losh/internal/core/product/models"
	"losh/internal/infra/dgraph/dgclient"
	"losh/internal/lib/log"
	"losh/internal/lib/net/download"
	"losh/internal/lib/net/request"
	"losh/internal/lib/util/stringutil"

	"github.com/aisbergg/go-errors/pkg/errors"
	"go.uber.org/zap"
)

type spdxLicenseFile struct {
	LicenseListVersion string        `json:"licenseListVersion"`
	Licenses           []spdxLicense `json:"licenses"`
}

type spdxLicense struct {
	LicenseID             string   `json:"licenseId"`
	Name                  string   `json:"name"`
	Reference             string   `json:"reference"`
	DetailsURL            string   `json:"detailsUrl"`
	SeeAlso               []string `json:"seeAlso"`
	ReferenceNumber       int      `json:"referenceNumber"`
	IsDeprecatedLicenseID bool     `json:"isDeprecatedLicenseId"`
	IsOSIApproved         bool     `json:"isOsiApproved"`
	IsFSFLibre            bool     `json:"isFsfLibre"`
	IsBlocked             bool     `json:"isBlocked"`
	Type                  string   `json:"type"`
}

type spdxLicenseDetails struct {
	LicenseText     string `json:"licenseText"`
	LicenseTextHTML string `json:"licenseTextHtml"`
}

//go:embed *.json
var assets embed.FS

// SpdxOrgProvider is a license provider that downloads SPDX licenses from the
// SPDX Workgroup website.
type SpdxOrgProvider struct {
	// https://raw.githubusercontent.com/spdx/license-list-data/master/json/licenses.json
	// use detailsUrl to get the license text
	dowloader   *download.Downloader
	licensesURL string
	log         *zap.SugaredLogger
}

// NewSpdxOrgProvider creates a new SpdxOrgProvider.
func NewSpdxOrgProvider(userAgent string) *SpdxOrgProvider {
	log := log.NewLogger("prov-spdxorg")
	requester := request.NewHTTPRequester(http.DefaultClient).SetLogger(log)
	downloader := download.NewDownloaderWithRequester(requester).SetUserAgent(userAgent)
	return &SpdxOrgProvider{
		dowloader:   downloader,
		licensesURL: "https://raw.githubusercontent.com/spdx/license-list-data/master/json/licenses.json",
		log:         log,
	}
}

// GetLicense returns the license with the given id.
func (p *SpdxOrgProvider) GetLicense(ctx context.Context, _, spdxID *string) (*models.License, error) {
	p.log.Debug("downloading base license information")
	licenses, err := p.getBaseLicenses(ctx)
	if err != nil {
		return nil, err
	}

	for _, l := range licenses {
		if l.Xid == spdxID {
			if l.DetailsURL != nil && *l.DetailsURL == "" {
				return l, nil
			}
			text, textHTML, err := p.getLicenseText(ctx, *l.DetailsURL)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to get license text for %s", l.Xid)
			}
			l.Text = &text
			l.TextHTML = &textHTML

			return l, nil
		}
	}

	return nil, errors.New("license not found")
}

// GetAllLicenses returns a list of all licenses
func (p *SpdxOrgProvider) GetAllLicenses(ctx context.Context) ([]*models.License, error) {
	p.log.Debug("downloading base license information")
	incompleteLicenses, err := p.getBaseLicenses(ctx)
	if err != nil {
		return nil, err
	}
	licenses := make([]*models.License, 0, len(incompleteLicenses))

	// download license texts
	for _, l := range incompleteLicenses {
		if l.DetailsURL == nil || *l.DetailsURL == "" {
			continue
		}
		p.log.Debugw("downloading license text", "spdxId", l.Xid)
		text, textHTML, err := p.getLicenseText(ctx, *l.DetailsURL)
		if err != nil {
			p.log.Errorw("failed to get license text", "spdxId", l.Xid)
			continue
		}
		l.Text = &text
		l.TextHTML = &textHTML
		licenses = append(licenses, l)
	}

	return licenses, nil
}

// getLicenseText returns the license text for the given url.
func (p *SpdxOrgProvider) getLicenseText(ctx context.Context, url string) (text string, textHTML string, err error) {
	if url == "" {
		return "", "", nil
	}

	// download the license details file
	detailsContent, err := p.dowloader.DownloadContent(ctx, url)
	if err != nil {
		return "", "", errors.Wrap(err, "failed to download")
	}

	// parse license details
	var details spdxLicenseDetails
	err = json.Unmarshal(detailsContent, &details)
	if err != nil {
		return "", "", errors.Wrapf(err, "failed to parse content, content was: %s", stringutil.Ellipses(strings.ReplaceAll(string(detailsContent), "\n", "\\n"), 60))
	}

	return details.LicenseText, details.LicenseTextHTML, nil
}

// getLicenses returns a list of all licenses without license text.
func (p *SpdxOrgProvider) getBaseLicenses(ctx context.Context) ([]*models.License, error) {
	// download the license list
	lcnt, err := p.dowloader.DownloadContent(ctx, p.licensesURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to download license list")
	}

	// parse licenses
	var licenseFile spdxLicenseFile
	err = json.Unmarshal(lcnt, &licenseFile)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse license file")
	}
	licenses := make([]*models.License, 0, len(licenseFile.Licenses))
	isSPDX := true
	lt := dgclient.LicenseTypeUnknown
	for i := 0; i < len(licenseFile.Licenses); i++ {
		l := licenseFile.Licenses[i]
		licenses = append(licenses, &models.License{
			Xid:           &l.LicenseID,
			Name:          &l.Name,
			ReferenceURL:  &l.Reference,
			DetailsURL:    &l.DetailsURL,
			IsSpdx:        &isSPDX,
			IsDeprecated:  &l.IsDeprecatedLicenseID,
			IsOsiApproved: &l.IsOSIApproved,
			IsFsfLibre:    &l.IsFSFLibre,
			Type:          &lt,
		})
	}

	// enrich with details
	lcnte, err := assets.ReadFile("spdx-licenses-extra.json")
	if err != nil {
		panic(err)
	}
	var licenseExtraFile spdxLicenseFile
	err = json.Unmarshal(lcnte, &licenseExtraFile)
	if err != nil {
		panic(err)
	}
	licensesMap := make(map[string]*models.License, len(licenses))
	for i := 0; i < len(licenses); i++ {
		// copy ptr of license into map
		licensesMap[*licenses[i].Xid] = licenses[i]
	}
	for _, lraw := range licenseExtraFile.Licenses {
		l, ok := licensesMap[lraw.LicenseID]
		if !ok {
			continue
		}
		isBlocked := lraw.IsBlocked
		l.IsBlocked = &isBlocked
		lt := models.AsLicenseType(lraw.Type)
		l.Type = &lt
	}

	return licenses, nil
}
