package module

import (
	"HOMEWORK-1/internal/models"
	"context"
)

// Список возвратов
func (m Module) ListRefund(ctx context.Context, page int, pageSize int) ([]models.Order, error) {
	refunds, err := m.Repository.ListRefund(ctx, page, pageSize)
	return refunds, err
}
