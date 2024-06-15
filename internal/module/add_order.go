package module

import (
	"HOMEWORK-1/internal/models"
	"context"
)

// Добавить заказ
func (m Module) AddOrder(ctx context.Context, Order models.Order) error {
	return m.Repository.AddOrder(ctx, Order)
}
