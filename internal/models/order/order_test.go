package order_test

import (
	"errors"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gulmix/orders/internal/models/order"
)

func TestNewItem(t *testing.T) {
	tests := []struct {
		name       string
		sku        string
		qty        int
		priceCents int64
		wantErr    error
	}{
		{name: "valid", sku: "soundtrack", qty: 1, priceCents: 249},
		{name: "valid with zero price", sku: "soundtrack", qty: 1, priceCents: 0},
		{name: "valid with max allowed price", sku: "soundtrack", qty: 1, priceCents: order.MaxPriceCents},
		{name: "valid with max allowed qty", sku: "soundtrack", qty: order.MaxQty, priceCents: 249},
		{name: "zero qty", sku: "soundtrack", qty: 0, priceCents: 249, wantErr: order.ErrInvalidQty},
		{name: "negative qty", sku: "soundtrack", qty: -123, priceCents: 249, wantErr: order.ErrInvalidQty},
		{name: "qty too large", sku: "soundtrack", qty: order.MaxQty + 1, priceCents: 249, wantErr: order.ErrQtyTooLarge},
		{name: "negative price", sku: "soundtrack", qty: 1, priceCents: -249, wantErr: order.ErrNegativePrice},
		{name: "price too large", sku: "soundtrack", qty: 1, priceCents: order.MaxPriceCents + 1, wantErr: order.ErrPriceTooLarge},
		{name: "empty sku", sku: "", qty: 1, priceCents: 249, wantErr: order.ErrEmptySKU},
		{name: "sku with spaces", sku: "  ", qty: 1, priceCents: 249, wantErr: order.ErrEmptySKU},
		{name: "sku with tabulation", sku: "\t", qty: 1, priceCents: 249, wantErr: order.ErrEmptySKU},
	}
	for _, tc := range tests {
		t.Run(tc.name,
			func(t *testing.T) {
				item, err := order.NewItem(tc.sku, tc.qty, tc.priceCents)
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("NewItem() error: %v, wantErr: %v", err, tc.wantErr)
				}
				if err != nil {
					if !errors.Is(err, order.ErrValidation) {
						t.Errorf("Error %v does not wrap ErrValidation", err)
					}
					if item != (order.Item{}) {
						t.Errorf("NewItem() returned non-empty item %+v with an error", item)
					}
					return
				}
				if got := item.SKU(); got != tc.sku {
					t.Errorf("SKU() = %q, want %q", got, tc.sku)
				}
				if got := item.Qty(); got != tc.qty {
					t.Errorf("Qty() = %d, want %d", got, tc.qty)
				}
				if got := item.PriceCents(); got != tc.priceCents {
					t.Errorf("PriceCents() = %d, want %d", got, tc.priceCents)
				}
			})
	}
}

func mustItem(t *testing.T, sku string, qty int, priceCents int64) order.Item {
	t.Helper()
	item, err := order.NewItem(sku, qty, priceCents)
	if err != nil {
		t.Fatalf("mustItem(%q, %d, %d): %v", sku, qty, priceCents, err)
	}
	return item
}

func TestNewOrder(t *testing.T) {
	testTime := time.Date(2006, time.October, 8, 5, 12, 13, 0, time.UTC)
	tests := []struct {
		name       string
		customerID string
		items      []order.Item
		createdAt  time.Time
		wantTotal  int64
		wantErr    error
	}{
		{name: "valid", customerID: "abcefg", items: []order.Item{mustItem(t, "soundtrack", 1, 249)}, createdAt: testTime, wantTotal: 249},
		{name: "empty order", customerID: "abcefg", items: []order.Item{}, createdAt: testTime, wantTotal: 249, wantErr: order.ErrEmptyOrder},
		{name: "nil order", customerID: "abcefg", items: nil, createdAt: testTime, wantTotal: 249, wantErr: order.ErrEmptyOrder},
		{name: "empty item", customerID: "abcefg", items: []order.Item{mustItem(t, "soundtrack", 1, 249), {}}, createdAt: testTime, wantTotal: 249, wantErr: order.ErrInvalidQty},
		{name: "empty customerID", customerID: " ", items: []order.Item{mustItem(t, "soundtrack", 1, 249)}, createdAt: testTime, wantTotal: 249, wantErr: order.ErrEmptyCustomerID},
		{name: "tabulation in customerID", customerID: "\t", items: []order.Item{mustItem(t, "soundtrack", 1, 249)}, createdAt: testTime, wantTotal: 249, wantErr: order.ErrEmptyCustomerID},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ord, err := order.NewOrder(tc.customerID, tc.items, tc.createdAt)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("NewOrder() error: %v, want %v", err, tc.wantErr)
			}
			if err != nil {
				if !errors.Is(err, order.ErrValidation) {
					t.Errorf("Error %v does not wrap ErrValidation", err)
				}
				if ord != nil {
					t.Errorf("NewOrder() returned non-empty order %+v with an error", ord)
				}
				return
			}
			if ord.ID() == uuid.Nil {
				t.Errorf("ID() is nil, want non-nil id")
			}
			if got := ord.Status(); got != order.StatusCreated {
				t.Errorf("Status() = %v, want %v", got, order.StatusCreated)
			}
			if got := strings.TrimSpace(ord.CustomerID()); got != tc.customerID {
				t.Errorf("CustomerID() = %v, want %v", got, tc.customerID)
			}
			if got := ord.CreatedAt(); got.Compare(tc.createdAt) != 0 {
				t.Errorf("CreatedAt() = %v, want %v", got, tc.createdAt)
			}
			if got := ord.Items(); !slices.Equal(got, tc.items) {
				t.Errorf("Items() = %v, want %v", got, tc.items)
			}
			if got := ord.Total(); got != tc.wantTotal {
				t.Errorf("Total() = %d, want %d", got, tc.wantTotal)
			}
		})
	}
}
