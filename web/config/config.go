package config

import (
	"fmt"
	"losh/web/app/controllers"
	"losh/web/lib/template/liquid"
	"losh/web/lib/utils"
	"os"

	"github.com/rotisserie/eris"
	"github.com/spf13/viper"

	"github.com/gofiber/fiber/v2"
)

type Config struct {
	*viper.Viper

	errorHandler fiber.ErrorHandler
	fiber        *fiber.Config
}

func New() *Config {
	config := &Config{
		Viper: viper.New(),
	}

	// Set default configurations
	config.setDefaults()

	// Select the .env file
	config.SetConfigName(".env")
	config.SetConfigType("dotenv")
	config.AddConfigPath(".")

	// Automatically refresh environment variables
	config.AutomaticEnv()

	// Read configuration
	if err := config.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Println("failed to read configuration:", err.Error())
			os.Exit(1)
		}
	}

	config.SetErrorHandler(controllers.ErrorHandler(config.GetBool("FIBER_VIEWS_DEBUG")))

	// TODO: Logger (Maybe a different zap object)

	// TODO: Add APP_KEY generation

	// TODO: Write changes to configuration file

	// Set Fiber configurations
	config.setFiberConfig()

	return config
}

func (config *Config) SetErrorHandler(errorHandler fiber.ErrorHandler) {
	config.errorHandler = errorHandler
}

func (config *Config) setDefaults() {
	// Set default App configuration
	config.SetDefault("APP_ADDR", ":8080")
	config.SetDefault("APP_ENV", "local")

	// Set default Fiber configuration
	config.SetDefault("FIBER_PREFORK", false)
	config.SetDefault("FIBER_SERVERHEADER", "")
	config.SetDefault("FIBER_STRICTROUTING", false)
	config.SetDefault("FIBER_CASESENSITIVE", false)
	config.SetDefault("FIBER_IMMUTABLE", false)
	config.SetDefault("FIBER_UNESCAPEPATH", false)
	config.SetDefault("FIBER_ETAG", false)
	config.SetDefault("FIBER_BODYLIMIT", 4194304)
	config.SetDefault("FIBER_CONCURRENCY", 262144)
	// config.SetDefault("FIBER_VIEWS", "html")
	config.SetDefault("FIBER_VIEWS", "django")

	config.SetDefault("FIBER_VIEWS_DEBUG", true)

	config.SetDefault("FIBER_READTIMEOUT", 0)
	config.SetDefault("FIBER_WRITETIMEOUT", 0)
	config.SetDefault("FIBER_IDLETIMEOUT", 0)
	config.SetDefault("FIBER_READBUFFERSIZE", 4096)
	config.SetDefault("FIBER_WRITEBUFFERSIZE", 4096)
	config.SetDefault("FIBER_COMPRESSEDFILESUFFIX", ".fiber.gz")
	config.SetDefault("FIBER_PROXYHEADER", "")
	config.SetDefault("FIBER_GETONLY", false)
	config.SetDefault("FIBER_DISABLEKEEPALIVE", false)
	config.SetDefault("FIBER_DISABLEDEFAULTDATE", false)
	config.SetDefault("FIBER_DISABLEDEFAULTCONTENTTYPE", false)
	config.SetDefault("FIBER_DISABLEHEADERNORMALIZING", false)
	config.SetDefault("FIBER_DISABLESTARTUPMESSAGE", false)
	config.SetDefault("FIBER_REDUCEMEMORYUSAGE", false)

	// Set default Custom Access Logger middleware configuration
	config.SetDefault("MW_ACCESS_LOGGER_ENABLED", true)
	config.SetDefault("MW_ACCESS_LOGGER_TYPE", "console")
	config.SetDefault("MW_ACCESS_LOGGER_FILENAME", "access.log")
	config.SetDefault("MW_ACCESS_LOGGER_MAXSIZE", 500)
	config.SetDefault("MW_ACCESS_LOGGER_MAXAGE", 28)
	config.SetDefault("MW_ACCESS_LOGGER_MAXBACKUPS", 3)
	config.SetDefault("MW_ACCESS_LOGGER_LOCALTIME", false)
	config.SetDefault("MW_ACCESS_LOGGER_COMPRESS", false)

	// Set default Force HTTPS middleware configuration
	config.SetDefault("MW_FORCE_HTTPS_ENABLED", false)

	// Set default HSTS middleware configuration
	config.SetDefault("MW_HSTS_ENABLED", false)
	config.SetDefault("MW_HSTS_MAXAGE", 31536000)
	config.SetDefault("MW_HSTS_INCLUDESUBDOMAINS", true)
	config.SetDefault("MW_HSTS_PRELOAD", false)

	// Set default Fiber Cache middleware configuration
	config.SetDefault("MW_FIBER_CACHE_ENABLED", false)
	config.SetDefault("MW_FIBER_CACHE_EXPIRATION", "1m")
	config.SetDefault("MW_FIBER_CACHE_CACHECONTROL", false)

	// Set default Fiber Compress middleware configuration
	config.SetDefault("MW_FIBER_COMPRESS_ENABLED", false)
	config.SetDefault("MW_FIBER_COMPRESS_LEVEL", 0)

	// Set default Fiber CORS middleware configuration
	config.SetDefault("MW_FIBER_CORS_ENABLED", false)
	config.SetDefault("MW_FIBER_CORS_ALLOWORIGINS", "*")
	config.SetDefault("MW_FIBER_CORS_ALLOWMETHODS", "GET,POST,HEAD,PUT,DELETE,PATCH")
	config.SetDefault("MW_FIBER_CORS_ALLOWHEADERS", "")
	config.SetDefault("MW_FIBER_CORS_ALLOWCREDENTIALS", false)
	config.SetDefault("MW_FIBER_CORS_EXPOSEHEADERS", "")
	config.SetDefault("MW_FIBER_CORS_MAXAGE", 0)

	// Set default Fiber CSRF middleware configuration
	config.SetDefault("MW_FIBER_CSRF_ENABLED", false)
	config.SetDefault("MW_FIBER_CSRF_TOKENLOOKUP", "header:X-CSRF-Token")
	config.SetDefault("MW_FIBER_CSRF_COOKIE_NAME", "_csrf")
	config.SetDefault("MW_FIBER_CSRF_COOKIE_SAMESITE", "Strict")
	config.SetDefault("MW_FIBER_CSRF_COOKIE_EXPIRES", "24h")
	config.SetDefault("MW_FIBER_CSRF_CONTEXTKEY", "csrf")

	// Set default Fiber ETag middleware configuration
	config.SetDefault("MW_FIBER_ETAG_ENABLED", false)
	config.SetDefault("MW_FIBER_ETAG_WEAK", false)

	// Set default Fiber Expvar middleware configuration
	config.SetDefault("MW_FIBER_EXPVAR_ENABLED", false)

	// Set default Fiber Favicon middleware configuration
	config.SetDefault("MW_FIBER_FAVICON_ENABLED", false)
	config.SetDefault("MW_FIBER_FAVICON_FILE", "")
	config.SetDefault("MW_FIBER_FAVICON_CACHECONTROL", "public, max-age=31536000")

	// Set default Fiber Limiter middleware configuration
	config.SetDefault("MW_FIBER_LIMITER_ENABLED", false)
	config.SetDefault("MW_FIBER_LIMITER_MAX", 5)
	config.SetDefault("MW_FIBER_LIMITER_DURATION", "1m")

	// Set default Fiber Monitor middleware configuration
	config.SetDefault("MW_FIBER_MONITOR_ENABLED", false)

	// Set default Fiber Pprof middleware configuration
	config.SetDefault("MW_FIBER_PPROF_ENABLED", false)

	// Set default Fiber Recover middleware configuration
	config.SetDefault("MW_FIBER_RECOVER_ENABLED", true)

	// Set default Fiber RequestID middleware configuration
	config.SetDefault("MW_FIBER_REQUESTID_ENABLED", false)
	config.SetDefault("MW_FIBER_REQUESTID_HEADER", "X-Request-ID")
	config.SetDefault("MW_FIBER_REQUESTID_CONTEXTKEY", "requestid")
}

func (config *Config) getFiberViewsEngine() fiber.Views {
	tmplPath, err := utils.ResolveExecRelPath("assets/templates")
	if err != nil {
		panic(eris.Wrapf(err, "missing template dir '%s'", tmplPath))
	}
	engine := liquid.New(tmplPath, ".html")
	engine.Layout("content")
	engine.Reload(config.GetBool("FIBER_VIEWS_DEBUG"))
	err = engine.Load()
	if err != nil {
		panic(err)
	}
	return engine
}

func (config *Config) setFiberConfig() {
	config.fiber = &fiber.Config{
		AppName:               "LOSH",
		Prefork:               config.GetBool("FIBER_PREFORK"),
		UnescapePath:          true,
		BodyLimit:             20 * 1024 * 1024,
		Concurrency:           config.GetInt("FIBER_CONCURRENCY"),
		Views:                 config.getFiberViewsEngine(),
		ViewsLayout:           "layouts/default2",
		PassLocalsToViews:     false,
		ReadTimeout:           config.GetDuration("FIBER_READTIMEOUT"),
		WriteTimeout:          config.GetDuration("FIBER_WRITETIMEOUT"),
		IdleTimeout:           config.GetDuration("FIBER_IDLETIMEOUT"),
		ReadBufferSize:        config.GetInt("FIBER_READBUFFERSIZE"),
		WriteBufferSize:       config.GetInt("FIBER_WRITEBUFFERSIZE"),
		CompressedFileSuffix:  config.GetString("FIBER_COMPRESSEDFILESUFFIX"),
		ProxyHeader:           config.GetString("FIBER_PROXYHEADER"),
		ErrorHandler:          config.errorHandler,
		DisableStartupMessage: config.GetBool("FIBER_DISABLESTARTUPMESSAGE"),
		ReduceMemoryUsage:     false,
	}
}

func (config *Config) GetFiberConfig() *fiber.Config {
	return config.fiber
}
