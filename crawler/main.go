package main

import (
	"fmt"
	"os"

	"losh/crawler/cmd"
	"losh/internal/errors"

	"github.com/gookit/gcli/v3"
)

func main() {
	// create the CLI application
	app := gcli.NewApp()
	app.Version = "0.0.1"
	app.Desc = "LOSH Crawler"

	// register error handler
	app.On(gcli.EvtAppRunError, errorHandler)

	// register commands
	app.Add(cmd.ConfigCommand)
	app.Add(cmd.ManageCommand)

	// run the application
	os.Exit(app.Run(nil))
}

func errorHandler(data ...interface{}) (stop bool) {
	if len(data) == 2 && data[1] != nil {
		if err, ok := data[1].(error); ok {
			fmt.Fprintln(os.Stderr, errors.FormatColorfulCLIMessage(err))
		}
	}
	return
}
