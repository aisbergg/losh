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

package main

import (
	"fmt"
	"os"

	"losh/crawler/intf/cmd"
	"losh/internal/lib/errors"

	"github.com/gookit/gcli/v3"
)

func main() {
	// create the CLI application
	app := gcli.NewApp()
	app.Version = "0.0.1"
	app.Desc = "LOSH Crawler"

	// register error handler
	app.On(gcli.EvtAppRunError, errorHandler)

	// register commands
	app.Add(cmd.DevCommand)
	app.Add(cmd.DiscoverCommand)
	app.Add(cmd.ConfigCommand)
	app.Add(cmd.ManageCommand)

	// run the application
	os.Exit(app.Run(nil))
}

func errorHandler(data ...interface{}) (stop bool) {
	if len(data) == 2 && data[1] != nil {
		if err, ok := data[1].(error); ok {
			fmt.Fprintln(os.Stderr, errors.FormatColorfulCLIMessage(err))
		}
	}
	return
}
