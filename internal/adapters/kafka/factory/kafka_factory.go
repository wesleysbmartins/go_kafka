package factory

import (
	"go_kafka/internal/adapters/kafka"

	"github.com/IBM/sarama"
)

type KafkaFactory struct{}

type IKafkaFactory interface {
	CreateProducer() (sarama.SyncProducer, error)
	CreateConsumer() (sarama.Consumer, error)
	CreateConsumerGroup() (sarama.ConsumerGroup, error)
}

func (o *KafkaFactory) CreateProducer() (sarama.SyncProducer, error) {
	producer, err := sarama.NewSyncProducerFromClient(kafka.Client)
	return producer, err
}

func (o *KafkaFactory) CreateConsumer() (sarama.Consumer, error) {
	consumer, err := sarama.NewConsumerFromClient(kafka.Client)
	return consumer, err
}

func (o *KafkaFactory) CreateConsumerGroup(groupId string) (sarama.ConsumerGroup, error) {
	consumerGroup, err := sarama.NewConsumerGroupFromClient(groupId, kafka.Client)
	return consumerGroup, err
}
