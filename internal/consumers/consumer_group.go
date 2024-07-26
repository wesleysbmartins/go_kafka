package consumers

import (
	"context"
	"fmt"
	"go_kafka/internal/consumers/handlers"
	"go_kafka/internal/services/kafka/operations"
	"go_kafka/internal/topics"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
)

type ConsumerGroup struct{}

type IConsumerGroup interface {
	create() error
	Listen()
}

var allTopics = []string{topics.Activity, "other-topic"}

const groupId = "activity-group"

var consumerGroup sarama.ConsumerGroup

func (c *ConsumerGroup) create() error {
	var err error

	if err == nil {
		operations := &operations.KafkaOperations{}
		consumerGroup, err = operations.CreateConsumerGroup(groupId)
		if err != nil {
			fmt.Println("Create Consumer Group Error!\n", err)
		}
	}

	return err
}

func (c *ConsumerGroup) Listen() {
	c.create()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	handler := handlers.ConsumerGroupHandlers{}

	go func() {
		for {

			if err := consumerGroup.Consume(ctx, allTopics, handler); err != nil {
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
