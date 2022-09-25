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
	"losh/internal/infra/dgraph"
	"losh/internal/lib/log"

	"github.com/aisbergg/go-errors/pkg/errors"
)

func initConfig(cfgPth string) (config.Config, error) {
	// configuration
	cfgSvc := config.NewService(cfgPth)
	cfg, err := cfgSvc.Get()
	if err != nil {
		return config.Config{}, errors.Wrap(err, "failed to load configuration")
	}

	// logging
	err = log.Initialize(cfg.Log)
	if err != nil {
		return config.Config{}, errors.Wrap(err, "failed to initialize logging")
	}

	return cfg, nil
}

func initConfigAndDatabase(cfgPth string) (config.Config, *dgraph.DgraphRepository, error) {
	cfg, err := initConfig(cfgPth)
	if err != nil {
		return config.Config{}, nil, err
	}

	// database
	db, err := dgraph.NewDgraphRepository(cfg.Database)
	if err != nil {
		return config.Config{}, nil, errors.Wrap(err, "failed to initialize Dgraph database connection")
	}
	if err = db.WaitUntilReachable(); err != nil {
		return config.Config{}, nil, errors.New("failed to connect to Dgraph database")
	}

	return cfg, db, nil
}
