package controllers

import (
	"losh/web/lib"

	"github.com/gofiber/fiber/v2"
)

func RegisterAboutRoute(router fiber.Router) {
	router.Get("/about/project", AboutProjectHandler())
	router.Get("/about/faq", FAQHandler())
}

func AboutProjectHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		bindings := lib.GetBindings()
		page := bindings["page"].(map[string]interface{})
		page["menu"] = "about.project"
		return ctx.Render("blank", bindings)
	}
}

func FAQHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		bindings := lib.GetBindings()
		page := bindings["page"].(map[string]interface{})
		page["menu"] = "about.faq"
		return ctx.Render("blank", bindings)
	}
}
