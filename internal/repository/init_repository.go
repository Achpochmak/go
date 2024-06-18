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
	orderColumns = []string{"id", "idReceiver", "storageTime", "delivered", "refund", "createdAt", "deliveredAt", "weightKg", "price", "packaging"}
	packagingMap = map[string]models.Packaging{
		"bag":  models.PackagingBag,
		"box":  models.PackagingBox,
		"film": models.PackagingFilm,
		"None": models.PackagingNone,
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
		packaging = models.PackagingNone
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
