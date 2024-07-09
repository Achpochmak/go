package sender_test

import (
	"testing"
	"time"

	"HOMEWORK-1/internal/infrastructure/app/sender"

	"github.com/stretchr/testify/assert"
)
var brokers = []string{
	"127.0.0.1:9091",
}
func TestSendMessage(t *testing.T) {
	topic := "test-topic"

	kafkaSender, err := sender.NewKafkaSender(brokers, topic)
	assert.NoError(t, err, "Expected no error from NewKafkaSender")

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
