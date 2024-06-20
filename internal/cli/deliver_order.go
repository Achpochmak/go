package cli

import (
	"HOMEWORK-1/internal/models/customErrors"
	"context"
	"flag"
	"fmt"
	"strconv"
	"strings"
)

// Доставить заказ
func (c *CLI) deliverOrder(ctx context.Context, args []string) error {
	orderIDs, idReceiver, err := c.parseDeliverOrder(args)
	if err != nil {
		return err
	}

	orders, err := c.Module.DeliverOrder(ctx, orderIDs, idReceiver)
	if err != nil {
		return err
	}

	for _, order := range orders {
		fmt.Printf("ID заказа: %d\nID получателя: %d\nВремя хранения: %s\n", order.ID, order.IDReceiver, order.StorageTime)
	}
	return nil
}

// Парсинг параметров доставки заказа
func (c *CLI) parseDeliverOrder(args []string) ([]int, int, error) {
	var idReceiver int
	var orderIdsUnparsed string
	fs := flag.NewFlagSet(deliverOrder, flag.ContinueOnError)
	fs.IntVar(&idReceiver, "idReceiver", 0, "use --id=1")
	fs.StringVar(&orderIdsUnparsed, "id", "", "use --id=1,2,3")

	if err := fs.Parse(args); err != nil {
		return nil, 0, err
	}

	if idReceiver == 0 {
		return nil, 0, customErrors.ErrReceiverNotFound
	}

	orderIDs := strings.Split(orderIdsUnparsed, ",")
	var ids []int
	for _, numStr := range orderIDs {
		id, err := strconv.Atoi(numStr)
		if err != nil {
			return nil, 0, customErrors.ErrIDNotFound
		}
		ids = append(ids, id)
	}
	return ids, idReceiver, nil
}
