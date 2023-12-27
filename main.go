package main

import (
	"log"
	"sync"

	// "example/tes-websocket/internal/ws"
	"example/tes-websocket/rabbitmq"
	// "example/tes-websocket/router"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		conn, ch, err := rabbitmq.InitBroker()
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()
		defer ch.Close()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		setupGinServer()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-stop
		fmt.Println("\nReceived an interrupt signal. Shutting down...")
		os.Exit(0)
	}()

	wg.Wait()
	fmt.Println("Application gracefully terminated")

	// conn, ch, err := rabbitmq.InitBroker()

	// if err != nil {
	// 	fmt.Printf("Error setting up RabbitMQ: %v\n", err)
	// 	return
	// }

	// r := gin.Default()
	// r.POST("/message", rabbitmq.GetRabbitMQBroker().SendMessage)
	// r.Run()

	// defer conn.Close()
	// defer ch.Close()

	// hub := ws.NewHub()
	// wsHandler := ws.NewHandler(hub)
	// go hub.Run()

	// router.InitRouter(wsHandler)
	// router.Start("localhost:5000")
}

func setupGinServer() {
	r := gin.Default()

	// Define your Gin routes here
	r.POST("/message", rabbitmq.GetRabbitMQBroker().SendMessageToQueueDB)

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Error starting Gin server: %v", err)
	}
}
