package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func EnvRabbitMQURI() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return os.Getenv("RABBITMQ_URI")
}

func EnvExchangeNameSend() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return os.Getenv("EXCHANGE_NAME_SEND")
}

func EnvExchangeNameReceive() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return os.Getenv("EXCHANGE_NAME_RECEIVE")
}

func EnvExchangeType() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return os.Getenv("EXCHANGE_TYPE")
}

func EnvQueueNameSend() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return os.Getenv("QUEUE_NAME_SEND")
}

func EnvQueueNameReceive() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return os.Getenv("QUEUE_NAME_RECEIVE")
}

func EnvRoutingKeySend() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return os.Getenv("ROUTING_KEY_SEND")
}

func EnvRoutingKeyReceive() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return os.Getenv("ROUTING_KEY_RECEIVE")
}

func EnvMongoURI() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return os.Getenv("MONGO_URI")
}

func EnvDatabase() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return os.Getenv("DATABASE")
}