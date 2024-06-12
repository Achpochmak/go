package module

import (
	"time"

	"HOMEWORK-1/internal/models"
	"HOMEWORK-1/internal/models/customErrors"
)

//Проверка параметров удаления заказа
func (m Module) validateDeleteOrder(Order models.Order) error {

	if time.Now().Before(Order.Storage_time) {
		return customErrors.ErrStorageTimeNotEnded
	}
	if Order.Delivered {
		return customErrors.ErrDelivered
	}
	return nil
}

//Проверка параметров доставки
func (m Module) validateDeliverOrder(order_ids []int, id_receiver int) ([]models.Order, error) {
	set := []models.Order{}
	for _, id := range order_ids {
		order, err := m.Storage.GetOrderByID(models.ID(id))
		if err != nil {
			return nil, customErrors.ErrOrderNotFound
		}

		if !time.Now().Before(order.Storage_time) {
			return nil, customErrors.ErrStorageTimeEnded
		}

		if order.Delivered {
			return nil, customErrors.ErrDelivered
		}

		if order.ID_receiver != models.ID(id_receiver) {
			return nil, customErrors.ErrWrongReceiver
		}

		set = append(set, order)
	}
	return set, nil
}

//Проверка параметров возврата
func (m Module) validateRefund(id int, id_receiver int) (models.Order, error) {
	order, err := m.Storage.GetOrderByID(models.ID(id))

	if err != nil {
		return models.Order{}, customErrors.ErrOrderNotFound
	}

	if !time.Now().Before(order.Delivered_time.Add(48 * time.Hour)) {
		return models.Order{}, customErrors.ErrRefundTimeEnded
	}

	if !order.Delivered {
		return models.Order{}, customErrors.ErrNotDelivered
	}

	if order.ID_receiver != models.ID(id_receiver) {
		return models.Order{}, customErrors.ErrWrongReceiver
	}
	return order, nil
}
