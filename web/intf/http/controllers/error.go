package controllers

import (
	losherrors "losh/internal/lib/errors"
	"losh/web/intf/http/controllers/binding"
	"losh/web/lib/template/liquid"

	"github.com/aisbergg/go-errors/pkg/errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/utils"
	"go.uber.org/zap"
)

// IsControllerError returns true if err is of type ControllerError.
func IsControllerError(err error) bool {
	_, ok := err.(interface{ IsControllerError() })
	return ok
}

// ControllerError is an error that is temporary and can be retried.
type ControllerError struct {
	losherrors.AppError
	RequestInfo *RequestInfo
}

// newControllerError wraps an error into ControllerError.
func newControllerError(err error, requestInfo *RequestInfo, format string, args ...interface{}) error {
	svrErr := &ControllerError{
		AppError:    *losherrors.NewAppErrorWrap(err, format, args...),
		RequestInfo: requestInfo,
	}
	svrErr.Add("reqURL", requestInfo.URL)
	svrErr.Add("reqMethod", requestInfo.Method)
	svrErr.Add("reqHeaders", requestInfo.Headers)
	svrErr.Add("reqParams", requestInfo.Params)
	svrErr.Add("reqQueryParams", requestInfo.QueryParams)
	return svrErr
}

// IsControllerError indicates the type of the error.
func (e *ControllerError) IsControllerError() {}

type errorInfo struct {
	Code    int    `json:"code" liquid:"code"`
	Reason  string `json:"reason" liquid:"reason"`
	Title   string `json:"title" liquid:"title"`
	Message string `json:"message" liquid:"message"`
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
		"internalControllerError",
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

type ErrorController struct {
	log       *zap.SugaredLogger
	tplBndPrv binding.TemplateBindingProvider
	debug     bool
}

// NewErrorController creates a new ErrorController.
func NewErrorController(debug bool, log *zap.SugaredLogger, tplBndPrv binding.TemplateBindingProvider) ErrorController {
	return ErrorController{
		debug:     debug,
		log:       log,
		tplBndPrv: tplBndPrv,
	}
}

// RegisterDebugRoute registers a simple debug route handler.
func RegisterDebugRoute(router fiber.Router) {
	router.Get("/error", func(ctx *fiber.Ctx) error {
		reqInf := parseRequestInfo(ctx, nil, nil)
		barErr := losherrors.NewAppError("bar").
			Add("foo", "1").
			Add("bar", "2")
		fooErr := newControllerError(barErr, reqInf, "foo")
		return fooErr
	})
}

func (c ErrorController) Handle(ctx *fiber.Ctx, err error) error {
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
		return c.handleClientError(ctx, code, message)
	}
	return c.handleControllerError(ctx, code, message, err)
}

// handleControllerError handles a 4XX error.
func (c ErrorController) handleClientError(ctx *fiber.Ctx, code int, message string) (err error) {
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
			c.log.Errorw("failed to render error.html", "error", err)
			err = ctx.SendString(info.Message)
		}
	}
	return
}

// handleControllerError handles a 5XX error.
func (c ErrorController) handleControllerError(ctx *fiber.Ctx, code int, message string, err error) error {
	// log the error
	c.logError(err)

	// render the error page
	errInf, ok := errorInfos[code]
	if !ok {
		errInf = errorInfo{Code: code, Reason: "unknown", Message: message}
	}
	tplBnd := c.tplBndPrv.Get()
	tplBnd["error"] = errInf

	// add more values in debug mode
	if c.debug {
		// add user values (locals set via ctx.Locals(k, v))
		locals := make(map[string]interface{})
		ctx.Context().VisitUserValues(func(key []byte, value interface{}) {
			locals[utils.UnsafeString(key)] = value
		})
		tplBnd["locals"] = locals

		// add template bind vars
		if tplErr, ok := err.(liquid.TemplatingError); ok {
			tplBnd["vars"] = tplErr.Bindings()
		}

		// add request information
		request := make(map[string]interface{})
		request["uri"] = utils.UnsafeString(ctx.Request().RequestURI())
		request["method"] = utils.UnsafeString(ctx.Request().Header.Method())
		tplBnd["request"] = request

		// add request method
		request["method"] = ctx.Method()

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

		// error
		errMsg := errors.ToString(err, false)
		var errStk interface{} = []string{}
		if terr, ok := err.(interface{ FullStack() errors.Stack }); ok {
			errStk = terr.FullStack()
		}
		errCtx := losherrors.GetFullContext(err)
		tplBnd["error"] = map[string]interface{}{
			"message": errMsg,
			"stack":   errStk,
			"context": errCtx,
		}
	}

	contentType := string(ctx.Request().Header.ContentType())
	switch contentType {
	case fiber.MIMEApplicationJSON, fiber.MIMEApplicationJSONCharsetUTF8:
		// TODO: json error https://cloud.google.com/storage/docs/json_api/v1/status-codes#403-forbidden
		return ctx.JSON(tplBnd)

	default:
		ctx.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
		if c.debug {
			err = ctx.Render("error-debug.html", tplBnd, "layouts/base")
		} else {
			err = ctx.Render("error.html", tplBnd, "layouts/base")
		}
		if err != nil {
			c.log.Errorw("failed to render error.html", "error", err)
			err = ctx.SendString(errInf.Message)
			return err
		}
	}

	return nil
}

func (c ErrorController) logError(err error) {
	errCnt := losherrors.GetFullContext(err)
	if terr, ok := err.(interface{ FullStack() errors.Stack }); ok {
		errCnt["stack"] = terr.FullStack()
	}
	ctxKV := make([]interface{}, 0, len(errCnt)*2)
	for k, v := range errCnt {
		ctxKV = append(ctxKV, k, v)
	}
	msg := errors.ToString(err, false)
	c.log.Errorw(msg, ctxKV...)
}
