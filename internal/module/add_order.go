package module

import (
	"context"
	"time"

	"HOMEWORK-1/internal/models"
	"HOMEWORK-1/internal/models/customErrors"

	"github.com/pkg/errors"
)

// Добавить заказ
func (m Module) AddOrder(ctx context.Context, Order models.Order) error {

	err := m.checkOrder(ctx, Order.ID, Order.StorageTime)
	if err != nil {
		return errors.Wrap(err, "некорректные данные")
	}
	
	Order.Price, err = m.packOrder(Order.WeightKg, Order.Price, Order.Packaging)
	if err != nil {
		return errors.Wrap(err, "не получилось упаковать заказ")
	}

	if err := m.Repository.AddOrder(ctx, Order); err != nil {
		return errors.Wrap(err, "не получилось добавить заказ")
	}
	return nil
}

// Упаковать заказ
func (m Module) packOrder(weight float64, price float64, packaging models.Packaging) (float64, error) {
	if packaging.MaxWeight > 0 && weight > packaging.MaxWeight {
		return 0, customErrors.ErrWeightIsTooBig
	}
	totalPrice := price + packaging.Price
	return totalPrice, nil
}

func (m Module) checkOrder(ctx context.Context, id models.ID, st time.Time) error {
	_, err := m.Repository.GetOrderByID(ctx, id)
	if err == nil {
		return customErrors.ErrOrderAlreadyExists
	}
	if time.Now().After(st) {
		return customErrors.ErrStorageTimeEnded
	}
	return nil
}
