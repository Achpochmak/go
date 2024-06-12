package module

import "HOMEWORK-1/internal/models"

type Storage interface {
	AddOrder(Order models.Order) error
	ListOrder() ([]models.Order, error)
	ReWrite(Orders []models.Order) error
	GetOrderByID(ID models.ID) (models.Order, error)
	UpdateOrder(Order models.Order) error
}

type Deps struct {
	Storage Storage
}

type Module struct {
	Deps
}
