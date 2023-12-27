package rabbitmq

import (
	"context"
	"encoding/json"
	"example/tes-websocket/config"
	"example/tes-websocket/database"
	"time"

	// "example/tes-websocket/internal/ws"
	"example/tes-websocket/utils"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Broker struct {
	PublisherQueue 			amqp091.Queue
	ReceiverQueue 			amqp091.Queue
	Channel        			*amqp091.Channel
	PublisherExchange		string
	ReceiverExchange		string
	PublisherRoutingKey string
	ReceiverRoutingKey  string
	ExchangeType	 			string
}

type Message struct {
	MessageID	 uuid.UUID `json:"messageid" bson:"messageid"`
	From			 string	`json:"from" bson:"from"`
	To				 string `json:"to" bson:"to"`
	Message    string `json:"message" bson:"message"`
	Timestamp	 int64 `json:"timestamp" bson:"timestamp"`
}

// type MessageDB struct {
// 	ID         string `json:"_id" bson:"_id"`
// 	MessageID  uuid.UUID `json:"messageid" bson:"messageid"`
// 	From			 string	`json:"from" bson:"from"`
// 	To				 string `json:"to" bson:"to"`
// 	Message    string `json:"message" bson:"message"`
// }

func (b *Broker) SetUp(channel *amqp091.Channel) error {
	exchangeNameSend := config.EnvExchangeNameSend()
	exchangeNameReceive := config.EnvExchangeNameReceive()
	exchangeType := config.EnvExchangeType()
	queueNameSend := config.EnvQueueNameSend()
	queueNameReceive := config.EnvQueueNameReceive()
	routingKeySend := config.EnvRoutingKeySend()
	routingKeyReceive := config.EnvRoutingKeyReceive()

	if queueNameSend == "" {
		return fmt.Errorf("queue name send is not set in environment variables")
	}

	if queueNameReceive == "" {
		return fmt.Errorf("queue name receive is not set in environment variables")
	}

	if routingKeySend == "" {
		return fmt.Errorf("routing key send is not set in environment variables")
	}

	if routingKeyReceive == "" {
		return fmt.Errorf("routing key receive is not set in environment variables")
	}

	err := channel.ExchangeDeclare(
		exchangeNameSend, // Name of the exchange
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
	fmt.Printf("Exchange %s declared\n", exchangeNameSend)

	err = channel.ExchangeDeclare(
		exchangeNameReceive, // Name of the exchange
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
	fmt.Printf("Exchange %s declared\n", exchangeNameReceive)

	queueSend, err := channel.QueueDeclare(
		queueNameSend, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	if err != nil {
		return fmt.Errorf("failed to declare a queue: %v", err)
	}
	fmt.Printf("Queue %s declared\n", queueNameSend)

	queueReceive, err := channel.QueueDeclare(
		queueNameReceive, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	if err != nil {
		return fmt.Errorf("failed to declare a queue: %v", err)
	}
	fmt.Printf("Queue %s declared\n", queueNameReceive)

	err = channel.QueueBind(
		queueNameSend,    // queue name
		routingKeySend,   // routing key
		exchangeNameSend, // exchange name
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to bind a queue: %v", err)
	}
	fmt.Println("Queue Bound")

	err = channel.QueueBind(
		queueNameReceive,    // queue name
		routingKeyReceive,   // routing key
		exchangeNameReceive, // exchange name
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to bind a queue: %v", err)
	}
	fmt.Println("Queue Bound")

	b.PublisherQueue = queueSend
	b.ReceiverQueue = queueReceive
	b.Channel = channel
	b.PublisherExchange = exchangeNameSend
	b.ReceiverExchange = exchangeNameReceive
	b.PublisherRoutingKey = routingKeySend
	b.ReceiverRoutingKey = routingKeyReceive
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

func (b *Broker) SendMessageToQueueDB(c *gin.Context) {
	var message1 Message
	message1.MessageID = utils.GenerateUUIDs()
	message1.Timestamp = utils.CurrentTime()
	c.ShouldBind(&message1)

	msg, _ := json.Marshal(message1)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

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

	err := b.Channel.PublishWithContext(ctx,
		b.PublisherExchange,// "chat_exchange",
		b.PublisherRoutingKey,// "chat",
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

	queueName := config.EnvQueueNameSend()

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

			var result Message

			filter := bson.D{{"messageid", receivedMessage.MessageID}}

			err = collection.FindOne(context.TODO(), filter).Decode(&result)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					return
				}
				panic(err)
			}
			fmt.Printf("Found document: %+v\n", result)

			jsonData, err := json.Marshal(result)
			if err != nil {
				fmt.Println("Error marshaling JSON:", err)
				return
			}

			message := amqp091.Publishing{
				ContentType: "application/json",
				Body: jsonData,
			}

			err = b.Channel.PublishWithContext(ctx,
				b.ReceiverExchange,// "chat_exchange",
				b.ReceiverRoutingKey,// "chat",
				false,
				false,
				message,
			)
			if err != nil {
				panic(err)
			}
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
