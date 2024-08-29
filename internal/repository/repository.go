package repository

import (
	"context"
	"errors"
	"wb_l0/internal/models"
)

var (
	ErrorNotFound = errors.New("order not found")
)

type Repo interface {
	AddOrder(ctx context.Context, order models.Order) error
	GetOrderByID(ctx context.Context, orderUID string) (models.Order, error)
	AllOrders(ctx context.Context) ([]models.Order, error)
}
