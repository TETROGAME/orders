package order_test

import (
	"errors"
	"testing"

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
