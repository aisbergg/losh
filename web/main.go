package main

import (
	"fmt"
	"os"

	"losh/internal/lib/errors"
	"losh/web/intf/cmd"

	"github.com/gookit/gcli/v3"
)

func main() {
	// create the CLI application
	app := gcli.NewApp()
	app.Version = "0.0.1"
	app.Desc = "LOSH Web"

	// register error handler
	app.On(gcli.EvtAppRunError, errorHandler)

	// register commands
	app.Add(cmd.RunCommand)
	app.Add(cmd.SearchCommand)
	app.Add(cmd.ConfigCommand)

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
