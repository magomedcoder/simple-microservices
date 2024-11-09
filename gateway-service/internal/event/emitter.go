package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type Emitter struct {
	Connection *amqp.Connection
}

func (e *Emitter) setup() error {
	channel, err := e.Connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()
	return declareExchange(channel)
}

func (e *Emitter) Push(event string, severity string) error {
	channel, err := e.Connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()
	log.Println("Отправка в канал...")
	err = channel.Publish(
		"logs_topic",
		severity,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(event),
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func NewEventEmitter(conn *amqp.Connection) (Emitter, error) {
	emitter := Emitter{
		Connection: conn,
	}
	if err := emitter.setup(); err != nil {
		return Emitter{}, err
	}

	return emitter, nil
}
