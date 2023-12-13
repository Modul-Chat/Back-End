package rabbitmq

import (
	"example/tes-websocket/config"
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

var br Broker

func InitBroker() (*amqp091.Connection, *amqp091.Channel, error) {
	rabbitMQURI := config.EnvRabbitMQURI()
	if rabbitMQURI == "" {
		return nil, nil, fmt.Errorf("RabbitMQ_URI not set in environment variables")
	}

	mqConn, err := amqp091.Dial(rabbitMQURI)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}
	fmt.Println("Connected to RabbitMQ")

	channel, err := mqConn.Channel()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open a channel: %v", err)
	}
	fmt.Println("Channel created")

	err = br.SetUp(channel)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to set up RabbitMQ: %v", err)
	}

	// err = br.PublishMessage()
	// if err != nil {
	// 	return nil, nil, fmt.Errorf("failed to publish message to RabbitMQ: %v", err)
	// }

	err = br.ConsumeMessage()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to consume message: %v", err)
	}

	return mqConn, channel, nil
}

func GetRabbitMQBroker() *Broker {
	return &br
}
