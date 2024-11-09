package internal

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"math"
	"os"
	"time"
)

type Config struct {
	Rabbit *amqp.Connection
}

func Connect() (*amqp.Connection, error) {
	var count int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection
	for {
		dsn := os.Getenv("DSN_RABBITMQ")
		c, err := amqp.Dial(dsn)
		if err != nil {
			fmt.Println("RabbitMQ ещё не готов...")
			count++
		} else {
			log.Println("Успешное подключение к RabbitMQ")
			connection = c
			break
		}
		if count > 5 {
			fmt.Println(err)
			return nil, err
		}
		backOff = time.Duration(math.Pow(float64(count), 2)) * time.Second
		log.Println("Задержка...")
		time.Sleep(backOff)
	}
	return connection, nil
}
