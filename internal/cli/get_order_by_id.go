package cli

import (
	"context"
	"fmt"

	"HOMEWORK-1/internal/models"

	"github.com/pkg/errors"
)

// Найти заказ
func (c *CLI) GetOrderByID(ctx context.Context, args []string) error {
	id, err := c.parseID(args)
	if err != nil {
		return errors.Wrap(err, "некорректный ввод")
	}

	order, err := c.Module.GetOrderByID(ctx, models.ID(id))
	if err != nil {
		return err
	}

	fmt.Printf("ID заказа: %d\nID получателя: %d\nВремя хранения: %s\n", order.ID, order.IDReceiver, order.StorageTime)
	return nil
}
