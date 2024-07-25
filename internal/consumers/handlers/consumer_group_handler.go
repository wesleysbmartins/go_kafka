package handlers

import (
	"fmt"

	"github.com/IBM/sarama"
)

type ConsumerGroupHandlers struct{}

func (ConsumerGroupHandlers) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (ConsumerGroupHandlers) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h ConsumerGroupHandlers) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {

		fmt.Printf("Received message - Topic: %q - Value: %s\n", msg.Topic, msg.Value)

		session.MarkMessage(msg, "Readed")
	}

	return nil
}
