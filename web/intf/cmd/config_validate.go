package cmd

import (
	"losh/web/core/config"

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
