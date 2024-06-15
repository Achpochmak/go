package module

import (
	"HOMEWORK-1/internal/models"
	"HOMEWORK-1/internal/models/customErrors"
	"context"
	"time"
)

// Удалить заказ
func (m Module) DeleteOrder(ctx context.Context, order models.Order) error {
	if err := m.validateDeleteOrder(order); err != nil {
		return err
	}
	return m.Repository.DeleteOrder(ctx, order.ID)
}

// Проверка параметров удаления заказа
func (m Module) validateDeleteOrder(Order models.Order) error {

	if time.Now().Before(Order.StorageTime) {
		return customErrors.ErrStorageTimeNotEnded
	}
	if Order.Delivered {
		return customErrors.ErrDelivered
	}
	return nil
}
