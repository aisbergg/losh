// Copyright 2022 André Lehmann
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
