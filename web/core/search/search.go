package search

import (
	"context"
	"strings"
	"unicode/utf8"

	searchmodels "losh/web/core/search/models"
	"losh/web/core/search/parser"

	"github.com/aisbergg/go-errors/pkg/errors"
)

func (s *Service) Search3(ctx context.Context, queryStr string, orderBy searchmodels.OrderBy, pagination searchmodels.Pagination) (searchmodels.Results, error) {
	query, operators, err := s.parseQuery(queryStr)
	if err != nil {
		return searchmodels.Results{}, err
	}

	prds, count, err := s.repo.SearchProducts(ctx, query, orderBy, pagination)
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
