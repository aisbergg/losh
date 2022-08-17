package cmd

import (
	"context"
	"fmt"
	"time"

	"losh/internal/core/product/models"
	"losh/internal/core/product/services"
	"losh/internal/infra/dgraph/dgclient"
	"losh/internal/lib/log"

	"github.com/aisbergg/go-errors/pkg/errors"
	"github.com/gookit/gcli/v3"
)

// DevUploadTestDataCommand is the CLI command to upload some test data to the
// database.
var DevUploadTestDataCommand = &gcli.Command{
	Name: "upload-test-data",
	Desc: "Upload test data to the database",
	Func: func(cmd *gcli.Command, args []string) error {
		_, db, err := initConfigAndDatabase(devOptions.ConfigPath)
		if err != nil {
			return err
		}

		log := log.NewLogger("cmd")
		log.Info("uploading test data now")

		// load licenses
		svc := services.NewService(db)
		err = svc.ReloadLicenseCache()
		if err != nil {
			return errors.Wrap(err, "failed to load licenses")
		}

		prd := getTestData()
		err = svc.SaveNode(context.Background(), prd)
		if err != nil {
			return errors.Wrap(err, "failed to save product")
		}

		log.Info("successfully uploaded test data")

		return nil
	},
}

func getTestData() *models.Product {
	host := &models.Host{
		Domain: p("example.org"),
		Name:   p("Example"),
	}
	licensor := &models.User{
		Xid:   p("example.org/johndoe"),
		Name:  p("John Doe"),
		Email: p("john.doe@example.org"),
		Host:  host,
	}
	license := &models.License{
		Xid: p("MIT"),
	}
	compSrc := &models.Repository{
		Xid:       p("example.org/johndoe/test-product/v1.0.0/okh.yml"),
		URL:       p("https://example.org/johndoe/test-product"),
		PermaURL:  p("https://example.org/johndoe/test-product"),
		Host:      host,
		Owner:     licensor,
		Name:      p("test-product"),
		Reference: p("v1.0.0"),
		Path:      p("okh.yml"),
	}
	release := &models.Component{
		DiscoveredAt:                p(time.Now()),
		LastIndexedAt:               p(time.Now()),
		DataSource:                  compSrc,
		Xid:                         compSrc.Xid,
		Name:                        p("Test Product"),
		Description:                 p("This is a test product"),
		Version:                     p("1.0.0"),
		CreatedAt:                   p(time.Now()),
		IsLatest:                    p(true),
		Repository:                  compSrc,
		License:                     license,
		TechnologyReadinessLevel:    p(dgclient.TechnologyReadinessLevelUndetermined),
		DocumentationReadinessLevel: p(dgclient.DocumentationReadinessLevelUndetermined),
		Licensor:                    licensor,
		DocumentationLanguage:       p("en"),
	}
	release.Releases = []*models.Component{release}

	forks := []*models.Product{}
	if forks == nil {
		fmt.Println("No forks")
	}

	product := &models.Product{
		DiscoveredAt:  p(time.Now()),
		LastIndexedAt: p(time.Now()),
		DataSource:    compSrc,

		Name:                  release.Name,
		Xid:                   p("example.org/johndoe/test-product"),
		Description:           release.Description,
		DocumentationLanguage: release.DocumentationLanguage,
		Version:               release.Version,
		License:               release.License,
		Licensor:              licensor,
		Release:               release,
		Releases:              release.Releases,
	}

	licensor.Products = []*models.Product{product}

	return product
}

func p[T any](v T) *T {
	return &v
}
