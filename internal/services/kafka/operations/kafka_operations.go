package operations

import (
	kafka_service "go_kafka/internal/services/kafka"

	"github.com/IBM/sarama"
)

type KafkaOperations struct{}

type IKafkaOperations interface {
	CreateProducer() (sarama.SyncProducer, error)
	CreateConsumer() (sarama.Consumer, error)
	CreateConsumerGroup() (sarama.ConsumerGroup, error)
}

func (o *KafkaOperations) CreateProducer() (sarama.SyncProducer, error) {
	producer, err := sarama.NewSyncProducerFromClient(kafka_service.Client)
	return producer, err
}

func (o *KafkaOperations) CreateConsumer() (sarama.Consumer, error) {
	consumer, err := sarama.NewConsumerFromClient(kafka_service.Client)
	return consumer, err
}

func (o *KafkaOperations) CreateConsumerGroup(groupId string) (sarama.ConsumerGroup, error) {
	consumerGroup, err := sarama.NewConsumerGroupFromClient(groupId, kafka_service.Client)
	return consumerGroup, err
}
