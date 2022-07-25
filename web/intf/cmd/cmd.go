package cmd

import (
	"losh/crawler/core/config"
	"losh/internal/infra/dgraph"
	"losh/internal/lib/log"

	"github.com/aisbergg/go-errors/pkg/errors"
)

func initConfigAndDatabase(cfgPth string) (config.Config, *dgraph.DgraphRepository, error) {
	// configuration
	cfgSvc := config.NewService(configInitOptions.Output)
	cfg, err := cfgSvc.Get()
	if err != nil {
		return config.Config{}, nil, errors.Wrap(err, "failed to load configuration")
	}

	// logging
	err = log.Initialize(cfg.Log)
	if err != nil {
		return config.Config{}, nil, errors.Wrap(err, "failed to initialize logging")
	}

	// database
	db, err := dgraph.NewDgraphRepository(cfg.Database)
	if err != nil {
		return config.Config{}, nil, errors.Wrap(err, "failed to initialize Dgraph database connection")
	}
	if !db.IsReachable() {
		return config.Config{}, nil, errors.New("failed to connect to Dgraph database")
	}

	return cfg, db, nil
}
