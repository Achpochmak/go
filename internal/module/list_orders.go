package module

import (
	"context"
	"HOMEWORK-1/internal/models"
)

// Список заказов
func (m Module) ListOrder(ctx context.Context) ([]models.Order, error) {
	return m.Repository.ListOrder(ctx)
}
