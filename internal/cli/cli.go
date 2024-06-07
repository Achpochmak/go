package cli

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
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
}

// NewCLI creates a command line interface
func NewCLI(d Deps) CLI {
	return CLI{
		Deps: d,
		commandList: []command{
			{
				name:        help,
				description: "справка",
			},
			{
				name:        addOrder,
				description: "добавить заказ: использование add --id=1 --id_receiver=8435432342 --storage_time=2023-06-15T15:04:05Z",
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
		},
	}
}

// Run ..
func (c CLI) Run() error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("> ")
	input, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	args := strings.Fields(strings.TrimSpace(input))
	if len(args) == 0 {
		return fmt.Errorf("command isn't set")
	}

	commandName := args[0]
	switch commandName {
	case help:
		c.help()
		return nil
	case addOrder:
		return c.addOrder(args[1:])
	case deleteOrder:
		return c.deleteOrder(args[1:])
	case deliverOrder:
		return c.deliverOrder(args[1:])
	case listOrder:
		return c.listOrder()
	case findOrder:
		return c.findOrder(args[1:])
	case OrdersByCustomer:
		return c.OrdersByCustomer(args[1:])
	case Refund:
		return c.Refund(args[1:])
	case listRefund:
		return c.listRefund(args[1:])
	case exit:
		fmt.Println("Exiting...")
		os.Exit(0)
	}
	return fmt.Errorf("command isn't set")
}

//Добавить заказ
func (c CLI) addOrder(args []string) error {
	var id, id_receiver int
	var storage_time string
	fs := flag.NewFlagSet(addOrder, flag.ContinueOnError)
	fs.IntVar(&id, "id", 0, "use --id=1")
	fs.IntVar(&id_receiver, "id_receiver", 0, "use --id_receiver=1")
	fs.StringVar(&storage_time, "storage_time", "", "use --storage_time=2023-06-15T15:04:05Z")

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
func (c CLI) listOrder() error {
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
func (c CLI) deleteOrder(args []string) error {
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
	return c.Module.DeleteOrder(models.Order(order))
}

//Доставить заказ
func (c CLI) deliverOrder(args []string) error {
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
		num, err := strconv.Atoi(numStr)
		if err != nil {
			return customErrors.ErrIdNotFound
		}
		ids = append(ids, num)
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
func (c CLI) findOrder(args []string) error {
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
func (c CLI) help() {
	fmt.Println("command list:")
	for _, cmd := range c.commandList {
		fmt.Println("", cmd.name, cmd.description)
	}
	return
}

//Получить список заказов по получателю
func (c CLI) OrdersByCustomer(args[]string) error {
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
func (c CLI) Refund(args []string) error {
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

	err := c.Module.Refund(id, id_receiver)
	if err != nil {
		return err
	}
	
	return nil
}

//Список возвратов
func (c CLI) listRefund(args []string) error {
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