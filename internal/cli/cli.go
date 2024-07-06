package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"HOMEWORK-1/internal/infrastructure/app/sender"
	"HOMEWORK-1/internal/infrastructure/kafka"
	"HOMEWORK-1/internal/infrastructure/outbox"

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
		outbox: outbox.OutboxRepo{
			Mu:     sync.RWMutex{},
			Outbox: make(map[int]*sender.Message),
		},
		outputKafka: false,
		kafkaConfig: KafkaConfig{
			Brokers: outbox.Brokers,
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

	defer cancel()
	for i := 0; i < c.numWorkers; i++ {
		c.wg.Add(1)
		go c.worker(ctx)
	}

	c.handleSignals(cancel)

	c.InitKafka(ctx)

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
		c.processCommand(commandName, input)

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
			if c.outputKafka {
				go c.ConsumeFromKafka(ctx)
			} else {
				c.ReadFromStdin(cancel)
			}
		} else {
			fmt.Println("Доступ закрыт")
		}
	}

	c.wg.Wait()
	fmt.Println("Все задачи завершены.")
	os.Exit(0)
	return nil
}

func (c *CLI) processCommand(commandName, input string) {
	go func() {
		ctxOutbox, cancelOutbox := context.WithCancel(context.Background())
		defer cancelOutbox()

		c.AnswerID++
		answerID := c.AnswerID
		c.outbox.CreateMessage(&sender.Message{Command: commandName, Request: input, AnswerID: int(answerID)})
		c.outbox.OutboxProcessor(ctxOutbox)
		<-ctxOutbox.Done()

			msg := sender.Message{
				Command:       commandName,
				Request:       input,
				AnswerID:      int(answerID),
				CreatedAt:     time.Now(),
				IsAquired:     true,
				IsProcessed:   true,
				Success:       true,
				ProcessedInOB: time.Now(),
			}
			c.kafkaSender.SendMessage(&msg)
		
	}()
}
func (c *CLI) ConsumeFromKafka(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			c.kafkaConsumer.Consume(ctx, func(message *sarama.ConsumerMessage) {
				command := strings.Fields(string(message.Value))
				if len(command) > 0 {
					commandName := command[0]
					c.taskQueue <- task{commandName: commandName, args: command[1:]}
				}
			})
		}
	}
}
func (c *CLI) ReadFromStdin(cancel context.CancelFunc) {
	reader := bufio.NewReader(os.Stdin)
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		args := strings.Fields(strings.TrimSpace(input))
		if len(args) == 0 {
			fmt.Println("command isn't set")
			continue
		}

		commandName := args[0]
		c.processCommand(commandName, input)

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

		c.taskQueue <- task{commandName: commandName, args: args[1:]}

	}
}

func (c *CLI) InitKafka(ctx context.Context) error {
	var err error
	c.kafkaConsumer, err = kafka.NewConsumer(c.kafkaConfig.Brokers)
	if err != nil {
		return fmt.Errorf("ошибка создания Kafka consumer: %w", err)
	}

	producer, err := kafka.NewProducer(c.kafkaConfig.Brokers)
	if err != nil {
		return fmt.Errorf("ошибка создания Kafka producer: %w", err)
	}

	c.kafkaSender = sender.NewKafkaSender(producer, c.kafkaConfig.Topic)

	consumerGroup := kafka.NewConsumerGroup()
	group, err := sarama.NewConsumerGroup(c.kafkaConfig.Brokers, "group-id", nil)
	if err != nil {
		return fmt.Errorf("ошибка создания Kafka consumer group: %w", err)
	}

	go func() {
		for {
			if err := group.Consume(ctx, []string{c.kafkaConfig.Topic}, consumerGroup); err != nil {
				fmt.Printf("Error from consumer: %v\n", err)
			}
			if ctx.Err() != nil {
				return
			}
			consumerGroup.Ready = make(chan *sarama.ConsumerMessage)
		}
	}()

	<-consumerGroup.IsReady()
	fmt.Println("Kafka consumer group is ready.")
	return nil
}