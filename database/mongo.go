package database

import (
	"context"
	"example/tes-websocket/config"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() (*mongo.Client, error) {
	mongoDBURI := config.EnvMongoURI()
	if mongoDBURI == "" {
		return nil, fmt.Errorf("MONGODB_URI not set in environment variables")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoDBURI))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}
	// defer client.Disconnect(ctx)

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	fmt.Println("Connected to MongoDB")

	return client, nil
}

func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	db := client.Database(config.EnvDatabase())
	collection := db.Collection(collectionName)
	
	return collection
}