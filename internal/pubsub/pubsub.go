package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	//Marshal the val to JSON bytes
	jsonBytes, err := json.Marshal(val)
	if err != nil {
		return err
	}

	return ch.PublishWithContext(
		context.Background(),
		exchange,
		key,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonBytes,
		},
	)
}

// an enum to represent "durable" or "transient"
type SimpleQueueType string

const (
	DurableQueue   SimpleQueueType = "durable"
	TransientQueue SimpleQueueType = "transient"
)

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
) (*amqp.Channel, amqp.Queue, error) {
	newChannel, err := conn.Channel()
	if err != nil {
		fmt.Printf("Failed to open a channel in DeclareAndBind: %s\n", err)
		return nil, amqp.Queue{}, err
	}

	newQueue, err := newChannel.QueueDeclare(queueName, queueType == DurableQueue, queueType == TransientQueue, queueType == TransientQueue, false, nil)
	if err != nil {
		fmt.Printf("Failed to declare a queue in DeclareAndBind: %s\n", err)
		return nil, amqp.Queue{}, err
	}

	// Bind the queue to the exchange
	err = newChannel.QueueBind(queueName, key, exchange, false, nil)
	if err != nil {
		fmt.Printf("Failed to bind queue to exchange in DeclareAndBind: %s\n", err)
		return nil, amqp.Queue{}, err
	}

	return newChannel, newQueue, nil
}
