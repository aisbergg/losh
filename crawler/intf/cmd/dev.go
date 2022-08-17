package cmd

import "github.com/gookit/gcli/v3"

var devOptions = struct {
	ConfigPath string
}{}

// DevCommand is the CLI command to run development tasks.
var DevCommand = &gcli.Command{
	Name: "dev",
	Desc: "Development tasks",
	Config: func(c *gcli.Command) {
		c.StrOpt(&devOptions.ConfigPath, "config", "c", "", "configuration file path")
	},
	Subs: []*gcli.Command{
		DevCrawlProductCommand,
		DevDownloadLicensesCommand,
		DevUploadFile,
		DevUploadLicensesCommand,
		DevUploadTestDataCommand,
	},
}
