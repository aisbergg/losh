package cmd

import (
	"encoding/json"
	"fmt"
	"losh/internal/core/product/models"
	"losh/internal/license"
	"time"

	"github.com/aisbergg/go-copier/pkg/copier"
	"github.com/aisbergg/go-errors/pkg/errors"
	"github.com/gookit/gcli/v3"
)

var DevUploadTestDataCommand = &gcli.Command{
	Name: "upload-test-data",
	Desc: "Upload test data to the database",
	Func: func(cmd *gcli.Command, args []string) error {
		_, db, err := initConfigAndDatabase(devOptions.Path)
		if err != nil {
			return err
		}

		// load licenses
		lcsC := license.NewLicenseCache(db)
		err = lcsC.Reload()
		if err != nil {
			return errors.Wrap(err, "failed to load licenses")
		}

		tstPrd := getTestData()
		js, err := json.MarshalIndent(tstPrd, "", "  ")
		if err != nil {
			return errors.Wrap(err, "failed to marshal test product")
		}
		fmt.Println(string(js))

		cpr := copier.New(copier.Options{AutoConvert: true, IgnoreEmpty: true})
		addPrd := &models.AddProductInput{}
		if err := cpr.CopyTo(tstPrd, addPrd); err != nil {
			return errors.Wrap(err, "failed to copy test product")
		}
		js, err = json.MarshalIndent(addPrd, "", "  ")
		if err != nil {
			return errors.Wrap(err, "failed to marshal test product")
		}
		fmt.Println()
		fmt.Println()
		fmt.Println(string(js))

		err = db.SaveProduct(tstPrd)
		if err != nil {
			return errors.Wrap(err, "failed to save product")
		}

		return nil
	},
}

func getTestData() *models.Product {
	host := models.Host{
		Domain: "example.org",
		Name:   "Example",
	}
	owner := models.User{
		Xid:   "example.org/johndoe",
		Name:  "John Doe",
		Email: sPtr("john.doe@example.org"),
		Host:  host,
	}
	license := models.License{
		Xid: "MIT",
	}
	compSrc := models.Repository{
		Xid:      "example.org/johndoe/test-product/v1.0.0/okh.yml",
		URL:      "https://example.org/johndoe/test-product",
		PermaURL: "https://example.org/johndoe/test-product",
		Host:     host,
		Owner:    &owner,
		Name:     sPtr("test-product"),
		Tag:      sPtr("v1.0.0"),
		Path:     sPtr("okh.yml"),
	}
	release := models.Component{
		DiscoveredAt:  time.Now(),
		LastIndexedAt: time.Now(),
		DataSource:    compSrc,

		Xid:                   compSrc.Xid,
		Name:                  "Test Product",
		Description:           "This is a test product",
		Owner:                 &owner,
		Version:               "1.0.0",
		CreatedAt:             time.Now(),
		Releases:              []*models.Component{},
		IsLatest:              true,
		Repository:            compSrc,
		License:               license,
		Licensor:              &owner,
		DocumentationLanguage: "en",
	}

	forks := []*models.Product{}
	if forks == nil {
		fmt.Println("No forks")
	}

	product := models.Product{
		DiscoveredAt:  time.Now(),
		LastIndexedAt: time.Now(),
		DataSource:    compSrc,

		Name:        release.Name,
		Xid:         "example.org/johndoe/test-product",
		Owner:       &owner,
		Description: release.Description,
		Version:     release.Version,
		Release:     &release,
		Releases:    []*models.Component{&release},
		Forks:       []*models.Product{},
		Tags:        []*models.Tag{},
	}

	return &product
}

func sPtr(s string) *string {
	return &s
}
