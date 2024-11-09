package main

import (
	"github.com/magomedcoder/simple-microservice/listener-service/internal"
	"github.com/magomedcoder/simple-microservice/listener-service/internal/event"
	"log"
	"os"
)

func main() {
	rabbitConn, err := internal.Connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()
	log.Println("Прослушивание и потребление сообщений из RabbitMQ...")
	consumer, err := event.NewConsumer(rabbitConn)
	if err != nil {
		panic(err)
	}
	if err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"}); err != nil {
		log.Println(err)
	}
}
