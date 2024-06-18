package module

import (
	"HOMEWORK-1/internal/models"
	"HOMEWORK-1/internal/models/customErrors"
	"context"
	"time"

	"github.com/pkg/errors"
)

// Удалить заказ
func (m Module) DeleteOrder(ctx context.Context, id models.ID) error {
	if err := m.validateDeleteOrder(ctx, id); err != nil {
		return errors.Wrap(err, "некорректные данные")
	}
	if err := m.Repository.DeleteOrder(ctx, id); err != nil {
		return errors.Wrap(err, "не получилось удалить заказ")
	}
	return nil
}

// Проверка параметров удаления заказа
func (m Module) validateDeleteOrder(ctx context.Context, id models.ID) error {
	order, err := m.Repository.GetOrderByID(ctx, id)
	if err != nil {
		return err
	}
	if time.Now().Before(order.StorageTime) {
		return customErrors.ErrStorageTimeNotEnded
	}
	if order.Delivered {
		return customErrors.ErrDelivered
	}
	return nil
}
