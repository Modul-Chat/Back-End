package main

import (
	"context"
	"fmt"
	"go-conversation/db"
	"go-conversation/handler"
	"go-conversation/repository"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewRepository(client *mongo.Client) *repository.Repository {
    return repository.NewRepository(client)
}

// Update main function to include the new endpoint
func main() {
	client := db.MgoConn()
	defer client.Disconnect(context.Background())

	repository := NewRepository(client)
	router := mux.NewRouter()

	router.HandleFunc("/user/register", handler.RegisterHandler(repository)).Methods("POST")
	router.HandleFunc("/user/login", handler.LoginHandler(repository)).Methods("POST")
	router.HandleFunc("/sendchat", handler.SendChatHandler(repository)).Methods("POST")
	router.HandleFunc("/user/conversations", handler.GetUserConversationHandler(repository)).Methods("GET")

	port := 8080
	fmt.Printf("Server is running on port %d...\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), router)
	
}