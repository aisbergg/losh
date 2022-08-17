package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/aisbergg/go-errors/pkg/errors"
	"github.com/gookit/event"
	"github.com/gookit/gcli/v3"

	"losh/internal/lib/log"
	loshapp "losh/web/intf/http"
)

var runOptions = struct {
	Path string
}{}

// RunCommand is the CLI command to run the web application.
var RunCommand = &gcli.Command{
	Name: "run",
	Desc: "Run the application",
	Config: func(c *gcli.Command) {
		c.StrOpt(&runOptions.Path, "config", "c", "", "configuration file path")
	},
	Func: func(cmd *gcli.Command, args []string) error {
		cfg, db, err := initConfigAndDatabase(runOptions.Path)
		if err != nil {
			return err
		}

		// flush logs and close log file after server shutdown
		event.On("server.stop", event.ListenerFunc(func(e event.Event) error {
			return log.Close()
		}))
		// rotate logs on usr1 signal
		event.On("signal.us1", event.ListenerFunc(func(e event.Event) error {
			return log.RotateLogFile()
		}))

		// server
		if err, _ = event.Fire("server.initialize", nil); err != nil {
			return err
		}
		server, err := loshapp.NewServer(&cfg, db)
		if err != nil {
			return errors.Wrap(err, "failed to initialize server")
		}

		var shutdownErr error
		done := make(chan struct{})
		defer close(done)
		sigChannel := make(chan os.Signal, 1)
		signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1, syscall.SIGUSR2)
		go func() {
			for {
				select {
				case sig := <-sigChannel:
					switch sig {
					case syscall.SIGINT, syscall.SIGTERM:
						event.Fire("signal.term", nil)
						// shutdown server on term and int signal
						shutdownErr = server.Shutdown()
					case syscall.SIGUSR1:
						event.Fire("signal.usr1", nil)
					case syscall.SIGUSR2:
						event.Fire("signal.usr2", nil)
					}
				case <-done:
					return
				}
			}
		}()

		// start listening on the specified address
		listenOn := fmt.Sprintf("%s:%d", cfg.Server.Interface, cfg.Server.Port)

		if err := server.Listen(listenOn); err != nil {
			shutdownErr = server.Shutdown()
			log.Close()
			return err
		}

		err = log.Close()
		if shutdownErr != nil {
			return errors.Wrap(err, "failed to shutdown server")
		}

		return err
	},
}
