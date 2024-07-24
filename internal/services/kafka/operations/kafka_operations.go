package operations

import (
	"encoding/json"
	"fmt"
	kafka_service "go_kafka/internal/services/kafka"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
)

type KafkaOperations struct{}

type IKafkaOperations interface {
	Send(topic string, key string, message interface{}) error
	Listen(topic string, partitionNumber int32)
}

func (o *KafkaOperations) Send(topic string, key string, message interface{}) error {

	value, _ := json.Marshal(message)

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(value),
	}

	_, _, err := kafka_service.Producer.SendMessage(msg)

	if err != nil {
		fmt.Println("Send message ERROR: ", err)
	} else {
		fmt.Println("Send message SUCCESS!")
	}

	return err
}

func (o *KafkaOperations) Listen(topic string, partitionNumber int32) {

	consumer, err := kafka_service.Consumer.ConsumePartition(topic, partitionNumber, sarama.OffsetOldest)

	if err != nil {
		fmt.Println("Consumer ERROR: ", err)
	}

	defer consumer.Close()

	ch := make(chan os.Signal)

	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	messages := consumer.Messages()

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
