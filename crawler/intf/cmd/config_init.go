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
	"losh/crawler/core/config"

	"github.com/aisbergg/go-errors/pkg/errors"
	"github.com/gookit/gcli/v3"
)

var configInitOptions = struct {
	Output string
}{}

// ConfigInitCommand is the CLI command to initialize a new default
// configuration.
var ConfigInitCommand = &gcli.Command{
	Name:    "init",
	Desc:    "Create configuration with default values",
	Aliases: []string{"ini", "i"},
	Config: func(c *gcli.Command) {
		c.StrOpt(&configInitOptions.Output, "output", "o", "", "Output file path; if non given, will print to stdout")
	},
	Func: func(cmd *gcli.Command, args []string) error {
		cfgSvc := config.NewService(configInitOptions.Output)
		err := cfgSvc.Init()
		if err != nil {
			return errors.Wrap(err, "failed to initialize configuration")
		}
		return nil
	},
}
