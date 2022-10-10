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
	"context"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"
	"unsafe"

	"losh/internal/core/product/models"
	"losh/internal/core/product/services"
	"losh/internal/lib/util/reflectutil"
	"losh/web/intf/http/controllers/binding"

	"github.com/aisbergg/go-errors/pkg/errors"
	"github.com/anglo-korean/rdf"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/bytebufferpool"
)

const (
	MIMETextHTML = "text/html"
	MIMEAny      = "*/*"
	MIMETurtle   = "text/turtle"
)

// RDFController is the controller for the rdf resource pages at '/rdf/*'.
type RDFController struct {
	Controller
	prdSvc  *services.Service
	baseURL string
}

// NewRDFController creates a new RDFController.
func NewRDFController(svc *services.Service, tplBndPrv binding.TemplateBindingProvider, baseURL string) RDFController {
	return RDFController{
		Controller: Controller{tplBndPrv},
		prdSvc:     svc,
		baseURL:    baseURL,
	}
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
	reqInfo, tplBnd := c.preprocessRequest(ctx, nil, nil)

	id := parseHexID(ctx.Params("id"))
	if id == "" {
		return fiber.ErrNotFound
	}

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

	page := tplBnd["page"].(map[string]interface{})
	page["title"] = "Resource"
	page["id"] = id
	page["resource"] = processNode(node.(models.Node))
	v, t := extractValueAndType(node.(models.Node))
	page["page-header"] = fmt.Sprintf("%s: %s", t, v)
	page["extraHeaders"] = []string{
		`<link rel="alternate" type="text/turtle" href="` + reqInfo.BaseURL + `/rdf/data/ttl/` + id[2:] + `" />`,
	}

	return ctx.Render("rdf-page", tplBnd)
}

// HandleData handles the request for the RDF data page.
func (c RDFController) HandleData(ctx *fiber.Ctx) error {
	id := parseHexID(ctx.Params("id"))
	if id == "" {
		return fiber.ErrNotFound
	}
	rdfType := ctx.Params("type")
	var encFmt rdf.Format
	switch rdfType {
	case "ntriples":
		encFmt = rdf.NTriples
	case "ttl":
		encFmt = rdf.Turtle
	}

	// get resource from database
	svcCtx, cancel := context.WithTimeout(ctx.Context(), dbTimeout)
	defer cancel()
	node, err := c.prdSvc.GetNode(svcCtx, id)
	if err != nil {
		return err
	}
	if node == nil {
		return fiber.ErrNotFound
	}

	triples := newRDFProcessor(c.baseURL).process(node.(models.Node))
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)
	enc := rdf.NewTripleEncoder(bb, encFmt)
	err = enc.EncodeAll(triples)
	if err != nil {
		return errors.CEWrap(err, "failed to encode RDF").
			Add("type", rdfType).
			Add("id", id)
	}
	enc.Close()
	return ctx.Send(bb.B)
}

func processNode(node models.Node) []map[string]interface{} {
	ret := make([]map[string]interface{}, 0, 30)
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

var timeType = reflect.TypeOf(time.Time{})

func formatValue(value interface{}, m map[string]interface{}) {
	if reflectutil.IsNil(value) {
		m["value"] = nil
		return
	}

	if n, ok := value.(models.Node); ok {
		m["link"] = "/rdf/resource/" + (*n.GetID())[2:]
		m["value"], m["type"] = extractValueAndType(n)
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
		if vd.Type() == timeType {
			m["type"] = "datetime"
		}
	}
}

func extractValueAndType(node models.Node) (value, typ string) {
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

type rdfProcessor struct {
	baseURL string
	triples []rdf.Triple
}

func newRDFProcessor(baseURL string) *rdfProcessor {
	return &rdfProcessor{
		baseURL: baseURL,
		triples: make([]rdf.Triple, 0, 40),
	}
}

func (p *rdfProcessor) process(node models.Node) []rdf.Triple {
	val := reflect.ValueOf(node)
	valDrf := reflectutil.Indirect(val)
	flds := reflectutil.GetStructFields(valDrf)

	keys := make([]string, 0, len(flds))
	for k := range flds {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	subj, _ := rdf.NewIRI(p.baseURL + "/rdf/resource/" + *node.GetID())

	for _, k := range keys {
		fldVal := valDrf.Field(flds[k].Index)
		fldInf := fldVal.Interface()

		p.nodeToTriple(subj, k, fldInf)
	}

	return p.triples
}

func (p *rdfProcessor) nodeToTriple(subj rdf.IRI, key string, value interface{}) {
	if reflectutil.IsNil(value) {
		return
	}

	if n, ok := value.(models.Node); ok {
		pred, _ := rdf.NewIRI(key)
		obj, _ := rdf.NewIRI(p.baseURL + "/rdf/resource/" + *n.GetID())
		p.triples = append(p.triples, rdf.Triple{subj, pred, obj})
		return
	}

	v := reflect.ValueOf(value)
	vd := reflectutil.Indirect(v)
	if !vd.IsValid() {
		return
	}
	vdInf := vd.Interface()
	pred, _ := rdf.NewIRI(key)
	var obj rdf.Object

	switch vd.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < vd.Len(); i++ {
			v := vd.Index(i).Interface()
			p.nodeToTriple(subj, key, v)
		}
		return
	case reflect.String:
		obj, _ = rdf.NewLiteral(vd.String())
	default:
		obj, _ = rdf.NewLiteral(vdInf)
		fmt.Println("value", value)
		fmt.Println("obj", obj)
	}

	p.triples = append(p.triples, rdf.Triple{subj, pred, obj})
}
