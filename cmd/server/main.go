package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")

	rabbitMqURL := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(rabbitMqURL)
	if err != nil {
		log.Fatalf("error establishing RabbitMQ connection: %e\n", err)
	}
	defer conn.Close()
	fmt.Println("RabbitMQ connection stablished successfully!")

	messageChannel, err := conn.Channel()
	if err != nil {
		log.Fatalf("error creating Channel: %e\n", err)
	}

	err = pubsub.PublishJSON(messageChannel, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: true})
	fmt.Println("Sending msg")
	if err != nil {
		log.Fatalf("error publishing message: %e\n", err)
	}

	gamelogic.PrintServerHelp()

REPL:
	for {
		input := gamelogic.GetInput()
		if len(input) != 1 {
			log.Print("Invalid number of arguments received\n")
			continue
		}
		switch input[0] {
		case "pause":
			log.Print("Sending pause message\n")
			err = pubsub.PublishJSON(
				messageChannel,
				routing.ExchangePerilDirect,
				routing.PauseKey,
				routing.PlayingState{IsPaused: true})
			if err != nil {
				log.Fatalf("error publishing message: %e\n", err)
			}
		case "resume":
			log.Print("Sending resume message\n")
			err = pubsub.PublishJSON(
				messageChannel,
				routing.ExchangePerilDirect,
				routing.PauseKey,
				routing.PlayingState{IsPaused: false})
			if err != nil {
				log.Fatalf("error publishing message: %e\n", err)
			}
		case "quit":
			log.Print("Exiting...\n")
			break REPL
		default:
			log.Printf("Invalid argument '%s' received\n", input[0])
			continue
		}
	}

	closeSig := make(chan os.Signal, 1)
	signal.Notify(closeSig, os.Interrupt)
	//<-closeSig

	for sigg := range closeSig {
		fmt.Printf("\nTerminating program: singal %s received\n", sigg)
		os.Exit(1)
	}
}
