//go:generate mockgen -source ./init.go -destination=./mocks/cli.go -package=mock_cli

package cli

import (
	"context"
	"sync"

	"HOMEWORK-1/internal/models"
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
	switchOutput        = "switch"
)

const (
	helpDescription                = "справка"
	addOrderDescription            = "добавить заказ: использование add --id=1 --idReceiver=1 --storageTime=2025-06-15T15:04:05Z --weightKg=1 --price=100(опционально --packaging=box|film|bag)"
	deleteOrderDescription         = "удалить заказ: использование delete --id=1"
	deliverOrderDescription        = "доставить заказ: использование deliver --id=1,2 --idReceiver=1"
	listOrderDescription           = "вывести список заказов: использование list"
	getOrdersByCustomerDescription = "вывести последние n заказов покупателя: использование customer --id=1 --n=1"
	GetOrderByIDDescription        = "найти заказ: использование find --id=1"
	refundDescription              = "вернуть заказ: использование refund --idReceiver=1 --id=1,2"
	listRefundDescription          = "вывести список возвратов: использование listrefund (опционально:--page=1 --pageSize=1)"
	exitDescription                = "exit"
	setWorkersDescription          = "вывести список возвратов: использование setWorkers --num=5"
	switchOutputDescription        = "изменить метод вывода: switch"
)

type command struct {
	name        string
	description string
}

type Module interface {
	AddOrder(context.Context, models.Order) error
	ListOrder(context.Context) ([]models.Order, error)
	DeleteOrder(context.Context, models.ID) error
	DeliverOrder(context.Context, []int, int) ([]*models.Order, error)
	GetOrderByID(context.Context, models.ID) (models.Order, error)
	GetOrdersByCustomer(context.Context, int, int) ([]models.Order, error)
	Refund(context.Context, int, int) error
	ListRefund(context.Context, int, int) ([]models.Order, error)
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
	wg            sync.WaitGroup
	mu            sync.Mutex
	taskQueueOpen bool
	outbox        OutboxRepo
	outputKafka   bool
	kafkaConfig   KafkaConfig
	AnswerID      int
}

type KafkaConfig struct {
	Brokers []string
	Topic   string
}

type task struct {
	commandName string
	args        []string
}

var brokers = []string{
	"127.0.0.1:9091",
	"127.0.0.1:9092",
	"127.0.0.1:9093",
}
