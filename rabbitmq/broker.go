package rabbitmq

import (
	"context"
	"encoding/json"
	"example/tes-websocket/config"
	"example/tes-websocket/internal/ws"
	"fmt"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type Broker struct {
	PublisherQueue amqp091.Queue
	Channel *amqp091.Channel
	Exchange string
	RoutingKey string
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
	b.Exchange = exchangeName
	b.RoutingKey = routingKey

	return nil
}

func (b *Broker) PublishMessage(message *ws.Message) error {
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error marshaling message: %s", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = b.Channel.PublishWithContext(ctx,
		b.Exchange,             // exchange
		b.RoutingKey, 					// routing key
		false,                 	// mandatory
		false,                 	// immediate
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		})

	cancel()

	if err != nil {
		return fmt.Errorf("PublishMessage Error occurred: %s", err)
	}

	return nil
}