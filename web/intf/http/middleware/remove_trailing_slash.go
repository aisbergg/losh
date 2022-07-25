package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

// RemoveTrailingSlash will remove a trailing slash.
func RemoveTrailingSlash() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		originalURL := utils.CopyString(ctx.OriginalURL())
		if len(originalURL) > 1 && strings.HasSuffix(originalURL, "/") {
			return ctx.Redirect(strings.TrimSuffix(originalURL, "/"))
		}
		return ctx.Next()
	}
}
