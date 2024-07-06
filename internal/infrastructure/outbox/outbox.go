package outbox

import (
	"HOMEWORK-1/internal/infrastructure/app/sender"
	"HOMEWORK-1/internal/infrastructure/kafka"
	"context"
	"fmt"
	"sync"
	"time"
)

var Brokers = []string{
	"127.0.0.1:9091",
	"127.0.0.1:9092",
	"127.0.0.1:9093",
}

type OutboxRepo struct {
	Mu     sync.RWMutex
	Outbox map[int]*sender.Message
}

func (o *OutboxRepo) CreateMessage(msg *sender.Message) {
	o.Mu.Lock()
	defer o.Mu.Unlock()
	msg.CreatedAt = time.Now()
	o.Outbox[msg.AnswerID] = msg
}

func (o *OutboxRepo) OutboxProcessor(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			message := o.getMessageFromRepoToOutbox()
			if message == nil {
				continue
			}

			err := producerExample(Brokers, message)
			if err != nil {
				o.rollback(message)
			} else {
				o.commit(message)
			}
		}
	}
}

func (o *OutboxRepo) commit(a *sender.Message) {
	o.Mu.Lock()
	defer o.Mu.Unlock()
	a.IsProcessed = true
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

func producerExample(brokers []string, payment *sender.Message) error {
	kafkaProducer, err := kafka.NewProducer(brokers)
	if err != nil {
		fmt.Println(err)
		return err
	}

	producer := sender.NewKafkaSender(kafkaProducer, "my-topic")
	err = producer.SendMessage(payment)
	if err != nil {
		return err
	}
	err = kafkaProducer.Close()
	if err != nil {
		fmt.Println("Close producers error ", err)
		return err
	}
	return nil
}
