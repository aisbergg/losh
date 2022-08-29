package controllers

import (
	"context"
	"math"
	"strconv"

	"losh/internal/infra/dgraph"
	"losh/internal/lib/util/mathutil"
	"losh/web/core/search"
	searchmodels "losh/web/core/search/models"
	"losh/web/intf/http/controllers/binding"

	gourl "net/url"

	"github.com/aisbergg/go-errors/pkg/errors"
	"github.com/gofiber/fiber/v2"
)

const (
	MIMECSVCharsetUTF8 = "text/csv; charset=utf-8"
	MIMETSVCharsetUTF8 = "text/tab-separated-values; charset=utf-8"
)

// SearchController is the controller for the search page at '/search'.
type SearchController struct {
	searchService *search.Service
	tplBndPrv     binding.TemplateBindingProvider
}

// NewSearchController creates a new SearchController.
func NewSearchController(db *dgraph.DgraphRepository, tplBndPrv binding.TemplateBindingProvider, debug bool) SearchController {
	return SearchController{search.NewService(db, debug), tplBndPrv}
}

// Register registers the controller with the given router.
func (c SearchController) Register(router fiber.Router) {
	router.Get("/search", c.Handle)
}

// Handle handles the request for the search page.
func (c SearchController) Handle(ctx *fiber.Ctx) error {
	// get request information (for template use)
	reqInfo := parseRequestInfo(ctx, c.parseQueryParams, nil)

	// tell client that hints about color scheme are accepted
	ctx.Set("Accept-CH", "Sec-CH-Prefers-Color-Scheme")
	ctx.Set("Vary", "Sec-CH-Prefers-Color-Scheme")
	ctx.Set("Critical-CH", "Sec-CH-Prefers-Color-Scheme")

	preferredColorScheme := ctx.Get("Sec-CH-Prefers-Color-Scheme")

	// parse search query
	svcCtx, cancel := context.WithTimeout(ctx.Context(), dbTimeout)
	defer cancel()
	params := reqInfo.QueryParams.(SearchQueryParams)
	first := mathutil.Max(1, params.ResultsPerPage)
	offset := mathutil.Max(0, (params.Page-1)*params.ResultsPerPage)
	results, err := c.searchService.Search3(svcCtx, params.Query, searchmodels.OrderByFromCombinedStr(params.Order), searchmodels.Pagination{First: first, Offset: offset})
	if err != nil {
		if serr, ok := err.(*search.Error); ok && serr.Type == search.ErrorLimitExceeded {
			tplBnd := c.tplBndPrv.Get()
			tplBnd["req"] = reqInfo
			page := tplBnd["page"].(map[string]interface{})
			page["title"] = "Search"
			page["menu"] = "search"
			page["error"] = serr.Error()
			page["results"] = 0
			page["curPage"] = 1
			page["numPages"] = 0
			if preferredColorScheme == "dark" {
				page["body-class"] = "theme-dark"
			} else {
				page["body-class"] = "theme-light"
			}
			return ctx.Render("search", tplBnd)
		}
		// if not a search error handle it as an internal server error
		return err
	}

	// export results
	if params.Export != "" {
		switch params.Export {
		case "csv":
			b, err := c.searchService.ExportUpTo300Results(searchmodels.ExportTypeCSV, svcCtx, params.Query, searchmodels.OrderByFromCombinedStr(params.Order))
			if err != nil {
				return errors.CEWrap(err, "failed to export results as CSV").
					Add("query", params.Query).
					Add("oder", params.Order).
					Add("page", params.Page)
			}
			ctx.Response().Header.SetContentType(MIMECSVCharsetUTF8)
			ctx.Response().Header.SetContentLength(len(b))
			ctx.Set("Content-Disposition", "attachment; filename=results.csv")
			return ctx.Send(b)

		case "tsv":
			b, err := c.searchService.ExportUpTo300Results(searchmodels.ExportTypeTSV, svcCtx, params.Query, searchmodels.OrderByFromCombinedStr(params.Order))
			if err != nil {
				return errors.CEWrap(err, "failed to export results as TSV").
					Add("query", params.Query).
					Add("oder", params.Order).
					Add("page", params.Page)
			}
			ctx.Response().Header.SetContentType(MIMETSVCharsetUTF8)
			ctx.Response().Header.SetContentLength(len(b))
			ctx.Set("Content-Disposition", "attachment; filename=results.tsv")
			return ctx.Send(b)
		}
		// invalid export type -> continue serving search page
	}

	// display results page

	// TODO: update query params with new values
	numPages := int(math.Ceil(float64(results.Count) / float64(params.ResultsPerPage)))
	if params.Page > numPages {
		params.Page = numPages
	}

	tplBnd := c.tplBndPrv.Get()
	tplBnd["req"] = reqInfo
	page := tplBnd["page"].(map[string]interface{})
	page["title"] = "Search"
	page["menu"] = "search"
	page["page-header"] = "Product Search"
	page["results"] = results
	page["curPage"] = params.Page
	page["numPages"] = numPages

	if preferredColorScheme == "dark" {
		page["body-class"] = "theme-dark"
	} else {
		page["body-class"] = "theme-light"
	}

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
	DisplayMode    string `query:"dm" json:"displayMode" liquid:"displayMode"`
	Export         string `query:"export" json:"export" liquid:"export"`
}

func (p SearchQueryParams) String() string {
	v := gourl.Values{}
	v.Set("q", p.Query)
	v.Set("o", p.Order)
	v.Set("page", strconv.Itoa(p.Page))
	v.Set("rpp", strconv.Itoa(p.ResultsPerPage))
	return v.Encode()
}
