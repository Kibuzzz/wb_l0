package inmemory

import (
	"context"
	"fmt"
	"sync"
	"wb_l0/internal/models"
	"wb_l0/internal/repository"
)

type InMemory struct {
	store map[string]models.Order
	db    repository.Repo
	mu    *sync.Mutex
}

func New(ctx context.Context, db repository.Repo) *InMemory {
	inMemory := &InMemory{
		store: make(map[string]models.Order),
		db:    db,
		mu:    &sync.Mutex{},
	}
	err := inMemory.loadDB(ctx)
	if err != nil {
		panic(err)
	}
	return inMemory
}

func (i *InMemory) AddOrder(ctx context.Context, order models.Order) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.store[order.OrderUID] = order
	return nil

}
func (i *InMemory) GetOrderByID(ctx context.Context, orderUID string) (models.Order, error) {
	i.mu.Lock()
	defer i.mu.Unlock()
	order, exist := i.store[orderUID]
	if !exist {
		return models.Order{}, repository.ErrorNotFound
	}
	return order, nil
}

func (i *InMemory) AllOrders(ctx context.Context) ([]models.Order, error) {
	orders := []models.Order{}
	for _, order := range i.store {
		orders = append(orders, order)
	}
	return orders, nil
}

func (i *InMemory) loadDB(ctx context.Context) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	orders, err := i.db.AllOrders(ctx)
	if err != nil {
		return fmt.Errorf("failed to load orders from db: %v", err)
	}
	for _, order := range orders {
		i.store[order.OrderUID] = order
	}
	return nil
}
