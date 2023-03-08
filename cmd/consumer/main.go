package main

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/mhrocha1997/tax_calculation/internal/infra/database"
	"github.com/mhrocha1997/tax_calculation/internal/usecase"
	"github.com/mhrocha1997/tax_calculation/pkg/kafka"
	"github.com/mhrocha1997/tax_calculation/pkg/rabbitmq"
	"github.com/rabbitmq/amqp091-go"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./orders.db")

	if err != nil {
		panic(err)
	}
	defer db.Close()

	repository := database.NewOrderRepository(db)

	usecase := usecase.CalculateFinalPrice{OrderRepository: repository}

	msgChanKafka := make(chan *ckafka.Message)

	topics := []string{"orders"}
	servers := "host.docker.internal:9094"

	go kafka.Consume(topics, servers, msgChanKafka)

	go kafkaWorker(msgChanKafka, usecase)

	ch, err := rabbitmq.OpenChannel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	msgRabbitmqChannel := make(chan amqp091.Delivery)
	go rabbitmq.Consume(ch, msgRabbitmqChannel)
	rabbitmqWorker(msgRabbitmqChannel, usecase)
}

func kafkaWorker(msgChan chan *ckafka.Message, uc usecase.CalculateFinalPrice) {
	fmt.Println("Kafka worker has started")
	for msg := range msgChan {

		var OrderInputDTO usecase.OrderInputDTO

		err := json.Unmarshal(msg.Value, &OrderInputDTO)
		if err != nil {
			panic(err)
		}
		outputDto, err := uc.Execute(&OrderInputDTO)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Kafka has processed order %s\n", outputDto.ID)
	}

}

func rabbitmqWorker(msgChan chan amqp091.Delivery, uc usecase.CalculateFinalPrice) {
	fmt.Println("Rabbitmq worker has started")
	for msg := range msgChan {

		var OrderInputDTO usecase.OrderInputDTO

		err := json.Unmarshal(msg.Body, &OrderInputDTO)
		if err != nil {
			panic(err)
		}
		outputDto, err := uc.Execute(&OrderInputDTO)
		if err != nil {
			panic(err)
		}
		msg.Ack(false)
		fmt.Printf("Rabbitmq has processed order %s	", outputDto.ID)
	}
}
