package pubsub

import (
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchangeName string,
	queueName string,
	key string,
	queueType SimpleQueueType,
	handler func(T),
) error {

	ch, _, err := DeclareAndBind(conn, exchangeName, queueName, key, queueType)
	if err != nil {
		return err
	}

	deliveryCh, err := ch.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		amqp.Table{},
	)
	if err != nil {
		return fmt.Errorf("error: could not create delivery channel for queue %s, %v", queueName, err)
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
