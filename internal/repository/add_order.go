package repository

import (
	"context"
	"HOMEWORK-1/internal/models"
)

func (r *Repository) AddOrder(ctx context.Context, order models.Order) error {
	db := r.QueryEngineProvider.GetQueryEngine(ctx)
	_, err := db.Exec(ctx, "INSERT INTO orders (id, id_receiver, storage_time, delivered, refund, created_at, delivered_at, weight_kg, price, packaging) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
		order.ID, order.IDReceiver, order.StorageTime, order.Delivered, order.Refund, order.CreatedAt, order.DeliveryTime, order.WeightKg, order.Price, order.Packaging.GetName())
	if err != nil {
		return err
	}
	return nil
}
