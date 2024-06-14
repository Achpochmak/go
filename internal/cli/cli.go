package cli

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"HOMEWORK-1/internal/models"
	"HOMEWORK-1/pkg/hash"
)

// NewCLI creates a command line interface
func NewCLI(d Deps) *CLI {
	cli := &CLI{
		Deps: d,
		commandList: []command{
			{
				name:        help,
				description: help_description,
			},
			{
				name:        addOrder,
				description: addOrder_description,
			},
			{
				name:        deleteOrder,
				description: deleteOrder_description,
			},
			{
				name:        deliverOrder,
				description: deliverOrder_description,
			},
			{
				name:        GetOrderByID,
				description: GetOrderByID_description,
			},
			{
				name:        listOrder,
				description: listOrder_description,
			},
			{
				name:        refund,
				description: refund_description,
			},
			{
				name:        listRefund,
				description: listRefund_description,
			},
			{
				name:        setWorkers,
				description: setWorkers_description,
			},
		},
		taskQueue:     make(chan task, 10),
		numWorkers:    2,
		workerPool:    make(chan struct{}, 2),
		orderLocks:    make(map[models.ID]*sync.Mutex),
		notifications: make(chan string, 10),
	}
	go cli.notificationHandler()
	return cli
}

// Run ..
func (c *CLI) Run() error {
	ctx := context.Background() 

	for i := 0; i < c.numWorkers; i++ {
		c.wg.Add(1)
		go c.worker(ctx)
	}

	c.handleSignals()

	reader := bufio.NewReader(os.Stdin)
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		args := strings.Fields(strings.TrimSpace(input))
		if len(args) == 0 {
			fmt.Println("command isn't set")
			continue
		}

		commandName := args[0]
		if commandName == exit {
			close(c.taskQueue)
			break
		}
		c.taskQueue <- task{commandName: commandName, args: args[1:]}
	}

	c.wg.Wait()
	fmt.Println("All tasks completed. Exiting...")
	os.Exit(0)
	return nil
}

// Обработка уведомлений
func (c *CLI) notificationHandler() {
	for msg := range c.notifications {
		fmt.Println(msg)
	}
}

func (c *CLI) worker(ctx context.Context) {
	defer c.wg.Done()
	for t := range c.taskQueue {
		var err error
		startMsg := fmt.Sprintf("Началась обработка команды: %s", t.commandName)
		endMsg := fmt.Sprintf("Завершилась обработка команды: %s", t.commandName)
		c.notifications <- startMsg
		switch t.commandName {

		case help:
			c.help()
		case addOrder:
			err = c.addOrder(ctx,t.args)
		case deleteOrder:
			err = c.deleteOrder(ctx,t.args)
		case deliverOrder:
			err = c.deliverOrder(ctx,t.args)
		case listOrder:
			err = c.listOrder(ctx)
		case GetOrderByID:
			err = c.GetOrderByID(ctx,t.args)
		case getOrdersByCustomer:
			err = c.getOrdersByCustomer(ctx,t.args)
		case refund:
			err = c.refund(ctx,t.args)
		case listRefund:
			err = c.listRefund(ctx,t.args)
		case setWorkers:
			err = c.setWorkers(t.args)
		case exit:
			fmt.Println("Exiting...")
			c.mu.Unlock()
			close(c.taskQueue)
			os.Exit(0)
		default:
			fmt.Println("command isn't set")
		}
		if err != nil {
			fmt.Println("Ошибка:", err)
		}
		c.notifications <- endMsg
	}
}

// Обработка сигналов
func (c *CLI) handleSignals() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigChan
		fmt.Printf("Получена команда %s. Exiting...\n", sig)
		close(c.taskQueue)
		c.wg.Wait()
		os.Exit(0)
	}()
}



// Измeнение числа рутин
func (c *CLI) setWorkers(args []string) error {
	num, err := c.parseSetWorkers(args)
	if err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if num > c.numWorkers {
		for i := c.numWorkers; i < num; i++ {
			c.wg.Add(1)
			go c.worker(context.Background())
		}
	} else if num < c.numWorkers {
		for i := num; i < c.numWorkers; i++ {
			c.taskQueue <- task{commandName: "exit"}
		}
	}

	c.numWorkers = num
	fmt.Printf("Число рутин %d\n", c.numWorkers)
	return nil
}

// Добавить заказ
func (c *CLI) addOrder(ctx context.Context,args []string) error {
	id, id_receiver, st, err := c.parseAddOrder(ctx, args)
	if err != nil {
		return err
	}


	return c.Module.AddOrder(ctx,models.Order{
		ID:           models.ID(id),
		ID_receiver:  models.ID(id_receiver),
		Storage_time: st,
		Delivered:    false,
		Created_at:   time.Now(),
		Hash:         hash.GenerateHash(),
	})
}

// Список заказов
func (c *CLI) listOrder(ctx context.Context) error {
	list, err := c.Module.ListOrder(ctx)
	if err != nil {
		return err
	}

	for _, order := range list {
		fmt.Printf("ID заказа: %d\nID получателя: %d\nВремя хранения: %s\n", order.ID, order.ID_receiver, order.Storage_time)
	}
	return nil
}

// Удалить заказ
func (c *CLI) deleteOrder(ctx context.Context,args []string) error {
	id, err := c.parseID(args)
	if err != nil {
		return err
	}

	order, err := c.Module.GetOrderByID(ctx, models.ID(id))

	if err != nil {
		return err
	}

	return c.Module.DeleteOrder(ctx, models.Order(order))
}

// Доставить заказ
func (c *CLI) deliverOrder(ctx context.Context,args []string) error {
	orderIDs, id_receiver, err := c.parseDeliverOrder(args)
	if err != nil {
		return err
	}

	orders, err := c.Module.DeliverOrder(ctx, orderIDs, id_receiver)
	if err != nil {
		return err
	}

	for _, order := range orders {
		fmt.Printf("ID заказа: %d\nID получателя: %d\nВремя хранения: %s\n", order.ID, order.ID_receiver, order.Storage_time)
	}
	return nil
}

// Найти заказ
func (c *CLI) GetOrderByID(ctx context.Context,args []string) error {
	id, err := c.parseID(args)
	if err != nil {
		return err
	}

	order, err := c.Module.GetOrderByID(ctx, models.ID(id))
	if err != nil {
		return err
	}

	fmt.Printf("ID заказа: %d\nID получателя: %d\nВремя хранения: %s\n", order.ID, order.ID_receiver, order.Storage_time)
	return nil
}

// Помощь
func (c *CLI) help() {
	fmt.Println("command list:")
	for _, cmd := range c.commandList {
		fmt.Println("", cmd.name, cmd.description)
	}
}

// Получить список заказов по получателю
func (c *CLI) getOrdersByCustomer(ctx context.Context,args []string) error {
	id_receiver, amount, err := c.parseGetOrdersByCustomer(args)
	if err != nil {
		return err
	}

	list, err := c.Module.GetOrdersByCustomer(ctx, id_receiver, amount)
	if err != nil {
		return err
	}

	for _, order := range list {
		fmt.Printf("ID заказа: %d\nID получателя: %d\nВремя хранения: %s\n", order.ID, order.ID_receiver, order.Storage_time)
	}
	return nil
}

// Вернуть заказ
func (c *CLI) refund(ctx context.Context,args []string) error {
	id, id_receiver, err := c.parseRefund(args)
	if err != nil {
		return err
	}

	err = c.Module.Refund(ctx, id, id_receiver)
	if err != nil {
		return err
	}

	return nil
}

// Список возвратов
func (c *CLI) listRefund(ctx context.Context,args []string) error {
	page, page_size, err := c.parseListRefund(args)
	if err != nil {
		return err
	}

	list, err := c.Module.ListRefund(ctx)
	if err != nil {
		return err
	}

	if page == 0 || page_size == 0 {
		for _, order := range list {
			fmt.Printf("ID заказа: %d\nID получателя: %d\nВремя хранения: %s\n", order.ID, order.ID_receiver, order.Storage_time)
		}
		return nil
	}

	start := (page - 1) * page_size
	end := start + page_size

	if end > len(list) || start < 0 {
		return errors.New("пустая страница")
	}

	for _, order := range list[start:end] {
		fmt.Printf("ID заказа: %d\nID получателя: %d\nВремя хранения: %s\n", order.ID, order.ID_receiver, order.Storage_time)
	}
	return nil
}
