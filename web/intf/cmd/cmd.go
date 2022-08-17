package cmd

import (
	"losh/internal/infra/dgraph"
	"losh/internal/lib/log"
	"losh/web/core/config"

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
	if !db.IsReachable() {
		return config.Config{}, nil, errors.New("failed to connect to Dgraph database")
	}

	return cfg, db, nil
}
