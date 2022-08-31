package controllers

import (
	"losh/web/intf/http/controllers/binding"

	"github.com/gofiber/fiber/v2"
)

// AboutController is the controller for the resource details page at '/details/:id'.
type AboutController struct {
	Controller
}

// NewAboutController creates a new AboutController.
func NewAboutController(tplBndPrv binding.TemplateBindingProvider) AboutController {
	return AboutController{Controller: Controller{tplBndPrv: tplBndPrv}}
}

// Register registers the controller with the given router.
func (c AboutController) Register(router fiber.Router) {
	router.Get("/about/project", c.AboutProjectHandler)
	router.Get("/about/faq", c.FAQHandler)
}

func (c AboutController) AboutProjectHandler(ctx *fiber.Ctx) error {
	_, tplBnd := c.preprocessRequest(ctx, nil, nil)
	page := tplBnd["page"].(map[string]interface{})
	page["title"] = "About the Project"
	page["page-header"] = "About the Project"
	page["menu"] = "about.project"
	return ctx.Render("blank", tplBnd)
}

func (c AboutController) FAQHandler(ctx *fiber.Ctx) error {
	_, tplBnd := c.preprocessRequest(ctx, nil, nil)
	page := tplBnd["page"].(map[string]interface{})
	page["title"] = "FAQ"
	page["page-header"] = "FAQ"
	page["menu"] = "about.faq"
	return ctx.Render("blank", tplBnd)
}
