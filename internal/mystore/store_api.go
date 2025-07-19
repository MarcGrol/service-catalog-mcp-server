package mystore

import (
	"context"
)

// ctxTransactionKey is an unexported type for context keys to avoid collisions.
type ctxTransactionKey struct{}

// Filter is a placeholder type for filtering operations in the Query method.
// Its actual structure would depend on the specific filtering logic required.
type Filter struct {
	Field string
	Value string
	Op    string // e.g., "eq", "gt", "lt"
}

//go:generate mockgen -source=api.go -package mystore -destination store_mock.go Store
type Store[T any] interface {
	RunInTransaction(c context.Context, f func(c context.Context) error) error
	Put(c context.Context, uid string, value T) error
	Get(c context.Context, uid string) (T, bool, error)
	List(c context.Context) ([]T, error)
	Query(c context.Context, filters []Filter, orderByField string) ([]T, error)
}

func New[T any](c context.Context) (Store[T], func(), error) {
	return NewInMemoryStore[T](c)
}
