// Package order — хранилище заказов.
//
// В Модуле 0 это карта в памяти. В Модуле 5 на её место встанет Postgres,
// и если границы проведены правильно, менять придётся только этот пакет.
package order

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"

	"github.com/gulmix/orders/internal/models/order"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("order already exists")
	ErrNilOrder      = errors.New("cannot save nil order")
)

type Repository struct {
	mu     sync.RWMutex
	orders map[uuid.UUID]*order.Order
}

func New() *Repository {
	return &Repository{
		orders: make(map[uuid.UUID]*order.Order),
	}
}

func (r *Repository) Get(ctx context.Context, targetID uuid.UUID) (*order.Order, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	ord, ok := r.orders[targetID]
	if !ok {
		return nil, ErrNotFound
	}
	return ord, nil
}

func (r *Repository) Save(ctx context.Context, orderToSave *order.Order) error {
	if orderToSave == nil {
		return ErrNilOrder
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	orderID := orderToSave.ID()
	if _, exists := r.orders[orderID]; exists {
		return fmt.Errorf("%w: UUID %v is already present in the database", ErrAlreadyExists, orderID)
	}
	r.orders[orderID] = orderToSave
	return nil
}
