package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"net/http"
)

type Consumer struct {
	Conn      *amqp.Connection
	QueueName string
}

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		Conn: conn,
	}

	if err := consumer.setup(); err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

func (c *Consumer) setup() error {
	channel, err := c.Conn.Channel()
	if err != nil {
		return err
	}
	return declareExchange(channel)
}

func (c *Consumer) Listen(topics []string) error {
	ch, err := c.Conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	q, err := declareRandomQueue(ch)
	if err != nil {
		return err
	}
	for _, s := range topics {
		ch.QueueBind(q.Name, s, "logs_topic", false, nil)
		if err != nil {
			return err
		}
	}
	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}
	forever := make(chan bool)
	go func() {
		for d := range messages {
			var payload Payload
			_ = json.Unmarshal(d.Body, &payload)

			go handlePayload(payload)
		}
	}()
	fmt.Printf("ожидание сообщения [exchange, queue] [logs_topic, %s]\n", q.Name)
	<-forever
	return nil
}

func handlePayload(payload Payload) {
	switch payload.Name {
	case "log", "event":
		if err := logEvent(payload); err != nil {
			log.Println(err)
		}
	case "auth":
	default:
		if err := logEvent(payload); err != nil {
			log.Println(err)
		}
	}
}

func logEvent(entry Payload) error {
	jsonData, _ := json.MarshalIndent(entry, "", "\t")
	request, err := http.NewRequest("POST", "http://logger-service/log", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusAccepted {
		return err
	}
	return nil
}
