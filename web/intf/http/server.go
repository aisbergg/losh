package app

import (
	"strings"
	"time"

	"losh/internal/core/product/services"
	"losh/internal/infra/dgraph"
	"losh/internal/lib/log"
	"losh/web/build/assets"
	"losh/web/core/config"
	"losh/web/intf/http/controllers"
	"losh/web/intf/http/controllers/binding"
	"losh/web/intf/http/middleware"
	"losh/web/lib/template/liquid"

	"github.com/aisbergg/go-errors/pkg/errors"
	"github.com/gofiber/fiber/v2"
	futils "github.com/gofiber/utils"
	"github.com/gookit/event"
	"go.uber.org/zap"

	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/expvar"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

const logSelector = "server"

// Server represents the web server.
type Server struct {
	*fiber.App
	config    *config.Config
	db        *dgraph.DgraphRepository
	prdSvc    *services.Service
	log       *zap.SugaredLogger
	tplBndPrv binding.TemplateBindingProvider
}

// Listen starts the server.
func (s *Server) Listen(addr string) error {
	s.log.Info("starting server to listen on ", addr)
	err, _ := event.Fire("server.start", nil)
	if err != nil {
		return errors.Wrap(err, "failed to execute server.start event")
	}

	err = s.App.Listen(addr)
	if err != nil {
		event.Fire("server.error", nil)
		return errors.Wrap(err, "failed to run server")
	}

	return nil
}

// Shutdown closes the server gracefully.
func (s *Server) Shutdown() error {
	s.log.Info("shutting down server")
	shutdownErr := s.App.Shutdown()

	err, _ := event.Fire("server.stop", nil)
	if shutdownErr != nil {
		return shutdownErr
	}
	if err != nil {
		return errors.Wrap(err, "failed to execute server.stop event")
	}

	return nil
}

// NewServer creates a new server instance.
func NewServer(config *config.Config, db *dgraph.DgraphRepository) (*Server, error) {
	log := log.NewLogger(logSelector)
	prdSvc := services.NewService(db)
	tplBndPrv := binding.NewTemplateBindingProvider(config)
	fiberConfig, err := createFiberConfig(config, log, tplBndPrv)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create fiber config")
	}
	app := Server{
		App:       fiber.New(*fiberConfig),
		config:    config,
		db:        db,
		prdSvc:    prdSvc,
		log:       log,
		tplBndPrv: tplBndPrv,
	}

	// register middleware
	if err := app.registerMiddlewares(); err != nil {
		return nil, errors.Wrap(err, "failed to register middlewares")
	}

	// register routes
	if err := app.registerRoutes(); err != nil {
		return nil, errors.Wrap(err, "failed to register routes")
	}

	return &app, nil
}

func createFiberConfig(config *config.Config, log *zap.SugaredLogger, tplBndPrv binding.TemplateBindingProvider) (*fiber.Config, error) {
	// fiber view engine
	// tmplPath, err := utils.ResolveExecRelPath("assets/templates")
	// if err != nil {
	// 	return nil, errors.Wrapf(err, "missing template dir '%s'", "assets/templates")
	// }
	// viewsEngine := liquid.New(tmplPath, ".html")
	viewsEngine := liquid.NewFileSystem(assets.AssetsHTTP, "templates", ".html")
	viewsEngine.Layout("content")                  // use 'content' as var name in layout
	viewsEngine.EnableReload(config.Debug.Enabled) // in debug mode reload templates on every request
	viewsEngine.EnableFrontmatter(true)            // enable frontmatter
	if err := viewsEngine.Load(); err != nil {
		return nil, errors.Wrap(err, "failed to load templates")
	}

	// fiber error handler
	errCtl := controllers.NewErrorController(config.Debug.Enabled, log, tplBndPrv)

	// https://docs.gofiber.io/api/fiber#config
	return &fiber.Config{
		AppName:                 "LOSH",
		Prefork:                 false,
		ReduceMemoryUsage:       false,
		UnescapePath:            true,
		BodyLimit:               20 * 1024 * 1024,
		Concurrency:             256 * 1024,
		Views:                   viewsEngine,
		ViewsLayout:             "layouts/default",
		PassLocalsToViews:       false,
		ReadTimeout:             0,
		WriteTimeout:            0,
		IdleTimeout:             0,
		ReadBufferSize:          8192,
		WriteBufferSize:         8192,
		EnableTrustedProxyCheck: true,
		TrustedProxies:          config.Server.TrustedDomains,
		ProxyHeader:             "X-Forwarded-For",
		ErrorHandler:            errCtl.Handle,
		DisableStartupMessage:   true,
	}, nil
}

func (s *Server) registerRoutes() error {
	// static files
	s.Use("/static", filesystem.New(filesystem.Config{
		Root:       assets.AssetsHTTP,
		PathPrefix: "static", // path within the FS
		Browse:     false,
	}))

	// web routes
	tplBndPrv := binding.NewTemplateBindingProvider(s.config)
	web := s.Group("")
	controllers.NewHomeController(tplBndPrv).Register(web)
	controllers.NewSearchController(s.db, tplBndPrv, s.config.Debug.Enabled).Register(web)
	controllers.NewDetailsController(s.prdSvc, tplBndPrv).Register(web)
	controllers.NewAboutController(tplBndPrv).Register(web)
	controllers.NewRDFController(s.prdSvc, tplBndPrv).Register(web)

	// api routes
	// api := app.Group("/api")
	// apiv1 := api.Group("/v1")
	// routes.RegisterAPI(apiv1, app.DB)

	// routes for debugging purposes
	if s.config.Debug.Enabled {
		controllers.RegisterDebugRoute(web)
	}

	// 404 handler (served when no other route matches)
	s.Use(func(c *fiber.Ctx) error {
		// main error handler takes care of it
		return fiber.ErrNotFound
	})

	return nil
}

func (s *Server) registerMiddlewares() error {
	// Access Logs
	if s.config.AccessLog.Enabled {
		lgMdlw, err := middleware.AccessLogger(s.config.AccessLog)
		if err != nil {
			return errors.Wrap(err, "failed to create access logger middleware")
		}
		s.Use(lgMdlw)
	}

	// remove trailing slash
	s.Use(middleware.RemoveTrailingSlash())

	// Cache
	if s.config.Server.Cache.Enabled {
		s.Use(cache.New(cache.Config{
			Expiration:   s.config.Server.Cache.Expiration,
			CacheControl: s.config.Server.Cache.CacheControl, // use client side caching
			KeyGenerator: func(ctx *fiber.Ctx) string {
				return ctx.OriginalURL()
			},
			MaxBytes: 100 * 1024 * 1024, // 100MB
		}))
	}

	// Compress
	if s.config.Server.Compress >= 0 {
		if s.config.Server.Compress > 2 {
			return errors.New("invalid compress level")
		}
		s.Use(compress.New(compress.Config{
			Level: compress.Level(s.config.Server.Compress),
		}))
	}

	// CORS
	allowedOrigins := []string{}
	for _, td := range s.config.Server.TrustedDomains {
		td = strings.TrimSpace(td)
		allowedOrigins = append(allowedOrigins, "http://"+td, "https://"+td)
	}
	s.Use(cors.New(cors.Config{
		AllowOrigins: strings.Join(allowedOrigins, ", "),
		MaxAge:       600, // 10 minutes
	}))

	// CSRF
	s.Use(csrf.New(csrf.Config{
		KeyLookup:      "header:X-Csrf-Token",
		CookieName:     "csrf_",
		CookieSameSite: "Strict",
		Expiration:     1 * time.Hour,
		KeyGenerator:   futils.UUID,
	}))

	// ETag
	s.Use(etag.New(etag.Config{
		Weak: true,
	}))

	// TODO: Favicon

	// RequestID
	s.Use(requestid.New(requestid.Config{
		Header:     "X-Request-ID",
		ContextKey: "requestid",
	}))

	// TODO: Filesystem

	// TODO: Limiter

	// TODO: Proxy

	// TODO: Timeout

	//
	// For Debugging
	//

	// Pprof
	if s.config.Debug.Pprof {
		s.Use(pprof.New())
	}

	// Expvar
	if s.config.Debug.Expvar {
		s.Use(expvar.New())
	}

	return nil
}
