package main

import (
	"go_kafka/internal/adapters/kafka"
	factory_consumer "go_kafka/internal/adapters/kafka/factory/consumer"
	factory_producer "go_kafka/internal/adapters/kafka/factory/producer"
	"go_kafka/internal/entities"
	"go_kafka/internal/usecases"
)

func init() {
	kafka_client := &kafka.Kafka{}
	kafka_client.Connect()
}

const (
	partition = 0
	groupId   = "activity-group"
	topic     = "activity-topic"
	key       = "activity"
)

func main() {
	activityProducer := &factory_producer.ProducerFactory{}
	activityProducer.Create(key, topic)

	// activityConsumer := &factory_consumer.ConsumerFactory{}
	// activityConsumer.Create(topic, partition)

	consumerGroup := &factory_consumer.ConsumerGroupFactory{}
	consumerGroup.Create(groupId, []string{topic})

	consumerGrouphandler := &usecases.ConsumerGroupUsecase{}

	consumerGroup.Listen(consumerGrouphandler)

	activities := []entities.Activity{
		{
			Id:          "1",
			Title:       "Example 01",
			Description: "KAFKA IMPLEMENTATION DESCRIPTION 01",
		}, {
			Id:          "2",
			Title:       "Example 02",
			Description: "KAFKA IMPLEMENTATION DESCRIPTION 02",
		},
		{
			Id:          "3",
			Title:       "Example 03",
			Description: "KAFKA IMPLEMENTATION DESCRIPTION 03",
		},
	}

	for _, activity := range activities {
		activityProducer.Send(activity)
	}

	// consumerUsecase := &usecases.ConsumerUsecase{}

	// activityConsumer.Listen(consumerUsecase)
}
