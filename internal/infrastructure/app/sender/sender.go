//go:generate mockgen -source ./sender.go -destination=./mocks/sender.go -package=mock_sender

	package sender

	import (
		"HOMEWORK-1/internal/infrastructure/kafka"
		"encoding/json"
		"fmt"
		"time"

		"github.com/IBM/sarama"
	)

	type Message struct {
		AnswerID      int
		Command       string
		Args          []string
		Success       bool
		CreatedAt     time.Time
		ProcessedInOB time.Time

		IsAquired   bool
		IsProcessed bool

		RetryCount uint
	}

	type KafkaSender struct {
		producer *kafka.Producer
		topic    string
	}

	func NewKafkaSender(producer *kafka.Producer, topic string) *KafkaSender {
		return &KafkaSender{
			producer,
			topic,
		}
	}

	func (s *KafkaSender) SendMessage(message *Message) error {
		kafkaMsg, err := s.buildMessage(*message)
		if err != nil {
			fmt.Println("Send message marshal error", err)
			return err
		}

		_, _, err = s.producer.SendSyncMessage(kafkaMsg)
		
		if err != nil {
			fmt.Println("Send message connector error", err)
			s.updateMessageStatus(message, false)
			return err
		}

		s.updateMessageStatus(message, true)

		return nil
	}

	func (s *KafkaSender) buildMessage(message Message) (*sarama.ProducerMessage, error) {
		msg, err := json.Marshal(message)

		if err != nil {
			fmt.Println("Send message marshal error", err)
			return nil, err
		}

		return &sarama.ProducerMessage{
			Topic:     s.topic,
			Value:     sarama.ByteEncoder(msg),
			Partition: -1,
			Key:       sarama.StringEncoder(fmt.Sprint(message.AnswerID)),
			Headers: []sarama.RecordHeader{ // например, в хедер можно записать версию релиза
				{
					Key:   []byte("test-header"),
					Value: []byte("test-value"),
				},
			},
		}, nil
	}

	func (s *KafkaSender) updateMessageStatus(message *Message, success bool) {
		message.Success = success
		message.ProcessedInOB = time.Now()
	}
