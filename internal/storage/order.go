package storage

import (
	"time"

	"HOMEWORK-1/internal/models"
)

type orderRecord struct {
	ID           int       `json:"id"`
	ID_receiver  int       `json:"receiver"`
	Storage_time time.Time `json:"storage_time"`
	Delivered    bool      `json:"delivered"`
	Refund       bool      `json:"refund"`
	CreatedAt    time.Time `json:"created_at"`
	DeliveredAt  time.Time `json:"delivered_at"`
	Hash         string    `json:"hash"`
}

//Преобразование из записи заказа в заказ
func (t orderRecord) toDomain() models.Order {
	return models.Order{
		ID:             models.ID(t.ID),
		ID_receiver:    models.ID(t.ID_receiver),
		Storage_time:   (t.Storage_time),
		Delivered:      (t.Delivered),
		Created_at:     (t.CreatedAt),
		Refund:         (t.Refund),
		Delivered_time: (t.DeliveredAt),
		Hash:           (t.Hash),
	}
}

//Преобразование из заказа в запись заказа
func transform(order models.Order) orderRecord {
	return orderRecord{
		ID:           int(order.ID),
		ID_receiver:  int(order.ID_receiver),
		Storage_time: (order.Storage_time),
		Delivered:    bool(order.Delivered),
		Refund:       bool(order.Refund),
		CreatedAt:    (order.Created_at),
		DeliveredAt:  (order.Delivered_time),
		Hash:         (order.Hash),
	}
}
