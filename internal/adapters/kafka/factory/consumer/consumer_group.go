package factory

import (
	"context"
	"fmt"
	"go_kafka/internal/adapters/kafka"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
)

type ConsumerGroupFactory struct {
	groupId  string
	topics   []string
	instance sarama.ConsumerGroup
}

type IConsumerGroupfactory interface {
	Create(groupId string, topics []string)
	Listen(handler IConsumerGroupHandler)
}

func (c *ConsumerGroupFactory) Create(groupId string, topics []string) {
	c.groupId = groupId
	c.topics = topics
	consumerGroup, err := sarama.NewConsumerGroupFromClient(groupId, kafka.Client)
	if err != nil {
		panic("Error to create Consumer Group!")
	} else {
		c.instance = consumerGroup
	}
}

func (c *ConsumerGroupFactory) Listen(handler IConsumerGroupHandler) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for {

			if err := c.instance.Consume(ctx, c.topics, handler); err != nil {
				panic(fmt.Sprintf("Consumer Group Error!\n%s", err.Error()))
			}

			if ctx.Err() != nil {
				return
			}
		}
	}()

	ch := make(chan os.Signal, 1)

	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case <-ctx.Done():
		fmt.Println("Context Cancelled!")
	case signal := <-ch:
		fmt.Printf("Signal Event: %v\n", signal)
	}
}
