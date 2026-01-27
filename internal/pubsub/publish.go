package pubsub

import (
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishJSON[T any](
	ch *amqp.Channel,
	exchangeName string,
	key string,
	msg T,
) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	ch.PublishWithContext(context.Background(), exchangeName, key, false, false, amqp.Publishing{ContentType: "application/json", Body: data})

	return nil
}
