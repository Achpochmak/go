package repository

import (
	"HOMEWORK-1/internal/models"
	"context"
)

func (r *Repository) UpdateOrder(ctx context.Context, order models.Order) error {
	db := r.QueryEngineProvider.GetQueryEngine(ctx)
	_, err := db.Query(ctx, "UPDATE orders SET  delivered = $1, deliveredAt=$2, refund=$3 WHERE id = $4", order.Delivered, order.DeliveryTime, order.Refund, order.ID)
	if err != nil {
		return err
	}

	return nil
}
