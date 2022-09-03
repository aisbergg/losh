package search

import (
	"context"

	"losh/internal/core/product/models"
	searchmodels "losh/web/core/search/models"
	"losh/web/core/search/parser"
)

type Repository interface {
	SearchProducts(ctx context.Context, query *parser.Query, order searchmodels.OrderBy, pagination searchmodels.Pagination) ([]*models.Product, uint64, error)
}

type Service struct {
	repo  Repository
	debug bool
}

func NewService(repo Repository, debug bool) *Service {
	return &Service{
		repo:  repo,
		debug: debug,
	}
}
