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

package dgraph

import (
	"losh/internal/core/product/models"
)

type PaginatableType interface {
	*models.Product
}

type ItemProvider[T PaginatableType] func(first int, offset int) (items []T, total uint64, err error)

type PaginatedList[T PaginatableType] struct {
	itemPrvFn ItemProvider[T]

	first  int
	offset int
	total  uint64
	i      int
	items  []T
}

// NewPaginatedList returns a new PaginatedList.
func NewPaginatedList[T PaginatableType](itemPrvFn ItemProvider[T], first, offset int) *PaginatedList[T] {
	return &PaginatedList[T]{
		itemPrvFn: itemPrvFn,
		first:     first,
		offset:    offset,
	}
}

// Next returns the next item in the list.
func (pl *PaginatedList[T]) Next() (T, error) {
	if pl.items == nil || pl.i >= len(pl.items) {
		if err := pl.nextPage(); err != nil {
			return nil, err
		}
		if len(pl.items) == 0 {
			return nil, nil
		}
	}
	item := pl.items[pl.i]
	pl.i++
	return item, nil
}

// HasNext returns true if there are more items in the list.
func (pl *PaginatedList[T]) HasNext() (bool, error) {
	if pl.items == nil || pl.i >= len(pl.items) {
		if err := pl.nextPage(); err != nil {
			return false, err
		}
	}
	return len(pl.items) > 0, nil
}

// nextPage returns the next page of items in the list.
func (pl *PaginatedList[T]) nextPage() error {
	// going past the end of the list results in all items being returned ->
	// check if we're already at the end
	if pl.offset > 0 && uint64(pl.offset) >= pl.total {
		pl.items = []T{}
		return nil
	}
	items, total, err := pl.itemPrvFn(pl.first, pl.offset)
	if err != nil {
		return err
	}
	pl.items = items
	pl.total = total
	pl.i = 0
	pl.offset += int(len(items))
	return nil
}
