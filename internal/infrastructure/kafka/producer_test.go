package kafka_test

import (
	"testing"

	"HOMEWORK-1/internal/infrastructure/kafka"

	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var brokers = []string{
	"127.0.0.1:9091",
	"127.0.0.1:9092",
	"127.0.0.1:9093",
}

func TestProducer(t *testing.T) {
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
}
