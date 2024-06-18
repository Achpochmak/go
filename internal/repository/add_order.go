package repository

import (
	"HOMEWORK-1/internal/models"
	"context"
)

func (r *Repository) AddOrder(ctx context.Context, order models.Order) error {
	db := r.QueryEngineProvider.GetQueryEngine(ctx)

	_, err := db.Query(ctx, "INSERT INTO orders (id, idReceiver, storageTime, delivered, refund, createdAt, deliveredAt, weightKg, price, packaging) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)", order.ID, order.IDReceiver, order.StorageTime, order.Delivered, order.Refund, order.CreatedAt, order.DeliveryTime, order.WeightKg, order.Price, order.Packaging.Name)
	if err != nil {
		return err
	}
	return nil
}
