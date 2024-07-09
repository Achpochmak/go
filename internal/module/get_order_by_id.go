package module

import (
	"HOMEWORK-1/internal/models"
	"context"
)

// Поиск заказа
func (m Module) GetOrderByID(ctx context.Context, n models.ID) (models.Order, error) {
	return m.Repository.GetOrderByID(ctx, n)
}
