package rabbitmq

import (
	"example/tes-websocket/config"
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

type Broker struct {
	PublisherQueue amqp091.Queue
	Channel *amqp091.Channel
}

func (b *Broker) SetUp(channel *amqp091.Channel) error {
	exchangeName := config.EnvExchangeName()
	exchangeType := config.EnvExchangeType()
	queueName := config.EnvQueueName()
	routingKey := config.EnvRoutingKey()

	if queueName == "" {
		return fmt.Errorf("queue name is not set in environment variables")
	}

	if routingKey == "" {
		return fmt.Errorf("routing key is not set in environment variables")
	}

	err := channel.ExchangeDeclare(
		exchangeName, // Name of the exchange
		exchangeType, // Type of the exchange: "direct", "fanout", "topic", etc.
		false,        // Durable
		false,        // AutoDelete
		false,        // Internal
		false,        // NoWait
		nil,          // Arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare an exchange: %v", err)
	}
	fmt.Printf("Exchange %s declared\n", exchangeName)

	queue, err := channel.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %v", err)
	}
	fmt.Printf("Queue %s declared\n", queueName)

	err = channel.QueueBind(
		queueName,    // queue name
		routingKey,   // routing key
		exchangeName, // exchange name
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to bind a queue: %v", err)
	}
	fmt.Println("Queue Bound")

	b.PublisherQueue = queue
	b.Channel = channel

	return nil
}