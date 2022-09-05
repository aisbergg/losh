package dgraph

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	productmodels "losh/internal/core/product/models"
	searchmodels "losh/web/core/search/models"
	"losh/web/core/search/parser"

	"github.com/aisbergg/go-errors/pkg/errors"
	"github.com/golang-module/carbon/v2"
)

const selectQueryFragment = `
q(func: uid(%s), first: $first, offset: $offset, %s) {
	uid
	CrawlerMeta.discoveredAt
	CrawlerMeta.lastIndexedAt
	Product.xid
	Product.name
	Product.website
	Product.state
	Product.lastUpdatedAt
	Product.renamedTo {uid}
	Product.renamedFrom {uid}
	Product.forkOf {
		uid Product.release {
			Component.repository {
				Repository.url
			}
		}
	}
	Product.forks {uid}
	Product.forkCount
	Product.starCount
	Product.releases {uid}
	Product.tags {
		uid
		Tag.name
		Tag.aliases {
			uid
			Tag.name
		}
		Tag.related {
			uid
			Tag.name
		}
	}
	Product.category {
		uid
		Category.xid
		Category.fullName
		Category.name
		Category.description
		Category.parent{
			uid
			Category.fullName
		}
		Category.children {
			uid
			Category.fullName
		}
		Category.products {
			uid
			Product.name
		}
	}
	Product.release {
		Component.id
		Component.xid
		Component.name
		Component.description
		Component.version
		Component.createdAt
		Component.releases {uid}
		Component.isLatest
		Component.documentationLanguage
		Component.technologyReadinessLevel
		Component.documentationReadinessLevel
		Component.attestation
		Component.publication
		Component.compliesWith {TechnicalStandard.name}
		Component.cpcPatentClass Component.tsdc {uid}
		Component.components {uid}
		Component.software {uid}
		Component.product {uid}
		Component.usedIn {uid}
		Component.organization {uid}
		Component.mass Component.material {uid}
		Component.manufacturingProcess {uid}
		Component.productionMetadata {uid}
		Component.repository {
			uid
			Repository.id
			Repository.xid
			Repository.url
			Repository.permaUrl
			Repository.host {
				uid
				Host.name
			}
			Repository.name
			Repository.reference
			Repository.path
		}
		Component.license {
			uid License.xid
			License.name
			License.isSpdx
			License.isDeprecated
			License.isOsiApproved
			License.isFsfLibre
			License.isBlocked
		}
		Component.additionalLicenses {
			uid
			License.xid
			License.name
			License.isSpdx
			License.isDeprecated
			License.isOsiApproved
			License.isFsfLibre
			License.isBlocked
		}
		Component.licensor {
			dgraph.type
			uid
			UserOrGroup.host {
				uid
				Host.name
			}
			UserOrGroup.name
			UserOrGroup.fullName
			UserOrGroup.email
			UserOrGroup.description
			UserOrGroup.avatar {
				uid
				File.path
			}
			UserOrGroup.url
			UserOrGroup.memberOf {
				uid
				Group.fullName
			}
			UserOrGroup.products {
				uid
				Product.name
			}
			User.locale
			Group.members {
				dgraph.type
				uid
			}
		}

		Component.image {
			uid
			CrawlerMeta.discoveredAt
			CrawlerMeta.lastIndexedAt
			File.name
			File.path
			File.mimeType
			File.url
			File.createdAt
		}

		Component.readme {
			uid
			CrawlerMeta.discoveredAt
			CrawlerMeta.lastIndexedAt
			File.name
			File.path
			File.mimeType
			File.url
			File.createdAt
		}

		Component.contributionGuide {
			uid
			CrawlerMeta.discoveredAt
			CrawlerMeta.lastIndexedAt
			File.name
			File.path
			File.mimeType
			File.url
			File.createdAt
		}

		Component.bom {
			uid
			CrawlerMeta.discoveredAt
			CrawlerMeta.lastIndexedAt
			File.name
			File.path
			File.mimeType
			File.url
			File.createdAt
		}

		Component.manufacturingInstructions {
			uid
			CrawlerMeta.discoveredAt
			CrawlerMeta.lastIndexedAt
			File.name
			File.path
			File.mimeType
			File.url
			File.createdAt
		}

		Component.userManual {
			uid
			CrawlerMeta.discoveredAt
			CrawlerMeta.lastIndexedAt
			File.name
			File.path
			File.mimeType
			File.url
			File.createdAt
		}

		Component.source {
			uid
			CrawlerMeta.discoveredAt
			CrawlerMeta.lastIndexedAt
			File.name
			File.path
			File.mimeType
			File.url
			File.createdAt
		}

		Component.export {
			uid
			CrawlerMeta.discoveredAt
			CrawlerMeta.lastIndexedAt
			File.name
			File.path
			File.mimeType
			File.url
			File.createdAt
		}

		Component.auxiliary {
			uid
			CrawlerMeta.discoveredAt
			CrawlerMeta.lastIndexedAt
			File.name
			File.path
			File.mimeType
			File.url
			File.createdAt
		}

		Component.outerDimensions {
			dgraph.type
			uid
			BoundingBoxDimensions.height
			BoundingBoxDimensions.width
			BoundingBoxDimensions.depth
			OpenSCADDimensions.openscad
			OpenSCADDimensions.unit
		}
	}
}`

func (dr *DgraphRepository) SearchProducts(ctx context.Context, query *parser.Query, order searchmodels.OrderBy, pagination searchmodels.Pagination) ([]*productmodels.Product, uint64, error) {
	dr.log.Debugw("search Products")

	q, v := createDQLQuery(query, order, pagination)
	// fmt.Println("q: ", q)
	// fmt.Println("v: ", v)

	rsp, err := dr.dgraphClient.NewTxn().QueryWithVars(ctx, q, v)
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to execute query")
	}

	// bind results to map - will be copied into data model at the end
	var rspData map[string]interface{}
	err = json.Unmarshal(rsp.Json, &rspData)
	if err != nil {
		return nil, 0, err
	}

	rawRes := rspData["q"].([]interface{})
	if len(rawRes) == 0 {
		return nil, 0, nil
	}
	ret := make([]*productmodels.Product, 0, len(rawRes))
	if err = dr.dqlCopier.CopyTo(rawRes, &ret); err != nil {
		panic(err)
	}

	return ret, rsp.Metrics.NumUids["Product.xid"], nil
}

// createDQLQuery creates a DQL query from a search query.
//
// I first tried to use github.com/fenos/dqlx to programmatically build a query. It was cumbersome but the actual deal breaker was its bugginess. Therefore I crafted a query manually.
func createDQLQuery(query *parser.Query, order searchmodels.OrderBy, pagination searchmodels.Pagination) (q string, v map[string]string) {
	encoder := newEncoder()
	lastVar := encoder.encodeQuery(query, "")
	if lastVar == "" {
		// if no root query was created, we need to create a query for returning all products
		lastVar = encoder.addVariableWithFilter("", "", "")
	}
	lastVar = encoder.appendOrderByVariable(order, lastVar)

	ordFrg := ""
	if order.Descending {
		// XXX: multiple sorting works only for predicates, not for variables (second term in this case will be ignored); if you change the order of those two terms, the query will fail telling you the fact I just stated
		ordFrg = "orderdesc: val(order), orderasc: Product.name"
	} else {
		ordFrg = "orderasc: val(order), orderasc: Product.name"
	}

	encoder.appendSelectQuery(selectQueryFragment, ordFrg, lastVar)

	encoder.addArg("$first", strconv.Itoa(pagination.First), "int")
	encoder.addArg("$offset", strconv.Itoa(pagination.Offset), "int")

	q = encoder.String()
	v = encoder.getArgs()

	return q, v
}

type encoder struct {
	buf *bytes.Buffer

	// encoder state
	args    [][3]string
	lastArg int
	lastVar int

	// limits
	nodes     int
	words     int
	wildcards int

	errors []string
}

// newEncoder returns a new encoder.
func newEncoder() *encoder {
	return &encoder{
		buf: new(bytes.Buffer),
	}
}

func (e *encoder) String() string {
	argsFmt := make([]string, 0, len(e.args))
	for _, arg := range e.args {
		argsFmt = append(argsFmt, fmt.Sprintf("%s: %s", arg[0], arg[2]))
	}
	return fmt.Sprintf(`query q(%s) {
%s
}`, strings.Join(argsFmt, ", "), e.buf.String())
}

func (e *encoder) encode(v interface{}) {
}

func (e *encoder) CrateArg(val, typ string) string {
	e.lastArg++
	argNme := fmt.Sprintf("$a%d", e.lastArg)
	e.args = append(e.args, [3]string{argNme, val, typ})
	return argNme
}

func (e *encoder) addArg(name, val, typ string) {
	e.args = append(e.args, [3]string{name, val, typ})
}

func (e *encoder) createVar() string {
	e.lastVar++
	varNme := fmt.Sprintf("v%d", e.lastVar)
	return varNme
}

func (e *encoder) getArgs() map[string]string {
	ret := make(map[string]string, len(e.args))
	for _, arg := range e.args {
		ret[arg[0]] = arg[1]
	}
	return ret
}

func (e *encoder) encodeQuery(query *parser.Query, parVar string) (curVar string) {
	if query == nil {
		// add an empty filter to return all results
		curVar = e.addVariableWithFilter("", "", parVar)
		return
	}
	curVar = parVar

	// flatten if only one OR condition
	if len(query.Or) == 1 {
		andCnds := query.Or[0].And
		for _, andCnd := range andCnds {
			lastVar := e.encodeAndCondition(andCnd, curVar)
			if lastVar != "" {
				curVar = lastVar
			}
		}
		return
	}

	// encode OR conditions
	vars := make([]string, 0, len(query.Or))
	for _, orCnd := range query.Or {
		for _, andCnd := range orCnd.And {
			v := e.encodeAndCondition(andCnd, curVar)
			if v != "" {
				vars = append(vars, v)
			}
		}
	}
	if len(vars) == 0 {
		return
	}
	// create union
	return e.appendUnionVariable(vars...)
}

func (e *encoder) encodeAndCondition(andCnd *parser.AndCondition, parVar string) (curVar string) {
	if andCnd.Not != nil {
		notVar := e.encodeAndCondition(andCnd.Not, parVar)
		return e.appendNegationVariable(parVar, notVar)
	}

	return e.encodeExpression(andCnd.Operand, parVar)
}

var multGlobPattern = regexp.MustCompile(`\*+`)

// encodeExpression encodes an expression into DQL variable and returns the
// variable name.
func (e *encoder) encodeExpression(expr *parser.Expression, parVar string) (curVar string) {
	if expr.Text != nil {
		// extract the value as text
		text := ""
		if expr.Text.Words != nil {
			text = *expr.Text.Words
		} else {
			text = *expr.Text.Exact
		}
		if text == "" {
			return
		}
		nameOpr := operators["name"]
		descOpr := operators["description"]
		tagOpr := operators["tag"]
		filter := e.generateOrFilter(
			e.generateTextFilter(nameOpr.Predicate, nameOpr.Type, text, false, false),
			e.generateTextFilter(descOpr.Predicate, descOpr.Type, text, false, false),
		)
		if filter == nil {
			return
		}
		filter = e.generateFilterExpression(filter)
		var1 := e.addVariableWithFilter(string(filter), "", parVar)

		tagfilter := e.generateTextFilter(tagOpr.Predicate, tagOpr.Type, text, false, false)
		tagfilter = e.generateFilterExpression(tagfilter)
		sel := fmt.Sprintf(`%s %s %s`, tagOpr.SelectionStart, string(tagfilter), tagOpr.SelectionEnd)
		var2 := e.addVariableWithFilter("", sel, parVar)

		curVar = e.appendUnionVariable(var1, var2)
		return

	} else if expr.Operator != nil {
		// operator search
		return e.encodeOperator(expr.Operator, parVar)

	} else if expr.Sub != nil {
		// more deeply nested expression
		return e.encodeQuery(expr.Sub, parVar)
	}

	return
}

func (e *encoder) encodeOperator(opr *parser.Operator, parVar string) (curVar string) {
	oprName := strings.ToLower(opr.Name)

	// special operators
	switch oprName {
	case "is":
		if opr.Value == nil {
			return ""
		}
		if opr.Value.Exact != nil {
			oprName = *opr.Value.Exact
		} else {
			oprName = *opr.Value.Words
		}
		oprName = strings.ToLower(oprName)
		o, ok := operators["is"+oprName]
		if !ok {
			return
		}
		if o.Type != booleanIsOperator {
			return
		}
		curVar = e.appendVariable2(o.IsRootFilter, func() { e.appendBooleanIsFilter(o.Predicate, o.Value, false) }, o.SelectionStart, o.SelectionEnd, parVar)
		return

	case "has":
		if opr.Value == nil {
			return
		}
		if opr.Value.Exact != nil {
			oprName = *opr.Value.Exact
		} else {
			oprName = *opr.Value.Words
		}
		oprName = strings.ToLower(oprName)
		o, ok := operators["has"+oprName]
		if !ok {
			return
		}
		if o.Type != booleanHasOperator {
			return ""
		}
		curVar = e.appendVariable2(o.IsRootFilter, func() { e.appendBooleanHasFilter(o.Predicate, false) }, o.SelectionStart, o.SelectionEnd, parVar)
		return
	}

	// named operators
	o, ok := operators[oprName]
	if !ok {
		return
	}
	switch o.Type {
	case uidOperator:
		text, _, _, _ := e.extractOperatorText(opr)
		if text == "" {
			return
		}
		filter := e.generateUIDFilter(text)
		if filter == nil || len(filter) == 0 {
			return
		}
		filter = e.generateFilterExpression(filter)
		if o.IsRootFilter {
			sel := fmt.Sprintf(`%s %s`, o.SelectionStart, o.SelectionEnd)
			curVar = e.addVariableWithFilter(string(filter), sel, parVar)
			return
		}

		sel := fmt.Sprintf(`%s %s %s`, o.SelectionStart, string(filter), o.SelectionEnd)
		curVar = e.addVariableWithFilter("", sel, parVar)
		return

	case booleanIsOperator, booleanHasOperator:
		// ignore; handled by special operators
		return

	case textFullContainsOperator, textTermExactOperator, textTermContainsOperator, textExactOperator:
		text, exact, match, not := e.extractOperatorText(opr)
		if text == "" {
			return
		}
		filter := e.generateTextFilter(o.Predicate, o.Type, text, exact, match)
		if filter == nil || len(filter) == 0 {
			return
		}
		if not {
			filter = e.generateNotFilter(filter)
		}
		filter = e.generateFilterExpression(filter)
		if o.IsRootFilter {
			sel := fmt.Sprintf(`%s %s`, o.SelectionStart, o.SelectionEnd)
			curVar = e.addVariableWithFilter(string(filter), sel, parVar)
			return
		}

		sel := fmt.Sprintf(`%s %s %s`, o.SelectionStart, string(filter), o.SelectionEnd)
		curVar = e.addVariableWithFilter("", sel, parVar)
		return

	case numberFloatOperator, numberIntOperator:
		return e.appendVariable2(o.IsRootFilter, func() { e.appendNumberFilter(o.Predicate, opr, o.Type == numberIntOperator) }, o.SelectionStart, o.SelectionEnd, parVar)

	case dateTimeOperator:
		return e.appendVariable2(o.IsRootFilter, func() { e.appendDateTimeFilter(o.Predicate, opr) }, o.SelectionStart, o.SelectionEnd, parVar)

	default:
		// should never happen unless we missed something
		panic("unsupported operator type")
	}
}

type operatorType int

const (
	uidOperator operatorType = iota
	// Matches either full-text or phrases contained in the text.
	// Requires indexes: fulltext, regexp, exacti
	//   opr:text or opr:`text` -> full-text-match (alloftext)
	//   opr:"text" -> contains-match (regexp: /.../i)
	//   opr:* or opr:`text*` -> wildcard-contains-match (regexp: /.../i and alloftext)
	//   opr:"*" -> wildcard-contains-match (regexp: /.../i)
	//   opr:==text -> exact-match (allof with exacti index)
	textFullContainsOperator
	// Matches either terms or the exact text.
	// Requires indexes: term, regexp, exacti
	//   opr:text or opr:`text` -> term-match (allofterms)
	//   opr:"text" or opr:==text -> exact-match (allof with exacti index)
	//   opr:* or opr:`text*` -> wildcard-contains-match (regexp: /.../i and allofterms)
	//   opr:"*" -> wildcard exact-match (regexp: /^...$/)
	textTermExactOperator
	// Matches either terms or phrases contained in the text.
	// Requires indexes: term, regexp, exacti
	//   opr:text or opr:`text` -> term-match (allofterms)
	//   opr:"text" -> contains-match (regexp: /.../i)
	//   opr:* or opr:`text*` -> wildcard-contains-match (regexp: /.../i and allofterms)
	//   opr:"*" -> wildcard-contains-match (regexp: /.../i)
	//   opr:==text -> exact-match (allof with exacti index)
	textTermContainsOperator
	// Matches the exact text (case insensitive).
	// Requires indexes: regexp, exacti
	//   opr:text or opr:`text` or opr:"text" or opr:==text -> exact-match (allof with exacti index)
	//   opr:* or opr:`*` or opr:"*" -> wildcard exact-match (regexp: /^...$/i)
	textExactOperator
	numberFloatOperator
	numberIntOperator
	booleanIsOperator
	booleanHasOperator
	dateTimeOperator
)

type operator struct {
	Type           operatorType
	IsRootFilter   bool
	Predicate      string
	SelectionStart string
	SelectionEnd   string

	// used for bool
	Value string
}

var operators = map[string]operator{
	//
	// Product
	//
	"name": {
		Type:           textFullContainsOperator,
		IsRootFilter:   true,
		Predicate:      "Product.name",
		SelectionStart: "uid",
	},
	"description": {
		Type:           textFullContainsOperator,
		IsRootFilter:   true,
		Predicate:      "Product.description",
		SelectionStart: "uid",
	},
	"documentationlanguage": {
		Type:           textExactOperator,
		Predicate:      "Component.documentationLanguage",
		SelectionStart: `Product.release`,
		SelectionEnd:   `{uid}`,
	},
	"language": {
		Type:           textExactOperator,
		Predicate:      "Component.documentationLanguage",
		SelectionStart: `Product.release`,
		SelectionEnd:   `{uid}`,
	},
	"version": {
		Type:           textExactOperator,
		IsRootFilter:   true,
		Predicate:      "Product.version",
		SelectionStart: "uid",
	},
	"website": {
		Type:           textFullContainsOperator,
		IsRootFilter:   true,
		Predicate:      "Product.website",
		SelectionStart: "uid",
	},
	"starcount": {
		Type:           numberIntOperator,
		IsRootFilter:   true,
		Predicate:      "Product.starCount",
		SelectionStart: "uid",
	},
	"forkcount": {
		Type:           numberIntOperator,
		IsRootFilter:   true,
		Predicate:      "Product.forkCount",
		SelectionStart: "uid",
	},
	"discoveredat": {
		Type:           dateTimeOperator,
		IsRootFilter:   true,
		Predicate:      "CrawlerMeta.discoveredAt",
		SelectionStart: `uid`,
	},
	"lastindexedat": {
		Type:           dateTimeOperator,
		IsRootFilter:   true,
		Predicate:      "CrawlerMeta.lastIndexedAt",
		SelectionStart: `uid`,
	},
	"lastupdatedat": {
		Type:           dateTimeOperator,
		IsRootFilter:   true,
		Predicate:      "Product.lastUpdatedAt",
		SelectionStart: "uid",
	},
	"releasecount": {
		Type:           numberIntOperator,
		IsRootFilter:   true,
		Predicate:      "count(Product.releases)",
		SelectionStart: "uid",
	},

	// state
	"isactive": {
		Type:           booleanIsOperator,
		IsRootFilter:   true,
		Predicate:      "Product.state",
		SelectionStart: "uid",
		Value:          "ACTIVE",
	},
	"isinactive": {
		Type:           booleanIsOperator,
		IsRootFilter:   true,
		Predicate:      "Product.state",
		SelectionStart: "uid",
		Value:          "INACTIVE",
	},
	"isarchived": {
		Type:           booleanIsOperator,
		IsRootFilter:   true,
		Predicate:      "Product.state",
		SelectionStart: "uid",
		Value:          "ARCHIVED",
	},
	"isdeprecated": {
		Type:           booleanIsOperator,
		IsRootFilter:   true,
		Predicate:      "Product.state",
		SelectionStart: "uid",
		Value:          "DEPRECATED",
	},
	"ismissing": {
		Type:           booleanIsOperator,
		IsRootFilter:   true,
		Predicate:      "Product.state",
		SelectionStart: "uid",
		Value:          "MISSING",
	},

	//
	// License
	//
	"haslicense": {
		Type:           booleanHasOperator,
		Predicate:      "Component.license",
		SelectionStart: `Product.release`,
		SelectionEnd:   `{uid}`,
	},
	"hasadditionallicenses": {
		Type:           booleanHasOperator,
		Predicate:      "Component.additionalLicenses",
		SelectionStart: `Product.release`,
		SelectionEnd:   `{uid}`,
	},
	"license": { // alias for licenseid
		Type:           textTermExactOperator,
		Predicate:      "License.xid",
		SelectionStart: `Product.release {Component.license`,
		SelectionEnd:   `{uid}}`,
	},
	"licenseid": {
		Type:           textTermExactOperator,
		Predicate:      "License.xid",
		SelectionStart: `Product.release {Component.license`,
		SelectionEnd:   `{uid}}`,
	},
	"licensename": {
		Type:           textTermExactOperator,
		Predicate:      "License.name",
		SelectionStart: `Product.release {Component.license`,
		SelectionEnd:   `{uid}}`,
	},
	"islicensespdx": {
		Type:           booleanIsOperator,
		Predicate:      "License.isSpdx",
		SelectionStart: `Product.release {Component.license`,
		SelectionEnd:   `{uid}}`,
	},
	"islicensedeprecated": {
		Type:           booleanIsOperator,
		Predicate:      "License.isDeprecated",
		SelectionStart: `Product.release {Component.license`,
		SelectionEnd:   `{uid}}`,
	},
	"islicenseosiapproved": {
		Type:           booleanIsOperator,
		Predicate:      "License.isOsiApproved",
		SelectionStart: `Product.release {Component.license`,
		SelectionEnd:   `{uid}}`,
	},
	"islicensefsflibre": {
		Type:           booleanIsOperator,
		Predicate:      "License.isFsfLibre",
		SelectionStart: `Product.release {Component.license`,
		SelectionEnd:   `{uid}}`,
	},
	"islicenseblocked": {
		Type:           booleanIsOperator,
		Predicate:      "License.isBlocked",
		SelectionStart: `Product.release {Component.license`,
		SelectionEnd:   `{uid}}`,
	},
	"islicensestrong": {
		Type:           booleanIsOperator,
		Predicate:      "License.type",
		SelectionStart: `Product.release {Component.license`,
		SelectionEnd:   `{uid}}`,
		Value:          "STRONG",
	},
	"islicenseweak": {
		Type:           booleanIsOperator,
		Predicate:      "License.type",
		SelectionStart: `Product.release {Component.license`,
		SelectionEnd:   `{uid}}`,
		Value:          "WEAK",
	},
	"islicensepermissive": {
		Type:           booleanIsOperator,
		Predicate:      "License.type",
		SelectionStart: `Product.release {Component.license`,
		SelectionEnd:   `{uid}}`,
		Value:          "PERMISSIVE",
	},

	//
	// Licensor
	//
	"licensoruid": {
		Type:           uidOperator,
		SelectionStart: `Product.licensor`,
		SelectionEnd:   `{uid}`,
	},
	"licensor": { // alias for licensorfullname
		Type:           textTermExactOperator,
		Predicate:      "UserOrGroup.fullName",
		SelectionStart: `Product.licensor`,
		SelectionEnd:   `{uid}`,
	},
	"licensorname": {
		Type:           textTermExactOperator,
		Predicate:      "UserOrGroup.name",
		SelectionStart: `Product.licensor`,
		SelectionEnd:   `{uid}`,
	},
	"licensorfullname": {
		Type:           textTermExactOperator,
		Predicate:      "UserOrGroup.fullName",
		SelectionStart: `Product.licensor`,
		SelectionEnd:   `{uid}`,
	},
	"islicensoruser": {
		Type:           booleanIsOperator,
		Predicate:      "dgraph.type",
		SelectionStart: `Product.licensor`,
		SelectionEnd:   `{uid}`,
		Value:          "User",
	},
	"islicensorgroup": {
		Type:           booleanIsOperator,
		Predicate:      "dgraph.type",
		SelectionStart: `Product.licensor`,
		SelectionEnd:   `{uid}`,
		Value:          "Group",
	},

	//
	// Categorization
	//
	"hastags": {
		Type:           booleanHasOperator,
		IsRootFilter:   true,
		Predicate:      "Product.tags",
		SelectionStart: "uid",
	},
	"tag": {
		Type:           textFullContainsOperator,
		Predicate:      "Tag.name",
		SelectionStart: `Product.tags`,
		SelectionEnd:   `{uid}`,
	},
	"tagcount": {
		Type:           numberIntOperator,
		IsRootFilter:   true,
		Predicate:      "count(Product.tags)",
		SelectionStart: "uid",
	},
	"hascategory": {
		Type:           booleanHasOperator,
		IsRootFilter:   true,
		Predicate:      "Product.category",
		SelectionStart: "uid",
	},
	"category": {
		Type:           textFullContainsOperator,
		Predicate:      "Category.name",
		SelectionStart: `Product.category`,
		SelectionEnd:   `{uid}`,
	},
	"categoryname": {
		Type:           textFullContainsOperator,
		Predicate:      "Category.name",
		SelectionStart: `Product.category`,
		SelectionEnd:   `{uid}`,
	},
	"categoryfullname": {
		Type:           textFullContainsOperator,
		Predicate:      "Category.fullName",
		SelectionStart: `Product.category`,
		SelectionEnd:   `{uid}`,
	},

	//
	// Component
	//
	"createdat": {
		Type:           dateTimeOperator,
		Predicate:      "Component.createdAt",
		SelectionStart: `Product.release`,
		SelectionEnd:   `{uid}`,
	},

	"repository": { // alias for repositoryhost
		Type:           textTermExactOperator,
		Predicate:      "Host.name",
		SelectionStart: `Product.release {Component.repository {Repository.host`,
		SelectionEnd:   `{uid}}}`,
	},
	"repositoryhost": {
		Type:           textTermExactOperator,
		Predicate:      "Host.name",
		SelectionStart: `Product.release {Component.repository {Repository.host`,
		SelectionEnd:   `{uid}}}`,
	},
	"repositoryowner": {
		Type:           textTermExactOperator,
		Predicate:      "UserOrGroup.fullName",
		SelectionStart: `Product.release {Component.repository {Repository.owner`,
		SelectionEnd:   `{uid}}}`,
	},
	"repositoryname": {
		Type:           textTermExactOperator,
		Predicate:      "Repository.name",
		SelectionStart: `Product.release {Component.repository`,
		SelectionEnd:   `{uid}}`,
	},

	"datasource": { // alias for datasourcehost
		Type:           textTermExactOperator,
		Predicate:      "Host.name",
		SelectionStart: `Product.release {Component.dataSource {Repository.host`,
		SelectionEnd:   `{uid}}}`,
	},
	"datasourcehost": {
		Type:           textTermExactOperator,
		Predicate:      "Host.name",
		SelectionStart: `Product.release {Component.dataSource {Repository.host`,
		SelectionEnd:   `{uid}}}`,
	},
	"datasourceowner": {
		Type:           textTermExactOperator,
		Predicate:      "UserOrGroup.fullName",
		SelectionStart: `Product.release {Component.dataSource {Repository.owner`,
		SelectionEnd:   `{uid}}}`,
	},
	"datasourcename": {
		Type:           textTermExactOperator,
		Predicate:      "Repository.name",
		SelectionStart: `Product.release {Component.dataSource`,
		SelectionEnd:   `{uid}}`,
	},

	"host": { // alias for repositoryhost
		Type:           textTermExactOperator,
		Predicate:      "Host.name",
		SelectionStart: `Product.release {Component.repository {Repository.host`,
		SelectionEnd:   `{uid}}}`,
	},

	// TODO: should be usable like this: technologyreadinesslevel:>=3
	// "technologyreadinesslevel": {
	// 	Type:           textTermExactOperator,
	// 	Predicate:      "Component.technologyReadinessLevel",
	// 	SelectionStart: `Product.release`,
	// 	SelectionEnd:   `{uid}`,
	// },
	// "documentationreadinesslevel": {
	// 	Type:           textTermExactOperator,
	// 	Predicate:      "Component.documentationReadinessLevel",
	// 	SelectionStart: `Product.release`,
	// 	SelectionEnd:   `{uid}`,
	// },

	"hasattestation": {
		Type:           booleanHasOperator,
		IsRootFilter:   true,
		Predicate:      "Component.attestation",
		SelectionStart: "Product.release",
		SelectionEnd:   "{uid}",
	},
	"attestation": {
		Type:           textTermExactOperator,
		Predicate:      "Component.attestation",
		SelectionStart: "Product.release",
		SelectionEnd:   "{uid}",
	},
	"haspublication": {
		Type:           booleanHasOperator,
		IsRootFilter:   true,
		Predicate:      "Component.publication",
		SelectionStart: "Product.release",
		SelectionEnd:   "{uid}",
	},
	"publication": {
		Type:           textTermExactOperator,
		Predicate:      "Component.publication",
		SelectionStart: "Product.release",
		SelectionEnd:   "{uid}",
	},
	"hasissuetracker": {
		Type:           booleanHasOperator,
		IsRootFilter:   true,
		Predicate:      "Component.issues",
		SelectionStart: "Product.release",
		SelectionEnd:   "{uid}",
	},
	"issuetracker": {
		Type:           textTermExactOperator,
		Predicate:      "Component.issues",
		SelectionStart: "Product.release",
		SelectionEnd:   "{uid}",
	},
	"hascomplieswith": {
		Type:           booleanHasOperator,
		IsRootFilter:   true,
		Predicate:      "Component.compliesWith",
		SelectionStart: "Product.release",
		SelectionEnd:   "{uid}",
	},
	"complieswith": {
		Type:           textTermExactOperator,
		Predicate:      "TechnicalStandard.name",
		SelectionStart: "Product.release { Component.compliesWith",
		SelectionEnd:   "{uid}}",
	},
	"hascpcpatentclass": {
		Type:           booleanHasOperator,
		IsRootFilter:   true,
		Predicate:      "Component.cpcPatentClass",
		SelectionStart: "Product.release",
		SelectionEnd:   "{uid}",
	},
	"cpcpatentclass": {
		Type:           textTermExactOperator,
		Predicate:      "Component.cpcPatentClass",
		SelectionStart: "Product.release",
		SelectionEnd:   "{uid}",
	},
	"hastsdc": {
		Type:           booleanHasOperator,
		IsRootFilter:   true,
		Predicate:      "Component.tsdc",
		SelectionStart: "Product.release",
		SelectionEnd:   "{uid}",
	},
	"tsdc": {
		Type:           textTermExactOperator,
		Predicate:      "TechnologySpecificDocumentationCriteria.name",
		SelectionStart: "Product.release { Component.tsdc",
		SelectionEnd:   "{uid}}",
	},

	// TODO: sub components operators
	// TODO: software operators
	"hassoftware": {
		Type:           booleanHasOperator,
		Predicate:      "Component.software",
		SelectionStart: `Product.release`,
		SelectionEnd:   `{uid}`,
	},

	//
	// Files
	//
	"hasimage": {
		Type:           booleanHasOperator,
		Predicate:      "Component.image",
		SelectionStart: `Product.release`,
		SelectionEnd:   `{uid}`,
	},
	"hasreadme": {
		Type:           booleanHasOperator,
		Predicate:      "Component.readme",
		SelectionStart: `Product.release`,
		SelectionEnd:   `{uid}`,
	},
	"hascontributionguide": {
		Type:           booleanHasOperator,
		Predicate:      "Component.contributionGuide",
		SelectionStart: `Product.release`,
		SelectionEnd:   `{uid}`,
	},
	"hasbom": {
		Type:           booleanHasOperator,
		Predicate:      "Component.bom",
		SelectionStart: `Product.release`,
		SelectionEnd:   `{uid}`,
	},
	"hasmanufacturinginstructions": {
		Type:           booleanHasOperator,
		Predicate:      "Component.manufacturingInstructions",
		SelectionStart: `Product.release`,
		SelectionEnd:   `{uid}`,
	},
	"hasusermanual": {
		Type:           booleanHasOperator,
		Predicate:      "Component.userManual",
		SelectionStart: `Product.release`,
		SelectionEnd:   `{uid}`,
	},
	"hassource": {
		Type:           booleanHasOperator,
		Predicate:      "Component.source",
		SelectionStart: `Product.release`,
		SelectionEnd:   `{uid}`,
	},
	"hasexport": {
		Type:           booleanHasOperator,
		Predicate:      "Component.export",
		SelectionStart: `Product.release`,
		SelectionEnd:   `{uid}`,
	},
	"hasauxiliary": {
		Type:           booleanHasOperator,
		Predicate:      "Component.auxiliary",
		SelectionStart: `Product.release`,
		SelectionEnd:   `{uid}`,
	},

	// TODO: more fields
}

func (e *encoder) extractOperatorText(opr *parser.Operator) (text string, exact, match, not bool) {
	var val *parser.Text

	// check if the operator is a range, comparison or a plain value
	if opr.Range != nil {
		// ignore
		return
	} else if opr.Value != nil {
		val = opr.Value
	} else if opr.Comparison != nil {
		val = opr.Comparison.Value
		match = true
		switch opr.Comparison.Operator {
		case parser.CompOpEq:
			not = false
		case parser.CompOpNe:
			not = true
		default:
			// ignore
			return
		}
	}
	if val == nil {
		return
	}

	// extract the value as text
	if val.Words != nil {
		text = *val.Words
	} else {
		exact = true
		text = *val.Exact
	}

	return
}

func (e *encoder) generateNotFilter(sub []byte) []byte {
	if sub == nil || len(sub) == 0 {
		return nil
	}
	var b bytes.Buffer
	b.WriteString(`NOT (`)
	b.Write(sub)
	b.WriteString(`)`)
	return b.Bytes()
}

func (e *encoder) generateOrFilter(subs ...[]byte) []byte {
	if subs == nil || len(subs) == 0 {
		return nil
	}
	if len(subs) == 1 {
		return subs[0]
	}
	var b bytes.Buffer
	first := true
	for _, sub := range subs {
		if sub == nil || len(sub) == 0 {
			continue
		}
		if !first {
			b.WriteString(`) OR `)
		}
		b.WriteRune('(')
		b.Write(sub)
		first = false
	}
	b.WriteRune(')')
	return b.Bytes()
}

func (e *encoder) generateFilterExpression(filter []byte) []byte {
	if filter == nil || len(filter) == 0 {
		return nil
	}
	var b bytes.Buffer
	b.WriteString(`@filter(`)
	b.Write(filter)
	b.WriteString(`)`)
	return b.Bytes()
}

var uidPattern = regexp.MustCompile(`^0x[a-f0-9]{1,16}$`)

func (e *encoder) generateUIDFilter(uid string) []byte {
	if len(uid) == 0 {
		return nil
	}
	if !uidPattern.MatchString(uid) {
		return nil
	}
	var b bytes.Buffer
	b.WriteString(`uid(`)
	arg := e.CrateArg(uid, "string")
	b.WriteString(arg)
	b.WriteString(`)`)
	return b.Bytes()
}

// generateTextFilter ...
//
// match indicates whether the text need to be matched (enclosed with quotes "")
func (e *encoder) generateTextFilter(predicate string, oprTyp operatorType, text string, exact, match bool) []byte {
	var b bytes.Buffer

	if oprTyp == textExactOperator { // only supports exact matches
		match = true
	}

	// check first if has to be matched exactly
	if match {
		if oprTyp == textExactOperator || oprTyp == textTermExactOperator { // only supports exact-matches
			exact = true
		}

		// handle wildcards (*)
		if strings.Contains(text, "*") {
			parts := multGlobPattern.Split(text, -1)
			for i, part := range parts {
				// quote special characters so it is save to use with regex
				parts[i] = regexp.QuoteMeta(part)
			}
			var barg strings.Builder
			barg.WriteRune('/')
			if exact {
				barg.WriteRune('^')
			}
			for j, part := range parts {
				if part == "" {
					if j == len(parts)-1 {
						// if glob is last -> match until end of string
						barg.WriteString(`.*?`)
					}
					continue
				}
				if j != 0 {
					// max distance between matches of 5
					barg.WriteString(`(\s*\S*){0,5}`)
				}
				barg.WriteString(part)
			}
			if exact {
				barg.WriteRune('$')
			}
			barg.WriteString(`/i`)
			arg := e.CrateArg(barg.String(), "string")

			b.WriteString(`regexp(`)
			b.WriteString(predicate)
			b.WriteString(`, `)
			b.WriteString(arg)
			b.WriteString(`)`)

			return b.Bytes()
		}

		// does not contain wildcards
		switch oprTyp {
		case textFullContainsOperator, textTermContainsOperator, textExactOperator, textTermExactOperator:
			// match (containing phrase) using "regexp" filter (don't have an
			// optimized filter for case insensitive exact match yet)
			// TODO: implement a custom strings.contains filter in database that allows case insensitive matching and use that instead of regex
			b.WriteString(`regexp(`)
			b.WriteString(predicate)
			b.WriteString(`, `)

			text = regexp.QuoteMeta(text)
			var barg strings.Builder
			barg.Grow(len(text) + 10)
			barg.WriteRune('/')
			if exact {
				barg.WriteRune('^')
			}
			barg.WriteString(text)
			if exact {
				barg.WriteRune('$')
			}
			barg.WriteString(`/i`)
			arg := e.CrateArg(barg.String(), "string")

			b.WriteString(arg)
			b.WriteString(`)`)

			// case textExactOperator, textTermExactOperator:
			// 	// full match using "allof" filter on "exacti" index
			// 	b.WriteString(`allof(`)
			// 	b.WriteString(predicate)

			// XXX: it seems I cannot add the exacti index via GraphQL schema...
			// 	b.WriteString(`, exacti, `)

			// 	arg := e.CrateArg(text, "string")
			// 	b.WriteString(arg)
			// 	b.WriteString(`)`)
		}

		return b.Bytes()
	}

	// not of type "match" -> terms or full text search
	fltName := "alloftext"
	if oprTyp == textTermExactOperator || oprTyp == textTermContainsOperator {
		fltName = "allofterms"
	}

	// handle wildcards (*)
	if strings.Contains(text, "*") {
		words := strings.Split(text, " ")
		nonWildcardWords := []string{}
		for _, word := range words {
			if word == "" {
				continue
			}

			if strings.Contains(word, "*") {
				parts := multGlobPattern.Split(word, -1)
				for i, part := range parts {
					parts[i] = regexp.QuoteMeta(part)
				}
				var barg strings.Builder
				barg.WriteString(`/`)
				for j, part := range parts {
					if part == "" {
						if j == len(parts)-1 {
							// if glob is last -> match until next word
							barg.WriteString(`(\s*\S*)?`)
						}
						continue
					}
					if j != 0 {
						// max distance between matches of 5
						barg.WriteString(`(\s*\S*){0,5}`)
					}

					barg.WriteString(part)
				}
				barg.WriteString(`/i`)
				arg := e.CrateArg(barg.String(), "string")

				if b.Len() > 0 {
					b.WriteString(` AND `)
				}
				b.WriteString(`regexp(`)
				b.WriteString(predicate)
				b.WriteString(`, `)
				b.WriteString(arg)
				b.WriteString(`)`)

			} else {
				nonWildcardWords = append(nonWildcardWords, word)
			}
		}
		// use full text or term search for the rest
		if len(nonWildcardWords) > 0 {
			s := strings.Join(nonWildcardWords, " ")
			arg := e.CrateArg(s, "string")
			b.WriteString(` AND `)
			b.WriteString(fltName)
			b.WriteRune('(')
			b.WriteString(predicate)
			b.WriteString(`, `)
			b.WriteString(arg)
			b.WriteRune(')')
		}

		return b.Bytes()
	}

	// without wildcards
	arg := e.CrateArg(text, "string")
	b.WriteString(fltName)
	b.WriteRune('(')
	b.WriteString(predicate)
	b.WriteString(`, `)
	b.WriteString(arg)
	b.WriteRune(')')
	return b.Bytes()
}

func (e *encoder) appendBooleanIsFilter(predicate string, value string, not bool) {
	e.buf.WriteString(`@filter(`)
	if not {
		e.buf.WriteString(`NOT `)
	}
	e.buf.WriteString(`eq(`)
	e.buf.WriteString(predicate)
	if value != "" {
		e.buf.WriteString(`, "`)
		// value is defined by developer, so we can use it directly
		e.buf.WriteString(value)
		e.buf.WriteString(`"))`)
	} else {
		e.buf.WriteString(`, true))`)
	}
}

func (e *encoder) appendBooleanHasFilter(predicate string, not bool) {
	e.buf.WriteString(`@filter(`)
	if not {
		e.buf.WriteString(`NOT `)
	}
	e.buf.WriteString(`has(`)
	e.buf.WriteString(predicate)
	e.buf.WriteString(`))`)
}

func extractNumberValue(rawVal *parser.Text) (number float64, ok bool) {
	if rawVal == nil {
		return
	}
	// extract the value as text
	var strVal *string
	if rawVal.Words != nil {
		strVal = rawVal.Words
	} else {
		strVal = rawVal.Exact
	}
	return parseNumberValue(strVal)
}

func parseNumberValue(strVal *string) (number float64, ok bool) {
	if strVal == nil {
		return
	}
	*strVal = strings.TrimSpace(*strVal)
	floatVal, err := strconv.ParseFloat(*strVal, 64)
	if err != nil {
		// invalid value -> ignore
		return
	}

	return floatVal, true
}

func (e *encoder) appendNumberFilter(predicate string, opr *parser.Operator, isInt bool) {
	var (
		txtVal *parser.Text
		cmpOpr parser.CompOperator
	)

	// check if the operator is a range, comparison or a plain value
	if opr.Range != nil {
		if opr.Range.OpenStart && opr.Range.OpenEnd {
			// both open -> ignore
			return
		}
		var b bytes.Buffer
		b.WriteString(`@filter(`)
		if !opr.Range.OpenStart {
			number, ok := parseNumberValue(opr.Range.Start)
			if !ok {
				return
			}
			b.WriteString(`ge(`)
			b.WriteString(predicate)
			b.WriteString(`, `)
			if isInt {
				b.WriteString(strconv.FormatInt(int64(number), 10))
			} else {
				b.WriteString(strconv.FormatFloat(number, 'f', -1, 64))
			}
			b.WriteString(`)`)
		}

		if !opr.Range.OpenEnd {
			if !opr.Range.OpenStart {
				b.WriteString(` AND `)
			}
			number, ok := parseNumberValue(opr.Range.End)
			if !ok {
				return
			}
			b.WriteString(`le(`)
			b.WriteString(predicate)
			b.WriteString(`, `)
			if isInt {
				b.WriteString(strconv.FormatInt(int64(number), 10))
			} else {
				b.WriteString(strconv.FormatFloat(number, 'f', -1, 64))
			}
			b.WriteString(`)`)
		}
		b.WriteRune(')')

		// append the value
		e.buf.Write(b.Bytes())
		return

	} else if opr.Value != nil {
		// treat as equality comparison
		txtVal = opr.Value
		cmpOpr = parser.CompOpEq

	} else if opr.Comparison != nil {
		txtVal = opr.Comparison.Value
		cmpOpr = opr.Comparison.Operator
	}
	number, ok := extractNumberValue(txtVal)
	if !ok {
		return
	}

	e.buf.WriteString(`@filter(`)
	switch cmpOpr {
	case parser.CompOpEq:
		e.buf.WriteString(`eq(`)
	case parser.CompOpNe:
		e.buf.WriteString(`ne(`)
	case parser.CompOpLt:
		e.buf.WriteString(`lt(`)
	case parser.CompOpLe:
		e.buf.WriteString(`le(`)
	case parser.CompOpGt:
		e.buf.WriteString(`gt(`)
	case parser.CompOpGe:
		e.buf.WriteString(`ge(`)
	default:
		panic("unknown comparison operator")
	}
	e.buf.WriteString(predicate)
	e.buf.WriteString(`, `)
	if isInt {
		e.buf.WriteString(strconv.FormatInt(int64(number), 10))
	} else {
		e.buf.WriteString(strconv.FormatFloat(number, 'f', -1, 64))
	}
	e.buf.WriteString(`))`)
}

func extractDateTimeValue(rawVal *parser.Text) (dt time.Time, isDuration, ok bool) {
	if rawVal == nil {
		return
	}
	// extract the value as text
	var strVal *string
	if rawVal.Words != nil {
		strVal = rawVal.Words
	} else {
		strVal = rawVal.Exact
	}
	return parseDateTimeValue(strVal)
}

// parseDateTimeValue parses the given string value and returns the time.Time.
func parseDateTimeValue(rawVal *string) (dt time.Time, isDuration, ok bool) {
	strVal := strings.TrimSpace(*rawVal)

	// try parsing as time duration
	duration, ok := parseDuration(strVal)
	if ok {
		return time.Now().Add(-1 * duration), true, true
	}

	// try parsing as date time
	c := carbon.Parse(strVal)
	if c.Error == nil {
		return c.Carbon2Time(), false, true
	}

	// invalid value -> ignore
	return
}

var timeDurationPattern = regexp.MustCompile(`^(?P<years>\d+[yY])?(?P<months>\d+[mM])?(?P<weeks>\d+[wW])?(?P<days>\d+[dD])?$`)

func parseDuration(str string) (time.Duration, bool) {
	str = strings.ReplaceAll(str, " ", "")
	matches := timeDurationPattern.FindStringSubmatch(str)
	if matches == nil {
		return 0, false
	}

	years := parseDurationMatch(matches[1])
	months := parseDurationMatch(matches[2])
	weeks := parseDurationMatch(matches[3])
	days := parseDurationMatch(matches[4])

	hour := int64(time.Hour)
	return time.Duration(years*24*365*hour + months*30*24*hour + weeks*7*24*hour + days*24*hour), true
}

func parseDurationMatch(value string) int64 {
	if len(value) == 0 {
		return 0
	}
	parsed, err := strconv.ParseInt(value[:len(value)-1], 10, 64)
	if err != nil {
		return 0
	}
	return parsed
}

func (e *encoder) appendDateTimeFilter(predicate string, opr *parser.Operator) {
	var (
		txtVal *parser.Text
		cmpOpr parser.CompOperator
	)

	// check if the operator is a range, comparison or a plain value
	if opr.Range != nil {
		if opr.Range.OpenStart && opr.Range.OpenEnd {
			// both open -> ignore
			return
		}
		var b bytes.Buffer
		b.WriteString(`@filter(`)
		if !opr.Range.OpenStart {
			dt, _, ok := parseDateTimeValue(opr.Range.Start)
			if !ok {
				return
			}
			b.WriteString(`ge(`)
			b.WriteString(predicate)
			b.WriteString(`, "`)
			b.WriteString(dt.Format(time.RFC3339))
			b.WriteString(`")`)
		}

		if !opr.Range.OpenEnd {
			if !opr.Range.OpenStart {
				b.WriteString(` AND `)
			}
			dt, _, ok := parseDateTimeValue(opr.Range.End)
			if !ok {
				return
			}
			b.WriteString(`le(`)
			b.WriteString(predicate)
			b.WriteString(`, "`)
			b.WriteString(dt.Format(time.RFC3339))
			b.WriteString(`")`)
		}
		b.WriteRune(')')

		// append the value
		e.buf.Write(b.Bytes())
		return

	} else if opr.Value != nil {
		// treat as equality comparison
		txtVal = opr.Value
		cmpOpr = parser.CompOpEq

	} else if opr.Comparison != nil {
		txtVal = opr.Comparison.Value
		cmpOpr = opr.Comparison.Operator
	}
	dt, isDuration, ok := extractDateTimeValue(txtVal)
	if !ok {
		return
	}

	e.buf.WriteString(`@filter(`)
	switch cmpOpr {
	case parser.CompOpEq:
		e.buf.WriteString(`eq(`)
	case parser.CompOpNe:
		e.buf.WriteString(`ne(`)
	case parser.CompOpLt:
		if isDuration {
			e.buf.WriteString(`gt(`)
		} else {
			e.buf.WriteString(`lt(`)
		}
	case parser.CompOpLe:
		if isDuration {
			e.buf.WriteString(`ge(`)
		} else {
			e.buf.WriteString(`le(`)
		}
	case parser.CompOpGt:
		if isDuration {
			e.buf.WriteString(`lt(`)
		} else {
			e.buf.WriteString(`gt(`)
		}
	case parser.CompOpGe:
		if isDuration {
			e.buf.WriteString(`le(`)
		} else {
			e.buf.WriteString(`ge(`)
		}
	default:
		panic("unknown comparison operator")
	}
	e.buf.WriteString(predicate)
	e.buf.WriteString(`, "`)
	e.buf.WriteString(dt.Format(time.RFC3339))
	e.buf.WriteString(`"))`)
}

func (e *encoder) addVariableWithFilter(rootFilter, selection, parVar string) string {
	buf := e.buf

	if selection == "" {
		selection = "{uid}"
	} else {
		selection = fmt.Sprintf(`@cascade {%s}`, selection)
	}

	// use type() filter
	if parVar == "" {
		curVar := e.createVar()
		s := `%s as var(func:type(Product)) %s %s`
		s = fmt.Sprintf(s, curVar, rootFilter, selection)
		buf.WriteString(s)
		buf.WriteString("\n")
		return curVar
	}

	// use uid() filter
	curVar := e.createVar()
	s := `%s as var(func:uid(%s)) %s %s`
	s = fmt.Sprintf(s, curVar, parVar, rootFilter, selection)
	buf.WriteString(s)
	buf.WriteString("\n")
	return curVar
}

func (e *encoder) appendSelection(start, end string, filter func()) {
	e.buf.WriteString(start)
	if filter != nil {
		e.buf.WriteRune(' ')
		filter()
	}
	e.buf.WriteString(end)
}

func (e *encoder) appendVariable(rootFilter, selection func(), parVar string) string {
	curVar := e.createVar()
	e.buf.WriteString(curVar)
	if parVar == "" { // use type() filter
		e.buf.WriteString(" as var(func:type(Product)) ")
	} else { // use uid() filter
		e.buf.WriteString(" as var(func:uid(")
		e.buf.WriteString(parVar)
		e.buf.WriteString(")) ")
	}
	if rootFilter != nil {
		rootFilter()
	}
	if selection != nil {
		e.buf.WriteString("@cascade {")
		selection()
		e.buf.WriteRune('}')
	} else {
		e.buf.WriteString("{uid}")
	}
	e.buf.WriteRune('\n')
	return curVar
}

func (e *encoder) appendVariable2(filterAtRoot bool, encodeFilterFn func(), selectionStart, selectionEnd, parVar string) string {
	curVar := e.createVar()
	e.buf.WriteString(curVar)
	if parVar == "" { // use type() filter
		e.buf.WriteString(" as var(func:type(Product)) ")
	} else { // use uid() filter
		e.buf.WriteString(" as var(func:uid(")
		e.buf.WriteString(parVar)
		e.buf.WriteString(")) ")
	}
	if filterAtRoot {
		encodeFilterFn()
		e.buf.WriteRune('{')
		e.buf.WriteString(selectionStart)
		e.buf.WriteString(selectionEnd)
		e.buf.WriteRune('}')
	} else {
		e.buf.WriteString("@cascade {")
		e.buf.WriteString(selectionStart)
		e.buf.WriteRune(' ')
		encodeFilterFn()
		e.buf.WriteString(selectionEnd)
		e.buf.WriteRune('}')
	}
	e.buf.WriteRune('\n')
	return curVar
}

func (e *encoder) appendUnionVariable(parVar ...string) string {
	curVar := e.createVar()
	e.buf.WriteString(curVar)
	e.buf.WriteString(" as var(func:uid(")
	for i, v := range parVar {
		e.buf.WriteString(v)
		if i < len(parVar)-1 {
			e.buf.WriteRune(',')
		}
	}
	e.buf.WriteString(")) {uid}\n")
	return curVar
}

func (e *encoder) appendNegationVariable(parVar, notVar string) string {
	curVar := e.createVar()
	e.buf.WriteString(curVar)
	e.buf.WriteString(" as var(func:")
	if parVar == "" {
		e.buf.WriteString("type(Product)")
	} else {
		e.buf.WriteString("uid(")
		e.buf.WriteString(parVar)
		e.buf.WriteString(")")
	}
	e.buf.WriteString(") @filter(NOT uid(")
	e.buf.WriteString(notVar)
	e.buf.WriteString(")) {uid}\n")
	return curVar
}

func (e *encoder) appendOrderByVariable(orderBy searchmodels.OrderBy, parVar string) string {
	e.buf.WriteString("var(func:uid(")
	e.buf.WriteString(parVar)
	e.buf.WriteString(")) {")

	switch orderBy.Field {
	case searchmodels.OrderByName:
		e.buf.WriteString(`order as Product.name`)
	case searchmodels.OrderByDocumentationLanguage:
		e.buf.WriteString(`order as Product.documentationLanguage`)
	case searchmodels.OrderByState:
		e.buf.WriteString(`order as Product.state`)
	case searchmodels.OrderByForkCount:
		e.buf.WriteString(`order as Product.forkCount`)
	case searchmodels.OrderByStarCount:
		e.buf.WriteString(`order as Product.starCount`)
	case searchmodels.OrderByVersion:
		e.buf.WriteString(`order as Product.version`)
	case searchmodels.OrderByWebsite:
		e.buf.WriteString(`order as Product.website`)
	case searchmodels.OrderByCreatedAt:
		e.buf.WriteString(`Product.release { O1 as Component.createdAt } order as min(val(O1))`)
	case searchmodels.OrderByDiscoveredAt:
		e.buf.WriteString(`Product.release { O1 as CrawlerMeta.discoveredAt } order as min(val(O1))`)
	case searchmodels.OrderByLastIndexedAt:
		e.buf.WriteString(`Product.release { O1 as CrawlerMeta.lastIndexedAt } order as min(val(O1))`)
	case searchmodels.OrderByLastUpdatedAt:
		e.buf.WriteString(`order as Product.lastUpdatedAt`)
	case searchmodels.OrderByHasAdditionalLicenses:
		// TODO:
		e.buf.WriteString(`Product.release { Component.license { O2 as License.hasAdditionalLicenses } O1 as min(val(O2)) } order as min(val(O1))`)
	case searchmodels.OrderByLicenseID:
		e.buf.WriteString(`Product.release { Component.license { O2 as License.xid } O1 as min(val(O2)) } order as min(val(O1))`)
	case searchmodels.OrderByLicenseName:
		e.buf.WriteString(`Product.release { Component.license { O2 as License.name } O1 as min(val(O2)) } order as min(val(O1))`)
	case searchmodels.OrderByIsLicenseSpdx:
		e.buf.WriteString(`Product.release { Component.license { O2 as License.isSpdx } O1 as min(val(O2)) } order as min(val(O1))`)
	case searchmodels.OrderByIsLicenseDeprecated:
		e.buf.WriteString(`Product.release { Component.license { O2 as License.isDeprecated } O1 as min(val(O2)) } order as min(val(O1))`)
	case searchmodels.OrderByIsLicenseOsiApproved:
		e.buf.WriteString(`Product.release { Component.license { O2 as License.isOsiApproved } O1 as min(val(O2)) } order as min(val(O1))`)
	case searchmodels.OrderByIsLicenseFsfLibre:
		e.buf.WriteString(`Product.release { Component.license { O2 as License.isFsfLibre } O1 as min(val(O2)) } order as min(val(O1))`)
	case searchmodels.OrderByIsLicenseBlocked:
		e.buf.WriteString(`Product.release { Component.license { O2 as License.isBlocked } O1 as min(val(O2)) } order as min(val(O1))`)
	case searchmodels.OrderByLicenseType:
		e.buf.WriteString(`Product.release { Component.license { O2 as License.type } O1 as min(val(O2)) } order as min(val(O1))`)
	case searchmodels.OrderByIsLicenseStrong:
		// TODO:
		e.buf.WriteString(`Product.release { Component.license { O2 as eq(License.type, "STRONG") } O1 as min(val(O2)) } order as min(val(O1))`)
	case searchmodels.OrderByIsLicenseWeak:
		// TODO:
		e.buf.WriteString(`Product.release { Component.license { O2 as eq(License.type, "WEAK") } O1 as min(val(O2)) } order as min(val(O1))`)
	case searchmodels.OrderByIsLicensePermissive:
		// TODO:
		e.buf.WriteString(`Product.release { Component.license { O2 as eq(License.type, "PERMISSIVE") } O1 as min(val(O2)) } order as min(val(O1))`)
	case searchmodels.OrderByLicensorFullName:
		e.buf.WriteString(`Product.release { Component.licensor { O2 as UserOrGroup.fullName } O1 as min(val(O2)) } order as min(val(O1))`)
	case searchmodels.OrderByLicensorName:
		e.buf.WriteString(`Product.release { Component.licensor { O2 as UserOrGroup.name } O1 as min(val(O2)) } order as min(val(O1))`)
	case searchmodels.OrderByRepositoryHost:
		e.buf.WriteString(`Product.release { Component.repository { Repository.host { O3 as Host.name } O2 as min(val(O3)) } O1 as min(val(O2)) } order as min(val(O1))`)
	case searchmodels.OrderByRepositoryOwner:
		e.buf.WriteString(`Product.release { Component.repository { Repository.owner { O3 as UserOrGroup.fullName } O2 as min(val(O3)) } } O1 as min(val(O2)) } order as min(val(O1))`)
	case searchmodels.OrderByRepositoryName:
		e.buf.WriteString(`Product.release { Component.repository { O2 as Repository.name } O1 as min(val(O2)) } order as min(val(O1))`)
	case searchmodels.OrderByDatasourceHost:
		e.buf.WriteString(`Product.release { CrawlerMeta.dataSource { Repository.host { O3 as Host.name } O2 as min(val(O3)) } O1 as min(val(O2)) } order as min(val(O1))`)
	case searchmodels.OrderByDatasourceOwner:
		e.buf.WriteString(`Product.release { CrawlerMeta.dataSource { Repository.owner { O3 as UserOrGroup.fullName } O2 as min(val(O3)) } } O1 as min(val(O2)) } order as min(val(O1))`)
	case searchmodels.OrderByDatasourceName:
		e.buf.WriteString(`Product.release { CrawlerMeta.dataSource { O2 as Repository.name } O1 as min(val(O2)) } order as min(val(O1))`)
	case searchmodels.OrderByHasAttestation:
		// TODO:
		e.buf.WriteString(`Product.release { O1 as Component.attestation } order as min(val(O1))`)
	case searchmodels.OrderByAttestation:
		e.buf.WriteString(`Product.release { O1 as Component.attestation } order as min(val(O1))`)
	case searchmodels.OrderByHasPublication:
		// TODO:
		e.buf.WriteString(`Product.release { O1 as Component.publication } order as min(val(O1))`)
	case searchmodels.OrderByPublication:
		e.buf.WriteString(`Product.release { O1 as Component.publication } order as min(val(O1))`)
	case searchmodels.OrderByHasIssueTracker:
		// TODO:
		e.buf.WriteString(`Product.release { O1 as Component.issues } order as min(val(O1))`)
	case searchmodels.OrderByIssueTracker:
		e.buf.WriteString(`Product.release { O1 as Component.issues } order as min(val(O1))`)
	case searchmodels.OrderByHasComplieswith:
		// TODO:
		e.buf.WriteString(`Product.release { O1 as Component.complieswith } order as min(val(O1))`)
	case searchmodels.OrderByComplieswith:
		e.buf.WriteString(`Product.release { O1 as Component.complieswith } order as min(val(O1))`)
	case searchmodels.OrderByHasCpcpatentclass:
		// TODO:
		e.buf.WriteString(`Product.release { O1 as Component.cpcPatentClass } order as min(val(O1))`)
	case searchmodels.OrderByCpcPatentClass:
		e.buf.WriteString(`Product.release { O1 as Component.cpcPatentClass } order as min(val(O1))`)
	case searchmodels.OrderByHasTsdc:
		// TODO:
		e.buf.WriteString(`Product.release { O1 as Component.tsdc } order as min(val(O1))`)
	case searchmodels.OrderByTsdc:
		e.buf.WriteString(`Product.release { O1 as Component.tsdc } order as min(val(O1))`)
	case searchmodels.OrderByHasImage:
		// TODO:
		e.buf.WriteString(`Product.release { O1 as Component.image } order as min(val(O1))`)
	case searchmodels.OrderByImage:
		e.buf.WriteString(`Product.release { O1 as Component.image } order as min(val(O1))`)
	case searchmodels.OrderByHasReadme:
		// TODO:
		e.buf.WriteString(`Product.release { O1 as Component.readme } order as min(val(O1))`)
	case searchmodels.OrderByReadme:
		e.buf.WriteString(`Product.release { O1 as Component.readme } order as min(val(O1))`)
	case searchmodels.OrderByHasContributionGuide:
		// TODO:
		e.buf.WriteString(`Product.release { O1 as Component.contributionGuide } order as min(val(O1))`)
	case searchmodels.OrderByContributionGuide:
		e.buf.WriteString(`Product.release { O1 as Component.contributionGuide } order as min(val(O1))`)
	case searchmodels.OrderByHasBom:
		// TODO:
		e.buf.WriteString(`Product.release { O1 as Component.bom } order as min(val(O1))`)
	case searchmodels.OrderByBom:
		e.buf.WriteString(`Product.release { O1 as Component.bom } order as min(val(O1))`)
	case searchmodels.OrderByHasManufacturingInstructions:
		// TODO:
		e.buf.WriteString(`Product.release { O1 as Component.manufacturingInstructions } order as min(val(O1))`)
	case searchmodels.OrderByManufacturingInstructions:
		e.buf.WriteString(`Product.release { O1 as Component.manufacturingInstructions } order as min(val(O1))`)
	case searchmodels.OrderByHasUserManual:
		// TODO:
		e.buf.WriteString(`Product.release { O1 as Component.userManual } order as min(val(O1))`)
	case searchmodels.OrderByUserManual:
		e.buf.WriteString(`Product.release { O1 as Component.userManual } order as min(val(O1))`)
	case searchmodels.OrderByHasSource:
		// TODO:
		e.buf.WriteString(`Product.release { O1 as Component.source } order as min(val(O1))`)
	case searchmodels.OrderBySource:
		e.buf.WriteString(`Product.release { O1 as Component.source } order as min(val(O1))`)
	case searchmodels.OrderByHasExport:
		// TODO:
		e.buf.WriteString(`Product.release { O1 as Component.export } order as min(val(O1))`)
	case searchmodels.OrderByExport:
		e.buf.WriteString(`Product.release { O1 as Component.export } order as min(val(O1))`)
	case searchmodels.OrderByHasAuxiliary:
		// TODO:
		e.buf.WriteString(`Product.release { O1 as Component.auxiliary } order as min(val(O1))`)
	case searchmodels.OrderByAuxiliary:
		e.buf.WriteString(`Product.release { O1 as Component.auxiliary } order as min(val(O1))`)
	default:
		panic("unsupported orderBy field")
	}

	e.buf.WriteString("}\n")
	return parVar
}

func (e *encoder) appendSelectQuery(query, order, parVar string) {
	query = fmt.Sprintf(query, parVar, order)
	e.buf.WriteString(query)
	e.buf.WriteString("\n")
}
