package cmd

import (
	"encoding/json"
	"os"

	"losh/crawler/core/config"

	"github.com/aisbergg/go-errors/pkg/errors"
	"github.com/gookit/gcli/v3"
)

var configShowOptions = struct {
	Path string
}{}

// ConfigShowCommand is the CLI command to show the effective configuration.
var ConfigShowCommand = &gcli.Command{
	Name:    "show",
	Desc:    "Show the effective configuration",
	Aliases: []string{"shw", "s"},
	Config: func(c *gcli.Command) {
		c.StrOpt(&configShowOptions.Path, "config", "c", "", "configuration file path")
	},
	Func: func(cmd *gcli.Command, args []string) error {
		// load configuration
		cfgSvc := config.NewService(configInitOptions.Output)
		cfg, err := cfgSvc.Get()
		if err != nil {
			return errors.Wrap(err, "failed to load configuration")
		}

		// serialize to JSON
		b, err := json.MarshalIndent(cfg, "", "  ")
		if err != nil {
			return errors.Wrap(err, "failed to marshal configuration")
		}
		os.Stdout.Write(b)
		return nil
	},
}
