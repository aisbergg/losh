package controllers

import (
	"context"
	"time"

	"losh/internal/core/product/models"
	"losh/internal/core/product/services"
	"losh/internal/infra/dgraph"
	"losh/web/core/search"
	"losh/web/intf/http/controllers/binding"

	"github.com/gofiber/fiber/v2"
)

var dbTimeout = 30 * time.Second

// DetailsController is the controller for the resource details page at '/details/:id'.
type DetailsController struct {
	Controller
	prdSvc        *services.Service
	searchService *search.Service
}

// NewDetailsController creates a new DetailsController.
func NewDetailsController(db *dgraph.DgraphRepository, prdSvc *services.Service, tplBndPrv binding.TemplateBindingProvider, debug bool) DetailsController {
	return DetailsController{
		Controller:    Controller{tplBndPrv: tplBndPrv},
		prdSvc:        prdSvc,
		searchService: search.NewService(db, debug),
	}
}

// Register registers the controller with the given router.
func (c DetailsController) Register(router fiber.Router) {
	router.Get("/details/:id", c.Handle)
}

// Handle handles the request for the resource details page.
func (c DetailsController) Handle(ctx *fiber.Ctx) error {
	// parse request info including params
	reqInfo, tplBnd := c.preprocessRequest(ctx, parseSearchQueryParams, parseDetailsParams)
	params := reqInfo.Params.(DetailsParams)
	queryParams := reqInfo.QueryParams.(SearchQueryParams)

	// return 404 if params are invalid
	if params.ID == "" {
		return fiber.ErrNotFound
	}

	// retrieve the resource with given ID from the database
	svcCtx, cancel := context.WithTimeout(ctx.Context(), dbTimeout)
	defer cancel()
	data, err := c.prdSvc.GetNode(svcCtx, params.ID)
	if err != nil {
		return newControllerError(err, reqInfo, "failed to render details page")
	}

	// prepare template context
	page := tplBnd["page"].(map[string]interface{})
	page["data"] = data

	// get the correct template for the resource type
	tplNme := ""
	switch data.(type) {
	case *models.Product:
		tplNme = "details-product.html"
		page["title"] = "Product Details"
		page["page-header"] = "Product Details"

	case *models.License:
		tplNme = "details-license.html"
		page["title"] = "License Details"
		page["page-header"] = "License Details"

	case *models.User, *models.Group:
		tplNme = "details-user-group.html"
		page := tplBnd["page"].(map[string]interface{})
		if u, ok := data.(*models.User); ok {
			queryParams.Query = "licensoruid:" + *u.ID
		} else if g, ok := data.(*models.Group); ok {
			queryParams.Query = "licensoruid:" + *g.ID
		}

		// export results
		if queryParams.Export != "" {
			return performExport(ctx, svcCtx, c.searchService, queryParams)
		}

		// perform search
		tplBnd["page"], err = performSearch(svcCtx, c.searchService, queryParams, page)
		if err != nil {
			return err
		}

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

func parseDetailsParams(ctx *fiber.Ctx) interface{} {
	params := DetailsParams{}

	// parse and check ID param
	params.ID = parseHexID(ctx.Params("id"))

	return params
}

type DetailsParams struct {
	ID string `liquid:"id"`
}
