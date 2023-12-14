package rabbitmq

import (
	"context"
	"encoding/json"
	"example/tes-websocket/config"
	"example/tes-websocket/database"
	"time"

	// "example/tes-websocket/internal/ws"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/rabbitmq/amqp091-go"
	// "go.mongodb.org/mongo-driver/bson"
)

type Broker struct {
	PublisherQueue amqp091.Queue
	Channel        *amqp091.Channel
	Exchange       string
	RoutingKey     string
	ExchangeType	 string
}

type Message struct {
	SenderID   string `json:"senderID"`
	ReceiverID string `json:"receiverID"`
	Message    string `json:"message"`
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
	b.ExchangeType = exchangeType

	return nil
}

// func (b *Broker) PublishMessage(c *gin.Context) {
// 	var message Message
// 	c.ShouldBind(message)

// 	body, err := json.Marshal(message)
// 	if err != nil {
// 		// return fmt.Errorf("error marshaling message: %s", err)
// 		panic(err)
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	err = b.Channel.ExchangeDeclare(
// 		"chat_exchange",
// 		"direct",
// 		false,
// 		false,
// 		false,
// 		false,
// 		nil,
// 	)
// 	if err != nil {
// 		panic(err)
// 	}

// 	err = b.Channel.PublishWithContext(ctx,
// 		"chat_exchange",//b.Exchange,   // exchange
// 		"chat",// b.RoutingKey, // routing key
// 		false,        // mandatory
// 		false,        // immediate
// 		amqp091.Publishing{
// 			ContentType: "application/json",
// 			Body:        body,
// 		})

// 	cancel()

// 	if err != nil {
// 		// return fmt.Errorf("PublishMessage Error occurred: %s", err)
// 		panic(err)
// 	}

// 	// return nil
// }

func (b *Broker) SendMessage(c *gin.Context) {
	var message1 Message
	c.ShouldBind(&message1)

	msg, _ := json.Marshal(message1)

	message := amqp091.Publishing{
		ContentType: "application/json",
		Body:        msg,
	}

	// err := b.Channel.ExchangeDeclare(
	// 	b.Exchange,// "chat_exchange",
	// 	b.ExchangeType,// "direct",
	// 	false,
	// 	false,
	// 	false,
	// 	false,
	// 	nil,
	// )
	// if err != nil {
	// 	panic(err)
	// }

	err := b.Channel.Publish(
		b.Exchange,// "chat_exchange",
		b.RoutingKey,// "chat",
		false,
		false,
		message,
	)
	if err != nil {
		panic(err)
	}
}

func (b *Broker) ConsumeMessage() error {
	if b.Channel == nil {
		return fmt.Errorf("channel is not set up for consuming messages")
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := database.ConnectDB()
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	defer db.Disconnect(ctx)

	collection := database.GetCollection(db, "userss")

	queueName := config.EnvQueueName()

	// Set up the consumer
	msgs, err := b.Channel.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)

	if err != nil {
		return fmt.Errorf("failed to register a consumer: %v", err)
	}

	var forever chan struct{}
	// Start consuming messages
	go func() {
		for msg := range msgs {
			// Process the received message
			fmt.Printf("Received a message: %s\n", msg.Body)

			// Unmarshal the message body into the Message struct
			var receivedMessage Message
			err := json.Unmarshal(msg.Body, &receivedMessage)
			if err != nil {
				fmt.Printf("Error unmarshaling message: %v\n", err)
				continue
			}

			_, err = collection.InsertOne(context.TODO(), receivedMessage)
			if err != nil {
				panic(err)
			}
			// Add your custom message processing logic here
		}
	}()

	fmt.Printf("Started consuming messages from queue %s\n", queueName)
	<-forever
	return nil

	// var forever chan struct{}

	// go func() {
	// 	for d := range msgs {
	// 		var order OrderHeader
	// 		json.Unmarshal(d.Body, &order)
	// 		go sendEmailNotification(order.CustomerEmail)
	// 	}
	// }()
	// log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	// <-forever
}
