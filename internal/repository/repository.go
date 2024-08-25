package repository

import (
	"context"
	"wb_l0/internal/models"
)

type Repo interface {
	AddOrder(ctx context.Context, order models.Order) error
	GetOrderByID(ctx context.Context, orderUID string) (models.Order, error)
	AllOrders(ctx context.Context) ([]models.Order, error)
}
