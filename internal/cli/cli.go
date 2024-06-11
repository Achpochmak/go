package cli

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"HOMEWORK-1/internal/models"
	"HOMEWORK-1/internal/models/customErrors"
	"HOMEWORK-1/pkg/hash"
)

type Module interface {
	AddOrder(order models.Order) error
	ListOrder() ([]models.Order, error)
	DeleteOrder(order models.Order) error
	DeliverOrder([]int, int) ([]models.Order,error)
	FindOrder(models.Id) (models.Order, error)
	OrdersByCustomer(int, int)([]models.Order, error)
	Refund(int, int)(error)
	ListRefund() ([]models.Order, error)
}

type Deps struct {
	Module Module
}

type CLI struct {
	Deps
	commandList []command
	taskQueue chan task
	notifications chan string
	workerPool    chan struct{}
	numWorkers int
	mu sync.Mutex
	wg sync.WaitGroup
	orderLocks  map[models.Id]*sync.Mutex
}

type task struct {
	commandName string
	args        []string
}

// NewCLI creates a command line interface
func NewCLI(d Deps) *CLI {
	cli:=&CLI{
		Deps: d,
		commandList: []command{
			{
				name:        help,
				description: "справка",
			},
			{
				name:        addOrder,
				description: "добавить заказ: использование add --id=1 --id_receiver=1 --storage_time=2025-06-15T15:04:05Z",
			},
			{
				name:        deleteOrder,
				description: "удалить заказ: использование delete --id=1",
			},
			{
				name:        deliverOrder,
				description: "доставить заказ: использование deliver --id=1,2,3 --id_receive=1",
			},
			{
				name:        findOrder,
				description: "найти заказ: использование find --id=1",
			},
			{
				name:        listOrder,
				description: "вывести список заказов: использование list",
			},
			{
				name:        Refund,
				description: "вернуть заказ: использование refund --id_receiver=1 --id=1",
			},
			{
				name:        listRefund,
				description: "вывести список возвратов: использование refund (опционально:--page=1 --page_size=1)",
			},
			{
				name:        setWorkers,
				description: "вывести список возвратов: использование setWorkers --num=5",
			},
		},
		taskQueue: make(chan task, 10),
		numWorkers: 2,
		workerPool: make(chan struct{}, 2),
		orderLocks: make(map[models.Id]*sync.Mutex),
		notifications: make(chan string, 10),
			
	}
	go cli.notificationHandler()
	return cli
}

// Run ..
func (c *CLI) Run() error {
	for i := 0; i < c.numWorkers; i++ {
		c.wg.Add(1)
		go c.worker()
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

//Обработка уведомлений
func (c *CLI) notificationHandler() {
	for msg := range c.notifications {
		fmt.Println(msg)
	}
}

func (c *CLI) worker() {
	defer c.wg.Done()
	for t := range c.taskQueue {
		startMsg := fmt.Sprintf("Началась обработка команды: %s", t.commandName)
		endMsg := fmt.Sprintf("Завершилась обработка команды: %s", t.commandName)
		c.notifications <- startMsg
		switch t.commandName {

		case help:
			c.help()
		case addOrder:
			if err := c.addOrder(t.args); err != nil {
				fmt.Println("Ошибка:", err)
			}
		case deleteOrder:
			if err := c.deleteOrder(t.args); err != nil {
				fmt.Println("Ошибка:", err)
			}
		case deliverOrder:
			if err := c.deliverOrder(t.args); err != nil {
				fmt.Println("Ошибка:", err)
			}		
		case listOrder:
			if err := c.listOrder(); err != nil {
				fmt.Println("Ошибка:", err)
			}
		case findOrder:
			if err := c.findOrder(t.args); err != nil {
				fmt.Println("Ошибка:", err)
			}
		case OrdersByCustomer:
			if err := c.OrdersByCustomer(t.args); err != nil {
				fmt.Println("Ошибка:", err)
			}		
		case Refund:
			if err := c.Refund(t.args); err != nil {
				fmt.Println("Ошибка:", err)
			}		
		case listRefund:
			if err := c.listRefund(t.args); err != nil {
				fmt.Println("Ошибка:", err)
			}
		case setWorkers:
			if err := c.setWorkers(t.args); err != nil {
				fmt.Println("Ошибка:", err)
			}		
		case exit:
			fmt.Println("Exiting...")
			c.mu.Unlock()
			close(c.taskQueue)
			os.Exit(0)
		default:
			fmt.Println("command isn't set")
		}
		c.notifications <- endMsg
		
	}
}

//Обработка сигналов
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

//Заблокировать заказ
func (c *CLI) lockOrder(id int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, exists := c.orderLocks[models.Id(id)]; !exists {
		c.orderLocks[models.Id(id)] = &sync.Mutex{}
	}
	c.orderLocks[models.Id(id)].Lock()
}

//Разблокировать заказ
func (c *CLI) unlockOrder(id int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if lock, exists := c.orderLocks[models.Id(id)]; exists {
		lock.Unlock()
	}
}

//Измнение числа рутин
func (c *CLI) setWorkers(args []string) error {
	var num int
	fs := flag.NewFlagSet("setWorkers", flag.ContinueOnError)
	fs.IntVar(&num, "num", c.numWorkers, "use --num=1")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if num < 1 {
		return errors.New("количество должно быть 1")
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if num > c.numWorkers {
		for i := c.numWorkers; i < num; i++ {
			c.wg.Add(1)
			go c.worker()
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

//Добавить заказ
func (c *CLI) addOrder(args []string) error {
	var id, id_receiver int
	var storage_time string
	fs := flag.NewFlagSet(addOrder, flag.ContinueOnError)
	fs.IntVar(&id, "id", 0, "use --id=1")
	fs.IntVar(&id_receiver, "id_receiver", 0, "use --id_receiver=1")
	fs.StringVar(&storage_time, "storage_time", "", "use --storage_time=2025-06-15T15:04:05Z")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if id == 0 {
		return customErrors.ErrIdNotFound
	}
	
	order, err := c.Module.FindOrder(models.Id(id))
	if err == nil {
		fmt.Printf("Id заказа: %d\nId получателя: %d\nВремя хранения: %s\n", order.Id, order.Id_receiver, order.Storage_time)
		return errors.New("заказ дублируется")
	}

	if id_receiver == 0 {
		return customErrors.ErrReceiverNotFound
	}

	st, err := time.Parse(time.RFC3339, storage_time)
	if err != nil {
		return errors.New("неправильный формат времени")
	}

	if !time.Now().Before(st){
		return errors.New("время хранения окончилось")
	}

	c.lockOrder(id)
	defer c.unlockOrder(id)

	return c.Module.AddOrder(models.Order{
		Id:      models.Id(id),
		Id_receiver: models.Id(id_receiver),
		Storage_time: st,
		Delivered: false,
		Created_at: time.Now(),
		Hash: hash.GenerateHash(),
	})
	}

//Список заказов
func (c *CLI) listOrder() error {
	list, err := c.Module.ListOrder()
	if err != nil {
		return err
	}

	for _, order := range list {
		fmt.Printf("Id заказа: %d\nId получателя: %d\nВремя хранения: %s\n", order.Id, order.Id_receiver, order.Storage_time)
	}
	return nil
}

//Удалить заказ
func (c *CLI) deleteOrder(args []string) error {
	var id int
	fs := flag.NewFlagSet(deleteOrder, flag.ContinueOnError)
	fs.IntVar(&id, "id", 0, "use --id=1")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if id == 0 {
		return customErrors.ErrIdNotFound
	}

	order, err := c.Module.FindOrder(models.Id(id))
	
	if err != nil {
		return err
	}
	c.lockOrder(id)
	defer c.unlockOrder(id)

	return c.Module.DeleteOrder(models.Order(order))
}

//Доставить заказ
func (c *CLI) deliverOrder(args []string) error {
	var id_receiver int
	var order_ids string
	fs := flag.NewFlagSet(deliverOrder, flag.ContinueOnError)
	fs.IntVar(&id_receiver, "id_receiver", 0, "use --id=1")
	fs.StringVar(&order_ids, "id", "", "use --id=1,2,3")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if id_receiver == 0 {
		return customErrors.ErrReceiverNotFound
	}

	orderIds := strings.Split(order_ids, ",")
	var ids []int
	for _, numStr := range orderIds {
		id, err := strconv.Atoi(numStr)
		if err != nil {
			return customErrors.ErrIdNotFound
		}
		ids = append(ids, id)
		c.lockOrder(id)
		defer c.unlockOrder(id)
	}

	orders, err := c.Module.DeliverOrder(ids, id_receiver)
	if err != nil {
		return err
	}

	for _,order:=range orders{
		fmt.Printf("Id заказа: %d\nId получателя: %d\nВремя хранения: %s\n", order.Id, order.Id_receiver, order.Storage_time)
	}
	return nil
}

//Найти заказ
func (c *CLI) findOrder(args []string) error {
	var id int
	fs := flag.NewFlagSet(findOrder, flag.ContinueOnError)
	fs.IntVar(&id, "id", 0, "use --id=1")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if id == 0 {
		return customErrors.ErrIdNotFound
	}

	order, err := c.Module.FindOrder(models.Id(id))
	if err != nil {
		return err
	}

	fmt.Printf("Id заказа: %d\nId получателя: %d\nВремя хранения: %s\n", order.Id, order.Id_receiver, order.Storage_time)
	return nil
}

//Помощь
func (c *CLI) help() {
	fmt.Println("command list:")
	for _, cmd := range c.commandList {
		fmt.Println("", cmd.name, cmd.description)
	}
}

//Получить список заказов по получателю
func (c *CLI) OrdersByCustomer(args[]string) error {
	var id_receiver, amount int
	fs := flag.NewFlagSet(OrdersByCustomer, flag.ContinueOnError)
	fs.IntVar(&id_receiver, "id_receiver", 0, "use --id_receiver=1")
	fs.IntVar(&amount, "n", 0, "use --n=1")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if id_receiver == 0 {
		return customErrors.ErrReceiverNotFound
	}

	list, err := c.Module.OrdersByCustomer(id_receiver, amount)
	if err != nil {
		return err
	}

	for _, order := range list {
		fmt.Printf("Id заказа: %d\nId получателя: %d\nВремя хранения: %s\n", order.Id, order.Id_receiver, order.Storage_time)
	}
	return nil
}

//Вернуть заказ
func (c *CLI) Refund(args []string) error {
	var id_receiver, id int
	fs := flag.NewFlagSet(Refund, flag.ContinueOnError)
	fs.IntVar(&id_receiver, "id_receiver", 0, "use --id_receiver=1")
	fs.IntVar(&id, "id", 0, "use --id=1")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if id_receiver == 0 {
		return customErrors.ErrReceiverNotFound
	}

	if id == 0 {
		return customErrors.ErrIdNotFound
	}

	c.lockOrder(id)
	defer c.unlockOrder(id)
	err := c.Module.Refund(id, id_receiver)
	if err != nil {
		return err
	}
	
	return nil
}

//Список возвратов
func (c *CLI) listRefund(args []string) error {
	var page, page_size int
	fs := flag.NewFlagSet(listRefund, flag.ContinueOnError)
	fs.IntVar(&page, "page", 0, "use --id_receiver=1")
	fs.IntVar(&page_size, "page_size", 0, "use --id=1")

	if err := fs.Parse(args); err != nil {
		return err
	}

	list, err := c.Module.ListRefund()
	if err != nil {
		return err
	}

	if page == 0 || page_size==0 {	
		for _, order := range list {
			fmt.Printf("Id заказа: %d\nId получателя: %d\nВремя хранения: %s\n", order.Id, order.Id_receiver, order.Storage_time)
		}
		return nil
	}

	start:=(page-1)*page_size
	end:= start+page_size

	if end>len(list) || start<0{
		return errors.New("пустая страница")
	}

	for _, order := range list[start: end] {
		fmt.Printf("Id заказа: %d\nId получателя: %d\nВремя хранения: %s\n", order.Id, order.Id_receiver, order.Storage_time)
	}
	return nil
}