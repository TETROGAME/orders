// Package order — доменная модель заказа.
//
// Здесь живут типы предметной области и правила, которые к ним относятся.
// Пакет ничего не знает про HTTP, JSON и базу данных: он не импортирует
// net/http и не содержит тегов сериализации. Это проверяется на приёмке.
package order

import (
	"errors"
	"github.com/google/uuid"
	// "time"
)

type Item struct {
	sku        string
	qty        int
	priceCents int
}

type Order struct {
	id         uuid.UUID
	customerID string
	items      []Item
	// created_at  time.Time
}

func (i *Item) PriceCents() int {
	return i.priceCents
}

func NewItem(sku string, qty, priceCents int) (*Item, error) {
	if qty <= 0 {
		return nil, errors.New("Negative quantity is not allowed")
	}
	if priceCents < 0 {
		return nil, errors.New("Negative price is not allowed")
	}
	if sku == "" {
		return nil, errors.New("Empty sku is not allowed")
	}
	return &Item{sku, qty, priceCents}, nil
}

func NewOrder(customerID string, items []Item) (*Order, error) {
	if customerID == "" {
		return nil, errors.New("Empty customerID is not allowed")
	}
	if len(items) == 0 {
		return nil, errors.New("Empty order is not allowed")
	}
	return &Order{uuid.New(), customerID, items}, nil
}
