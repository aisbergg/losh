package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"losh/crawler/core/wikifactory"
	"losh/internal/core/product/models"
	"losh/internal/core/product/services"

	"github.com/aisbergg/go-errors/pkg/errors"
	"github.com/aisbergg/go-pathlib/pkg/pathlib"
	"github.com/gookit/gcli/v3"
)

var devCrawlProductOptions = struct {
	OutputPath string
}{}

// DevCrawlProductCommand is the CLI command to crawl a bunch of products by
// URL.
var DevCrawlProductCommand = &gcli.Command{
	Name: "crawl-product",
	Desc: "Crawl a product from a given URL",
	Config: func(c *gcli.Command) {
		c.AddArg("url", "URLs to crawl", true, true)
		c.StrOpt(&devCrawlProductOptions.OutputPath, "output", "o", "", "output file path; if defined, the output will be written to this file as json instead of saving to the database")
	},
	Func: func(cmd *gcli.Command, args []string) error {
		urls := cmd.Arg("url").Strings()
		outPath := pathlib.NewPath(devCrawlProductOptions.OutputPath)
		prdIDs := make([]models.ProductID, 0, len(urls))
		for _, url := range urls {
			prdID, err := models.NewProductIDFromURL(url)
			if err != nil {
				return errors.Wrap(err, "invalid or unsupported URL")
			}
			prdIDs = append(prdIDs, prdID)
		}

		cfg, db, err := initConfigAndDatabase(devOptions.ConfigPath)
		if err != nil {
			return err
		}

		// load licenses
		svc := services.NewService(db)
		err = svc.ReloadLicenseCache()
		if err != nil {
			return errors.Wrap(err, "failed to load licenses")
		}

		// setup crawler
		crw := wikifactory.NewWikifactoryCrawler(svc, cfg.Crawler.UserAgent)

		// crawl products
		prds := make([]models.Node, 0, len(prdIDs))
		ctx := context.Background()
		for _, prdID := range prdIDs {
			fmt.Println("Crawling product:", prdID)

			prd, err := crw.GetProduct(ctx, prdID)
			if err != nil {
				return errors.Wrap(err, "failed to get product")
			}
			prds = append(prds, prd)
		}

		if outPath.String() != "." {
			// save to json file
			for i := 0; i < len(prds); i++ {
				// turn into tree structure to remove circular references
				prds[i] = models.AsTree(prds[i])
			}
			js, err := json.MarshalIndent(prds, "", "  ")
			err = outPath.WriteFile(js)
			if err != nil {
				return errors.Wrap(err, "failed to write products to file")
			}

		} else {
			// save to database
			err = svc.SaveNodes(ctx, prds)
			if err != nil {
				return errors.Wrap(err, "failed to save product")
			}
		}

		return nil
	},
}
