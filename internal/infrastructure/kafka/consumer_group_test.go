package kafka_test

import (
	"context"
	"testing"
	"time"

	"HOMEWORK-1/internal/infrastructure/kafka"

	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConsumerGroup(t *testing.T) {
	consumer, err := kafka.NewConsumer(brokers)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	handler := func(message *sarama.ConsumerMessage) {
		assert.Equal(t, "test-topic", message.Topic)
		assert.Equal(t, "test message", string(message.Value))
		cancel()
	}

	go consumer.Consume(ctx, handler)

	producer, err := kafka.NewProducer(brokers)
	require.NoError(t, err)

	defer func() {
		err := producer.Close()
		require.NoError(t, err)
	}()

	msg := &sarama.ProducerMessage{
		Topic: "test-topic",
		Value: sarama.StringEncoder("test message"),
	}

	_, _, err = producer.SendSyncMessage(msg)
	assert.NoError(t, err)

	<-ctx.Done()
}
