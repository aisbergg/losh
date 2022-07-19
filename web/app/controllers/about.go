package controllers

import (
	"losh/web/app/controllers/binding"

	"github.com/gofiber/fiber/v2"
)

// AboutController is the controller for the resource details page at '/details/:id'.
type AboutController struct {
	tplBndPrv binding.TemplateBindingProvider
}

// NewAboutController creates a new AboutController.
func NewAboutController(tplBndPrv binding.TemplateBindingProvider) AboutController {
	return AboutController{tplBndPrv}
}

// Register registers the controller with the given router.
func (c AboutController) Register(router fiber.Router) {
	router.Get("/details/:id", c.AboutProjectHandler)
	router.Get("/details/:id", c.FAQHandler)
}

func (c AboutController) AboutProjectHandler(ctx *fiber.Ctx) error {
	bindings := c.tplBndPrv.Get()
	page := bindings["page"].(map[string]interface{})
	page["title"] = "About the Project"
	page["menu"] = "about.project"
	return ctx.Render("blank", bindings)
}

func (c AboutController) FAQHandler(ctx *fiber.Ctx) error {
	bindings := c.tplBndPrv.Get()
	page := bindings["page"].(map[string]interface{})
	page["title"] = "FAQ"
	page["menu"] = "about.faq"
	return ctx.Render("blank", bindings)
}
