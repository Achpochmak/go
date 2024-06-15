package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
)

// Список возвратов
func (c *CLI) listRefund(ctx context.Context, args []string) error {
	page, pageSize, err := c.parseListRefund(args)
	if err != nil {
		return err
	}

	list, err := c.Module.ListRefund(ctx, page, pageSize)
	if err != nil {
		return err
	}

	if page == 0 || pageSize == 0 {
		for _, order := range list {
			fmt.Printf("ID заказа: %d\nID получателя: %d\nВремя хранения: %s\n", order.ID, order.IDReceiver, order.StorageTime)
		}
		return nil
	}

	start := (page - 1) * pageSize
	end := start + pageSize

	if end > len(list) || start < 0 {
		return errors.New("пустая страница")
	}

	for _, order := range list[start:end] {
		fmt.Printf("ID заказа: %d\nID получателя: %d\nВремя хранения: %s\n", order.ID, order.IDReceiver, order.StorageTime)
	}
	return nil
}

// Парсинг параметров листа возвратов
func (c *CLI) parseListRefund(args []string) (int, int, error) {
	var page, pageSize int
	fs := flag.NewFlagSet(listRefund, flag.ContinueOnError)
	fs.IntVar(&page, "page", 0, "use --idReceiver=1")
	fs.IntVar(&pageSize, "pageSize", 0, "use --id=1")

	if err := fs.Parse(args); err != nil {
		return 0, 0, err
	}
	return page, pageSize, nil
}
