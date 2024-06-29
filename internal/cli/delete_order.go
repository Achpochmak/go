package cli

import (
	"context"

	"HOMEWORK-1/internal/models"

	"github.com/pkg/errors"
)

// Удалить заказ
func (c *CLI) deleteOrder(ctx context.Context, args []string) error {
	id, err := c.parseID(args)
	if err != nil {
		return errors.Wrap(err, "некорректный ввод")
	}

	return c.Module.DeleteOrder(ctx, models.ID(id))
}
