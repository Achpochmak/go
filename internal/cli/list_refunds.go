package cli

import (
	"context"
	"flag"
	"fmt"

	"github.com/pkg/errors"
)

// Список возвратов
func (c *CLI) ListRefund(ctx context.Context, args []string) error {
	page, pageSize, err := c.parseListRefund(args)
	if err != nil {
		return errors.Wrap(err, "некорректный ввод")
	}

	list, err := c.Module.ListRefund(ctx, page, pageSize)
	if err != nil {
		return err
	}

	if len(list) < 1 {
		return errors.New("пустая страница")
	}

	for _, order := range list {
		fmt.Printf("ID заказа: %d\nID получателя: %d\nВремя хранения: %s\n", order.ID, order.IDReceiver, order.StorageTime)
	}

	return nil
}

// Парсинг параметров листа возвратов
func (c *CLI) parseListRefund(args []string) (int, int, error) {
	var page, pageSize int
	fs := flag.NewFlagSet(listRefund, flag.ContinueOnError)
	fs.IntVar(&page, "page", 0, "use --page=1")
	fs.IntVar(&pageSize, "pageSize", 0, "use --pageSize=1")

	if err := fs.Parse(args); err != nil {
		return 0, 0, err
	}
	return page, pageSize, nil
}
