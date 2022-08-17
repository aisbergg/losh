package search

import (
	"context"
	"strings"
	"unicode/utf8"

	"losh/internal/core/product/models"
	"losh/internal/infra/dgraph/dgclient"
)

const (
	maxQueryStringLength = 256
)

type Repository interface {
	SearchProducts(ctx context.Context, filter *dgclient.ProductFilter, order *dgclient.ProductOrder, first *int64, offset *int64) ([]*models.Product, int64, error)
}

type Results struct {
	Count int64             `json:"count" liquid:"count"`
	Items []*models.Product `json:"items" liquid:"items"`
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Search(ctx context.Context, queryStr string, orderBy string, orderAscending bool, first, offset int64) (Results, error) {
	queryStr = strings.TrimSpace(queryStr)

	if utf8.RuneCountInString(queryStr) > maxQueryStringLength {
		return Results{}, &Error{"query too long"}
	}

	// if queryStr == "" {
	// 	return Results{}, nil
	// }
	// parse query
	// TODO: implement https://discuss.dgraph.io/t/rfc-nested-filters-in-graphql/13560
	query, err := s.parseQuery(queryStr)
	if err != nil {
		return Results{}, &Error{"invalid search query: " + err.Error()}
	}

	// query does not contain any search terms -> return all products

	var filter *dgclient.ProductFilter
	if query.query != "" {
		// convert query to graphQL
		graphQLQuery := s.searchQuery2GraphQL(query)

		// execute graphQL query
		filter = &dgclient.ProductFilter{
			Or: []*dgclient.ProductFilter{
				&dgclient.ProductFilter{
					Name: &dgclient.StringFullTextFilterStringHashFilterStringRegExpFilter{
						Anyoftext: &graphQLQuery,
					},
				},
				&dgclient.ProductFilter{
					Description: &dgclient.StringFullTextFilterStringRegExpFilter{
						Anyoftext: &graphQLQuery,
					},
				},
			},
		}
	}
	ordMde, ok := orderByLookup[orderBy]
	if !ok {
		ordMde = orderByLookup["name"]
	}
	order := &dgclient.ProductOrder{}
	if orderAscending {
		order.Asc = &ordMde
	} else {
		order.Desc = &ordMde
	}
	prds, count, err := s.repo.SearchProducts(ctx, filter, order, &first, &offset)
	if err != nil {
		return Results{}, err
	}

	return Results{Count: count, Items: prds}, nil
}

// TODO: need better ways to order results
var orderByLookup = map[string]dgclient.ProductOrderable{
	"discoveredAt":  dgclient.ProductOrderableDiscoveredAt,
	"lastIndexedAt": dgclient.ProductOrderableLastIndexedAt,
	"xid":           dgclient.ProductOrderableXid,
	"name":          dgclient.ProductOrderableName,
	"website":       dgclient.ProductOrderableWebsite,
	"version":       dgclient.ProductOrderableVersion,
}

// parseQuery parses a query string and returns a query object.
func (s *Service) parseQuery(queryStr string) (query Query, err error) {
	return Query{query: queryStr}, nil
}

// searchQuery2GraphQL converts a query object to a graphQL query.
func (s *Service) searchQuery2GraphQL(query Query) string {
	// TODO: implement
	return query.query
}

// Query represents a search query.
type Query struct {
	// TODO: implement
	query string
}
