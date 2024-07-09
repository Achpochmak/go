package integration_tests

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"HOMEWORK-1/internal/cli"
	"HOMEWORK-1/internal/infrastructure/app/receiver"
	"HOMEWORK-1/internal/infrastructure/app/sender"
	"HOMEWORK-1/internal/infrastructure/kafka"

	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
)

func TestKafkaIntegration(t *testing.T) {
	var err error
	topic := "test"
	brokers := []string{"127.0.0.1:9091"}

	consumer, err := kafka.NewConsumer(brokers)
	if err != nil {
		assert.NoError(t, err)
		return
	}

	producer, err := kafka.NewProducer(brokers)
	if err != nil {
		assert.NoError(t, err)
		return
	}

	messageChannel := make(chan *sender.Message)

	handlers := map[string]receiver.HandleFunc{
		topic: func(message *sarama.ConsumerMessage) {
			msg := sender.Message{}
			err = json.Unmarshal(message.Value, &msg)
			if err != nil {
				fmt.Println("Consumer error", err)
				return
			}
			messageChannel <- &msg
		},
	}

	testSender := sender.NewKafkaSender(producer, topic)
	testReceiver := receiver.NewReceiver(consumer, handlers)
	outbox := cli.OutboxRepo{
		Mu:     sync.RWMutex{},
		Outbox: make(map[int]*sender.Message),
		Sender: testSender,
	}

	ctxOutbox, cancelOutbox := context.WithCancel(context.Background())
	defer cancelOutbox()

	go outbox.OutboxProcessor(ctxOutbox)

	err = testReceiver.Subscribe(topic)
	if err != nil {
		assert.NoError(t, err)
		return
	}

	args := strings.Fields(strings.TrimSpace("add --id=2 --idReceiver=1 --storageTime=2025-06-15T15:04:05Z --weightKg=1 --price=100"))
	msg := sender.Message{
		Command:     args[0],
		Args:        args,
		AnswerID:    0,
		CreatedAt:   time.Now(),
		Success:     false,
		IsAquired:   false,
		IsProcessed: false,
	}

	err = testSender.SendMessage(&msg)
	if err != nil {
		assert.NoError(t, err)
		return
	}

	select {
	case receivedMsg := <-messageChannel:
		fmt.Println(receivedMsg)
		assert.NotNil(t, receivedMsg)
		assert.Equal(t, msg.Command, receivedMsg.Command)
		assert.Equal(t, msg.Args, receivedMsg.Args)
	case <-time.After(10 * time.Second):
		t.Fatal("Timeout waiting for message")
	}

	cancelOutbox()
}
