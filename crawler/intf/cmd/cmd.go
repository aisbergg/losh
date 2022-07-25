package cmd

import (
	"losh/crawler/core/config"
	"losh/internal/infra/dgraph"
	"losh/internal/lib/log"
	"losh/internal/lib/util/configutil"
	"strings"

	"github.com/aisbergg/go-errors/pkg/errors"
	"github.com/aisbergg/go-pathlib/pkg/pathlib"
)

func initConfigAndDatabase(cfgPth string) (config.Config, *dgraph.DgraphRepository, error) {
	// configuration
	path := pathlib.NewPath(strings.TrimSpace(cfgPth))
	cfg := config.DefaultConfig()
	err := configutil.Load(path, &cfg)
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
