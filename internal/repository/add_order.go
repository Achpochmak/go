package repository

import (
	"context"
	"HOMEWORK-1/internal/models"
)

func (r *Repository) AddOrder(ctx context.Context, order models.Order) error {
	db := r.QueryEngineProvider.GetQueryEngine(ctx)
	
	_, err := db.Query(ctx, "INSERT INTO orders (id, idReceiver, storageTime, delivered, refund, createdAt, deliveredAt) VALUES($1, $2, $3, $4, $5, $6, $7)", order.ID, order.IDReceiver, order.StorageTime, order.Delivered, order.Refund, order.CreatedAt, order.DeliveryTime)
	if err != nil {
		return err
	}

	return nil
}
