package controllers

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"unsafe"

	"losh/internal/core/product/models"
	"losh/internal/core/product/services"
	"losh/internal/lib/util/reflectutil"
	"losh/web/intf/http/controllers/binding"

	"github.com/gofiber/fiber/v2"
)

const (
	MIMETextHTML = "text/html"
	MIMEAny      = "*/*"
	MIMETurtle   = "text/turtle"
)

// RDFController is the controller for the rdf resource pages at '/rdf/*'.
type RDFController struct {
	prdSvc    *services.Service
	tplBndPrv binding.TemplateBindingProvider
}

// NewRDFController creates a new RDFController.
func NewRDFController(svc *services.Service, tplBndPrv binding.TemplateBindingProvider) RDFController {
	return RDFController{svc, tplBndPrv}
}

// Register registers the controller with the given router.
func (c RDFController) Register(router fiber.Router) {
	rdfRoutes := router.Group("/rdf")
	rdfRoutes.Get("/resource/:id", c.HandleResource)
	rdfRoutes.Get("/data/:type/:id", c.HandleData)
	rdfRoutes.Get("/page/:id", c.HandlePage)
}

// HandleResource handles the request for the RDF resource.
func (c RDFController) HandleResource(ctx *fiber.Ctx) error {
	id := parseID(ctx.Params("id"))
	if id == "" {
		return fiber.ErrNotFound
	}

	//  ctx.Accepts(MIMETextHTML, MIMEAny)
	accept := ctx.Get(fiber.HeaderAccept)
	if accept == "" {
		accept = MIMEAny
	}
	if ctx.Accepts(MIMETextHTML, MIMEAny) != "" {
		red := make([]byte, 0, len(id)+10)
		red = append(red, "/rdf/page/"...)
		red = append(red, id...)
		return ctx.Redirect(*(*string)(unsafe.Pointer(&red)), fiber.StatusSeeOther)

	} else if ctx.Accepts(MIMETurtle) != "" {
		red := make([]byte, 0, len(id)+14)
		red = append(red, "/rdf/data/ttl/"...)
		red = append(red, id...)
		return ctx.Redirect(*(*string)(unsafe.Pointer(&red)), fiber.StatusSeeOther)
	}

	return ctx.Status(fiber.StatusNotAcceptable).SendString("Not Acceptable")
}

// HandlePage handles the request for the RDF HTML page.
func (c RDFController) HandlePage(ctx *fiber.Ctx) error {
	id := parseHexID(ctx.Params("id"))
	if id == "" {
		return fiber.ErrNotFound
	}

	// tell client that hints about color scheme are accepted
	ctx.Set("Accept-CH", "Sec-CH-Prefers-Color-Scheme")
	ctx.Set("Vary", "Sec-CH-Prefers-Color-Scheme")
	ctx.Set("Critical-CH", "Sec-CH-Prefers-Color-Scheme")
	preferredColorScheme := ctx.Get("Sec-CH-Prefers-Color-Scheme")

	// get node from database
	svcCtx, cancel := context.WithTimeout(ctx.Context(), dbTimeout)
	defer cancel()
	node, err := c.prdSvc.GetNode(svcCtx, id)
	if err != nil {
		return err
	}
	if node == nil {
		return fiber.ErrNotFound
	}

	tplBnd := c.tplBndPrv.Get()
	page := tplBnd["page"].(map[string]interface{})
	page["title"] = "Resource"
	page["resource"] = processNode(node.(models.Node))
	v, t := nodeValueType(node.(models.Node))
	page["page-header"] = fmt.Sprintf("%s: %s", t, v)
	if preferredColorScheme == "dark" {
		page["body-class"] = "theme-dark"
	} else {
		page["body-class"] = "theme-light"
	}

	return ctx.Render("resource", tplBnd)
}

// HandleData handles the request for the RDF data page.
func (c RDFController) HandleData(ctx *fiber.Ctx) error {
	id := parseHexID(ctx.Params("id"))
	if id == "" {
		return fiber.ErrNotFound
	}
	rdfType := ctx.Params("type")
	switch rdfType {
	case "ttl":
		// get resource from database
		svcCtx, cancel := context.WithTimeout(ctx.Context(), dbTimeout)
		defer cancel()
		node, err := c.prdSvc.GetNode(svcCtx, id)
		if err != nil {
			return err
		}
		if node == nil {
			return ctx.SendStatus(fiber.StatusNotFound)
		}

		// TODO: implement
		return ctx.JSON("something")
		// return ctx.SendFile(c.prdSvc.GetProductRDF(id), MIMETurtle)
	}

	return fiber.ErrNotFound
}

func processNode(node models.Node) []map[string]interface{} {
	ret := make([]map[string]interface{}, 0, 40)
	val := reflect.ValueOf(node)
	valDrf := reflectutil.Indirect(val)
	flds := reflectutil.GetStructFields(valDrf)

	keys := make([]string, 0, len(flds))
	for k := range flds {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fldVal := valDrf.Field(flds[k].Index)
		fldInf := fldVal.Interface()

		m := make(map[string]interface{})
		m["name"] = k
		formatValue(fldInf, m)
		ret = append(ret, m)
	}

	return ret
}

func formatValue(value interface{}, m map[string]interface{}) {
	if reflectutil.IsNil(value) {
		m["value"] = nil
		return
	}

	n, ok := value.(models.Node)
	if ok {
		m["link"] = "/rdf/resource/" + (*n.GetID())[2:]
		m["value"], m["type"] = nodeValueType(n)
		return
	}

	v := reflect.ValueOf(value)
	vd := reflectutil.Indirect(v)
	if !vd.IsValid() {
		m["value"] = nil
		return
	}

	switch vd.Kind() {
	case reflect.String:
		s := vd.String()
		m["value"] = s
		if strings.HasPrefix(s, "http") {
			m["link"] = s
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		m["value"] = fmt.Sprintf("%v", vd.Int())
		m["type"] = "integer"
	case reflect.Float32, reflect.Float64:
		m["value"] = fmt.Sprintf("%v", vd.Float())
		m["type"] = "float"
	case reflect.Complex64, reflect.Complex128:
		m["value"] = fmt.Sprintf("%v", vd.Complex())
		m["type"] = "complex"
	case reflect.Bool:
		if vd.Bool() {
			m["value"] = "true"
		} else {
			m["value"] = "false"
		}
		m["type"] = "bool"
	case reflect.Slice, reflect.Array:
		l := make([]map[string]interface{}, 0, vd.Len())
		for i := 0; i < vd.Len(); i++ {
			v := vd.Index(i).Interface()
			m := make(map[string]interface{})
			formatValue(v, m)
			l = append(l, m)
		}
		m["value"] = l
		m["type"] = "list"
	default:
		m["value"] = fmt.Sprintf("%v", vd.Interface())
	}
}

func nodeValueType(node models.Node) (value, typ string) {
	switch t := node.(type) {
	case *models.BoundingBoxDimensions:
		value = fmt.Sprintf("%fm x %fm x %fm", *t.Width, *t.Height, *t.Depth)
		typ = "BoundingBoxDimensions"
	case *models.Category:
		value = *t.FullName
		typ = "Category"
	case *models.Component:
		value = *t.Name
		typ = "Component"
	case *models.File:
		value = *t.Path
		typ = "File"
	case *models.FloatV:
		value = fmt.Sprintf("%f", *t.Value)
		typ = "FloatV"
	case *models.Group:
		value = *t.Name
		typ = "Group"
	case *models.Host:
		value = *t.Name
		typ = "Host"
	case *models.KeyValue:
		switch tv := t.Value.(type) {
		case *models.FloatV:
			value = fmt.Sprintf("%s: %f", *t.Key, *tv.Value)
			typ = "FloatV"
		case *models.StringV:
			value = fmt.Sprintf("%s: %s", *t.Key, *tv.Value)
			typ = "StringV"
		}
	case *models.License:
		value = *t.Xid
		typ = "License"
	case *models.ManufacturingProcess:
		value = *t.Name
		typ = "ManufacturingProcess"
	case *models.Material:
		value = *t.Name
		typ = "Material"
	case *models.OpenSCADDimensions:
		value = fmt.Sprintf("%s %s", *t.Openscad, *t.Unit)
		typ = "OpenSCADDimensions"
	case *models.Product:
		value = *t.Name
		typ = "Product"
	case *models.Repository:
		value = *t.Name
		typ = "Repository"
	case *models.Software:
		value = *t.ID
		typ = "Software"
	case *models.StringV:
		value = *t.Value
		typ = "StringV"
	case *models.Tag:
		value = *t.Name
		typ = "Tag"
	case *models.TechnicalStandard:
		value = *t.Name
		typ = "TechnicalStandard"
	case *models.TechnologySpecificDocumentationCriteria:
		value = *t.Name
		typ = "TechnologySpecificDocumentationCriteria"
	case *models.User:
		if t.FullName != nil {
			value = *t.FullName
		} else {
			value = *t.Name
		}
		typ = "User"
	default:
		panic(fmt.Sprintf("unsupported node type: %T", node))
	}
	return
}
