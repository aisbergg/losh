package search

import (
	"context"
	"strings"
	"unicode/utf8"

	searchmodels "losh/web/core/search/models"
	"losh/web/core/search/parser"

	"github.com/aisbergg/go-errors/pkg/errors"
)

// func (s *Service) Search(ctx context.Context, queryStr string, orderBy string, orderAscending bool, first, offset int64) (searchmodels.Results, error) {
// 	queryStr = strings.TrimSpace(queryStr)

// 	if utf8.RuneCountInString(queryStr) > maxQueryStringLength {
// 		return searchmodels.Results{}, &Error{"query too long", ErrorLimitExceeded}
// 	}

// 	// if queryStr == "" {
// 	// 	return searchmodels.Results{}, nil
// 	// }
// 	// parse query
// 	// TODO: implement https://discuss.dgraph.io/t/rfc-nested-filters-in-graphql/13560
// 	query, err := s.parseQuery(queryStr)
// 	if err != nil {
// 		return searchmodels.Results{}, &Error{"invalid search query: " + err.Error(), ErrorInvalidQuery}
// 	}

// 	// query does not contain any search terms -> return all products

// 	var filter *dgclient.ProductFilter
// 	if query.query != "" {
// 		// convert query to graphQL
// 		graphQLQuery := s.searchQuery2GraphQL(query)

// 		// execute graphQL query
// 		filter = &dgclient.ProductFilter{
// 			Or: []*dgclient.ProductFilter{
// 				&dgclient.ProductFilter{
// 					Name: &dgclient.StringFullTextFilterStringHashFilterStringRegExpFilter{
// 						Anyoftext: &graphQLQuery,
// 					},
// 				},
// 				&dgclient.ProductFilter{
// 					Description: &dgclient.StringFullTextFilterStringRegExpFilter{
// 						Anyoftext: &graphQLQuery,
// 					},
// 				},
// 			},
// 		}
// 	}

// 	ordMde, ok := orderByLookup[orderBy]
// 	if !ok {
// 		ordMde = orderByLookup["name"]
// 	}
// 	order := &dgclient.ProductOrder{}
// 	if orderAscending {
// 		order.Asc = &ordMde
// 	} else {
// 		order.Desc = &ordMde
// 	}
// 	prds, count, err := s.repo.SearchProducts(ctx, filter, order, &first, &offset)
// 	if err != nil {
// 		return searchmodels.Results{}, err
// 	}

// 	return searchmodels.Results{Count: count, Items: prds}, nil
// }

// func (s *Service) Search2(ctx context.Context, queryStr string, orderBy string, orderAscending bool, first, offset int64) (searchmodels.Results, error) {
// 	queryStr = strings.TrimSpace(queryStr)

// 	if utf8.RuneCountInString(queryStr) > maxQueryStringLength {
// 		return searchmodels.Results{}, &Error{"query too long", ErrorLimitExceeded}
// 	}

// 	// TODO: return all results
// 	if queryStr == "" {
// 		return searchmodels.Results{}, nil
// 	}

// 	// parse query
// 	query, err := parser.Parse(queryStr)
// 	if err != nil {
// 		return searchmodels.Results{}, err
// 	}
// 	repr.Println(query, repr.Indent("  "), repr.OmitEmpty(false))

// 	// translate query to graphQL
// 	prdFlt, lcsFlt := translateQuery(query)
// 	repr.Println(prdFlt, repr.Indent("  "), repr.OmitEmpty(false))
// 	_ = lcsFlt

// 	ordMde, ok := orderByLookup[orderBy]
// 	if !ok {
// 		ordMde = orderByLookup["name"]
// 	}
// 	order := &dgclient.ProductOrder{}
// 	if orderAscending {
// 		order.Asc = &ordMde
// 	} else {
// 		order.Desc = &ordMde
// 	}
// 	prds, count, err := s.repo.SearchProducts(ctx, prdFlt, order, &first, &offset)
// 	if err != nil {
// 		return searchmodels.Results{}, err
// 	}

// 	return searchmodels.Results{Count: count, Items: prds}, nil
// }

func (s *Service) Search3(ctx context.Context, queryStr string, orderBy searchmodels.OrderBy, pagination searchmodels.Pagination) (searchmodels.Results, error) {
	query, operators, err := s.parseQuery(queryStr)
	if err != nil {
		return searchmodels.Results{}, err
	}

	prds, count, err := s.repo.SearchProductsDQL(ctx, query, orderBy, pagination)
	if err != nil {
		return searchmodels.Results{}, err
	}

	return searchmodels.Results{
		Count:     int64(count),
		Items:     prds,
		Operators: operators,
	}, nil
}

// parseQuery parses the query string into a Query object and makes sure the
// various limits are not exceeded.
func (s *Service) parseQuery(queryStr string) (query *parser.Query, operators []string, err error) {
	queryStr = strings.TrimSpace(queryStr)

	// check max length limit
	if utf8.RuneCountInString(queryStr) > maxQueryStringLength {
		return nil, nil, &Error{"query too long", ErrorLimitExceeded}
	}

	// parse query
	if queryStr != "" {
		query, err = parser.Parse(queryStr)
		if err != nil && s.debug {
			return nil, nil, &Error{errors.ToString(err, false), ErrorInvalidQuery}
		}
	}
	// repr.Println(query, repr.Indent("  "), repr.OmitEmpty(false))

	// check other limits
	limiter := newLimiter()
	// return (&limiter{}).check(query)
	// err = checkLimits(query)
	if err = limiter.check(query); err != nil {
		return nil, nil, err
	}

	return query, limiter.getOperators(), nil
}

// func translateQuery(query *parser.Query) (*dgclient.ProductFilter, *dgclient.LicenseFilter) {
// 	var prdFlt *dgclient.ProductFilter
// 	var lcsFlt *dgclient.LicenseFilter

// 	var prdOrFlts []*dgclient.ProductFilter
// 	for _, orCnd := range query.Or {
// 		var prdAndFlts []*dgclient.ProductFilter

// 		for _, andCnd := range orCnd.And {
// 			pf, lf := translateAndCondition(andCnd)
// 			if pf != nil {
// 				if prdAndFlts == nil {
// 					prdAndFlts = []*dgclient.ProductFilter{}
// 				}
// 				prdAndFlts = append(prdAndFlts, pf)
// 			}
// 			_ = lf
// 		}

// 		if prdAndFlts != nil {
// 			if prdOrFlts == nil {
// 				prdOrFlts = []*dgclient.ProductFilter{}
// 			}
// 			prdOrFlts = append(prdOrFlts, &dgclient.ProductFilter{
// 				And: prdAndFlts,
// 			})
// 		}
// 	}
// 	if prdOrFlts != nil {
// 		prdFlt = &dgclient.ProductFilter{
// 			Or: prdOrFlts,
// 		}
// 	}

// 	return prdFlt, lcsFlt
// }

// func translateAndCondition(andCnd *parser.AndCondition) (*dgclient.ProductFilter, *dgclient.LicenseFilter) {
// 	if andCnd.Not != nil {
// 		pf, lf := translateAndCondition(andCnd.Not)
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
// 		return translateExpression(andCnd.Operand)
// 	}
// }

// func translateExpression(expr *parser.Expression) (pf *dgclient.ProductFilter, lf *dgclient.LicenseFilter) {
// 	if expr.Text != nil {
// 		if expr.Text.Words != nil {
// 			if strings.Contains(*expr.Text.Words, "*") {
// 				// regex search
// 				words := strings.Split(*expr.Text.Words, " ")
// 				asFullText := []string{}
// 				pf = &dgclient.ProductFilter{
// 					And: []*dgclient.ProductFilter{},
// 				}

// 				for _, word := range words {
// 					if strings.Contains(word, "*") {
// 						// regex for words with globs (*)
// 						parts := strings.Split(word, "*")
// 						for i, part := range parts {
// 							parts[i] = regexp.QuoteMeta(part)
// 						}
// 						var b strings.Builder
// 						b.WriteString(`/.*?`)
// 						for _, part := range parts {
// 							if part == "" {
// 								continue
// 							}
// 							b.WriteString(part)
// 							b.WriteString(`.*?`)
// 						}
// 						b.WriteString(`/i`)
// 						s := b.String()
// 						pf.And = append(pf.And, &dgclient.ProductFilter{
// 							Or: []*dgclient.ProductFilter{
// 								&dgclient.ProductFilter{
// 									Name: &dgclient.StringFullTextFilterStringHashFilterStringRegExpFilter{
// 										Regexp: &s,
// 									},
// 								},
// 								&dgclient.ProductFilter{
// 									Description: &dgclient.StringFullTextFilterStringRegExpFilter{
// 										Regexp: &s,
// 									},
// 								},
// 							},
// 						})

// 					} else {
// 						asFullText = append(asFullText, word)
// 					}
// 				}
// 				// use full text search for the rest
// 				if len(asFullText) > 0 {
// 					s := strings.Join(asFullText, " ")
// 					pf.And = append(pf.And, &dgclient.ProductFilter{
// 						Or: []*dgclient.ProductFilter{
// 							&dgclient.ProductFilter{
// 								Name: &dgclient.StringFullTextFilterStringHashFilterStringRegExpFilter{
// 									Alloftext: &s,
// 								},
// 							},
// 							&dgclient.ProductFilter{
// 								Description: &dgclient.StringFullTextFilterStringRegExpFilter{
// 									Alloftext: &s,
// 								},
// 							},
// 						},
// 					})
// 				}

// 			} else {
// 				// full text search
// 				pf = &dgclient.ProductFilter{
// 					Or: []*dgclient.ProductFilter{
// 						&dgclient.ProductFilter{
// 							Name: &dgclient.StringFullTextFilterStringHashFilterStringRegExpFilter{
// 								Alloftext: expr.Text.Words,
// 							},
// 						},
// 						&dgclient.ProductFilter{
// 							Description: &dgclient.StringFullTextFilterStringRegExpFilter{
// 								Alloftext: expr.Text.Words,
// 							},
// 						},
// 					},
// 				}
// 			}

// 		} else if expr.Text.Exact != nil {
// 			// exact text search
// 			var b strings.Builder
// 			b.WriteString(`/.*?`)
// 			b.WriteString(regexp.QuoteMeta(*expr.Text.Exact))
// 			b.WriteString(`.*/i`)
// 			s := b.String()
// 			pf = &dgclient.ProductFilter{
// 				// TODO: replace with custom strings.contains filter instead of regex
// 				Or: []*dgclient.ProductFilter{
// 					&dgclient.ProductFilter{
// 						Name: &dgclient.StringFullTextFilterStringHashFilterStringRegExpFilter{
// 							Regexp: &s,
// 						},
// 					},
// 					&dgclient.ProductFilter{
// 						Description: &dgclient.StringFullTextFilterStringRegExpFilter{
// 							Regexp: &s,
// 						},
// 					},
// 				},
// 			}
// 		}

// 	} else if expr.Operator != nil {
// 		// operator search
// 		switch expr.Operator.Name {
// 		case "license":
// 			lf = &dgclient.LicenseFilter{
// 				Or: []*dgclient.LicenseFilter{},
// 			}
// 			if expr.Operator.Value != nil {

// 			} else if expr.Operator.Range != nil {
// 				// not supported

// 			} else if expr.Operator.Comparison != nil {
// 				if expr.Operator.Comparison.Operator == parser.CompOpEq {
// 					// lf = &dgclient.LicenseFilter{
// 					// 	// TODO: other comparison
// 					// 	Xid: &dgclient.StringHashFilter{
// 					// 		Eq: expr.Operator.Comparison.Value.Words,
// 					// 	},
// 					// }
// 				}
// 				// else if expr.Operator.Comparison.Operator == parser.CompOpNeq {

// 				// }
// 			}
// 		}

// 	} else if expr.Sub != nil {
// 		// more deeply nested expression
// 		return translateQuery(expr.Sub)
// 	}

// 	return
// }
