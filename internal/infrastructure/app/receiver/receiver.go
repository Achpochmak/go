package receiver

import (
	"HOMEWORK-1/internal/infrastructure/kafka"
	"errors"
	"fmt"

	"github.com/IBM/sarama"
)

type HandleFunc func(message *sarama.ConsumerMessage)

type KafkaReceiver struct {
	consumer *kafka.Consumer
	handlers map[string]HandleFunc
}

func NewReceiver(brokers []string, handlers map[string]HandleFunc) (*KafkaReceiver, error) {
	consumer, err := kafka.NewConsumer(brokers)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания Kafka consumer: %w", err)
	}
	return &KafkaReceiver{
		consumer: consumer,
		handlers: handlers,
	}, nil
}

func (r *KafkaReceiver) Subscribe(topic string) error {
	handler, ok := r.handlers[topic]

	if !ok {
		return errors.New("can not find handler")
	}

	// получаем все партиции топика
	partitionList, err := r.consumer.SingleConsumer.Partitions(topic)

	if err != nil {
		return err
	}
	initialOffset := sarama.OffsetNewest

	for _, partition := range partitionList {
		pc, err := r.consumer.SingleConsumer.ConsumePartition(topic, partition, initialOffset)

		if err != nil {
			return err
		}

		go func(pc sarama.PartitionConsumer, partition int32) {
			for message := range pc.Messages() {
				handler(message)
			}
		}(pc, partition)
	}

	return nil
}
