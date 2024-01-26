package repository

import (
	"context"
	"fmt"
	"go-conversation/model"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

// Repository represents the database repository.
type Repository struct {
	client *mongo.Client
}

// NewRepository creates a new repository instance.
func NewRepository(client *mongo.Client) *Repository {
	return &Repository{client: client}
}

// CreateUser inserts a new user into the database.
func (r *Repository) CreateUser(ctx context.Context, newUser model.User) (string, error) {
	newUser.CreatedAt = time.Now()

	// Hash the user's password before storing it in the database
	hashedPassword, err := hashPassword(newUser.Password)
	if err != nil {
		return "", fmt.Errorf("error hashing password: %v", err)
	}

	newUser.Password = hashedPassword

	result, err := r.client.Database("ChatGo").Collection("users").InsertOne(ctx, newUser)
	if err != nil {
		return "", fmt.Errorf("error creating user: %v", err)
	}
	return fmt.Sprintf("%v", result.InsertedID), nil
}

// hashPassword hashes the given password using a strong hashing algorithm (bcrypt).
func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// generateJWTToken generates a JWT token for the given user.
func generateJWTToken(user model.User, secretKey string) (string, error) {
	claims := jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token expiration time (1 day)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

// GetUserByUsername retrieves a user by username from the database.
func (r *Repository) GetUserByUsername(ctx context.Context, username string) (model.User, error) {
	var user model.User
	err := r.client.Database("ChatGo").Collection("users").
		FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

// AddMessage inserts a new message into the conversations collection.
func (r *Repository) AddMessage(ctx context.Context, conversation model.Conversation) (string, error) {
	result, err := r.client.Database("ChatGo").Collection("conversations").InsertOne(ctx, conversation)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", result.InsertedID), nil
}

// GetUserConversations retrieves all conversations for a user from the database.
func (r *Repository) GetUserConversations(ctx context.Context, userID string) ([]model.Conversation, error) {
	var conversations []model.Conversation

	cursor, err := r.client.Database("ChatGo").Collection("conversations").
		Find(ctx, bson.M{"$or": []bson.M{{"sender": userID}, {"receiver": userID}}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var conversation model.Conversation
		if err := cursor.Decode(&conversation); err != nil {
			return nil, err
		}
		conversations = append(conversations, conversation)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return conversations, nil
}

// GetUserConversationsWithTimeFilter retrieves conversations for a user based on user_id and optional time filters.
func (r *Repository) GetUserConversationsWithTimeFilter(ctx context.Context, userID string, afterUnixtime, beforeUnixtime int64, limit int) ([]model.Conversation, error) {
	var conversations []model.Conversation

	filter := bson.M{
		"$or": []bson.M{{"sender": userID}, {"receiver": userID}},
	}

	// Tambahkan filter waktu jika diberikan
	if afterUnixtime > 0 || beforeUnixtime > 0 {
    	timeFilter := bson.M{}
    	if afterUnixtime > 0 {
        	timeFilter["$gte"] = time.Unix(afterUnixtime, 0)
    	}
    	if beforeUnixtime > 0 {
        	timeFilter["$lte"] = time.Unix(beforeUnixtime, 0)
    	}
    	// Ubah agar kedua filter diterapkan bersamaan
    	filter["$and"] = []bson.M{filter, {"dateTime": timeFilter}}
	}

	// Tambahkan limit jika diberikan
	options := options.Find()
	if limit > 0 {
		options.SetLimit(int64(limit))
	}

	cursor, err := r.client.Database("ChatGo").Collection("conversations").Find(ctx, filter, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var conversation model.Conversation
		if err := cursor.Decode(&conversation); err != nil {
			return nil, err
		}
		conversations = append(conversations, conversation)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return conversations, nil
}