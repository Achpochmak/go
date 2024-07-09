package cli_test

import (
	"testing"
	"time"

	"HOMEWORK-1/internal/cli"
	"HOMEWORK-1/internal/infrastructure/app/sender"

	"context"

	"github.com/stretchr/testify/assert"
)

func TestOutboxRepo(t *testing.T) {
	brokers := []string{"localhost:9091"}
	topic := "test-topic"

	kafkaSender, err := sender.NewKafkaSender(brokers, topic)
	assert.NoError(t, err, "Expected no error from NewKafkaSender")

	outbox := cli.OutboxRepo{
		Outbox: make(map[int]*sender.Message),
		Sender: kafkaSender,
	}

	msg := &sender.Message{
		AnswerID: 1,
		Command:  "test-command",
		Args:     []string{"arg1", "arg2"},
		Success:  false,
	}

	outbox.CreateMessage(msg)
	assert.NotNil(t, outbox.Outbox[msg.AnswerID])

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go outbox.OutboxProcessor(ctx)
	//Время на обработку
	time.Sleep(2 * time.Second)

	processedMsg := outbox.Outbox[msg.AnswerID]
	assert.True(t, processedMsg.IsProcessed)
	assert.NotEqual(t, time.Now(), processedMsg.ProcessedInOB)

	cancel()
}
