package order_test

import (
	"context"
	"testing"
	"time"

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

// func TestRepository_Get_NotFound()
