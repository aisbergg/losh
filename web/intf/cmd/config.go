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
	"github.com/gookit/gcli/v3"
)

// ConfigCommand is the CLI command to run configuration tasks.
var ConfigCommand = &gcli.Command{
	Name: "config",
	Desc: "Initialize, validate and show configuration",
	Subs: []*gcli.Command{
		ConfigInitCommand,
		ConfigShowCommand,
		ConfigValidateCommand,
	},
	Aliases: []string{"cfg", "c"},
}
