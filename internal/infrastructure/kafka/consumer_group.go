package kafka

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

type Message struct {
	AnswerID int
	Command  string
	Args     string
}

type ConsumerGroup struct {
	Ready chan *sarama.ConsumerMessage
}


func NewConsumerGroup() *ConsumerGroup {
	return &ConsumerGroup{
		Ready: make(chan *sarama.ConsumerMessage),
	}
}

func (consumer *ConsumerGroup) IsReady() <-chan *sarama.ConsumerMessage {
	return consumer.Ready
}

// Setup Начинаем новую сессию, до ConsumeClaim
func (consumer *ConsumerGroup) Setup(_ sarama.ConsumerGroupSession) error {
	close(consumer.Ready)

	return nil
}

// Cleanup завершает сессию, после того, как все ConsumeClaim завершатся
func (consumer *ConsumerGroup) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim читаем до тех пор пока сессия не завершилась
func (consumer *ConsumerGroup) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():

			pm := Message{}
			err := json.Unmarshal(message.Value, &pm)
			if err != nil {
				fmt.Println("Consumer group error", err)
			}

			log.Printf("Message claimed: value = %v, timestamp = %v, topic = %s",
				pm,
				message.Timestamp,
				message.Topic,
			)

			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
}

