package kafka

import (
	"time"

	"github.com/IBM/sarama"
)

type Consumer struct {
	SingleConsumer sarama.Consumer
}

func NewConsumer(brokers []string) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = false
	//Так как не требуется точного контроля над смещенями в связи с довольно простой логикой обработки,
	// а также для упрощения кода было решено использовать autocommit
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer, err := sarama.NewConsumer(brokers, config)

	if err != nil {
		return nil, err
	}

	return &Consumer{
		SingleConsumer: consumer,
	}, err
}
