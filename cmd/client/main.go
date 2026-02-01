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
	fmt.Println("Starting Peril client...")

	rabbitMqURL := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(rabbitMqURL)
	if err != nil {
		log.Fatalf("error establishing RabbitMQ connection: %e\n", err)
	}
	defer conn.Close()

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("error getting username: %e\n", err)
	}

	queueName := fmt.Sprintf("%s.%s", routing.PauseKey, username)
	fmt.Printf("Using queueName %s.\n", queueName)
	_, _, err = pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilDirect,
		queueName,
		routing.PauseKey,
		pubsub.Transient)
	if err != nil {
		log.Fatalf("error declaring and binding channel/queue-exchange: %e\n", err)
	}

	closeSig := make(chan os.Signal, 1)
	signal.Notify(closeSig, os.Interrupt)
	//<-closeSig

	for sigg := range closeSig {
		fmt.Printf("\nTerminating program: singal %s received\n", sigg)
		os.Exit(1)
	}
}
