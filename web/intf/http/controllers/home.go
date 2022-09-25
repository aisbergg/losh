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
