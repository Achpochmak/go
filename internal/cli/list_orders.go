package cli

import (
	"context"
	"fmt"
)

// Список заказов
func (c *CLI) listOrder(ctx context.Context) error {
	list, err := c.Module.ListOrder(ctx)
	if err != nil {
		return err
	}

	for _, order := range list {
		fmt.Printf("ID заказа: %d\nID получателя: %d\nВремя хранения: %s\n", order.ID, order.IDReceiver, order.StorageTime)
	}
	return nil
}
