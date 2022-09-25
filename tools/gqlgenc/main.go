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
	"context"
	"flag"
	"fmt"
	"os"

	"losh/tools/gqlgenc/clientgenlosh"

	"github.com/99designs/gqlgen/api"
	"github.com/Yamashou/gqlgenc/clientgen"
	"github.com/Yamashou/gqlgenc/config"
	"github.com/Yamashou/gqlgenc/generator"
)

func main() {
	workDir := flag.String("d", "", "Switch to the specified working directory")
	flag.Parse()
	if *workDir != "" {
		if err := os.Chdir(*workDir); err != nil {
			fmt.Fprintf(os.Stderr, "failed to change working directory: %v\n", err)
			os.Exit(1)
		}
	}

	// load config file with default or custom name
	var cfg *config.Config
	var err error
	if flag.NArg() > 0 {
		cfgPath := flag.Arg(0)
		cfg, err = config.LoadConfig(cfgPath)
	} else {
		cfg, err = config.LoadConfigFromDefaultLocations()
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v", err.Error())
		os.Exit(2)
	}

	clientGen := api.AddPlugin(clientgen.New(cfg.Query, cfg.Client, cfg.Generate))
	if cfg.Generate != nil {
		if cfg.Generate.ClientV2 {
			clientGen = api.AddPlugin(clientgenlosh.New(cfg.Query, cfg.Client, cfg.Generate))
		}
	}

	ctx := context.Background()
	if err := generator.Generate(ctx, cfg, clientGen); err != nil {
		fmt.Fprintf(os.Stderr, "%+v", err.Error())
		os.Exit(4)
	}
}
