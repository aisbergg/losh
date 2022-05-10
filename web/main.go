package main

import (
	loshApp "losh/web/app"
	loshConfig "losh/web/config"

	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
)

type App struct {
	*fiber.App
}

func main() {
	// configuration
	config := loshConfig.New()
	app := loshApp.NewApp(config)

	// close any connections on interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		app.Shutdown()
	}()

	// start listening on the specified address
	if err := app.Listen(config.GetString("APP_ADDR")); err != nil {
		app.Shutdown()
	}
}
