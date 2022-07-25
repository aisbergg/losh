package cmd

import (
	"github.com/aisbergg/go-errors/pkg/errors"
	"github.com/gookit/gcli/v3"
)

var ConfigValidateCommand = &gcli.Command{
	Name:    "validate",
	Desc:    "Validate the configuration",
	Aliases: []string{"vld", "v"},
	Func: func(cmd *gcli.Command, args []string) error {
		return errors.New("not implemented")
	},
}
