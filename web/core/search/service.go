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
