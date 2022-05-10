package controllers

import (
	"losh/web/lib"

	"github.com/gofiber/fiber/v2"
)

func RegisterSearchRoute(router fiber.Router) {
	router.Get("/search", SearchHandler())
}

func SearchHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		bindings := lib.GetBindings()
		page := bindings["page"].(map[string]interface{})
		page["menu"] = "search"
		return ctx.Render("blank", bindings)
	}
}
