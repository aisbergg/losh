package cmd

import (
	"encoding/json"
	"losh/internal/util/configutil"
	"losh/web/config"
	"os"
	"strings"

	"github.com/aisbergg/go-errors/pkg/errors"
	"github.com/aisbergg/go-pathlib/pkg/pathlib"
	"github.com/gookit/gcli/v3"
)

var configShowOptions = struct {
	Path string
}{}

var ConfigShowCommand = &gcli.Command{
	Name:    "show",
	Desc:    "Show the effective configuration",
	Aliases: []string{"shw", "s"},
	Config: func(c *gcli.Command) {
		c.StrOpt(&configShowOptions.Path, "config", "c", "", "configuration file path")
	},
	Func: func(cmd *gcli.Command, args []string) error {
		path := pathlib.NewPath(strings.TrimSpace(configShowOptions.Path))
		config := config.DefaultConfig()
		err := configutil.Load(path, &config)
		if err != nil {
			return errors.Wrap(err, "failed to load configuration")
		}
		b, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			return errors.Wrap(err, "failed to marshal configuration")
		}
		os.Stdout.Write(b)
		return nil
	},
}
