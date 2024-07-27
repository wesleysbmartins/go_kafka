package usecases

import (
	"fmt"

	"github.com/IBM/sarama"
)

type ConsumerGroupUsecase struct{}

func (ConsumerGroupUsecase) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (ConsumerGroupUsecase) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h ConsumerGroupUsecase) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {

		fmt.Printf("Consumer Group Usecase Received Message - Topic: %q - Value: %s\n", msg.Topic, msg.Value)

		session.MarkMessage(msg, "Readed")
	}

	return nil
}
