package cli

import (
	"context"
	"flag"
	"fmt"
	
	"HOMEWORK-1/internal/models/customErrors"

	"github.com/pkg/errors"
)

// Получить список заказов по получателю
func (c *CLI) GetOrdersByCustomer(ctx context.Context, args []string) error {
	idReceiver, amount, err := c.parseGetOrdersByCustomer(args)
	if err != nil {
		return errors.Wrap(err, "некорректный ввод")
	}

	list, err := c.Module.GetOrdersByCustomer(ctx, idReceiver, amount)
	if err != nil {
		return err
	}

	for _, order := range list {
		fmt.Printf("ID заказа: %d\nID получателя: %d\nВремя хранения: %s\n", order.ID, order.IDReceiver, order.StorageTime)
	}
	return nil
}

// Парсинг получения параметров получения заказа по клиенту
func (c *CLI) parseGetOrdersByCustomer(args []string) (int, int, error) {
	var idReceiver, amount int
	fs := flag.NewFlagSet(getOrdersByCustomer, flag.ContinueOnError)
	fs.IntVar(&idReceiver, "idReceiver", 0, "use --idReceiver=1")
	fs.IntVar(&amount, "n", 0, "use --n=1")

	if err := fs.Parse(args); err != nil {
		return 0, 0, err
	}

	if idReceiver == 0 {
		return 0, 0, customErrors.ErrReceiverNotFound
	}
	return idReceiver, amount, nil
}
