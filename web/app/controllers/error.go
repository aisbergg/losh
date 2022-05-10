package controllers

import (
	"fmt"
	"losh/web/lib/template/liquid"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/utils"
)

type errorInfo struct {
	Code    int    `liquid:"code"`
	Reason  string `liquid:"reason"`
	Title   string `liquid:"title"`
	Message string `liquid:"message"`
}

var errorInfos = map[int]errorInfo{
	400: errorInfo{
		400,
		"badRequest",
		"Bad Request",
		"We are sorry but your request contains bad syntax and cannot be fulfilled",
	},
	401: errorInfo{
		401,
		"unauthorized",
		"Unauthorized",
		"We are sorry but you are not authorized to access this page",
	},
	403: errorInfo{
		403,
		"forbidden",
		"Forbidden",
		"We are sorry but you do not have permission to access this page",
	},
	404: errorInfo{
		404,
		"notFound",
		"Not Found",
		"We are sorry but the page you are looking for was not found",
	},
	429: errorInfo{
		429,
		"tooManyRequests",
		"Too Many Requests",
		"Rate limit was exceeded. Throttle your client's requests and use exponential backoff",
	},
	500: errorInfo{
		500,
		"internalServerError",
		"Internal Server Error",
		"We are sorry but our server encountered an internal error",
	},
	503: errorInfo{
		503,
		"serviceUnavailable",
		"Service Unavailable",
		"We are sorry but our service is currently not available",
	},
}

func ErrorHandler(debug bool) fiber.ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {
		// 500 as default error code
		code := fiber.StatusInternalServerError
		message := err.Error()

		// retrieve the custom status code if it's an fiber.*Error
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
			message = e.Message
		}

		// set response status code
		ctx.Status(code)

		// handle client or server error
		if 400 <= code && code < 500 {
			return handleClientError(ctx, code, message)
		}
		return handleServerError(ctx, code, message, debug)
	}
}

// handleServerError handles a 4XX error.
func handleClientError(ctx *fiber.Ctx, code int, message string) (err error) {
	info, ok := errorInfos[code]
	if !ok {
		info = errorInfo{Code: code, Message: message}
	}
	bind := map[string]interface{}{
		"error": info,
	}
	contentType := string(ctx.Request().Header.ContentType())
	if contentType == fiber.MIMEApplicationJSON || contentType == fiber.MIMEApplicationJSONCharsetUTF8 {
		// TODO: json error https://cloud.google.com/storage/docs/json_api/v1/status-codes#403-forbidden
		ctx.Set("Content-Type", "application/json")
		err = ctx.JSON(bind)
	} else {
		ctx.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
		err = ctx.Render("error.html", bind, "layouts/base")
		if err != nil {
			err = ctx.SendString(info.Message)
		}
	}
	return
}

// handleServerError handles a 5XX error.
func handleServerError(ctx *fiber.Ctx, code int, message string, debug bool) (err error) {
	info, ok := errorInfos[code]
	if !ok {
		info = errorInfo{Code: code, Reason: "unknown", Message: message}
	}
	bind := map[string]interface{}{
		"error": info,
	}

	// add more values in debug mode
	if debug {
		// add user values (locals set via ctx.Locals(k, v))
		locals := make(map[string]interface{})
		ctx.Context().VisitUserValues(func(key []byte, value interface{}) {
			locals[utils.UnsafeString(key)] = value
		})
		bind["locals"] = locals

		// add template bind vars
		if tplErr, ok := err.(liquid.TemplatingError); ok {
			bind["vars"] = tplErr.Bindings()
		}

		// add request information
		request := make(map[string]interface{})
		request["uri"] = utils.UnsafeString(ctx.Request().RequestURI())
		request["method"] = utils.UnsafeString(ctx.Request().Header.Method())
		bind["request"] = request

		// // add request method
		// requestMethod := make(map[string]interface{})
		// ctx.Request().Header.VisitAllInOrder(func(key, value []byte) {
		// 	requestMethod[utils.UnsafeString(key)] = utils.UnsafeString(value)
		// })
		// request["method"] = requestMethod

		// add request headers
		requestHeaders := make(map[string]interface{})
		ctx.Request().Header.VisitAllInOrder(func(key, value []byte) {
			requestHeaders[utils.UnsafeString(key)] = utils.UnsafeString(value)
		})
		request["headers"] = requestHeaders

		// add request cookies
		requestCookies := make(map[string]interface{})
		ctx.Request().Header.VisitAllCookie(func(key []byte, value []byte) {
			requestCookies[utils.UnsafeString(key)] = utils.UnsafeString(value)
		})
		request["cookies"] = requestCookies
	}

	// TODO: proper error logging; sentry?
	fmt.Printf("error: %d - %s", code, message)

	contentType := string(ctx.Request().Header.ContentType())
	if contentType == fiber.MIMEApplicationJSON || contentType == fiber.MIMEApplicationJSONCharsetUTF8 {
		// TODO: json error https://cloud.google.com/storage/docs/json_api/v1/status-codes#403-forbidden
		ctx.Set("Content-Type", "application/json")
		err = ctx.JSON(bind)
	} else {
		ctx.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
		if debug {
			err = ctx.Render("error-debug.html", bind, "layouts/base")
			fmt.Println(err)
		} else {
			err = ctx.Render("error.html", bind, "layouts/base")
		}
		if err != nil {
			err = ctx.SendString(info.Message)
		}
	}

	return
}
