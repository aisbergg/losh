package cmd

import "github.com/gookit/gcli/v3"

var devOptions = struct {
	Path string
}{}

// DevCommand is the CLI command to run development tasks.
var DevCommand = &gcli.Command{
	Name: "dev",
	Desc: "Development tasks",
	Config: func(c *gcli.Command) {
		c.StrOpt(&devOptions.Path, "config", "c", "", "configuration file path")
	},
	Subs: []*gcli.Command{
		DevCrawlProductCommand,
		DevUploadTestDataCommand,
	},
}
