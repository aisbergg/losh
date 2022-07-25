package cmd

import (
	"losh/web/core/config"

	"github.com/aisbergg/go-errors/pkg/errors"
	"github.com/gookit/gcli/v3"
)

var configInitOptions = struct {
	Output string
}{}

// ConfigInitCommand is the CLI command to initialize a new default
// configuration.
var ConfigInitCommand = &gcli.Command{
	Name:    "init",
	Desc:    "Create configuration with default values",
	Aliases: []string{"ini", "i"},
	Config: func(c *gcli.Command) {
		c.StrOpt(&configInitOptions.Output, "output", "o", "", "Output file path; if non given, will print to stdout")
	},
	Func: func(cmd *gcli.Command, args []string) error {
		cfgSvc := config.NewService(configInitOptions.Output)
		err := cfgSvc.Init()
		if err != nil {
			return errors.Wrap(err, "failed to initialize configuration")
		}
		return nil
	},
}
