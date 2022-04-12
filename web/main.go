package main

import (
	"losh/web/app/middleware"
	configuration "losh/web/config"
	"losh/web/database"
	"losh/web/routes"
	"time"

	"os"
	"os/signal"

	hashing "github.com/thomasvvugt/fiber-hashing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/utils"

	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/expvar"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	"github.com/gofiber/fiber/v2/middleware/session"
)

type App struct {
	*fiber.App

	DB      *database.Database
	Hasher  hashing.Driver
	Session *session.Store
}

func main() {
	config := configuration.New()

	app := App{
		App:     fiber.New(*config.GetFiberConfig()),
		Hasher:  hashing.New(config.GetHasherConfig()),
		Session: session.New(config.GetSessionConfig()),
	}

	app.registerMiddlewares(config)

	var err error
	// // Initialize database
	// db, err := database.New(&database.DatabaseConfig{
	// 	Driver:   config.GetString("DB_DRIVER"),
	// 	Host:     config.GetString("DB_HOST"),
	// 	Username: config.GetString("DB_USERNAME"),
	// 	Password: config.GetString("DB_PASSWORD"),
	// 	Port:     config.GetInt("DB_PORT"),
	// 	Database: config.GetString("DB_DATABASE"),
	// })

	// // Auto-migrate database models
	// if err != nil {
	// 	fmt.Println("failed to connect to database:", err.Error())
	// } else {
	// 	if db == nil {
	// 		fmt.Println("failed to connect to database: db variable is nil")
	// 	} else {
	// 		app.DB = db
	// 		err := app.DB.AutoMigrate(&models.Role{})
	// 		if err != nil {
	// 			fmt.Println("failed to automigrate role model:", err.Error())
	// 			return
	// 		}
	// 		err = app.DB.AutoMigrate(&models.User{})
	// 		if err != nil {
	// 			fmt.Println("failed to automigrate user model:", err.Error())
	// 			return
	// 		}
	// 	}
	// }

	// Register web routes
	web := app.Group("")
	routes.RegisterHome(web, app.Session, config.GetString("SESSION_LOOKUP"), app.DB, app.Hasher)

	// Register application API routes (using the /api/v1 group)
	// api := app.Group("/api")
	// apiv1 := api.Group("/v1")
	// routes.RegisterAPI(apiv1, app.DB)

	// Register static routes for the public directory
	app.Static("/", "./public")

	// Custom 404 Handler
	app.Use(func(c *fiber.Ctx) error {
		if err := c.SendStatus(fiber.StatusNotFound); err != nil {
			panic(err)
		}
		if err := c.Render("errors/404", fiber.Map{}); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		return err
	})

	// Close any connections on interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		app.exit()
	}()

	// Start listening on the specified address
	err = app.Listen(config.GetString("APP_ADDR"))
	if err != nil {
		app.exit()
	}
}

func (app *App) registerMiddlewares(config *configuration.Config) {

	// Custom Access Logger based on zap
	if config.GetBool("MW_ACCESS_LOGGER_ENABLED") {
		app.Use(middleware.AccessLogger(&middleware.AccessLoggerConfig{
			Type:        config.GetString("MW_ACCESS_LOGGER_TYPE"),
			Environment: config.GetString("APP_ENV"),
			Filename:    config.GetString("MW_ACCESS_LOGGER_FILENAME"),
			MaxSize:     config.GetInt("MW_ACCESS_LOGGER_MAXSIZE"),
			MaxAge:      config.GetInt("MW_ACCESS_LOGGER_MAXAGE"),
			MaxBackups:  config.GetInt("MW_ACCESS_LOGGER_MAXBACKUPS"),
			LocalTime:   config.GetBool("MW_ACCESS_LOGGER_LOCALTIME"),
			Compress:    config.GetBool("MW_ACCESS_LOGGER_COMPRESS"),
		}))
	}

	// Force HTTPS
	if config.GetBool("MW_FORCE_HTTPS_ENABLED") {
		app.Use(middleware.ForceHTTPS())
	}

	// Force trailing slash
	if config.GetBool("MW_FORCE_TRAILING_SLASH_ENABLED") {
		app.Use(middleware.RemoveTrailingSlash())
	}

	// Recover
	if config.GetBool("MW_FIBER_RECOVER_ENABLED") {
		app.Use(recover.New())
	}

	// Cache
	if config.GetBool("MW_FIBER_CACHE_ENABLED") {
		app.Use(cache.New(cache.Config{
			Expiration:   config.GetDuration("MW_FIBER_CACHE_EXPIRATION"),
			CacheControl: config.GetBool("MW_FIBER_CACHE_CACHECONTROL"),
		}))
	}

	// Compress
	if config.GetBool("MW_FIBER_COMPRESS_ENABLED") {
		lvl := compress.Level(config.GetInt("MW_FIBER_COMPRESS_LEVEL"))
		app.Use(compress.New(compress.Config{
			Level: lvl,
		}))
	}

	// CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "https://gofiber.io, https://gofiber.net",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// CSRF
	app.Use(csrf.New(csrf.Config{
		KeyLookup:      "header:X-Csrf-Token",
		CookieName:     "csrf_",
		CookieSameSite: "Strict",
		Expiration:     1 * time.Hour,
		KeyGenerator:   utils.UUID,
	}))

	// ETag
	app.Use(etag.New(etag.Config{
		Weak: true,
	}))

	// Favicon
	if config.GetBool("MW_FIBER_FAVICON_ENABLED") {
		app.Use(favicon.New(favicon.Config{
			File:         config.GetString("MW_FIBER_FAVICON_FILE"),
			CacheControl: config.GetString("MW_FIBER_FAVICON_CACHECONTROL"),
		}))
	}

	// TODO: Filesystem

	// Limiter
	if config.GetBool("MW_FIBER_LIMITER_ENABLED") {
		app.Use(limiter.New(limiter.Config{
			Max:      config.GetInt("MW_FIBER_LIMITER_MAX"),
			Duration: config.GetDuration("MW_FIBER_LIMITER_DURATION"),
			// TODO: Key
			// TODO: LimitReached
		}))
	}

	// TODO: Proxy

	// RequestID
	app.Use(requestid.New(requestid.Config{
		Header:     "X-Request-ID",
		ContextKey: "requestid",
	}))

	// TODO: Timeout

	//
	// For Debugging
	//

	// Pprof
	if config.GetBool("MW_FIBER_PPROF_ENABLED") {
		app.Use(pprof.New())
	}

	// Expvar
	if config.GetBool("MW_FIBER_EXPVAR_ENABLED") {
		app.Use(expvar.New())
	}
}

// Stop the Fiber application
func (app *App) exit() {
	_ = app.Shutdown()
}
