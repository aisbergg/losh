package controllers

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unsafe"

	"losh/web/intf/http/controllers/binding"

	gourl "net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/utils"
)

type Controller struct {
	tplBndPrv binding.TemplateBindingProvider
}

func (c Controller) preprocessRequest(ctx *fiber.Ctx, queryParser paramParser, paramParser paramParser) (*RequestInfo, map[string]interface{}) {
	// collect request information (mostly for template use)
	reqInfo := parseRequestInfo(ctx, queryParser, paramParser)

	// tell client that hints about color scheme are accepted
	ctx.Set("Accept-CH", "Sec-CH-Prefers-Color-Scheme")
	ctx.Set("Vary", "Sec-CH-Prefers-Color-Scheme")
	ctx.Set("Critical-CH", "Sec-CH-Prefers-Color-Scheme")

	preferredColorScheme := ctx.Get("Sec-CH-Prefers-Color-Scheme")

	tplBnd := c.tplBndPrv.Get()
	tplBnd["req"] = reqInfo
	page := tplBnd["page"].(map[string]interface{})
	if preferredColorScheme == "dark" {
		page["body-class"] = "theme-dark"
	} else {
		page["body-class"] = "theme-light"
	}

	return reqInfo, tplBnd
}

type RequestInfo struct {
	BaseURL  string            `json:"baseUrl" liquid:"baseUrl"`
	URL      string            `json:"url" liquid:"url"`
	Scheme   string            `json:"scheme" liquid:"scheme"`
	Hostname string            `json:"hostname" liquid:"hostname"`
	Port     uint16            `json:"port" liquid:"port"`
	Netloc   string            `json:"netloc" liquid:"netloc"`
	Path     string            `json:"path" liquid:"path"`
	FullPath string            `json:"fullPath" liquid:"fullPath"`
	Method   string            `json:"method" liquid:"method"`
	Headers  map[string]string `json:"headers" liquid:"headers"`
	Cookies  map[string]string `json:"cookies" liquid:"cookies"`
	IP       string            `json:"ip" liquid:"ip"`
	IPs      []string          `json:"ips" liquid:"ips"`

	// Page specific query params
	Params      interface{} `liquid:"params"`
	QueryParams interface{} `liquid:"queryParams"`
}

// paramParser is a function that parses the request params or query params.
type paramParser func(ctx *fiber.Ctx) interface{}

func parseRequestInfo(ctx *fiber.Ctx, queryParser paramParser, paramParser paramParser) *RequestInfo {
	req := &RequestInfo{}
	req.BaseURL = ctx.BaseURL()
	req.URL = ctx.BaseURL() + ctx.OriginalURL()
	parsedURL, _ := gourl.ParseRequestURI(req.URL)
	req.Scheme = parsedURL.Scheme
	req.Hostname = parsedURL.Hostname()
	port, _ := strconv.ParseUint(parsedURL.Port(), 10, 16)
	req.Port = uint16(port)
	req.Netloc = fmt.Sprintf("%s:%d", req.Hostname, req.Port)
	req.Path = parsedURL.Path
	req.FullPath = utils.SafeString(ctx.OriginalURL())
	req.Method = ctx.Method()
	reqHeader := make(map[string]string)
	ctx.Request().Header.VisitAllInOrder(func(key, value []byte) {
		reqHeader[utils.UnsafeString(key)] = utils.UnsafeString(value)
	})
	req.Headers = reqHeader
	reqCookies := make(map[string]string)
	ctx.Request().Header.VisitAllCookie(func(key []byte, value []byte) {
		reqCookies[utils.UnsafeString(key)] = utils.UnsafeString(value)
	})
	req.Cookies = reqCookies
	req.IP = ctx.IP()
	req.IPs = ctx.IPs()
	if queryParser != nil {
		req.QueryParams = queryParser(ctx)
	}
	if paramParser != nil {
		req.Params = paramParser(ctx)
	}
	return req
}

var idPattern = regexp.MustCompile(`^[a-f0-9]{1,16}$`)

func parseID(id string) string {
	if id == "" {
		return ""
	}
	id = strings.ToLower(strings.TrimSpace(id))
	if idPattern.MatchString(id) {
		return id
	}
	return ""
}

func parseHexID(id string) string {
	if id == "" {
		return ""
	}
	id = strings.ToLower(strings.TrimSpace(id))
	if idPattern.MatchString(id) {
		buf := make([]byte, 0, len(id)+2)
		buf = append(buf, "0x"...)
		buf = append(buf, id...)
		return *(*string)(unsafe.Pointer(&buf))
	}
	return ""
}
