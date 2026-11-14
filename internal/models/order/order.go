// Package order — доменная модель заказа.
//
// Здесь живут типы предметной области и правила, которые к ним относятся.
// Пакет ничего не знает про HTTP, JSON и базу данных: он не импортирует
// net/http и не содержит тегов сериализации. Это проверяется на приёмке.
package order

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Status string

const StatusCreated Status = "created"

const (
	MaxQty        int   = 10000
	MaxPriceCents int64 = 100 * 5 * 1_000_000
	MaxItems      int   = 1_000_000
)

var (
	ErrValidation      = errors.New("validation failed")
	ErrInvalidQty      = fmt.Errorf("%w: negative or zero quantity", ErrValidation)
	ErrQtyTooLarge     = fmt.Errorf("%w: quantity too large, max %d", ErrValidation, MaxQty)
	ErrNegativePrice   = fmt.Errorf("%w: negative price", ErrValidation)
	ErrPriceTooLarge   = fmt.Errorf("%w: price too large, max %d cents", ErrValidation, MaxPriceCents)
	ErrEmptySKU        = fmt.Errorf("%w: empty sku", ErrValidation)
	ErrEmptyCustomerID = fmt.Errorf("%w: empty customer id", ErrValidation)
	ErrEmptyOrder      = fmt.Errorf("%w: empty order", ErrValidation)
	ErrTooManyItems    = fmt.Errorf("%w: too many items in the order, max %d", ErrValidation, MaxItems)
)

type Item struct {
	sku        string
	qty        int
	priceCents int64
}

func (i Item) PriceCents() int64 {
	return i.priceCents
}

func (i Item) Qty() int {
	return i.qty
}

func (i Item) SKU() string {
	return i.sku
}

func (i Item) validate() error {
	if i.qty <= 0 {
		return ErrInvalidQty
	}
	if i.qty > MaxQty {
		return ErrQtyTooLarge
	}
	if i.priceCents < 0 {
		return ErrNegativePrice
	}
	if i.priceCents > MaxPriceCents {
		return ErrPriceTooLarge
	}
	if strings.TrimSpace(i.sku) == "" {
		return ErrEmptySKU
	}
	return nil
}

type Order struct {
	id         uuid.UUID
	customerID string
	items      []Item
	status     Status
	createdAt  time.Time
}

func (o *Order) ID() uuid.UUID {
	return o.id
}

func (o *Order) Items() []Item {
	return slices.Clone(o.items)
}

func (o *Order) CustomerID() string {
	return o.customerID
}

func (o *Order) Status() Status {
	return o.status
}

func (o *Order) CreatedAt() time.Time {
	return o.createdAt
}

func (o *Order) Total() int64 {
	var total int64
	for _, item := range o.items {
		total += int64(item.qty) * item.priceCents
	}
	return total
}

func NewItem(sku string, qty int, priceCents int64) (Item, error) {
	item := Item{sku, qty, priceCents}
	if err := item.validate(); err != nil {
		return Item{}, err
	}
	return item, nil
}

func NewOrder(customerID string, items []Item, createdAt time.Time) (*Order, error) {
	if strings.TrimSpace(customerID) == "" {
		return nil, ErrEmptyCustomerID
	}
	if len(items) > MaxItems {
		return nil, ErrTooManyItems
	}
	for i, item := range items {
		if err := item.validate(); err != nil {
			return nil, fmt.Errorf("item %d: %w", i+1, err)
		}
	}
	if len(items) == 0 {
		return nil, ErrEmptyOrder
	}
	return &Order{uuid.New(), customerID, slices.Clone(items), StatusCreated, createdAt}, nil
}
