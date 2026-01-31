package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	valJSON, err := json.Marshal(val)
	if err != nil {
		return fmt.Errorf("error marshalling val msg to JSON: %e", err)
	}
	err = ch.PublishWithContext(context.Background(), exchange, key, false, false,
		amqp.Publishing{ContentType: "application/json", Body: valJSON})
	if err != nil {
		return fmt.Errorf("error publishing message to channel: %e", err)
	}
	return nil
}
