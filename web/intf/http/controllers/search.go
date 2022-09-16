package controllers

import (
	"context"
	"math"
	gourl "net/url"
	"strconv"
	"strings"

	"losh/internal/infra/dgraph"
	"losh/internal/lib/util/mathutil"
	"losh/web/core/search"
	searchmodels "losh/web/core/search/models"
	"losh/web/intf/http/controllers/binding"

	"github.com/aisbergg/go-errors/pkg/errors"
	"github.com/gofiber/fiber/v2"
)

const (
	MIMECSVCharsetUTF8 = "text/csv; charset=utf-8"
	MIMETSVCharsetUTF8 = "text/tab-separated-values; charset=utf-8"
)

// SearchController is the controller for the search page at '/search'.
type SearchController struct {
	Controller
	searchService *search.Service
}

// NewSearchController creates a new SearchController.
func NewSearchController(db *dgraph.DgraphRepository, tplBndPrv binding.TemplateBindingProvider, debug bool) SearchController {
	return SearchController{
		Controller:    Controller{tplBndPrv},
		searchService: search.NewService(db, debug),
	}
}

// Register registers the controller with the given router.
func (c SearchController) Register(router fiber.Router) {
	router.Get("/search", c.Handle)
}

// Handle handles the request for the search page.
func (c SearchController) Handle(ctx *fiber.Ctx) error {
	reqInfo, tplBnd := c.preprocessRequest(ctx, parseSearchQueryParams, nil)

	// parse search query
	svcCtx, cancel := context.WithTimeout(ctx.Context(), dbTimeout)
	defer cancel()
	params := reqInfo.QueryParams.(SearchQueryParams)

	// export results
	if params.Export != "" {
		return performExport(ctx, svcCtx, c.searchService, params)
	}

	// search and display results page
	var err error
	page := tplBnd["page"].(map[string]interface{})
	page["title"] = "Search"
	page["menu"] = "search"
	page["page-header"] = "Product Search"
	tplBnd["page"], err = performSearch(svcCtx, c.searchService, params, page)
	if err != nil {
		return err
	}

	return ctx.Render("search", tplBnd)
}

func parseSearchQueryParams(ctx *fiber.Ctx) interface{} {
	params := SearchQueryParams{}
	ctx.QueryParser(&params)
	if params.Page < 1 {
		params.Page = 1
	}
	if params.ResultsPerPage < 1 || params.ResultsPerPage > 100 {
		params.ResultsPerPage = 10
	}
	params.Order = strings.ToLower(strings.TrimSpace(params.Order))
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

func performSearch(ctx context.Context, searchService *search.Service, queryParams SearchQueryParams, pageBinding map[string]interface{}) (map[string]interface{}, error) {
	first := mathutil.Max(1, queryParams.ResultsPerPage)
	offset := mathutil.Max(0, (queryParams.Page-1)*queryParams.ResultsPerPage)
	results, err := searchService.Search(
		ctx,
		queryParams.Query,
		searchmodels.OrderByFromCombinedStr(queryParams.Order),
		searchmodels.Pagination{First: first, Offset: offset},
	)

	if err != nil {
		if serr, ok := err.(*search.Error); ok && serr.Type == search.ErrorLimitExceeded {
			pageBinding["error"] = serr.Error()
		} else {
			// if not a search error handle it as an internal server error
			return pageBinding, err
		}
	}

	numPages := int(math.Ceil(float64(results.Count) / float64(queryParams.ResultsPerPage)))
	if queryParams.Page > numPages {
		queryParams.Page = numPages
	}
	pageBinding["results"] = results
	pageBinding["curPage"] = queryParams.Page
	pageBinding["numPages"] = numPages

	return pageBinding, nil
}

func performExport(fbrCtx *fiber.Ctx, svcCtx context.Context, searchService *search.Service, queryParams SearchQueryParams) error {
	switch queryParams.Export {
	case "csv":
		b, err := searchService.ExportUpTo300Results(searchmodels.ExportTypeCSV, svcCtx, queryParams.Query, searchmodels.OrderByFromCombinedStr(queryParams.Order))
		if err != nil {
			return errors.CEWrap(err, "failed to export results as CSV").
				Add("query", queryParams.Query).
				Add("oder", queryParams.Order).
				Add("page", queryParams.Page)
		}
		fbrCtx.Response().Header.SetContentType(MIMECSVCharsetUTF8)
		fbrCtx.Response().Header.SetContentLength(len(b))
		fbrCtx.Set("Content-Disposition", "attachment; filename=results.csv")
		return fbrCtx.Send(b)

	case "tsv":
		b, err := searchService.ExportUpTo300Results(searchmodels.ExportTypeTSV, svcCtx, queryParams.Query, searchmodels.OrderByFromCombinedStr(queryParams.Order))
		if err != nil {
			return errors.CEWrap(err, "failed to export results as TSV").
				Add("query", queryParams.Query).
				Add("oder", queryParams.Order).
				Add("page", queryParams.Page)
		}
		fbrCtx.Response().Header.SetContentType(MIMETSVCharsetUTF8)
		fbrCtx.Response().Header.SetContentLength(len(b))
		fbrCtx.Set("Content-Disposition", "attachment; filename=results.tsv")
		return fbrCtx.Send(b)

	default:
		return fiber.ErrBadRequest
	}
}
