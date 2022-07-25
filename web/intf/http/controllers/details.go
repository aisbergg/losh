package controllers

import (
	"regexp"
	"strings"

	"losh/internal/core/product/models"
	"losh/internal/infra/dgraph"
	"losh/web/intf/http/controllers/binding"

	"github.com/gofiber/fiber/v2"
)

// DetailsController is the controller for the resource details page at '/details/:id'.
type DetailsController struct {
	db        *dgraph.DgraphRepository
	tplBndPrv binding.TemplateBindingProvider
}

// NewDetailsController creates a new DetailsController.
func NewDetailsController(db *dgraph.DgraphRepository, tplBndPrv binding.TemplateBindingProvider) DetailsController {
	return DetailsController{db, tplBndPrv}
}

// Register registers the controller with the given router.
func (c DetailsController) Register(router fiber.Router) {
	router.Get("/details/:id", c.Handle)
}

// Handle handles the request for the resource details page.
func (c DetailsController) Handle(ctx *fiber.Ctx) error {
	// parse request info including params
	reqInfo := parseRequestInfo(ctx, parseParamNoop, c.parseParams)

	// return 404 if params are invalid
	params := reqInfo.Params.(DetailsParams)
	if params.ID == "" {
		return fiber.ErrNotFound
	}

	// retrieve the resource with given ID from the database
	data, err := c.db.GetNode(params.ID)
	if err != nil {
		return newControllerError(err, reqInfo, "failed to render details page")
	}

	// prepare template context
	tplBnd := c.tplBndPrv.Get()
	page := tplBnd["page"].(map[string]interface{})
	page["data"] = data

	// get the correct template for the resource type
	tplNme := ""
	switch data.(type) {
	case *models.Product:
		tplNme = "details-product.html"
		page["title"] = "Product Details"
		page["page-header"] = "Product Details"
	case *models.Component:
		tplNme = "details-component.html"
		page["title"] = "Component Details"
		page["page-header"] = "Component Details"
	case *models.License:
		tplNme = "details-license.html"
		page["title"] = "License Details"
		page["page-header"] = "License Details"
	case *models.User:
		tplNme = "details-user.html"
		page["title"] = "User Details"
		page["page-header"] = "User Details"
	case *models.Group:
		tplNme = "details-group.html"
		page["title"] = "Group Details"
		page["page-header"] = "Group Details"
	default:
		return fiber.ErrNotFound
	}

	// render the template
	err = ctx.Render(tplNme, tplBnd)
	if err != nil {
		return newControllerError(err, reqInfo, "failed to render details page")
	}
	return nil
}

var idPattern = regexp.MustCompile(`^[a-z0-9]{1,16}$`)

func (DetailsController) parseParams(ctx *fiber.Ctx) interface{} {
	params := DetailsParams{}

	// parse and check ID param
	rawID := strings.ToLower(strings.TrimSpace(ctx.Params("id")))
	if idPattern.MatchString(rawID) {
		params.ID = "0x" + rawID
	}

	return params
}

type DetailsParams struct {
	ID string `liquid:"id"`
}
