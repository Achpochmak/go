package cli

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"HOMEWORK-1/internal/infrastructure/app/receiver"
	"HOMEWORK-1/internal/infrastructure/app/sender"
	"HOMEWORK-1/internal/infrastructure/kafka"

	"github.com/IBM/sarama"
)

// NewCLI creates a command line interface
func NewCLI(d Deps) *CLI {
	cli := &CLI{
		Deps: d,
		commandList: []command{
			{
				name:        help,
				description: helpDescription,
			},
			{
				name:        addOrder,
				description: addOrderDescription,
			},
			{
				name:        deleteOrder,
				description: deleteOrderDescription,
			},
			{
				name:        deliverOrder,
				description: deliverOrderDescription,
			},
			{
				name:        GetOrderByID,
				description: GetOrderByIDDescription,
			},
			{
				name:        listOrder,
				description: listOrderDescription,
			},
			{
				name:        refund,
				description: refundDescription,
			},
			{
				name:        listRefund,
				description: listRefundDescription,
			},
			{
				name:        setWorkers,
				description: setWorkersDescription,
			},
			{
				name:        switchOutput,
				description: switchOutputDescription,
			},
		},
		taskQueue:     make(chan task, 10),
		numWorkers:    2,
		workerPool:    make(chan struct{}, 2),
		notifications: make(chan string, 10),
		taskQueueOpen: true,
		wg:            sync.WaitGroup{},
		outputKafka:   true,
		outbox: OutboxRepo{
			Mu:     sync.RWMutex{},
			Outbox: make(map[int]*sender.Message),
		},
		kafkaConfig: KafkaConfig{
			Brokers: brokers,
			Topic:   "my-topic",
		},
		AnswerID: 0,
	}

	return cli
}

// Run ..
func (c *CLI) Run() error {
	go c.notificationHandler()
	defer close(c.notifications)
	ctx, cancel := context.WithCancel(context.Background())
	c.InitKafka(ctx)

	defer cancel()
	for i := 0; i < c.numWorkers; i++ {
		c.wg.Add(1)
		go c.worker(ctx)
	}

	c.handleSignals(cancel)

	ctxOutbox, cancelOutbox := context.WithCancel(context.Background())

	//Запускаем запись в outbox
	go c.outbox.OutboxProcessor(ctxOutbox)

	reader := bufio.NewReader(os.Stdin)
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			cancelOutbox()
			return err
		}

		args := strings.Fields(strings.TrimSpace(input))
		if len(args) == 0 {
			fmt.Println("command isn't set")
			continue
		}

		commandName := args[0]
		c.ProcessCommand(commandName, args[1:])

		if commandName == exit {
			c.mu.Lock()
			if c.taskQueueOpen {
				c.taskQueueOpen = false
				close(c.taskQueue)
			}
			c.mu.Unlock()

			go func() {
				time.Sleep(5 * time.Second)
				cancel()
			}()
			break
		}

		if c.taskQueueOpen {
			if !c.outputKafka {
				//если вывод не через кафку, записываем в очередь выполнения то, что получили из консоли
				c.taskQueue <- task{commandName: args[0], args: args[1:]}
			}
		} else {
			fmt.Println("Доступ закрыт")
		}
	}
	cancelOutbox()

	c.wg.Wait()
	fmt.Println("Все задачи завершены.")
	os.Exit(0)
	return nil
}

//Запись сообщения в outbox
func (c *CLI) ProcessCommand(commandName string, args []string) {
	c.AnswerID++
	answerID := c.AnswerID

	msg := sender.Message{
		Command:     commandName,
		Args:        args,
		AnswerID:    int(answerID),
		CreatedAt:   time.Now(),
		Success:     false,
		IsAquired:   false,
		IsProcessed: false,
	}

	c.outbox.CreateMessage(&msg)
}

func (c *CLI) InitKafka(ctx context.Context) error {
	consumer, err := kafka.NewConsumer(c.kafkaConfig.Brokers)
	if err != nil {
		return fmt.Errorf("ошибка создания Kafka consumer: %w", err)
	}

	producer, err := kafka.NewProducer(c.kafkaConfig.Brokers)
	if err != nil {
		return fmt.Errorf("ошибка создания Kafka producer: %w", err)
	}

	handlers := map[string]receiver.HandleFunc{
		c.kafkaConfig.Topic: func(message *sarama.ConsumerMessage) {
			msg := sender.Message{}
			err = json.Unmarshal(message.Value, &msg)
			if err != nil {
				fmt.Println("Consumer error", err)
			}
			//Если вывод через кафку, то записываем в очередь выполения
			if c.taskQueueOpen && c.outputKafka {
				c.taskQueue <- task{commandName: msg.Command, args: msg.Args}
			}
		},
	}

	c.KafkaSender = sender.NewKafkaSender(producer, c.kafkaConfig.Topic)
	c.KafkaReceiver = receiver.NewReceiver(consumer, handlers)

	c.outbox.Sender = c.KafkaSender
	c.KafkaReceiver.Subscribe(c.kafkaConfig.Topic)
	return nil
}
