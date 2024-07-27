package usecases

import (
	"fmt"

	"github.com/IBM/sarama"
)

type ConsumerUsecase struct{}

func (c *ConsumerUsecase) Run(message sarama.ConsumerMessage) {
	fmt.Printf("Consumer Usecase Received Message - Topic: %q - Value: %s\n", message.Topic, message.Value)
}
