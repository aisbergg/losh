package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

// RemoveTrailingSlash will add a trailing slash (`/`) if this is not present is the client's request.
// This also takes the ability for the client to request a file extension into account.
func RemoveTrailingSlash() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		originalUrl := utils.ImmutableString(ctx.OriginalURL())

		if strings.HasSuffix(originalUrl, "/") {
			return ctx.Redirect(strings.TrimSuffix(originalUrl, "/"))
		}
		return ctx.Next()
	}
}
