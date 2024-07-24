# Golang Kafka Integration
Este é um exemplo de como uma aplicação Golang pode trabalhar integrada ao Kafka.

## Kafka
Apache Kafka é uma plataforma de streaming de eventos distribuída e altamente escalável, projetada para processar e gerenciar grandes volumes de dados em tempo real. Originalmente desenvolvido pela LinkedIn e posteriormente aberto como um projeto de código aberto pela Apache Software Foundation, o Kafka é amplamente utilizado em diversas indústrias para várias aplicações.

> **OBS:** Para entender melhor o que é o Kafka, como instalar e rodar em seu ambiente, acesse o meu [repositório](https://github.com/wesleysbmartins/kafka) onde registrei meus estudos sobre o tema.

## Implementação
Para utilizar o Kafka em uma aplicação Go é necessário usar bibliotecas como o [Sarama](https://pkg.go.dev/github.com/IBM/sarama). Existem outras abordagens para utilizar mas, eu estou usando esta.

Após iniciar sua aplicação Go e instalar a biblioteca do Sarama, você pode criar sua conexão, criando seu Producer ou seu Consumer.

Utilizei abordagens de padrões de projeto como **Singleton**, que será útil economizando recursos e melhorando a performance da aplicação e da comunicação com o Kafka.

### Producer
Entidades que publicam dados (mensagens) em tópicos do Kafka.
```go
package kafka

import (
	"fmt"

	"github.com/IBM/sarama"
)

// singleton
var Producer sarama.SyncProducer

func Connect() {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

    // credenciais como host e porta onde seu kafka esta disponível
	addr := []string{fmt.Sprintf("%s:%s", "localhost", "9092")}

	producer, err := sarama.NewSyncProducer(addr, config)

	if err != nil {
		fmt.Println("Producer instance ERROR: ", err)
	} else {
		Producer = producer
	}
}
```

### Consumer
Entidades que lêem e processam dados de tópicos do Kafka.
```go
package kafka

import (
	"fmt"

	"github.com/IBM/sarama"
)

var Consumer sarama.Consumer

func Connect() {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	addr := []string{fmt.Sprintf("%s:%s", "localhost", "9092")}

	consumer, err := sarama.NewConsumer(addr, config)

	if err != nil {
		fmt.Println("Cconsumer instance ERROR: ", err)
	} else {
		Consumer = consumer
	}
}
```

### Operations
Implementação de funções para enviar e ler mensagens dos tópicos Kafka.
### Send
Função Send que espera parametros como o tópico a ser enviado mensagens, a chave e o valor da mensagem.
a mensagem pode ser uma struct.
```go
package operations

import (
	"encoding/json"
	"fmt"
	kafka "go_kafka/internal/services/kafka"
	"os"

	"github.com/IBM/sarama"
)

func Send(topic string, key string, message interface{}) error {

	value, _ := json.Marshal(message)

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(value),
	}

	_, _, err := kafka.Producer.SendMessage(msg)

	if err != nil {
		fmt.Println("Send message ERROR: ", err)
	} else {
		fmt.Println("Send message SUCCESS!")
	}

	return err
}
```

### Listen
Função Listen que espera como parametro o nome e número/posição do tópico em que deverá ler as mensagens.
```go
package operations

import (
	"fmt"
	kafka "go_kafka/internal/services/kafka"
	"os"
	"os/signal"
	"syscall"
)

func Listen(topic string, partitionNumber int32) {

	consumer, err := kafka.Consumer.ConsumePartition(topic, partitionNumber, sarama.OffsetOldest)

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
```

### Entidade
A entidade de exemplo que representa as mensagens a serem enviadas e lidas.
```go
package entities

type Activity struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
```

### Main
Utilizando código implementado.
```go
package main

import (
	"go_kafka/internal/entities"
	kafka "go_kafka/internal/services/kafka"
	"go_kafka/internal/services/kafka/operations"
)

func init() {
	kafka.Connect()
}

const topic = "activity-topic"
const key = "user-activity"
const position = 0

func main() {

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
		operation.Send(topic, key, activity)
	}

	operations.Listen(topic, position)

}
```

### Resultado
Executando o main, você deve obter o seguinte resultado:
```shell
Send message SUCCESS!
Send message SUCCESS!
Send message SUCCESS!
Received message - Topic: activity-topic - Value: {"id":"1","title":"Example 01","description":"KAFKA IMPLEMENTATION DESCRIPTION 01"}
Received message - Topic: activity-topic - Value: {"id":"2","title":"Example 02","description":"KAFKA IMPLEMENTATION DESCRIPTION 02"}
Received message - Topic: activity-topic - Value: {"id":"3","title":"Example 03","description":"KAFKA IMPLEMENTATION DESCRIPTION 03"}
```

Para parar a aplicação basta pressionar **CTRL + C** em seu terminal:
```shell
Reveived message SIGNALL: interrupt
```

Assim voce tem uma aplicação Golang integrada ao Kafka, sendo capaz de realizar operações de envio ou leitura de mensagens ou eventos.
