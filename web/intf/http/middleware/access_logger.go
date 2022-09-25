// Copyright 2022 André Lehmann
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package middleware

import (
	"strconv"
	"strings"
	"sync"

	"losh/internal/lib/log"
	"losh/web/core/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gookit/event"

	"go.uber.org/zap"
)

// AccessLogger is a middleware that logs request to the server.
func AccessLogger(accessLogConfig config.AccessLogConfig) (fiber.Handler, error) {
	logConfig := log.Config{
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
	accessLoggerManager, err := log.NewLoggerManager(logConfig)
	if err != nil {
		return nil, err
	}
	logger := accessLoggerManager.NewLogger("")

	// flush logs and close log file after server shutdown
	event.On("server.stop", event.ListenerFunc(func(e event.Event) error {
		return accessLoggerManager.Close()
	}))

	// console logger
	if accessLogConfig.Format == "console" {
		pool := sync.Pool{
			New: func() interface{} { return &strings.Builder{} },
		}
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
			msgBuilder := pool.Get().(*strings.Builder)
			msgBuilder.Reset()
			msgBuilder.WriteString(ctx.IP())
			msgBuilder.WriteString(" - ")
			msgBuilder.WriteString(ctx.Method())
			msgBuilder.WriteString(" ")
			msgBuilder.WriteString(ctx.OriginalURL())
			msgBuilder.WriteString(" - ")
			msgBuilder.WriteString(strconv.Itoa(ctx.Response().StatusCode()))
			logger.Info(msgBuilder.String())
			pool.Put(msgBuilder)

			return err
		}, nil
	}

	// json logger
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
		logger.Infow("",
			zap.String("remote_addr", ctx.IP()),
			zap.String("request_uri", ctx.OriginalURL()),
			zap.String("request_protocol", ctx.Protocol()),
			zap.Int("http_status", ctx.Response().StatusCode()),
			zap.String("http_host", ctx.Hostname()),
			zap.String("http_method", ctx.Method()),
			zap.String("http_user_agent", ctx.Get(fiber.HeaderUserAgent)),
			zap.String("http_referer", ctx.Get(fiber.HeaderReferer)),
		)

		return err
	}, nil
}
