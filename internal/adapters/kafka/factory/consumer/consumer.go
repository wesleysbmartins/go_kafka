package factory

import (
	"fmt"
	"go_kafka/internal/adapters/kafka"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
)

type ConsumerFactory struct {
	topic     string
	partition int32
	instance  sarama.Consumer
}

type IConsumerFactory interface {
	Create(topic string, partition int32)
	Listen() error
}

func (c *ConsumerFactory) Create(topic string, partition int32) {
	c.topic = topic
	c.partition = partition
	consumer, err := sarama.NewConsumerFromClient(kafka.Client)

	if err != nil {
		panic("Error to create new Consumer!")
	} else {
		c.instance = consumer
	}
}

func (c *ConsumerFactory) Listen(handler IConsumerHandler) {
	consumerPartition, err := c.instance.ConsumePartition(c.topic, c.partition, sarama.OffsetOldest)

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
			handler.Run(*msg)
		}
	}
}
