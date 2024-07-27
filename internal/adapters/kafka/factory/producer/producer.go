package factory

import (
	"encoding/json"
	"fmt"
	"go_kafka/internal/adapters/kafka"

	"github.com/IBM/sarama"
)

type ProducerFactory struct {
	key      string
	topic    string
	instance sarama.SyncProducer
}

type IProducerFactory interface {
	Create(key string, topic string)
	SendMessage(message string) error
}

func (p *ProducerFactory) Create(key string, topic string) {
	p.key = key
	p.topic = topic
	producer, err := sarama.NewSyncProducerFromClient(kafka.Client)

	if err != nil {
		panic("Error to create new Producer!")
	} else {
		p.instance = producer
	}
}

func (p *ProducerFactory) Send(message interface{}) error {
	value, _ := json.Marshal(message)

	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.StringEncoder(p.key),
		Value: sarama.StringEncoder(value),
	}

	_, _, err := p.instance.SendMessage(msg)

	if err != nil {
		fmt.Println("Send message ERROR: ", err)
	} else {
		fmt.Println("Send message SUCCESS!")
	}

	return err
}
