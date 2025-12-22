package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")

	connectionString := "amqp://guest:guest@localhost:5672/"

	connection, err := amqp.Dial(connectionString)
	if err != nil {
		fmt.Printf("Failed to connect to RabbitMq: %s\n", err)
		return
	}
	defer connection.Close()

	newChannel, err := connection.Channel()
	if err != nil {
		fmt.Printf("Failed to open a channel: %s\n", err)
		return
	}
	defer newChannel.Close()

	// Use my internal/pubsub
	err = pubsub.PublishJSON(newChannel, routing.ExchangePerilDirect, string(routing.PauseKey), routing.PlayingState{IsPaused: true})
	if err != nil {
		fmt.Printf("Failed to publish message: %s\n", err)
		return
	}

	fmt.Println("Connected to RabbitMQ successfully!")

	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan

	fmt.Println("Peril server shutting down...")
}

// I'll continue this week
