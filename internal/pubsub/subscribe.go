package pubsub

import (
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

func subscribeJSON[T any](
	conn *amqp.Connection,
	exchangeName string,
	queueName string,
	key string,
	queueType SimpleQueueType,
	handler func(T),
) error {
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	ch, _, err = DeclareAndBind(conn, exchangeName, queueName, key, queueType)
	if err != nil {
		return err
	}

	deliveryCh, err := ch.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		amqp.Table{},
	)
	if err != nil {
		return err
	}

	receiveMessages := func() {
		for d := range deliveryCh {
			var msg T
			err := json.Unmarshal(d.Body, &msg)
			if err != nil {
				continue
			}
			handler(msg)
		}
	}
	go receiveMessages()

	return nil

}
