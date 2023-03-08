package main

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/mhrocha1997/tax_calculation/internal/infra/database"
	"github.com/mhrocha1997/tax_calculation/internal/usecase"
	"github.com/mhrocha1997/tax_calculation/pkg/kafka"

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

	kafkaWorker(msgChanKafka, usecase)
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
