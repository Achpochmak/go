package kafka

import (
	"context"
	"time"

	"github.com/IBM/sarama"
)

type Consumer struct {
	brokers        []string
	SingleConsumer sarama.Consumer
	ConsumerGroup  *ConsumerGroup
}

func NewConsumer(brokers []string) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = false
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = 5 * time.Second
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumerGroup := NewConsumerGroup()

	return &Consumer{ConsumerGroup: consumerGroup, brokers: brokers}, nil
}

func (c *Consumer) Consume(ctx context.Context, handler func(*sarama.ConsumerMessage)) {
	for {
		select {
		case <-ctx.Done():
			return
		case message := <-c.ConsumerGroup.Ready:
			if message == nil {
				continue
			}
			handler(message)
		}
	}
}
