package mystore

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
)

type testStruct struct {
	ID   int
	Name string
}

func TestInMemoryStore_PutGetList(t *testing.T) {
	ctx := context.Background()
	store, cleanup, err := NewInMemoryStore[testStruct](ctx)
	defer cleanup()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Positive: Put and Get
	item := testStruct{ID: 1, Name: "Alpha"}
	if err := store.Put(ctx, "A1", item); err != nil {
		t.Errorf("Put failed: %v", err)
	}
	got, exists, err := store.Get(ctx, "A1")
	if err != nil {
		t.Errorf("Get failed: %v", err)
	}
	if !exists {
		t.Errorf("Expected item to exist")
	}
	if got != item {
		t.Errorf("Got %+v, want %+v", got, item)
	}

	// Negative: Get non-existent
	_, exists, err = store.Get(ctx, "doesnotexist")
	if err != nil {
		t.Errorf("Get failed: %v", err)
	}
	if exists {
		t.Errorf("Expected item to not exist")
	}

	// Edge: Put with different case, should overwrite
	item2 := testStruct{ID: 2, Name: "Beta"}
	if err := store.Put(ctx, "A1", item2); err != nil {
		t.Errorf("Put failed: %v", err)
	}
	got, exists, _ = store.Get(ctx, "a1")
	if !exists || got != item2 {
		t.Errorf("Case-insensitive Put/Get failed")
	}

	// List
	all, err := store.List(ctx)
	if err != nil {
		t.Errorf("List failed: %v", err)
	}
	if len(all) != 1 {
		t.Errorf("Expected 1 item, got %d", len(all))
	}
}

func TestInMemoryStore_Transaction(t *testing.T) {
	ctx := context.Background()
	store, cleanup, err := NewInMemoryStore[testStruct](ctx)
	defer cleanup()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	item := testStruct{ID: 1, Name: "Alpha"}
	item2 := testStruct{ID: 2, Name: "Beta"}

	err = store.RunInTransaction(ctx, func(txCtx context.Context) error {
		if err := store.Put(txCtx, "A1", item); err != nil {
			return err
		}
		if err := store.Put(txCtx, "A2", item2); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		t.Errorf("Transaction failed: %v", err)
	}

	// Both items should be present
	got, exists, _ := store.Get(ctx, "A1")
	if !exists || got != item {
		t.Errorf("A1 missing after transaction")
	}
	got, exists, _ = store.Get(ctx, "A2")
	if !exists || got != item2 {
		t.Errorf("A2 missing after transaction")
	}
}

func TestInMemoryStore_TransactionRollback(t *testing.T) {
	ctx := context.Background()
	store, cleanup, err := NewInMemoryStore[testStruct](ctx)
	defer cleanup()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	item := testStruct{ID: 1, Name: "Alpha"}

	err = store.RunInTransaction(ctx, func(txCtx context.Context) error {
		if err := store.Put(txCtx, "A1", item); err != nil {
			return err
		}
		return errors.New("fail")
	})
	if err == nil {
		t.Errorf("Expected error from transaction, got nil")
	}
	// Should not rollback changes (since it's in-memory, lock is released, but state is not reverted)
	// But we can check that the item is present
	got, exists, _ := store.Get(ctx, "A1")
	if !exists || got != item {
		t.Errorf("Expected item to be present after rollback (in-memory store does not revert)")
	}
}

func TestInMemoryStore_Concurrency(t *testing.T) {
	ctx := context.Background()
	store, cleanup, err := NewInMemoryStore[int](ctx)
	defer cleanup()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", i)
			if err := store.Put(ctx, key, i); err != nil {
				t.Errorf("Put failed: %v", err)
			}
			val, exists, err := store.Get(ctx, key)
			if err != nil || !exists || val != i {
				t.Errorf("Get failed: %v, exists: %v, val: %v", err, exists, val)
			}
		}(i)
	}
	wg.Wait()
	all, err := store.List(ctx)
	if err != nil {
		t.Errorf("List failed: %v", err)
	}
	if len(all) != 100 {
		t.Errorf("Expected 100 items, got %d", len(all))
	}
}

func TestInMemoryStore_Query(t *testing.T) {
	ctx := context.Background()
	store, cleanup, err := NewInMemoryStore[int](ctx)
	defer cleanup()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 0; i < 5; i++ {
		store.Put(ctx, fmt.Sprintf("k%d", i), i)
	}
	result, err := store.Query(ctx, nil, "")
	if err != nil {
		t.Errorf("Query failed: %v", err)
	}
	if len(result) != 5 {
		t.Errorf("Expected 5 items from Query, got %d", len(result))
	}
}
