package module

import (
	"context"

	"HOMEWORK-1/internal/models"
)

// Поиск заказов по получателю
func (m Module) GetOrdersByCustomer(ctx context.Context, idReceiver int, amount int) ([]models.Order, error) {

	set, err := m.Repository.GetOrdersByCustomer(ctx, models.ID(idReceiver), amount)
	if err != nil {
		return nil, err
	}

	return set, nil
}
