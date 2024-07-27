# Golang [![My Skills](https://skillicons.dev/icons?i=golang)](https://skillicons.dev) Kafka Integration [![My Skills](https://skillicons.dev/icons?i=kafka)](https://skillicons.dev)
Este é um exemplo de como uma aplicação Golang pode trabalhar integrada ao Kafka.

## Kafka


Apache Kafka é uma plataforma de streaming de eventos distribuída e altamente escalável, projetada para processar e gerenciar grandes volumes de dados em tempo real. Originalmente desenvolvido pela LinkedIn e posteriormente aberto como um projeto de código aberto pela Apache Software Foundation, o Kafka é amplamente utilizado em diversas indústrias para várias aplicações.

> **OBS:** Para entender melhor o que é o Kafka, como instalar e rodar em seu ambiente, acesse o meu [repositório](https://github.com/wesleysbmartins/kafka) onde registrei meus estudos sobre o tema.

## Hands-On
Neste momento iremos abordar de forma simples como integra o Kafka a sua aplicação Golang usando a biblioteca [Sarama](https://pkg.go.dev/github.com/IBM/sarama), utilizando padrões de desenvolvimento como **Singleton**, **Factory**, **Entities** e **Usecases**.

Após iniciar sua aplicação Go e instalar a biblioteca do Sarama, você pode criar sua conexão, criando seu Producer ou seu Consumer.

### Client Kafka
Após inicar sua aplicação podemos iniciar o desenvolvimento da aplicação, neste primeiro momento vamos criar o client do Kafka para sua aplicação e instancia-lo utilizando um Singleton.

```go
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
```

### Factory Kafka
Para a criação dos Producers e Consumers do Kafka vamos utilizar a abordagem de Factory.

### Producer Factory
Para o factory dos producers retemos os métodos de criação do producer e o de envio de mensagens, onde esperam os parametros necessários para realizar tais operações como tópico, chave da mensagem e valor da mensagem.
```go
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
```

### Consumers
Temos dois tipos de implementações de consumers, um consumer de um único tópico e os consumers de grupos de tópicos.

Além disse foi criado interfaces de handlers para lidarem com as mensagens recebidas, o que representaria sua regra de negócio.

### Consumer Handlers Interfaces
```go
package factory

import "github.com/IBM/sarama"

type IConsumerHandler interface {
	Run(message sarama.ConsumerMessage)
}

type IConsumerGroupHandler interface {
	Setup(session sarama.ConsumerGroupSession) error
	Cleanup(session sarama.ConsumerGroupSession) error
	ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error
}
```

### Consumer Factory
```go
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
```

### Consumer Group Factory
```go
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
```

## Entities
Para exemplificar uma mensagem foi criado uma struct **Activity**.
```go
package entities

type Activity struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
```

### Usecases
Structs e métodos responsáveis por implementar as interfaces esperadas pelos consumers como handler e também, sua regra de negócio.

### Consumer Usecase
```go
package usecases

import (
	"fmt"

	"github.com/IBM/sarama"
)

type ConsumerUsecase struct{}

func (c *ConsumerUsecase) Run(message sarama.ConsumerMessage) {
	fmt.Printf("Consumer Usecase Received Message - Topic: %q - Value: %s\n", message.Topic, message.Value)
}
```

### Consumer Group Usecase
```go
package usecases

import (
	"fmt"

	"github.com/IBM/sarama"
)

type ConsumerGroupUsecase struct{}

func (ConsumerGroupUsecase) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (ConsumerGroupUsecase) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h ConsumerGroupUsecase) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {

		fmt.Printf("Consumer Group Usecase Received Message - Topic: %q - Value: %s\n", msg.Topic, msg.Value)

		session.MarkMessage(msg, "Readed")
	}

	return nil
}
```
### Main
Sendo assim, basta brincar com o código desenvolvido até aqui, use isto no seu main, execute e veja a mágica acontecer:
```go
package main

import (
	"go_kafka/internal/adapters/kafka"
	factory_consumer "go_kafka/internal/adapters/kafka/factory/consumer"
	factory_producer "go_kafka/internal/adapters/kafka/factory/producer"
	"go_kafka/internal/entities"
	"go_kafka/internal/usecases"
)

func init() {
	kafka_client := &kafka.Kafka{}
	kafka_client.Connect()
}

const (
	partition = 0
	groupId   = "activity-group"
	topic     = "activity-topic"
	key       = "activity"
)

func main() {
	activityProducer := &factory_producer.ProducerFactory{}
	activityProducer.Create(key, topic)

	// activityConsumer := &factory_consumer.ConsumerFactory{}
	// activityConsumer.Create(topic, partition)

	consumerGroup := &factory_consumer.ConsumerGroupFactory{}
	consumerGroup.Create(groupId, []string{topic})

	consumerGrouphandler := &usecases.ConsumerGroupUsecase{}

	consumerGroup.Listen(consumerGrouphandler)

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
		activityProducer.Send(activity)
	}

	// consumerUsecase := &usecases.ConsumerUsecase{}

	// activityConsumer.Listen(consumerUsecase)
}
```

### Resultado
Executando o projeto, você deve obter um resultado semelhante a este:
```shell
Consumer Group Usecase Received Message - Topic: "activity-topic" - Value: {"id":"1","title":"Example 01","description":"KAFKA IMPLEMENTATION DESCRIPTION 01"}
Consumer Group Usecase Received Message - Topic: "activity-topic" - Value: {"id":"2","title":"Example 02","description":"KAFKA IMPLEMENTATION DESCRIPTION 02"}
Consumer Group Usecase Received Message - Topic: "activity-topic" - Value: {"id":"3","title":"Example 03","description":"KAFKA IMPLEMENTATION DESCRIPTION 03"}
Consumer Group Usecase Received Message - Topic: "activity-topic" - Value: {"id":"1","title":"Example 01","description":"KAFKA IMPLEMENTATION DESCRIPTION 01"}
Consumer Group Usecase Received Message - Topic: "activity-topic" - Value: {"id":"2","title":"Example 02","description":"KAFKA IMPLEMENTATION DESCRIPTION 02"}
Consumer Group Usecase Received Message - Topic: "activity-topic" - Value: {"id":"3","title":"Example 03","description":"KAFKA IMPLEMENTATION DESCRIPTION 03"}
Signal Event: interrupt
Send message SUCCESS!
Send message SUCCESS!
Send message SUCCESS!
```

Para parar a aplicação basta pressionar **CTRL + C** em seu terminal:
```shell
Reveived message SIGNALL: interrupt
```

Assim voce tem uma aplicação Golang integrada ao Kafka, sendo capaz de realizar operações de envio ou leitura de mensagens ou eventos.
