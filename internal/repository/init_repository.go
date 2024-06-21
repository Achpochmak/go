package repository

import (
	"HOMEWORK-1/internal/models"
	"HOMEWORK-1/internal/repository/schema"
	"HOMEWORK-1/internal/repository/transactor"
)

const (
	orderTable = "orders"
)

var (
	orderColumns = []string{"id", "id_receiver", "storage_time", "delivered", "refund", "created_at", "delivered_at", "weight_kg", "price", "packaging"}
	packagingMap = map[string]models.Packaging{
		"bag":  models.NewBag(),
		"box":  models.NewBox(),
		"film": models.NewFilm(),
		"none": models.NewNoPackaging(),
	}
)

type Repository struct {
	transactor.QueryEngineProvider
}

func NewRepository(provider transactor.QueryEngineProvider) *Repository {
	return &Repository{provider}
}

func toDomainOrder(order schema.OrderInfo) models.Order {
	packaging, ok := packagingMap[order.Packaging]
	if !ok {
		packaging = models.NewNoPackaging()
	}

	return models.Order{
		ID:           models.ID(order.ID),
		IDReceiver:   models.ID(order.IDReceiver),
		Delivered:    order.Delivered,
		DeliveryTime: order.DeliveredAt.Time,
		Refund:       order.Refund,
		CreatedAt:    order.CreatedAt.Time,
		StorageTime:  order.StorageTime.Time,
		WeightKg:     order.WeightKg,
		Price:        order.Price,
		Packaging:    packaging,
	}
}
