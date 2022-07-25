package controllers

import (
	"losh/internal/infra/dgraph"
	"losh/web/intf/http/controllers/binding"
	"strconv"

	gourl "net/url"

	"github.com/gofiber/fiber/v2"
)

// SearchController is the controller for the search page at '/search'.
type SearchController struct {
	db        *dgraph.DgraphRepository
	tplBndPrv binding.TemplateBindingProvider
}

// NewSearchController creates a new SearchController.
func NewSearchController(db *dgraph.DgraphRepository, tplBndPrv binding.TemplateBindingProvider) SearchController {
	return SearchController{db, tplBndPrv}
}

func (c SearchController) Register(router fiber.Router) {
	router.Get("/search", c.Handle)
}

func (c SearchController) Handle(ctx *fiber.Ctx) error {
	// get request information (for template use)
	reqInfo := parseRequestInfo(ctx, c.parseQueryParams, parseParamNoop)

	// params := new(rawQueryParams)
	// ctx.QueryParser(params)

	// if err := ctx.QueryParser(params); err != nil {
	// 	// ignore errors
	// 	// return err
	// }
	// fmt.Println("params:", params)

	// pageInfo := map[string]interface{}{
	// 	"title":                   "Search",
	// 	"menu":                    "search",
	// 	"layout-navbar-condensed": true,
	// 	"container-centered":      true,
	// }

	tplBnd := c.tplBndPrv.Get()
	page := tplBnd["page"].(map[string]interface{})
	page["title"] = "Search"
	page["menu"] = "search"
	page["page-header"] = "Product Search"
	tplBnd["req"] = reqInfo

	return ctx.Render("search", tplBnd)
}

func (SearchController) parseQueryParams(ctx *fiber.Ctx) interface{} {
	params := SearchQueryParams{}
	ctx.QueryParser(&params)

	if params.Page < 1 {
		params.Page = 1
	}

	if params.ResultsPerPage < 1 || params.ResultsPerPage > 100 {
		params.ResultsPerPage = 10
	}

	return params
}

type SearchQueryParams struct {
	Query          string `query:"q" json:"query" liquid:"query"`
	Order          string `query:"o" json:"order" liquid:"order"`
	Page           int    `query:"page" json:"page" liquid:"page"`
	ResultsPerPage int    `query:"rpp" json:"resultsPerPage" liquid:"resultsPerPage"`
}

func (p SearchQueryParams) String() string {
	v := gourl.Values{}
	v.Set("q", p.Query)
	v.Set("o", p.Order)
	v.Set("page", strconv.Itoa(p.Page))
	v.Set("rpp", strconv.Itoa(p.ResultsPerPage))
	return v.Encode()
}
