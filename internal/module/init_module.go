package module

import (
	"HOMEWORK-1/internal/models"
	"context"
)

type Repository interface {
	AddOrder(context.Context, models.Order) error
	ListOrder(context.Context) ([]models.Order, error)
	GetOrderByID(context.Context,  models.ID) (models.Order, error)
	UpdateOrder(context.Context,  models.Order) error
	DeleteOrder(context.Context,  models.ID) error
	ListRefund(context.Context)([]models.Order, error)
	GetOrdersByCustomer(context.Context,  models.ID)([]models.Order, error)
}

type Deps struct {
	Repository Repository
}

type Module struct {
	Deps
}
