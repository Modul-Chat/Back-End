package main

import (
	// "example/tes-websocket/internal/ws"
	"example/tes-websocket/rabbitmq"
	// "example/tes-websocket/router"
	"fmt"
)

func main() {
	conn, ch, err := rabbitmq.InitBroker()
	
	if err != nil {
		fmt.Printf("Error setting up RabbitMQ: %v\n", err)
		return
	}
	
	defer conn.Close()
	defer ch.Close()

	// hub := ws.NewHub()
	// wsHandler := ws.NewHandler(hub)
	// go hub.Run()

	// router.InitRouter(wsHandler)
	// router.Start("localhost:5000")
}