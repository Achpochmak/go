package cli

import (
	"HOMEWORK-1/internal/models"
	"context"
)

// Удалить заказ
func (c *CLI) deleteOrder(ctx context.Context, args []string) error {
	id, err := c.parseID(args)
	if err != nil {
		return err
	}

	order, err := c.Module.GetOrderByID(ctx, models.ID(id))

	if err != nil {
		return err
	}

	return c.Module.DeleteOrder(ctx, models.Order(order))
}
