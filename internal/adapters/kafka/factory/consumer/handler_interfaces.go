package factory

import "github.com/IBM/sarama"

type IConsumerHandler interface {
	Run(message sarama.ConsumerMessage)
}

type IConsumerGroupHandler interface {
	Setup(session sarama.ConsumerGroupSession) error
	Cleanup(session sarama.ConsumerGroupSession) error
	ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error
}
