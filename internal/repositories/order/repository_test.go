package order_test

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"

	morder "github.com/gulmix/orders/internal/models/order"
	rorder "github.com/gulmix/orders/internal/repositories/order"
)

func mustOrder(t *testing.T, idx int) *morder.Order {
	t.Helper()

	item, err := morder.NewItem(fmt.Sprintf("test-sku-%v", idx), 1, 100)
	if err != nil {
		t.Fatalf("mustOrder: failed to create item: %v", err)
	}

	ord, err := morder.NewOrder(fmt.Sprintf("test-customer-%v", idx), []morder.Item{item}, time.Now())
	if err != nil {
		t.Fatalf("mustOrder: failed to create order: %v", err)
	}

	return ord
}

func TestRepository_SaveAndGet(t *testing.T) {
	repo := rorder.New()
	testOrder := mustOrder(t, 13)
	errSave := repo.Save(context.Background(), testOrder)
	if errSave != nil {
		t.Fatalf("unexpected error while saving order: %v", errSave)
	}
	gotOrder, errGet := repo.Get(context.Background(), testOrder.ID())
	if errGet != nil {
		t.Fatalf("unexpected error while getting order: %v", errGet)
	}
	if gotOrder != testOrder {
		t.Errorf("Get(): got %+v, want %+v", gotOrder, testOrder)
	}
}

func TestRepository_Get_NotFound(t *testing.T) {
	repo := rorder.New()
	_, errGet := repo.Get(context.Background(), uuid.New())
	if !errors.Is(errGet, rorder.ErrNotFound) {
		t.Fatalf("Get(): got %v, want %v", errGet, rorder.ErrNotFound)
	}
}

func TestRepository_Save_AlreadyExists(t *testing.T) {
	repo := rorder.New()
	testOrder := mustOrder(t, 13)
	errFirstSave := repo.Save(context.Background(), testOrder)
	if errFirstSave != nil {
		t.Fatalf("unexpected error while saving first order: %v", errFirstSave)
	}
	errSecondSave := repo.Save(context.Background(), testOrder)
	if !errors.Is(errSecondSave, rorder.ErrAlreadyExists) {
		t.Fatalf("Save(): got %v, want %v", errSecondSave, rorder.ErrAlreadyExists)
	}
}

func TestRepository_Save_NilOrder(t *testing.T) {
	repo := rorder.New()
	err := repo.Save(context.Background(), nil)
	if !errors.Is(err, rorder.ErrNilOrder) {
		t.Fatalf("Save(): got %v, want %v", err, rorder.ErrNilOrder)
	}
}

func TestRepository_Save_ContextCancel(t *testing.T) {
	repo := rorder.New()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := repo.Save(ctx, mustOrder(t, 13))
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Save(): got %v, want %v", err, context.Canceled)
	}
}

func TestRepository_Get_ContextCancel(t *testing.T) {
	repo := rorder.New()
	testOrder := mustOrder(t, 13)
	errSave := repo.Save(context.Background(), testOrder)
	if errSave != nil {
		t.Fatalf("unexpected error when saving order to repository: %v", errSave)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, errGet := repo.Get(ctx, testOrder.ID())
	if !errors.Is(errGet, context.Canceled) {
		t.Fatalf("Get(): got %v, want %v", errGet, context.Canceled)
	}
}

func TestRepository_Race(t *testing.T) {
	repo := rorder.New()
	const goroutines int = 50
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := range goroutines {
		testOrder := mustOrder(t, i)
		go func() {
			defer wg.Done()
			errSave := repo.Save(context.Background(), testOrder)
			if errSave != nil {
				t.Errorf("unexpected error when saving order to repo: %v", errSave)
				return
			}
			gotOrder, errGet := repo.Get(context.Background(), testOrder.ID())
			if errGet != nil {
				t.Errorf("unexpected error when getting order from repo: %v", errGet)
				return
			}
			if testOrder != gotOrder {
				t.Errorf("Get(): got %+v, want %+v", gotOrder, testOrder)
				return
			}
		}()
	}
	wg.Wait()
}
