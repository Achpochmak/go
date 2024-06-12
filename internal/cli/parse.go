package cli

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
	"time"

	"HOMEWORK-1/internal/models"
	"HOMEWORK-1/internal/models/customErrors"
)

// Парсинг параметров изменения количества рутин
func (c *CLI) parseSetWorkers(args []string) (int, error) {
	var num int
	fs := flag.NewFlagSet("setWorkers", flag.ContinueOnError)
	fs.IntVar(&num, "num", c.numWorkers, "use --num=1")

	if err := fs.Parse(args); err != nil {
		return 0, err
	}

	if num < 1 {
		return 0, customErrors.ErrWorkersLessThanOne
	}
	return num, nil
}

// Парсинг параметров добавления заказа
func (c *CLI) parseAddOrder(args []string) (int, int, time.Time, error) {
	var id, id_receiver int
	var storage_time string
	fs := flag.NewFlagSet(addOrder, flag.ContinueOnError)
	fs.IntVar(&id, "id", 0, "use --id=1")
	fs.IntVar(&id_receiver, "id_receiver", 0, "use --id_receiver=1")
	fs.StringVar(&storage_time, "storage_time", "", "use --storage_time=2025-06-15T15:04:05Z")

	if err := fs.Parse(args); err != nil {
		return 0, 0, time.Now(), err
	}

	if id == 0 {
		return 0, 0, time.Now(), customErrors.ErrIDNotFound
	}

	order, err := c.Module.GetOrderByID(models.ID(id))
	if err == nil {
		fmt.Printf("ID заказа: %d\nID получателя: %d\nВремя хранения: %s\n", order.ID, order.ID_receiver, order.Storage_time)
		return 0, 0, time.Now(), customErrors.ErrOrderAlreadyExists
	}

	if id_receiver == 0 {
		return 0, 0, time.Now(), customErrors.ErrReceiverNotFound
	}

	st, err := time.Parse(time.RFC3339, storage_time)
	if err != nil {
		return 0, 0, time.Now(), customErrors.ErrWrongTimeFormat
	}

	if time.Now().After(st) {
		return 0, 0, time.Now(), customErrors.ErrStorageTimeEnded
	}
	return id, id_receiver, st, nil
}

// Парсинг ID заказа
func (c *CLI) parseID(args []string) (int, error) {
	var id int
	fs := flag.NewFlagSet(deleteOrder, flag.ContinueOnError)
	fs.IntVar(&id, "id", 0, "use --id=1")

	if err := fs.Parse(args); err != nil {
		return 0, err
	}

	if id == 0 {
		return 0, customErrors.ErrIDNotFound
	}
	return id, nil
}

// Парсинг параметров доставки заказа
func (c *CLI) parseDeliverOrder(args []string) ([]int, int, error) {
	var id_receiver int
	var orderIdsUnparsed string
	fs := flag.NewFlagSet(deliverOrder, flag.ContinueOnError)
	fs.IntVar(&id_receiver, "id_receiver", 0, "use --id=1")
	fs.StringVar(&orderIdsUnparsed, "id", "", "use --id=1,2,3")

	if err := fs.Parse(args); err != nil {
		return nil, 0, err
	}

	if id_receiver == 0 {
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
	return ids, id_receiver, nil
}

// Парсинг получения параметров получения заказа по клиенту
func (c *CLI) parseGetOrdersByCustomer(args []string) (int, int, error) {
	var id_receiver, amount int
	fs := flag.NewFlagSet(getOrdersByCustomer, flag.ContinueOnError)
	fs.IntVar(&id_receiver, "id_receiver", 0, "use --id_receiver=1")
	fs.IntVar(&amount, "n", 0, "use --n=1")

	if err := fs.Parse(args); err != nil {
		return 0, 0, err
	}

	if id_receiver == 0 {
		return 0, 0, customErrors.ErrReceiverNotFound
	}
	return id_receiver, amount, nil
}

// Парсинг параметров возврата
func (c *CLI) parseRefund(args []string) (int, int, error) {
	var id_receiver, id int
	fs := flag.NewFlagSet(refund, flag.ContinueOnError)
	fs.IntVar(&id_receiver, "id_receiver", 0, "use --id_receiver=1")
	fs.IntVar(&id, "id", 0, "use --id=1")

	if err := fs.Parse(args); err != nil {
		return 0, 0, err
	}

	if id_receiver == 0 {
		return 0, 0, customErrors.ErrReceiverNotFound
	}

	if id == 0 {
		return 0, 0, customErrors.ErrIDNotFound
	}
	return id, id_receiver, nil
}

// Парсинг параметров листа возвратов
func (c *CLI) parseListRefund(args []string) (int, int, error) {
	var page, page_size int
	fs := flag.NewFlagSet(listRefund, flag.ContinueOnError)
	fs.IntVar(&page, "page", 0, "use --id_receiver=1")
	fs.IntVar(&page_size, "page_size", 0, "use --id=1")

	if err := fs.Parse(args); err != nil {
		return 0, 0, err
	}
	return page, page_size, nil
}
