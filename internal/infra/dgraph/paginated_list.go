package dgraph

import (
	"losh/internal/core/product/models"
)

type PaginatableType interface {
	*models.Product
}

type ItemProvider[T PaginatableType] func(first int64, offset int64) ([]T, error)

type PaginatedList[T PaginatableType] struct {
	itemPrvFn ItemProvider[T]

	first  int64
	offset int64
	i      int
	items  []T
}

// NewPaginatedList returns a new PaginatedList.
func NewPaginatedList[T PaginatableType](itemPrvFn func(first int64, offset int64) ([]T, error), batchSize int64) *PaginatedList[T] {
	return &PaginatedList[T]{
		itemPrvFn: itemPrvFn,
		first:     batchSize,
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
	items, err := pl.itemPrvFn(pl.first, pl.offset)
	if err != nil {
		return err
	}
	pl.items = items
	pl.i = 0
	pl.offset += int64(len(items))
	return nil
}
