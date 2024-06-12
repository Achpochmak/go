package cli

import (
	"HOMEWORK-1/internal/models"
	"sync"
)

const (
	help                = "help"
	addOrder            = "add"
	deleteOrder         = "delete"
	deliverOrder        = "deliver"
	listOrder           = "list"
	getOrdersByCustomer = "customer"
	GetOrderByID        = "find"
	refund              = "refund"
	listRefund          = "listrefund"
	exit                = "exit"
	setWorkers          = "setworkers"
)

const (
	help_description                = "справка"
	addOrder_description            = "добавить заказ: использование add --id=1 --id_receiver=1 --storage_time=2025-06-15T15:04:05Z"
	deleteOrder_description         = "удалить заказ: использование delete --id=1"
	deliverOrder_description        = "доставить заказ: использование deliver --id=1,2,3 --id_receive=1"
	listOrder_description           = "вывести список заказов: использование list"
	getOrdersByCustomer_description = "customer"
	GetOrderByID_description        = "найти заказ: использование find --id=1"
	refund_description              = "вернуть заказ: использование refund --id_receiver=1 --id=1"
	listRefund_description          = "вывести список возвратов: использование refund (опционально:--page=1 --page_size=1)"
	exit_description                = "exit"
	setWorkers_description          = "вывести список возвратов: использование setWorkers --num=5"
)

type command struct {
	name        string
	description string
}
type Module interface {
	AddOrder(order models.Order) error
	ListOrder() ([]models.Order, error)
	DeleteOrder(order models.Order) error
	DeliverOrder([]int, int) ([]models.Order, error)
	GetOrderByID(models.ID) (models.Order, error)
	GetOrdersByCustomer(int, int) ([]models.Order, error)
	Refund(int, int) error
	ListRefund() ([]models.Order, error)
}

type Deps struct {
	Module Module
}

type CLI struct {
	Deps
	commandList   []command
	taskQueue     chan task
	notifications chan string
	workerPool    chan struct{}
	numWorkers    int
	mu            sync.Mutex
	wg            sync.WaitGroup
	orderLocks    map[models.ID]*sync.Mutex
}

type task struct {
	commandName string
	args        []string
}
