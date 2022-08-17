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
