package kafka_service

import (
	"fmt"

	"github.com/IBM/sarama"
)

type KafkaService struct {
	host string
	port string
}

type IKafikaService interface {
	Connect()
}

var Producer sarama.SyncProducer

var Consumer sarama.Consumer

func (k *KafkaService) Connect() {
	credentials := KafkaService{
		host: "localhost",
		port: "9092",
	}

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Consumer.Return.Errors = true

	addr := []string{fmt.Sprintf("%s:%s", credentials.host, credentials.port)}

	consumer, err := sarama.NewConsumer(addr, config)

	if err != nil {
		fmt.Println("Cconsumer instance ERROR: ", err)
	} else {
		Consumer = consumer
	}

	producer, err := sarama.NewSyncProducer(addr, config)

	if err != nil {
		fmt.Println("Producer instance ERROR: ", err)
	} else {
		Producer = producer
	}
}
