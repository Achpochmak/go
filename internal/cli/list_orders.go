package cli

import (
	"context"
	"fmt"
	"time"
)

// Список заказов
func (c *CLI) listOrder(ctx context.Context) error {
	list, err := c.Module.ListOrder(ctx)
	if err != nil {
		return err
	}
	time.Sleep(time.Second * 5) // имитируем долгую работу

	for _, order := range list {
		fmt.Printf("ID заказа: %d\nID получателя: %d\nВремя хранения: %s\nВес: %.2f кг\nЦена: %.2f руб\n\n",
			order.ID, order.IDReceiver, order.StorageTime, order.WeightKg, order.Price)
	}
	return nil
}
