package module

import (
	"context"
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
func (m Module) AddOrder(ctx context.Context, Order models.Order) error {
	return m.Repository.AddOrder(ctx, Order)
}

// Список заказов
func (m Module) ListOrder(ctx context.Context) ([]models.Order, error) {
	return m.Repository.ListOrder(ctx)
}

// Удалить заказ
func (m Module) DeleteOrder(ctx context.Context, order models.Order) error {
	if err := m.validateDeleteOrder(order); err != nil {
		return err
	}
	return m.Repository.DeleteOrder(ctx, order.ID)
}

// Доставка заказа
func (m Module) DeliverOrder(ctx context.Context, order_ids []int, id_receiver int) ([]models.Order, error) {
	set, err := m.validateDeliverOrder(ctx, order_ids, id_receiver)

	if err != nil {
		return nil, err
	}
	for _, order := range set {
		order.Delivered = true
		order.Delivered_time = time.Now()
		if err := m.Repository.UpdateOrder(ctx, order); err != nil {
			return nil, customErrors.ErrNotUpdated
		}
	}
	return set, nil
}

// Поиск заказов по получателю
func (m Module) GetOrdersByCustomer(ctx context.Context, id_receiver int, amount int) ([]models.Order, error) {

	set, err := m.Repository.GetOrdersByCustomer(ctx, models.ID(id_receiver))
	if err != nil {
		return nil, err
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
func (m Module) GetOrderByID(ctx context.Context, n models.ID) (models.Order, error) {
	return m.Repository.GetOrderByID(ctx, n)
}

// Возврат заказа
func (m Module) Refund(ctx context.Context, id int, id_receiver int) error {
	order, err := m.validateRefund(ctx, id, id_receiver)
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

// Список возвратов
func (m Module) ListRefund(ctx context.Context) ([]models.Order, error) {
	refunds, err := m.Repository.ListRefund(ctx)
	return refunds, err

}
