package module

import (
	"HOMEWORK-1/internal/models"
	"HOMEWORK-1/internal/models/customErrors"
	"HOMEWORK-1/pkg/hash"
	"context"
	"time"
)

// Возврат заказа
func (m Module) Refund(ctx context.Context, id int, idReceiver int) error {
	order, err := m.validateRefund(ctx, id, idReceiver)
	if err != nil {
		return err
	}

	order.Delivered = false
	order.Refund = true
	order.Hash = hash.GenerateHash()

	if err := m.Repository.UpdateOrder(ctx, order); err != nil {
		return customErrors.ErrNotUpdated
	}

	return nil
}

// Проверка параметров возврата
func (m Module) validateRefund(ctx context.Context, id int, idReceiver int) (models.Order, error) {
	order, err := m.Repository.GetOrderByID(ctx, models.ID(id))

	if err != nil {
		return models.Order{}, customErrors.ErrOrderNotFound
	}

	if !time.Now().Before(order.DeliveryTime.Add(refundTime)) {
		return models.Order{}, customErrors.ErrRefundTimeEnded
	}

	if !order.Delivered {
		return models.Order{}, customErrors.ErrNotDelivered
	}

	if order.IDReceiver != models.ID(idReceiver) {
		return models.Order{}, customErrors.ErrWrongReceiver
	}
	return order, nil
}
