package module

import (
	"context"

	"HOMEWORK-1/internal/models"

	"github.com/pkg/errors"
)

// Поиск заказов по получателю
func (m Module) GetOrdersByCustomer(ctx context.Context, idReceiver int, amount int) ([]models.Order, error) {

	set, err := m.Repository.GetOrdersByCustomer(ctx, models.ID(idReceiver), amount)
	if err != nil {
		return nil, errors.Wrap(err, "ошибка поиска")
	}

	return set, nil
}
