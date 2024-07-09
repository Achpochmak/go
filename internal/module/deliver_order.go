package module

import (
	"context"
	"time"

	"HOMEWORK-1/internal/models"
	"HOMEWORK-1/internal/models/customErrors"

	"github.com/pkg/errors"
)

// Доставка заказа
func (m Module) DeliverOrder(ctx context.Context, ordersID []int, idReceiver int) ([]*models.Order, error) {
	set, err := m.validateDeliverOrder(ctx, ordersID, idReceiver)

	if err != nil {
		return nil, errors.Wrap(err, "некорректные данные")
	}
	for _, order := range set {
		order.Delivered = true
		order.DeliveryTime = time.Now().UTC()
		if err := m.Repository.UpdateOrder(ctx, order); err != nil {
			return nil, errors.Wrap(err, "не получилось обновить заказ")
		}
	}
	return set, nil
}

// Проверка параметров доставки
func (m Module) validateDeliverOrder(ctx context.Context, ordersID []int, idReceiver int) ([]*models.Order, error) {
	set := []*models.Order{}
	for _, id := range ordersID {
		order, err := m.Repository.GetOrderByID(ctx, models.ID(id))
		if err != nil {
			return nil, err
		}

		if !time.Now().Before(order.StorageTime) {
			return nil, customErrors.ErrStorageTimeEnded
		}

		if order.Delivered {
			return nil, customErrors.ErrDelivered
		}

		if order.IDReceiver != models.ID(idReceiver) {
			return nil, customErrors.ErrWrongReceiver
		}

		set = append(set, &order)
	}
	return set, nil
}
