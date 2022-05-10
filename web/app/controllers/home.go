package controllers

import (
	"losh/web/lib"

	"github.com/gofiber/fiber/v2"
)

func RegisterHomeRoute(router fiber.Router) {
	router.Get("/", HomeHandler())
}

func HomeHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Bind data to template
		bindings := lib.GetBindings()
		// Render template

		// out, err := yaml.Marshal(bindings)
		// if err != nil {
		// 	return ctx.SendString(err.Error())
		// }
		// return ctx.SendString(string(out))

		// return ctx.SendString(fmt.Sprintf("%s", bindings))
		return ctx.Render("index", bindings)
	}
}
