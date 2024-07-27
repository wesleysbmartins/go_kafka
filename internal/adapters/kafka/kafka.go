package kafka

import (
	"fmt"
	"time"

	"github.com/IBM/sarama"
)

type Kafka struct {
	host string
	port string
}

type IKafka interface {
	Connect()
}

var Client sarama.Client

func (k *Kafka) Connect() {
	if Client == nil {
		credentials := Kafka{
			host: "localhost",
			port: "9092",
		}

		config := sarama.NewConfig()
		config.Version = sarama.MaxVersion
		config.Producer.Return.Successes = true
		config.Consumer.Return.Errors = true
		config.Consumer.Offsets.AutoCommit.Enable = true
		config.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second

		broker := []string{fmt.Sprintf("%s:%s", credentials.host, credentials.port)}

		client, err := sarama.NewClient(broker, config)

		if err != nil {
			fmt.Println("Client Intance Error!\n", err)
		} else {
			Client = client
		}
	}
}
