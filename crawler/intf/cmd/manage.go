package cmd

import "github.com/gookit/gcli/v3"

var manageOptions = struct {
	Path string
}{}

// ManageCommand is the CLI command to run management tasks.
var ManageCommand = &gcli.Command{
	Name: "manage",
	Desc: "Management tasks",
	Config: func(c *gcli.Command) {
		c.StrOpt(&manageOptions.Path, "config", "c", "", "configuration file path")
	},
	Subs: []*gcli.Command{
		ManageUpdateLicensesCommand,
	},
	Aliases: []string{"mng", "m"},
}
