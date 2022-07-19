package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/aisbergg/go-errors/pkg/errors"
	"github.com/aisbergg/go-pathlib/pkg/pathlib"
	"github.com/gookit/event"
	"github.com/gookit/gcli/v3"

	"losh/internal/logging"
	"losh/internal/repository/dgraph"
	"losh/internal/util/configutil"
	loshapp "losh/web/app"
	"losh/web/config"
)

var runOptions = struct {
	Path string
}{}

var RunCommand = &gcli.Command{
	Name: "run",
	Desc: "Run the application",
	Config: func(c *gcli.Command) {
		c.StrOpt(&configShowOptions.Path, "config", "c", "", "configuration file path")
	},
	Func: func(cmd *gcli.Command, args []string) error {
		// configuration
		path := pathlib.NewPath(strings.TrimSpace(configShowOptions.Path))
		config := config.DefaultConfig()
		err := configutil.Load(path, &config)
		if err != nil {
			return errors.Wrap(err, "failed to load configuration")
		}

		// logging
		err = logging.Initialize(config.Log)
		if err != nil {
			return errors.Wrap(err, "failed to initialize logging")
		}
		// flush logs and close log file after server shutdown
		event.On("server.stop", event.ListenerFunc(func(e event.Event) error {
			return logging.Close()
		}))
		// rotate logs on usr1 signal
		event.On("signal.us1", event.ListenerFunc(func(e event.Event) error {
			return logging.RotateLogFile()
		}))

		// database
		db, err := dgraph.NewDgraphRepository(config.Database)
		if err != nil {
			return errors.Wrap(err, "failed to initialize Dgraph database connection")
		}
		if !db.IsReachable() {
			return errors.New("failed to connect to Dgraph database")
		}

		// server
		if err, _ = event.Fire("server.initialize", nil); err != nil {
			return err
		}
		server, err := loshapp.NewServer(&config, db)
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
		listenOn := fmt.Sprintf("%s:%d", config.Server.Interface, config.Server.Port)

		if err := server.Listen(listenOn); err != nil {
			shutdownErr = server.Shutdown()
			logging.Close()
			return err
		}

		err = logging.Close()
		if shutdownErr != nil {
			return errors.Wrap(err, "failed to shutdown server")
		}

		return err
	},
}
