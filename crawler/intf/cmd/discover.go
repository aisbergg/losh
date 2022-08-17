package cmd

import (
	"context"

	"losh/crawler/core/wikifactory"
	"losh/internal/core/product/services"
	"losh/internal/lib/log"

	"github.com/aisbergg/go-errors/pkg/errors"
	"github.com/gookit/gcli/v3"
)

var discoverOptions = struct {
	ConfigPath string
}{}

// DiscoverCommand is the CLI command to discover products and save them to the database.
var DiscoverCommand = &gcli.Command{
	Name: "discover",
	Desc: "Discover products and save them to the database",
	Config: func(c *gcli.Command) {
		c.StrOpt(&discoverOptions.ConfigPath, "config", "c", "", "configuration file path")
	},
	Func: func(cmd *gcli.Command, args []string) error {
		cfg, db, err := initConfigAndDatabase(discoverOptions.ConfigPath)
		if err != nil {
			return err
		}
		log := log.NewLogger("cmd")

		// setup crawler
		svc := services.NewService(db)
		svc.ReloadLicenseCache()
		crwl := wikifactory.NewWikifactoryCrawler(svc, cfg.Crawler.UserAgent)

		// discover products
		err = crwl.DiscoverProducts(context.Background())
		if err != nil {
			return errors.Wrap(err, "failed to discover products")
		}

		log.Info("successfully discovered products")
		return nil
	},
}
