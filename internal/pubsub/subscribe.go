package pubsub

import (
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type AckType int

const (
	Ack AckType = iota
	NackRequeue
	NackDiscard
)

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchangeName string,
	queueName string,
	key string,
	queueType SimpleQueueType,
	handler func(T) AckType,
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
			ackType := handler(msg)
			switch ackType {
			case Ack:
				d.Ack(false)
			case NackRequeue:
				d.Nack(false, true)
			case NackDiscard:
				d.Nack(false, false)
			}
		}
	}
	go receiveMessages()

	return nil

}
