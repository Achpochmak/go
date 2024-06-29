package cli

import (
	"context"
	"flag"
	
	"HOMEWORK-1/internal/models/customErrors"

	"github.com/pkg/errors"
)

// Вернуть заказ
func (c *CLI) Refund(ctx context.Context, args []string) error {
	id, idReceiver, err := c.parseRefund(args)
	if err != nil {
		return errors.Wrap(err, "некорректный ввод")
	}

	err = c.Module.Refund(ctx, id, idReceiver)
	if err != nil {
		return err
	}

	return nil
}

// Парсинг параметров возврата
func (c *CLI) parseRefund(args []string) (int, int, error) {
	var idReceiver, id int
	fs := flag.NewFlagSet(refund, flag.ContinueOnError)
	fs.IntVar(&idReceiver, "idReceiver", 0, "use --idReceiver=1")
	fs.IntVar(&id, "id", 0, "use --id=1")

	if err := fs.Parse(args); err != nil {
		return 0, 0, err
	}

	if idReceiver == 0 {
		return 0, 0, customErrors.ErrReceiverNotFound
	}

	if id == 0 {
		return 0, 0, customErrors.ErrIDNotFound
	}
	return id, idReceiver, nil
}
