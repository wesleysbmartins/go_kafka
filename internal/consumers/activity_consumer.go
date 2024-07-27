package consumers

import (
	"fmt"
	"go_kafka/internal/adapters/kafka/factory"
	"go_kafka/internal/topics"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
)

type ActivityConsumer struct{}

type IActivityConsumer interface {
	create() error
	Linsten()
}

var consumer sarama.Consumer

const (
	topic     = topics.Activity
	partition = 0
)

func (c *ActivityConsumer) create() error {
	var err error

	if consumer == nil {
		factory := &factory.KafkaFactory{}
		consumer, err = factory.CreateConsumer()
		if err != nil {
			fmt.Println("Create Consumer to Activity Error!\n", err)
		}
	}

	return err
}

func (c *ActivityConsumer) Listen() {
	c.create()

	consumerPartition, err := consumer.ConsumePartition(topic, partition, sarama.OffsetOldest)

	if err != nil {
		fmt.Println("Consumer ERROR: ", err)
	}

	defer consumerPartition.Close()

	ch := make(chan os.Signal, 1)

	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	messages := consumerPartition.Messages()

	for {
		select {
		case msg := <-ch:
			fmt.Printf("Reveived message SIGNALL: %v\n", msg)
			return
		case msg := <-messages:
			fmt.Printf("Received message - Topic: %s - Value: %s\n", msg.Topic, msg.Value)
		}
	}
}
