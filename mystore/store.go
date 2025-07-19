package mystore

import (
	"context"
	"log"
	"strings"
	"sync"
)

type InMemoryStore[T any] struct {
	sync.Mutex
	Items map[string]T
}

func NewInMemoryStore[T any](c context.Context) (*InMemoryStore[T], func(), error) {
	return &InMemoryStore[T]{
		Items: make(map[string]T),
	}, func() {}, nil
}

func (s *InMemoryStore[T]) RunInTransaction(c context.Context, f func(c context.Context) error) error {
	// Start transaction: acquire lock for the duration of the transaction
	s.Lock()
	defer s.Unlock() // Ensure the lock is released when the function exits

	ctx := context.WithValue(c, ctxTransactionKey{}, true)

	// Within this block everything is transactional
	log.Printf("Func %p with context %p", f, ctx)
	err := f(ctx)
	if err != nil {
		// Rollback: lock is released by defer
		return err
	}

	// Commit: lock is released by defer
	return nil
}

func (s *InMemoryStore[T]) Put(c context.Context, uid string, value T) error {
	nonTransactional := c.Value(ctxTransactionKey{}) == nil

	if nonTransactional {
		s.Lock()
		defer s.Unlock() // Acquire and release lock if not in a transaction
	}

	uid = strings.ToLower(uid)
	s.Items[uid] = value

	return nil
}

func (s *InMemoryStore[T]) Get(c context.Context, uid string) (T, bool, error) {
	nonTransactional := c.Value(ctxTransactionKey{}) == nil

	if nonTransactional {
		s.Lock()
		defer s.Unlock() // Acquire and release lock if not in a transaction
	}
	uid = strings.ToLower(uid)
	result, exists := s.Items[uid]

	return result, exists, nil
}

func (s *InMemoryStore[T]) List(c context.Context) ([]T, error) {
	nonTransactional := c.Value(ctxTransactionKey{}) == nil

	if nonTransactional {
		s.Lock()
		defer s.Unlock() // Acquire and release lock if not in a transaction
	}

	result := make([]T, 0, len(s.Items))
	for _, v := range s.Items {
		result = append(result, v)
	}

	return result, nil
}

func (s *InMemoryStore[T]) Query(c context.Context, filters []Filter, orderByField string) ([]T, error) {
	// This method currently just calls List, effectively ignoring filters and orderByField.
	// A full implementation would apply the filtering and sorting logic here.
	return s.List(c)
}
