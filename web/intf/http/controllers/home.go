package controllers

import (
	"losh/web/intf/http/controllers/binding"

	"github.com/gofiber/fiber/v2"
)

// HomeController is the controller for the homepage at '/'.
type HomeController struct {
	tplBndPrv binding.TemplateBindingProvider
}

// NewHomeController creates a new HomeController.
func NewHomeController(tplBndPrv binding.TemplateBindingProvider) HomeController {
	return HomeController{tplBndPrv}
}

func (c HomeController) Register(router fiber.Router) {
	router.Get("/", c.Handle)
}

func (c HomeController) Handle(ctx *fiber.Ctx) error {
	// Bind data to template
	bindings := c.tplBndPrv.Get()
	// Render template

	// out, err := yaml.Marshal(bindings)
	// if err != nil {
	// 	return ctx.SendString(err.Error())
	// }
	// return ctx.SendString(string(out))

	// return ctx.SendString(fmt.Sprintf("%s", bindings))
	return ctx.Render("index", bindings)
}
