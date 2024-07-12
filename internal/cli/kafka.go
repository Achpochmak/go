package cli

import (
	"HOMEWORK-1/internal/infrastructure/app/receiver"
	"HOMEWORK-1/internal/infrastructure/app/sender"
	"context"
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
)

func (c *CLI) InitKafka(ctx context.Context) error {
	var err error

	KafkaSender, err := sender.NewKafkaSender(c.kafkaConfig.Brokers, c.kafkaConfig.Topic)
	if err != nil {
		return err
	}

	KafkaReceiver, err := receiver.NewReceiver(c.kafkaConfig.Brokers, c.getKafkaHandlers())
	if err != nil {
		return err
	}

	c.outbox.Sender = KafkaSender
	KafkaReceiver.Subscribe(c.kafkaConfig.Topic)
	return nil
}

func (c *CLI) getKafkaHandlers() map[string]receiver.HandleFunc {
	return map[string]receiver.HandleFunc{
		c.kafkaConfig.Topic: func(message *sarama.ConsumerMessage) {
			var msg sender.Message
			if err := json.Unmarshal(message.Value, &msg); err != nil {
				fmt.Println("Consumer error", err)
				return
			}
			//Вывод из кафки(задание 4)
			if c.outputKafka{
				fmt.Println("Kafka: ", msg.Command, " status: ", msg.Status)
			}
		},
	}
}
