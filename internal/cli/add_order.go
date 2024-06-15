package cli

import (
	"context"
	"flag"
	"fmt"
	"time"

	"HOMEWORK-1/internal/models"
	"HOMEWORK-1/internal/models/customErrors"
	"HOMEWORK-1/pkg/hash"
)

// Добавить заказ
func (c *CLI) addOrder(ctx context.Context, args []string) error {
	id, idReceiver, st, err := c.parseAddOrder(args)
	if err != nil {
		return err
	}
	err = c.checkOrder(ctx, id, st)
	if err != nil {
		return err
	}

	return c.Module.AddOrder(ctx, models.Order{
		ID:          models.ID(id),
		IDReceiver:  models.ID(idReceiver),
		StorageTime: st,
		Delivered:   false,
		CreatedAt:   time.Now(),
		Hash:        hash.GenerateHash(),
	})
}

// Парсинг параметров добавления заказа
func (c *CLI) parseAddOrder(args []string) (int, int, time.Time, error) {
	var id, idReceiver int
	var storageTime string
	fs := flag.NewFlagSet(addOrder, flag.ContinueOnError)
	fs.IntVar(&id, "id", 0, "use --id=1")
	fs.IntVar(&idReceiver, "idReceiver", 0, "use --idReceiver=1")
	fs.StringVar(&storageTime, "storageTime", "", "use --storageTime=2025-06-15T15:04:05Z")

	if err := fs.Parse(args); err != nil {
		return 0, 0, time.Now(), err
	}

	if id == 0 {
		return 0, 0, time.Now(), customErrors.ErrIDNotFound
	}

	if idReceiver == 0 {
		return 0, 0, time.Now(), customErrors.ErrReceiverNotFound
	}

	st, err := time.Parse(time.RFC3339, storageTime)
	if err != nil {
		return 0, 0, time.Now(), customErrors.ErrWrongTimeFormat
	}

	return id, idReceiver, st, nil
}

func (c *CLI) checkOrder(ctx context.Context, id int, st time.Time) error {
	order, err := c.Module.GetOrderByID(ctx, models.ID(id))
	if err == nil {
		fmt.Printf("ID заказа: %d\nID получателя: %d\nВремя хранения: %s\n", order.ID, order.IDReceiver, order.StorageTime)
		return customErrors.ErrOrderAlreadyExists
	}
	if time.Now().After(st) {
		return customErrors.ErrStorageTimeEnded
	}
	return nil
}
