package cmd

import (
	"context"
	"fmt"

	"losh/crawler/core/wikifactory"
	"losh/internal/core/product/models"
	"losh/internal/license"

	"github.com/aisbergg/go-errors/pkg/errors"
	"github.com/gookit/gcli/v3"
)

// DevCrawlProductCommand is the CLI command to crawl a bunch of products by
// URL.
var DevCrawlProductCommand = &gcli.Command{
	Name: "crawl-product",
	Desc: "Crawl a product from a given URL",
	Config: func(c *gcli.Command) {
		c.AddArg("url", "URLs to crawl", true, true)
	},
	Func: func(cmd *gcli.Command, args []string) error {
		urls := cmd.Arg("url").Strings()
		prdIDs := make([]models.ProductID, 0, len(urls))
		for _, url := range urls {
			prdID, err := models.NewProductIDFromURL(url)
			if err != nil {
				return errors.Wrap(err, "invalid or unsupported URL")
			}
			prdIDs = append(prdIDs, prdID)
		}

		cfg, db, err := initConfigAndDatabase(devOptions.Path)
		if err != nil {
			return err
		}

		// load licenses
		lcsC := license.NewLicenseCache(db)
		err = lcsC.Reload()
		if err != nil {
			return errors.Wrap(err, "failed to load licenses")
		}

		// setup crawler
		crw := wikifactory.NewWikifactoryCrawler(db, lcsC, cfg.Crawler.UserAgent)

		// crawl products
		for _, prdID := range prdIDs {
			fmt.Println("Crawling product:", prdID)

			ctx := context.Background()
			prd, err := crw.GetProduct(ctx, prdID)
			if err != nil {
				return errors.Wrap(err, "failed to get product")
			}

			// prdJ, err := json.MarshalIndent(prd, "", "  ")
			// if err != nil {
			// 	return errors.Wrap(err, "failed to marshal product")
			// }
			// fmt.Println(string(prdJ))

			// save product
			err = db.SaveProduct(prd)
			if err != nil {
				return errors.Wrap(err, "failed to save product")
			}
		}

		return nil
	},
}
