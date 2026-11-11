// Package order — доменная модель заказа.
//
// Здесь живут типы предметной области и правила, которые к ним относятся.
// Пакет ничего не знает про HTTP, JSON и базу данных: он не импортирует
// net/http и не содержит тегов сериализации. Это проверяется на приёмке.
package order

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Item struct {
	sku        string
	qty        int
	priceCents int
}

func (i *Item) PriceCents() int {
	return i.priceCents
}

func (i *Item) Qty() int {
	return i.qty
}

func (i *Item) Sku() string {
	return i.sku
}

type Order struct {
	id         uuid.UUID
	customerID string
	items      []Item
	status     string
	created_at time.Time
}

func (o *Order) Id() uuid.UUID {
	return o.id
}

func (o *Order) Items() []Item {
	return o.items
}

func (o *Order) CustomerID() string {
	return o.customerID
}

func (o *Order) Status() string {
	return o.status
}

func (o *Order) CreatedAt() time.Time {
	return o.created_at
}

func (o *Order) Total() int64 {
	var total int64
	for _, item := range o.items {
		total += int64(item.qty) * int64(item.priceCents)
	}
	return total
}

var (
	ErrValidation      = errors.New("validation failed")
	ErrInvalidQty      = fmt.Errorf("%w: negative or zero quantity", ErrValidation)
	ErrNegativePrice   = fmt.Errorf("%w: negative price", ErrValidation)
	ErrEmptySKU        = fmt.Errorf("%w: empty sku", ErrValidation)
	ErrEmptyCustomerID = fmt.Errorf("%w: empty customer id", ErrValidation)
	ErrEmptyOrder      = fmt.Errorf("%w: empty order", ErrValidation)
)

func NewItem(sku string, qty, priceCents int) (*Item, error) {
	if qty <= 0 {
		return nil, ErrInvalidQty
	}
	if priceCents < 0 {
		return nil, ErrNegativePrice
	}
	if sku == "" {
		return nil, ErrEmptySKU
	}
	return &Item{sku, qty, priceCents}, nil
}

func NewOrder(customerID string, items []Item, created_at time.Time) (*Order, error) {
	if customerID == "" {
		return nil, ErrEmptyCustomerID
	}
	if len(items) == 0 {
		return nil, ErrEmptyOrder
	}
	return &Order{uuid.New(), customerID, items, "created", created_at}, nil
}
