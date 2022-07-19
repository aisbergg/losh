package spdxorg

import (
	"context"
	"embed"
	"encoding/json"
	"losh/crawler/logging"
	"losh/internal/models"
	"losh/internal/net/download"
	"losh/internal/net/request"
	"net/http"

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
	LicenseText string `json:"licenseText"`
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
	log := logging.NewLogger("prov-spdxorg")
	requester := request.NewHTTPRequester(http.DefaultClient).SetLogger(log)
	downloader := download.NewDownloaderWithRequester(requester).SetUserAgent(userAgent)
	return &SpdxOrgProvider{
		dowloader:   downloader,
		licensesURL: "https://raw.githubusercontent.com/spdx/license-list-data/master/json/licenses.json",
		log:         log,
	}
}

// GetLicense returns the license with the given id.
func (p *SpdxOrgProvider) GetLicense(_, spdxID *string) (*models.License, error) {
	p.log.Debug("downloading base license information")
	licenses, err := p.getBaseLicenses()
	if err != nil {
		return nil, err
	}

	for _, l := range licenses {
		if l.Xid == *spdxID {
			if l.DetailsURL != nil && *l.DetailsURL == "" {
				return l, nil
			}
			text, err := p.getLicenseText(*l.DetailsURL)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to get license text for %s", l.Xid)
			}
			l.Text = &text

			return l, nil
		}
	}

	return nil, errors.New("license not found")
}

// GetAllLicenses returns a list of all licenses
func (p *SpdxOrgProvider) GetAllLicenses() ([]*models.License, error) {
	p.log.Debug("downloading base license information")
	licenses, err := p.getBaseLicenses()
	if err != nil {
		return nil, err
	}

	// download license texts
	for _, l := range licenses {
		if l.DetailsURL == nil || *l.DetailsURL == "" {
			continue
		}
		p.log.Debugw("downloading license text", "spdxId", l.Xid)
		text, err := p.getLicenseText(*l.DetailsURL)
		if err != nil {
			p.log.Errorw("failed to get license text", "spdxId", l.Xid)
			continue
		}
		l.Text = &text
	}

	return licenses, nil
}

// getLicenseText returns the license text for the given url.
func (p *SpdxOrgProvider) getLicenseText(url string) (string, error) {
	if url == "" {
		return "", nil
	}

	// download the license details file
	ctx := context.Background()
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

	return details.LicenseText, nil
}

// getLicenses returns a list of all licenses without license text.
func (p *SpdxOrgProvider) getBaseLicenses() ([]*models.License, error) {
	// download the license list
	ctx := context.Background()
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
	for i := 0; i < len(licenseFile.Licenses); i++ {
		l := licenseFile.Licenses[i]
		licenses = append(licenses, &models.License{
			Xid:           l.LicenseID,
			Name:          l.Name,
			ReferenceURL:  &l.Reference,
			DetailsURL:    &l.DetailsURL,
			IsSpdx:        true,
			IsDeprecated:  l.IsDeprecatedLicenseID,
			IsOsiApproved: l.IsOSIApproved,
			IsFsfLibre:    l.IsFSFLibre,
			Type:          models.LicenseTypeUnknown,
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
	licensesMap := make(map[string]*models.License, len(licenseExtraFile.Licenses))
	for i := 0; i < len(licenses); i++ {
		// copy ptr of license into map
		licensesMap[licenses[i].Xid] = licenses[i]
	}
	for _, lraw := range licenseExtraFile.Licenses {
		l, ok := licensesMap[lraw.LicenseID]
		if !ok {
			continue
		}
		l.IsBlocked = lraw.IsBlocked
		l.Type = models.AsLicenseType(lraw.Type)
	}

	return licenses, nil
}
