package sender_test

import (
	"testing"
	"time"

	"HOMEWORK-1/internal/infrastructure/app/sender"
	"HOMEWORK-1/internal/infrastructure/kafka"

	"github.com/stretchr/testify/assert"
)

func TestSendMessage(t *testing.T) {
	brokers := []string{"localhost:9091"}
	topic := "test-topic"

	producer, err := kafka.NewProducer(brokers)
	if err != nil {
		t.Fatalf("Failed to create Kafka producer: %v", err)
	}

	kafkaSender := sender.NewKafkaSender(producer, topic)

	testMessage := &sender.Message{
		AnswerID:      1,
		Command:       "test-command",
		Args:          []string{"arg1", "arg2"},
		Success:       false,
		CreatedAt:     time.Now(),
		ProcessedInOB: time.Time{},
		IsAquired:     false,
		IsProcessed:   false,
		RetryCount:    0,
	}

	err = kafkaSender.SendMessage(testMessage)

	assert.NoError(t, err, "Expected no error from SendMessage")
	assert.True(t, testMessage.Success, "Expected Success to be true after SendMessage")
}
