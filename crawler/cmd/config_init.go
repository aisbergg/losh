package cmd

import (
	"io"
	"losh/web/config"
	"os"
	"strings"

	"github.com/aisbergg/go-errors/pkg/errors"
	"github.com/aisbergg/go-pathlib/pkg/pathlib"
	"github.com/gookit/gcli/v3"
	"gopkg.in/yaml.v3"
)

var configInitOptions = struct {
	Output string
}{}

var ConfigInitCommand = &gcli.Command{
	Name:    "init",
	Desc:    "Create configuration with default values",
	Aliases: []string{"ini", "i"},
	Config: func(c *gcli.Command) {
		c.StrOpt(&configInitOptions.Output, "output", "o", "", "Output file path; if non given, will print to stdout")
	},
	Func: func(cmd *gcli.Command, args []string) error {
		// create file to write to
		path := pathlib.NewPath(strings.TrimSpace(configInitOptions.Output))
		var writeTo io.Writer
		if path.String() != "." {
			exists, err := path.Exists()
			if err != nil {
				return err
			}
			if exists {
				isFile, err := path.IsFile()
				if err != nil {
					return err
				}
				if !isFile {
					return errors.Errorf("path is not a file: %s", path.String())
				}
			}
			file, err := path.OpenFile(os.O_RDWR | os.O_CREATE)
			if err != nil {
				return err
			}
			defer file.Close()
			writeTo = file

		} else {
			// write to stdout
			writeTo = os.Stdout
		}

		// encode as yaml
		yamlEncoder := yaml.NewEncoder(writeTo)
		yamlEncoder.SetIndent(2)
		if err := yamlEncoder.Encode(config.DefaultConfig); err != nil {
			return err
		}

		return nil
	},
}
