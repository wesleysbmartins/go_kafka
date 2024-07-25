package producers

import (
	"encoding/json"
	"fmt"
	"go_kafka/internal/entities"
	"go_kafka/internal/services/kafka/operations"
	"go_kafka/internal/topics"

	"github.com/IBM/sarama"
)

type ActivityProducer struct{}

type IActivityProducer interface {
	create() error
	Send(activity entities.Activity) error
}

var producer sarama.SyncProducer

const (
	topic = topics.Activity
	key   = "activity"
)

func (p *ActivityProducer) create() error {
	var err error
	if producer == nil {
		operations := operations.KafkaOperations{}
		producer, err = operations.CreateProducer()

		if err != nil {
			fmt.Println("Create Producer to Activity Error!\n", err)
		}
	}

	return err
}

func (p *ActivityProducer) Send(activity entities.Activity) error {
	p.create()

	value, _ := json.Marshal(activity)

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(value),
	}

	_, _, err := producer.SendMessage(msg)

	if err != nil {
		fmt.Println("Send message ERROR: ", err)
	} else {
		fmt.Println("Send message SUCCESS!")
	}

	return err
}
