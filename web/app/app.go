package app

import (
	"losh/web/app/controllers"
	"losh/web/app/middleware"
	configuration "losh/web/config"
	"losh/web/lib/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	futils "github.com/gofiber/utils"
	"github.com/rotisserie/eris"

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
)

type App struct {
	*fiber.App
}

func NewApp(config *configuration.Config) *App {
	// app
	app := App{
		App: fiber.New(*config.GetFiberConfig()),
	}

	// middelwares
	app.registerMiddlewares(config)

	// static files
	staticPath, err := utils.ResolveExecRelPath("assets/static")
	if err != nil {
		panic(eris.Wrapf(err, "missing static files '%s'", staticPath))
	}
	app.Static("/static", staticPath)
	app.Static("/", staticPath)

	// web routes
	web := app.Group("")
	controllers.RegisterHomeRoute(web)
	controllers.RegisterSearchRoute(web)
	controllers.RegisterAboutRoute(web)

	// api routes
	// api := app.Group("/api")
	// apiv1 := api.Group("/v1")
	// routes.RegisterAPI(apiv1, app.DB)

	// TODO: remove later
	app.Get("error", func(c *fiber.Ctx) error {
		return fiber.ErrInternalServerError
	})

	// 404 handler
	app.Use(func(c *fiber.Ctx) error {
		// main error handler takes care of it
		return fiber.ErrNotFound
	})

	return &app
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

	// force HTTPS
	if config.GetBool("MW_FORCE_HTTPS_ENABLED") {
		app.Use(middleware.ForceHTTPS())
	}

	// remove trailing slash
	app.Use(middleware.RemoveTrailingSlash())

	// recover from panics
	app.Use(recover.New())

	// cache
	if config.GetBool("MW_FIBER_CACHE_ENABLED") {
		app.Use(cache.New(cache.Config{
			Expiration:   config.GetDuration("MW_FIBER_CACHE_EXPIRATION"),
			CacheControl: config.GetBool("MW_FIBER_CACHE_CACHECONTROL"),
		}))
	}

	// compress
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
		KeyGenerator:   futils.UUID,
	}))

	// ETag
	app.Use(etag.New(etag.Config{
		Weak: true,
	}))

	// favicon
	if config.GetBool("MW_FIBER_FAVICON_ENABLED") {
		app.Use(favicon.New(favicon.Config{
			File:         config.GetString("MW_FIBER_FAVICON_FILE"),
			CacheControl: config.GetString("MW_FIBER_FAVICON_CACHECONTROL"),
		}))
	}

	// TODO: Filesystem

	// limiter
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
