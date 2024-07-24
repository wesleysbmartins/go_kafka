package main

import (
	"go_kafka/internal/entities"
	kafka_service "go_kafka/internal/services/kafka"
	"go_kafka/internal/services/kafka/operations"
)

func init() {
	kafka := kafka_service.KafkaService{}
	kafka.Connect()
}

const topic = "activity-topic"
const key = "user-activity"
const position = 0

func main() {

	kafka := operations.KafkaOperations{}

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
		kafka.Send(topic, key, activity)
	}

	kafka.Listen(topic, position)

}
