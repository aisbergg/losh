package controllers

import (
	"losh/web/intf/http/controllers/binding"

	"github.com/gofiber/fiber/v2"
)

// HomeController is the controller for the homepage at '/'.
type HomeController struct {
	Controller
}

// NewHomeController creates a new HomeController.
func NewHomeController(tplBndPrv binding.TemplateBindingProvider) HomeController {
	return HomeController{Controller: Controller{tplBndPrv: tplBndPrv}}
}

func (c HomeController) Register(router fiber.Router) {
	router.Get("/", c.Handle)
}

func (c HomeController) Handle(ctx *fiber.Ctx) error {
	_, tplBnd := c.preprocessRequest(ctx, nil, nil)
	return ctx.Render("home", tplBnd)
}
