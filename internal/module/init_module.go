package module

import (
	"HOMEWORK-1/internal/models"
	"context"
	"time"
)
const refundTime = 48 * time.Hour

type Repository interface {
	AddOrder(context.Context, models.Order) error
	ListOrder(context.Context) ([]models.Order, error)
	GetOrderByID(context.Context,  models.ID) (models.Order, error)
	UpdateOrder(context.Context,  models.Order) error
	DeleteOrder(context.Context,  models.ID) error
	ListRefund(context.Context, int, int)([]models.Order, error)
	GetOrdersByCustomer(context.Context,  models.ID, int)([]models.Order, error)
}

type Deps struct {
	Repository Repository
}

type Module struct {
	Deps
}

// NewModule
func NewModule(d Deps) Module {
	return Module{Deps: d}
}