package main

import (
	"go_kafka/internal/consumers"
	"go_kafka/internal/entities"
	"go_kafka/internal/producers"
	kafka_service "go_kafka/internal/services/kafka"
)

func init() {
	kafka := &kafka_service.KafkaService{}
	kafka.Connect()
}

func main() {
	activityProducer := &producers.ActivityProducer{}
	// activityConsumer := &consumers.ActivityConsumer{}
	consumerGroup := &consumers.ConsumerGroup{}

	consumerGroup.Listen()

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

	// activityConsumer.Listen()
}
