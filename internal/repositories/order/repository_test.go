package order_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	morder "github.com/gulmix/orders/internal/models/order"
	rorder "github.com/gulmix/orders/internal/repositories/order"
)

func mustOrder(t *testing.T) *morder.Order {
	t.Helper()

	item, err := morder.NewItem("test-sku", 1, 100)
	if err != nil {
		t.Fatalf("mustOrder: failed to create item: %v", err)
	}

	ord, err := morder.NewOrder("test-customer", []morder.Item{item}, time.Now())
	if err != nil {
		t.Fatalf("mustOrder: failed to create order: %v", err)
	}

	return ord
}

func TestRepository_SaveAndGet(t *testing.T) {
	repo := rorder.New()
	testOrder := mustOrder(t)
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
	testOrder := mustOrder(t)
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
