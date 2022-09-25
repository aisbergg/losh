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
	"encoding/json"
	"os"

	"losh/web/core/config"

	"github.com/aisbergg/go-errors/pkg/errors"
	"github.com/gookit/gcli/v3"
)

var configShowOptions = struct {
	Path string
}{}

// ConfigShowCommand is the CLI command to show the effective configuration.
var ConfigShowCommand = &gcli.Command{
	Name:    "show",
	Desc:    "Show the effective configuration",
	Aliases: []string{"shw", "s"},
	Config: func(c *gcli.Command) {
		c.StrOpt(&configShowOptions.Path, "config", "c", "", "configuration file path")
	},
	Func: func(cmd *gcli.Command, args []string) error {
		// load configuration
		cfgSvc := config.NewService(configShowOptions.Path)
		cfg, err := cfgSvc.Get()
		if err != nil {
			return errors.Wrap(err, "failed to load configuration")
		}

		// serialize to JSON
		b, err := json.MarshalIndent(cfg, "", "  ")
		if err != nil {
			return errors.Wrap(err, "failed to marshal configuration")
		}
		os.Stdout.Write(b)
		return nil
	},
}
