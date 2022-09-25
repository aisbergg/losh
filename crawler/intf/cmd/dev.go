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

import "github.com/gookit/gcli/v3"

var devOptions = struct {
	ConfigPath string
}{}

// DevCommand is the CLI command to run development tasks.
var DevCommand = &gcli.Command{
	Name: "dev",
	Desc: "Development tasks",
	Config: func(c *gcli.Command) {
		c.StrOpt(&devOptions.ConfigPath, "config", "c", "", "configuration file path")
	},
	Subs: []*gcli.Command{
		DevCrawlProductCommand,
		DevDownloadLicensesCommand,
		DevUploadFile,
		DevUploadLicensesCommand,
		DevUploadTestDataCommand,
	},
}
