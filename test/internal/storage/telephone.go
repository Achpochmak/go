package storage

import (
	"time"

	"HOMEWORK-1/internal/models"
)

type orderRecord struct {
	Id      int    `json:"id"`
	Id_receiver int    `json:"receiver"`
	Storage_time time.Time    `json:"storage_time"`
	Delivered bool `json:"delivered"`
	Refund bool `json:"refund"`
	CreatedAt time.Time `json:"created_at"`
	DeliveredAt time.Time `json:"delivered_at"`
	Hash string `json:"hash"`
}

func (t orderRecord) toDomain() models.Order {
	return models.Order{
		Id:      models.Id(t.Id),
		Id_receiver: models.Id(t.Id_receiver),
		Storage_time: (t.Storage_time),
		Delivered: (t.Delivered),
		Created_at: (t.CreatedAt),
		Refund: (t.Refund),
		Delivered_time: (t.DeliveredAt),
		Hash: (t.Hash),
	}
}

func transform(order models.Order) orderRecord {
	return orderRecord{
		Id:      int(order.Id),
		Id_receiver: int(order.Id_receiver),
		Storage_time: (order.Storage_time),
		Delivered: bool(order.Delivered),
		Refund: bool(order.Refund),
		CreatedAt: (order.Created_at),
		DeliveredAt: (order.Delivered_time),
		Hash: (order.Hash),
	}
}

