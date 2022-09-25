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
