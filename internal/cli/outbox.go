package cli

import (
	"HOMEWORK-1/internal/infrastructure/app/sender"
	"context"
	"sync"
	"time"
)

type OutboxRepo struct {
	Mu     sync.RWMutex
	Outbox map[int]*sender.Message
	Sender *sender.KafkaSender
}

func (o *OutboxRepo) CreateMessage(msg *sender.Message) {
	o.Mu.Lock()
	defer o.Mu.Unlock()
	o.Outbox[msg.AnswerID] = msg
}

func (o *OutboxRepo) OutboxProcessor(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			message := o.getMessageFromRepoToOutbox()
			if message != nil {
				err := o.Sender.SendMessage(message)
				if err != nil {
					o.rollback(message)
				} else {
					o.commit(message)
				}
			}
		}
	}
}

func (o *OutboxRepo) commit(a *sender.Message) {
	o.Mu.Lock()
	defer o.Mu.Unlock()
	a.IsProcessed = true
	a.ProcessedInOB = time.Now()
	o.Outbox[a.AnswerID] = a
}

func (o *OutboxRepo) rollback(a *sender.Message) {
	o.Mu.Lock()
	defer o.Mu.Unlock()
	a.IsAquired = false
	a.RetryCount++
	o.Outbox[a.AnswerID] = a
}

func (o *OutboxRepo) getMessageFromRepoToOutbox() *sender.Message {
	o.Mu.Lock()
	defer o.Mu.Unlock()
	for _, value := range o.Outbox {
		if !value.IsAquired && !value.IsProcessed && value.RetryCount < 5 {
			value.IsAquired = true
			value.RetryCount++
			return value
		}
	}
	return nil
}
