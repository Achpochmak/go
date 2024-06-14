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
	orderColumns = []string{"id", "id_receiver", "storage_time", "delivered", "refund", "created_at", "delivered_at"}
)

type Repository struct {
	transactor.QueryEngineProvider
}

func toDomainOrder(order schema.OrderInfo) models.Order {
	return models.Order{
		ID:             models.ID(order.ID),
		ID_receiver:    models.ID(order.ID_receiver),
		Delivered:      order.Delivered,
		Delivered_time: order.DeliveredAt.Time,
		Refund:         order.Refund,
		Created_at:     order.CreatedAt.Time,
		Storage_time:   order.Storage_time.Time,
	}
}
