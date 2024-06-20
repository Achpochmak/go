package module

import (
	"context"
	"time"

	"HOMEWORK-1/internal/models"
	"HOMEWORK-1/internal/repository/transactor"
)

const refundTime = 48 * time.Hour

type Repository interface {
	AddOrder(context.Context, models.Order) error
	ListOrder(context.Context) ([]models.Order, error)
	GetOrderByID(context.Context, models.ID) (models.Order, error)
	UpdateOrder(context.Context, models.Order) error
	DeleteOrder(context.Context, models.ID) error
	ListRefund(context.Context, int, int) ([]models.Order, error)
	GetOrdersByCustomer(context.Context, models.ID, int) ([]models.Order, error)
}

type Deps struct {
	Repository Repository
	Transactor *transactor.TransactionManager
}

type Module struct {
	Deps
}

// NewModule
func NewModule(d Deps) Module {
	return Module{Deps: d}
}
