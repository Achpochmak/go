package module

import (
	"sort"
	"time"

	"HOMEWORK-1/internal/models"
	"HOMEWORK-1/internal/models/customErrors"
	"HOMEWORK-1/pkg/hash"
)

// NewModule
func NewModule(d Deps) Module {
	return Module{Deps: d}
}

// Добавить заказ
func (m Module) AddOrder(Order models.Order) error {
	return m.Storage.AddOrder(Order)
}

// Список заказов
func (m Module) ListOrder() ([]models.Order, error) {
	return m.Storage.ListOrder()
}

// Удалить заказ
func (m Module) DeleteOrder(Order models.Order) error {
	if err := m.validateDeleteOrder(Order); err != nil {
		return err
	}

	orders, err := m.Storage.ListOrder()
	if err != nil {
		return err
	}

	set := make(map[models.ID]models.Order, len(orders))
	for _, order := range orders {
		set[order.ID] = order
	}

	_, ok := set[Order.ID]
	if !ok {
		return nil
	}

	delete(set, Order.ID)

	newOrders := make([]models.Order, 0, len(set))
	for _, value := range set {
		newOrders = append(newOrders, value)
	}
	return m.Storage.ReWrite(newOrders)
}

// Доставка заказа
func (m Module) DeliverOrder(order_ids []int, id_receiver int) ([]models.Order, error) {
	set, err := m.validateDeliverOrder(order_ids, id_receiver)
	if err != nil {
		return nil, err
	}

	for _, order := range set {
		order.Delivered = true
		order.Hash = hash.GenerateHash()
		order.Delivered_time = time.Now()
		if err := m.Storage.UpdateOrder(order); err != nil {
			return nil, customErrors.ErrNotUpdated
		}
	}
	return set, nil
}

// Поиск заказов по получателю
func (m Module) GetOrdersByCustomer(id_receiver int, amount int) ([]models.Order, error) {
	orders, err := m.Storage.ListOrder()
	if err != nil {
		return nil, err
	}

	set := []models.Order{}
	for _, order := range orders {
		if order.ID_receiver == models.ID(id_receiver) {
			set = append(set, order)
		}
	}

	sort.Slice(set, func(i, j int) bool {
		return set[j].Created_at.Before(set[i].Created_at)
	})

	if amount > 0 {
		return set[0:amount], nil
	}
	return set, nil
}

// Поиск заказа
func (m Module) GetOrderByID(n models.ID) (models.Order, error) {
	return m.Storage.GetOrderByID(n)
}

// Возврат заказа
func (m Module) Refund(id int, id_receiver int) error {
	order, err := m.validateRefund(id, id_receiver)
	if err != nil {
		return err
	}

	order.Delivered = false
	order.Refund = true
	order.Hash = hash.GenerateHash()

	if err := m.Storage.UpdateOrder(order); err != nil {
		return customErrors.ErrNotUpdated
	}

	return nil
}

// Список возвратов
func (m Module) ListRefund() ([]models.Order, error) {
	orders, err := m.Storage.ListOrder()
	refunds := []models.Order{}
	for _, order := range orders {
		if order.Refund {
			refunds = append(refunds, order)
		}
	}
	return refunds, err

}
