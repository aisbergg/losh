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

// ConfigValidateCommand is the CLI command to validate the configuration.
var ConfigValidateCommand = &gcli.Command{
	Name:    "validate",
	Desc:    "Validate the configuration",
	Aliases: []string{"vld", "v"},
	Func: func(cmd *gcli.Command, args []string) error {
		// load and validate configuration
		cfgSvc := config.NewService(configInitOptions.Output)
		_, err := cfgSvc.Get()
		if err != nil {
			return errors.Wrap(err, "failed to load configuration")
		}
		return nil
	},
}
