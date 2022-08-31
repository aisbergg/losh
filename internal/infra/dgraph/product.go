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
	"losh/internal/infra/dgraph/dgclient"
	searchmodels "losh/web/core/search/models"
	"losh/web/core/search/parser"

	"github.com/aisbergg/go-errors/pkg/errors"
	"github.com/golang-module/carbon/v2"
)

// SearchProducts returns a list of `Product` objects matching the filter criteria.
func (dr *DgraphRepository) SearchProducts(ctx context.Context, filter *dgclient.ProductFilter, order *dgclient.ProductOrder, first *int64, offset *int64) ([]*productmodels.Product, int64, error) {
	dr.log.Debugw("search Products")
	rsp, err := dr.client.SearchProducts(ctx, filter, order, first, offset)
	if err != nil {
		return nil, 0, WrapRepoError(err, errGetProductStr)
	}
	ret := make([]*productmodels.Product, 0, len(rsp.QueryProduct))
	if err = dr.copier.CopyTo(rsp.QueryProduct, &ret); err != nil {
		panic(err)
	}
	return ret, *rsp.AggregateProduct.Count, nil
}

// func (dr *DgraphRepository) SearchProductsDQL(ctx context.Context, filter *dgclient.ProductFilter, order *dgclient.ProductOrder, first *int64, offset *int64) ([]*models.Product, uint64, error) {
// 	dr.log.Debugw("search Products")

// 	// bind results to map - will be copied into data model at the end
// 	var rspData []map[string]interface{}

// 	// search products
// 	searchQuery := dqlx.Variable(dqlx.TypeFn("Product")).As("searchResults")
// 	// Filter(dqlx.And()).

// 	// order by
// 	orderBy := OrderBy{"license", true}
// 	switch orderBy.Field {
// 	case "license":
// 		// searchQuery = searchQuery.
// 		// 	Select(dqlx.As("order", dqlx.Min("O1"))).
// 		// 	Edge("Product.release", dqlx.As("O1", dqlx.Min("O2"))). // variable needs to be propagated through levels
// 		// 	Edge("Product.release->Component.license", dqlx.Select("O2 as License.xid"))
// 		// searchQuery = searchQuery.
// 		searchQuery = searchQuery.Select(dqlx.As("order", dqlx.Min("O1")))
// 		searchQuery = searchQuery.Edge("Product.release", dqlx.As("O1", dqlx.Min("O2"))) // variable needs to be propagated through levels
// 		searchQuery = searchQuery.Edge("Product.release->Component.license", dqlx.Select("O2 as License.xid"))
// 	default:
// 		// order by name
// 	}

// 	// collect results and select fields
// 	paginatedQuery := dr.dqlxClient.Query(dqlx.UIDFn(dqlx.P("searchResults"))).
// 		Variable(searchQuery).
// 		Raw(`
// 			uid
// 			Product.xid
// 			Product.name
// 			Product.website
// 			Product.renamedTo {uid}
// 			Product.renamedFrom {uid}
// 			Product.forkOf {
// 				uid
// 				Product.release {
// 					Component.repository {
// 						Repository.url
// 					}
// 				}
// 			}
// 			Product.forks {uid}
// 			Product.releases {uid}
// 		`).
// 		EdgeFromQuery(TagDQLFragment.Name("Product.tags")).
// 		EdgeFromQuery(CategoryDQLFragment.Name("Product.category")).
// 		EdgeFn("Product.release", func(builder dqlx.QueryBuilder) dqlx.QueryBuilder {
// 			return builder.Raw(`
// 					Component.id
// 					Component.xid
// 					Component.name
// 					Component.description
// 					Component.version
// 					Component.createdAt
// 					Component.releases {uid}
// 					Component.isLatest
// 					Component.documentationLanguage
// 					Component.technologyReadinessLevel
// 					Component.documentationReadinessLevel
// 					Component.attestation
// 					Component.publication
// 					Component.compliesWith {name}
// 					Component.cpcPatentClass
// 					Component.tsdc {uid}
// 					Component.components {uid}
// 					Component.software {uid}
// 					Component.product {uid}
// 					Component.usedIn {uid}
// 					Component.organization {uid}
// 					Component.mass
// 					Component.material {uid}
// 					Component.manufacturingProcess {uid}
// 					Component.productionMetadata {uid}
// 				`).
// 				EdgeFromQuery(RepositoryDQLFragment.Name("Component.repository")).
// 				EdgeFromQuery(LicenseBasicDQLFragment.Name("Component.license")).
// 				EdgeFromQuery(LicenseBasicDQLFragment.Name("Component.additionalLicenses")).
// 				EdgeFromQuery(UserOrGroupFullDQLFragment.Name("Component.licensor")).
// 				EdgeFromQuery(FileDQLFragment.Name("Component.image")).
// 				EdgeFromQuery(FileDQLFragment.Name("Component.readme")).
// 				EdgeFromQuery(FileDQLFragment.Name("Component.contributionGuide")).
// 				EdgeFromQuery(FileDQLFragment.Name("Component.bom")).
// 				EdgeFromQuery(FileDQLFragment.Name("Component.manufacturingInstructions")).
// 				EdgeFromQuery(FileDQLFragment.Name("Component.userManual")).
// 				EdgeFromQuery(FileDQLFragment.Name("Component.source")).
// 				EdgeFromQuery(FileDQLFragment.Name("Component.export")).
// 				EdgeFromQuery(FileDQLFragment.Name("Component.auxiliary")).
// 				EdgeFromQuery(OuterDimensionsDQLFragment.Name("Component.outerDimensions"))
// 		}).
// 		UnmarshalInto(&rspData)

// 	// order results
// 	paginatedQuery = paginatedQuery.OrderAsc(dqlx.Val("order"))

// 	// paginate results
// 	if first != nil {
// 		paginatedQuery = paginatedQuery.Paginate(dqlx.Cursor{
// 			// TODO: revert
// 			// First: int(*first),
// 			First:  1,
// 			Offset: int(*offset),
// 			After:  "",
// 		})
// 	}

// 	s, m, e := paginatedQuery.ToDQL()
// 	fmt.Println("query:", s)
// 	fmt.Println("m:", m)
// 	fmt.Println("e:", e)

// 	rsp, err := paginatedQuery.Execute(ctx, dqlx.WithReadOnly(true))
// 	if err != nil {
// 		return nil, 0, err
// 	}

// 	// fmt.Println("here")
// 	// fmt.Println(string(rsp.Raw.Json))
// 	// fmt.Println(rsp.Raw.Metrics.NumUids["uid"])

// 	ret := make([]*models.Product, 0, len(rspData))
// 	if err = dr.dqlCopier.CopyTo(rspData, &ret); err != nil {
// 		panic(err)
// 	}

// 	// fmt.Println("ret:", rspData[0])
// 	// fmt.Println("ret:", ret[0].Name)
// 	// fmt.Println("metrics:", rsp.Raw.Metrics)

// 	return ret, rsp.Raw.Metrics.NumUids["Product.xid"], nil
// }

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
		Tag.xid
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

func (dr *DgraphRepository) SearchProductsDQL(ctx context.Context, query *parser.Query, order searchmodels.OrderBy, pagination searchmodels.Pagination) ([]*productmodels.Product, uint64, error) {
	dr.log.Debugw("search Products")

	q, v := createDQLQuery(query, order, pagination)
	fmt.Println("q: ", q)
	fmt.Println("v: ", v)

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

// func (e *encoder) encodeAndConditions(andCnds []*parser.AndCondition) (curVar string) {
// 	for _, andCnd := range andCnds {
// 		pf, lf := encodeAndCondition(andCnd)
// 		if pf != nil {
// 			if prdAndFlts == nil {
// 				prdAndFlts = []*dgclient.ProductFilter{}
// 			}
// 			prdAndFlts = append(prdAndFlts, pf)
// 		}
// 		_ = lf
// 	}

// 	if andCnd.Not != nil {
// 		pf, lf := encodeAndCondition(andCnd.Not)
// 		if pf != nil {
// 			pf = &dgclient.ProductFilter{
// 				Not: pf,
// 			}
// 		}
// 		if lf != nil {
// 			lf = &dgclient.LicenseFilter{
// 				Not: lf,
// 			}
// 		}
// 		return pf, lf

// 	} else {
// 		return encodeExpression(andCnd.Operand)
// 	}
// }

func (e *encoder) encodeAndCondition(andCnd *parser.AndCondition, parVar string) (curVar string) {
	if andCnd.Not != nil {
		notVar := e.encodeAndCondition(andCnd.Not, parVar)
		return e.appendNegationVariable(parVar, notVar)
	}

	return e.encodeExpression(andCnd.Operand, parVar)
}

// encodeExpression encodes an expression into DQL variable and returns the
// variable name.
// func (e *encoder) encodeExpression(expr *parser.Expression, parVar string) (curVar string) {
// 	if expr.Text != nil {
// 		if expr.Text.Words != nil {
// 			arg := e.addArg(*expr.Text.Words, "string")
// 			filter := fmt.Sprintf(`@filter(
// 				alloftext(Product.name, %s)
// 				OR alloftext(Product.description, %s)
// 			)`, arg, arg)
// 			curVar = e.addVariableWithFilter(filter, "", parVar)
// 		}
// 	}

// 	return
// }

var multGlobPattern = regexp.MustCompile(`\*+`)

// encodeExpression encodes an expression into DQL variable and returns the
// variable name.
func (e *encoder) encodeExpression(expr *parser.Expression, parVar string) (curVar string) {
	if expr.Text != nil {
		if expr.Text.Words != nil {
			if strings.Contains(*expr.Text.Words, "*") {
				// regex search
				words := strings.Split(*expr.Text.Words, " ")
				asFullText := []string{}

				var fb strings.Builder
				wildcardsCount := 0
				for i, word := range words {
					if strings.Contains(word, "*") && wildcardsCount <= 3 { // allow up to 3 wildcards
						word = multGlobPattern.ReplaceAllString(word, "*")
						// regex for words with globs (*)
						parts := strings.Split(word, "*")
						for i, part := range parts {
							parts[i] = regexp.QuoteMeta(part)
						}
						var b strings.Builder
						b.WriteString(`/.*?`)
						for j, part := range parts {
							if wildcardsCount > 3 {
								// discard remaining parts
								break
							}
							if part == "" {
								if j == len(parts)-1 {
									// if glob is last -> match until next word
									b.WriteString(`(\s*\S*)?`)
								}
								continue
							}
							if j != 0 {
								// max distance between matches of 5
								b.WriteString(`(\s*\S*){0,5}`)
							}

							b.WriteString(part)
							wildcardsCount++
						}
						b.WriteString(`/i`)
						arg := e.CrateArg(b.String(), "string")

						if i != 0 {
							fb.WriteString(` AND `)
						}
						fb.WriteString(`(regexp(Product.name, `)
						fb.WriteString(arg)
						fb.WriteString(`) OR regexp(Product.description, `)
						fb.WriteString(arg)
						fb.WriteString(`))`)

					} else {
						asFullText = append(asFullText, word)
					}
				}
				// use full text search for the rest
				if len(asFullText) > 0 {
					s := strings.Join(asFullText, " ")
					arg := e.CrateArg(s, "string")
					fb.WriteString(` AND (alloftext(Product.name, `)
					fb.WriteString(arg)
					fb.WriteString(`) OR alloftext(Product.description, `)
					fb.WriteString(arg)
					fb.WriteString(`))`)
				}
				filter := fmt.Sprintf(`@filter(%s)`, fb.String())
				curVar = e.addVariableWithFilter(filter, "", parVar)

			} else {
				// full text search
				arg := e.CrateArg(*expr.Text.Words, "string")
				filter := fmt.Sprintf(`@filter(
	alloftext(Product.name, %s)
	OR alloftext(Product.description, %s)
)`, arg, arg)
				curVar = e.addVariableWithFilter(filter, "", parVar)
			}

		} else if expr.Text.Exact != nil {
			// TODO: might also contain wildcards

			// exact text search
			// TODO: implement a custom strings.contains filter in database that allows case insensitive matching and use that instead of regex
			exact := regexp.QuoteMeta(*expr.Text.Exact)
			var b strings.Builder
			b.Grow(len(exact) + 10)
			b.WriteString(`/.*?`)
			b.WriteString(exact)
			b.WriteString(`.*/i`)

			arg := e.CrateArg(b.String(), "string")
			filter := fmt.Sprintf(`@filter(
	regexp(Product.name, %s)
	OR regexp(Product.description, %s)
)`, arg, arg)
			curVar = e.addVariableWithFilter(filter, "", parVar)
		}

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
	case booleanIsOperator, booleanHasOperator:
		// ignore; handled by special operators
		return

	case textFullContainsOperator, textTermExactOperator, textTermContainsOperator, textExactOperator:
		// text, full, match, not := e.extractOperatorText(opr)
		// if text == "" {
		// 	return
		// }
		// filter := ""
		// if full || match {
		// 	filter = e.exactTextFilter(o.Predicate, text, not, match)
		// } else {
		// 	filter = e.fullTextFilter(o.Predicate, text, not, (o.Type == textTermExactOperator || o.Type == textExactOperator))
		// }
		// if o.IsRootFilter {
		// 	sel := fmt.Sprintf(`%s %s`, o.SelectionStart, o.SelectionEnd)
		// 	curVar = e.addVariableWithFilter(filter, sel, parVar)
		// 	return
		// }

		// sel := fmt.Sprintf(`%s %s %s`, o.SelectionStart, filter, o.SelectionEnd)
		// curVar = e.addVariableWithFilter("", sel, parVar)
		// return

		text, exact, match, not := e.extractOperatorText(opr)
		if text == "" {
			return
		}
		filter := e.encodeTextFilter(o.Predicate, o.Type, text, exact, match)
		if filter == nil || len(filter) == 0 {
			return
		}
		if not {
			filter = e.encodeNotFilter(filter)
		}
		filter = e.encodeFilterExpression(filter)
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
	// Matches either full-text or phrases contained in the text.
	// Requires indexes: fulltext, regexp, exacti
	//   opr:text or opr:`text` -> full-text-match (alloftext)
	//   opr:"text" -> contains-match (regexp: /.../i)
	//   opr:* or opr:`text*` -> wildcard-contains-match (regexp: /.../i and alloftext)
	//   opr:"*" -> wildcard-contains-match (regexp: /.../i)
	//   opr:==text -> exact-match (allof with exacti index)
	textFullContainsOperator operatorType = iota
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
	// Product Operators
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
	"lastupdatedat": {
		Type:           dateTimeOperator,
		IsRootFilter:   true,
		Predicate:      "Product.lastUpdatedAt",
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
	// Component Operators
	//
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
	"discoveredat": {
		Type:           dateTimeOperator,
		Predicate:      "CrawlerMeta.discoveredAt",
		SelectionStart: `Product.release`,
		SelectionEnd:   `{uid}`,
	},
	"lastindexedat": {
		Type:           dateTimeOperator,
		Predicate:      "CrawlerMeta.lastIndexedAt",
		SelectionStart: `Product.release`,
		SelectionEnd:   `{uid}`,
	},
	"createdat": {
		Type:           dateTimeOperator,
		Predicate:      "Component.createdAt",
		SelectionStart: `Product.release`,
		SelectionEnd:   `{uid}`,
	},

	//
	// Licensor Operators
	//
	"licensor": {
		Type:           textTermExactOperator,
		Predicate:      "UserOrGroup.name",
		SelectionStart: `Product.release {Component.licensor`,
		SelectionEnd:   `{uid}}`,
	},
	"licensorname": {
		Type:           textTermExactOperator,
		Predicate:      "UserOrGroup.name",
		SelectionStart: `Product.release {Component.licensor`,
		SelectionEnd:   `{uid}}`,
	},
	"licensorfullname": {
		Type:           textTermExactOperator,
		Predicate:      "UserOrGroup.fullName",
		SelectionStart: `Product.release {Component.licensor`,
		SelectionEnd:   `{uid}}`,
	},

	//
	// Image Operators
	//
	"hasimage": {
		Type:           booleanHasOperator,
		Predicate:      "Component.image",
		SelectionStart: `Product.release`,
		SelectionEnd:   `{uid}`,
	},

	//
	// License Operators
	//
	"haslicense": {
		Type:           booleanHasOperator,
		Predicate:      "Component.license",
		SelectionStart: `Product.release`,
		SelectionEnd:   `{uid}`,
	},
	"license": {
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

// func (e *encoder) extractOperatorBoolean(opr *parser.Operator) (value, ok bool) {
// 	var val *parser.Text
// 	negate := false

// 	// check if the operator is a range, comparison or a plain value
// 	if opr.Range != nil {
// 		// ignore
// 		return
// 	} else if opr.Value != nil {
// 		val = opr.Value
// 	} else if opr.Comparison != nil {
// 		val = opr.Comparison.Value
// 		switch opr.Comparison.Operator {
// 		case parser.CompOpEq:
// 			break
// 		case parser.CompOpNe:
// 			negate = true
// 		default:
// 			// ignore
// 			return
// 		}
// 	}
// 	if val == nil {
// 		return
// 	}

// 	// extract the value as text
// 	text := ""
// 	if val.Words != nil {
// 		text = *val.Words
// 	} else {
// 		text = *val.Exact
// 	}

// 	// convert to boolean
// 	b, ok := stringutil.ParseBool(text)
// 	if !ok {
// 		return false, false
// 	}
// 	if negate {
// 		b = !b
// 	}
// 	return b, true
// }

// exactTextFilter returns a filter for exact text search. If match is true, the
// whole string needs to match.
func (e *encoder) exactTextFilter(predicate, text string, not, match bool) string {
	if strings.Contains(text, "*") {
		text = multGlobPattern.ReplaceAllString(text, "*")
		parts := strings.Split(text, "*")
		if len(parts) <= 3 { // allow up to 3 wildcards
			for i, part := range parts {
				parts[i] = regexp.QuoteMeta(part)
			}
			var b strings.Builder
			b.WriteRune('/')
			if match {
				b.WriteRune('^')
			}
			for j, part := range parts {
				if part == "" {
					if j == len(parts)-1 {
						// if glob is last -> match until end of string
						b.WriteString(`.*?`)
					}
					continue
				}
				if j != 0 {
					// max distance between matches of 5
					b.WriteString(`(\s*\S*){0,5}`)
				}
				b.WriteString(part)
			}
			if match {
				b.WriteRune('$')
			}
			b.WriteString(`/i`)
			arg := e.CrateArg(b.String(), "string")
			if not {
				return fmt.Sprintf(`@filter(NOT regexp(%s, %s))`, predicate, arg)
			}
			return fmt.Sprintf(`@filter(regexp(%s, %s))`, predicate, arg)
		}
	}

	// without wildcards
	text = regexp.QuoteMeta(text)
	// TODO: implement a custom strings.contains filter in database that allows case insensitive matching and use that instead of regex
	var b strings.Builder
	b.Grow(len(text) + 10)
	b.WriteRune('/')
	if match {
		b.WriteRune('^')
	}
	b.WriteString(text)
	if match {
		b.WriteRune('$')
	}
	b.WriteString(`/i`)
	arg := e.CrateArg(b.String(), "string")
	if not {
		return fmt.Sprintf(`@filter(NOT regexp(%s, %s))`, predicate, arg)
	}
	return fmt.Sprintf(`@filter(regexp(%s, %s))`, predicate, arg)
}

// fullTextFilter returns a filter for full text search. If terms is set,
// `allofterms` will be used instead of `alloftext`.
func (e *encoder) fullTextFilter(predicate, text string, not, terms bool) string {
	fltName := "alloftext"
	if terms {
		fltName = "allofterms"
	}

	// regex for words with globs (*)
	if strings.Contains(text, "*") {
		words := strings.Split(text, " ")
		asFullText := []string{}
		var fb strings.Builder
		for _, word := range words {
			if word == "" {
				continue
			}

			if strings.Contains(word, "*") {
				word = multGlobPattern.ReplaceAllString(word, "*")

				parts := strings.Split(word, "*")
				for i, part := range parts {
					parts[i] = regexp.QuoteMeta(part)
				}
				var b strings.Builder
				b.WriteString(`/`)
				for j, part := range parts {
					if part == "" {
						if j == len(parts)-1 {
							// if glob is last -> match until next word
							b.WriteString(`(\s*\S*)?`)
						}
						continue
					}
					if j != 0 {
						// max distance between matches of 5
						b.WriteString(`(\s*\S*){0,5}`)
					}

					b.WriteString(part)
				}
				b.WriteString(`/i`)
				arg := e.CrateArg(b.String(), "string")

				if fb.Len() > 0 {
					fb.WriteString(` AND `)
				}
				fb.WriteString(`regexp(`)
				fb.WriteString(predicate)
				fb.WriteString(`, `)
				fb.WriteString(arg)
				fb.WriteString(`)`)

			} else {
				asFullText = append(asFullText, word)
			}
		}
		// use full text search for the rest
		if len(asFullText) > 0 {
			s := strings.Join(asFullText, " ")
			arg := e.CrateArg(s, "string")
			fb.WriteString(` AND `)
			fb.WriteString(fltName)
			fb.WriteString(`(`)
			fb.WriteString(predicate)
			fb.WriteString(`, `)
			fb.WriteString(arg)
			fb.WriteString(`)`)
		}
		if not {
			return fmt.Sprintf(`@filter(NOT (%s))`, fb.String())
		}
		return fmt.Sprintf(`@filter(%s)`, fb.String())
	}

	// without wildcards
	arg := e.CrateArg(text, "string")
	if not {
		return fmt.Sprintf(`@filter(NOT %s(%s, %s))`, fltName, predicate, arg)
	}
	return fmt.Sprintf(`@filter(%s(%s, %s))`, fltName, predicate, arg)
}

func (e *encoder) encodeNotFilter(sub []byte) []byte {
	if sub == nil || len(sub) == 0 {
		return nil
	}
	var b bytes.Buffer
	b.WriteString(`NOT (`)
	b.Write(sub)
	b.WriteString(`)`)
	return b.Bytes()
}

func (e *encoder) encodeOrFilter(subs [][]byte) []byte {
	if subs == nil || len(subs) == 0 {
		return nil
	}
	var b bytes.Buffer
	first := true
	for _, sub := range subs {
		if sub == nil || len(sub) == 0 {
			continue
		}
		if !first {
			b.WriteString(` OR `)
		}
		b.Write(sub)
		first = false
	}
	return b.Bytes()
}

func (e *encoder) encodeFilterExpression(filter []byte) []byte {
	if filter == nil || len(filter) == 0 {
		return nil
	}
	var b bytes.Buffer
	b.WriteString(`@filter(`)
	b.Write(filter)
	b.WriteString(`)`)
	return b.Bytes()
}

// encodeTextFilter ...
//
// match indicates whether the text need to be matched (enclosed with quotes "")
func (e *encoder) encodeTextFilter(predicate string, oprTyp operatorType, text string, exact, match bool) []byte {
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

func extractDateTimeValue(rawVal *parser.Text) (dt time.Time, ok bool) {
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

func parseDateTimeValue(rawVal *string) (dt time.Time, ok bool) {
	strVal := strings.TrimSpace(*rawVal)

	// try parsing as time duration
	duration, ok := parseDuration(strVal)
	if ok {
		return time.Now().Add(-1 * duration), true
	}

	// try parsing as date time
	c := carbon.Parse(strVal)
	if c.Error == nil {
		return c.Carbon2Time(), true
	}

	// invalid value -> ignore
	return
}

var timeDurationPattern = regexp.MustCompile(`^(?P<years>\d+y)?(?P<months>\d+M)?(?P<weeks>\d+w)?(?P<days>\d+d)?T?(?P<hours>\d+h)?(?P<minutes>\d+m)?(?P<seconds>\d+s)?$`)

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
	hours := parseDurationMatch(matches[5])
	minutes := parseDurationMatch(matches[6])
	seconds := parseDurationMatch(matches[7])

	hour := int64(time.Hour)
	minute := int64(time.Minute)
	second := int64(time.Second)
	return time.Duration(years*24*365*hour + months*30*24*hour + weeks*7*24*hour + days*24*hour + hours*hour + minutes*minute + seconds*second), true
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
			dt, ok := parseDateTimeValue(opr.Range.Start)
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
			dt, ok := parseDateTimeValue(opr.Range.End)
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
	dt, ok := extractDateTimeValue(txtVal)
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
	case searchmodels.OrderByCreatedAt:
		e.buf.WriteString(`Product.release { O1 as Component.createdAt } order as min(val(O1))`)
	case searchmodels.OrderByDiscoveredAt:
		e.buf.WriteString(`order as CrawlerMeta.discoveredAt`)
	case searchmodels.OrderByLastIndexedAt:
		e.buf.WriteString(`order as CrawlerMeta.lastIndexedAt`)
	case searchmodels.OrderByDocumentationLanguage:
		e.buf.WriteString(`order as Product.documentationLanguage`)
	case searchmodels.OrderByState:
		e.buf.WriteString(`order as Product.state`)
	case searchmodels.OrderByForkCount:
		e.buf.WriteString(`order as Product.forkCount`)
	case searchmodels.OrderByStarCount:
		e.buf.WriteString(`order as Product.starCount`)
	case searchmodels.OrderByLicense:
		e.buf.WriteString(`Product.release { Component.license { O2 as License.xid } O1 as min(val(O2)) } order as min(val(O1))`)
	default:
		panic("unsupported order by field")
	}

	e.buf.WriteString("}\n")
	return parVar
}

func (e *encoder) appendSelectQuery(query, order, parVar string) {
	query = fmt.Sprintf(query, parVar, order)
	e.buf.WriteString(query)
	e.buf.WriteString("\n")
}
