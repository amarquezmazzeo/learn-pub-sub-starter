package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type SimpleQueueType int

const (
	Durable SimpleQueueType = iota
	Transient
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

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType,
) (*amqp.Channel, amqp.Queue, error) {
	chann, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("error creating channel: %e", err)
	}

	queue, err := chann.QueueDeclare(
		queueName,
		queueType == 0,
		queueType != 0,
		queueType != 0,
		false,
		nil)
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("error creating queue: %e", err)
	}

	err = chann.QueueBind(
		queueName,
		key,
		exchange,
		false, nil)
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("error binding queue to exchange: %e", err)
	}

	return chann, queue, nil
}
