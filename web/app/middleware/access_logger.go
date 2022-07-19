package middleware

import (
	"losh/internal/logging"
	"losh/web/config"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gookit/event"

	"go.uber.org/zap"
)

// AccessLogger is a middleware that logs request to the server.
func AccessLogger(accessLogConfig config.AccessLogConfig) (fiber.Handler, error) {
	logConfig := logging.Config{
		Level:       "info",
		Format:      accessLogConfig.Format,
		Filename:    accessLogConfig.Filename,
		MaxSize:     accessLogConfig.MaxSize,
		MaxAge:      accessLogConfig.MaxAge,
		MaxBackups:  accessLogConfig.MaxBackups,
		LocalTime:   accessLogConfig.LocalTime,
		Compress:    accessLogConfig.Compress,
		Permissions: accessLogConfig.Permissions,
	}
	accessLoggerManager, err := logging.NewLoggerManager(logConfig)
	if err != nil {
		return nil, err
	}
	logger := accessLoggerManager.NewLogger("")

	// flush logs and close log file after server shutdown
	event.On("server.stop", event.ListenerFunc(func(e event.Event) error {
		return accessLoggerManager.Close()
	}))

	return func(ctx *fiber.Ctx) error {
		// Handle the request to calculate the number of bytes sent
		err := ctx.Next()

		// Chained error
		if err != nil {
			if chainErr := ctx.App().Config().ErrorHandler(ctx, err); chainErr != nil {
				_ = ctx.SendStatus(fiber.StatusInternalServerError)
			}
		}

		// Send structured information message to the logger
		msgBuilder := strings.Builder{}
		msgBuilder.Grow(256) // TODO: optimize value?
		msgBuilder.WriteString(ctx.IP())
		msgBuilder.WriteString(" - ")
		msgBuilder.WriteString(ctx.Method())
		msgBuilder.WriteString(" ")
		msgBuilder.WriteString(ctx.OriginalURL())
		msgBuilder.WriteString(" - ")
		msgBuilder.WriteString(strconv.Itoa(ctx.Response().StatusCode()))
		msgBuilder.WriteString(" - ")
		msgBuilder.WriteString(strconv.Itoa(len(ctx.Response().Body())))
		logger.Info(msgBuilder.String(),

			zap.String("ip", ctx.IP()),
			zap.String("hostname", ctx.Hostname()),
			zap.String("method", ctx.Method()),
			zap.String("path", ctx.OriginalURL()),
			zap.String("protocol", ctx.Protocol()),
			zap.Int("status", ctx.Response().StatusCode()),

			zap.String("x-forwarded-for", ctx.Get(fiber.HeaderXForwardedFor)),
			zap.String("user-agent", ctx.Get(fiber.HeaderUserAgent)),
			zap.String("referer", ctx.Get(fiber.HeaderReferer)),

			zap.Int("bytes_received", len(ctx.Request().Body())),
			zap.Int("bytes_sent", len(ctx.Response().Body())),
		)

		return err
	}, nil
}
