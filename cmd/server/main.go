package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")

	amqpConn := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(amqpConn)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	fmt.Println("RabbitMQ connection stablished successfully!")

	closeSig := make(chan os.Signal, 1)
	signal.Notify(closeSig, os.Interrupt)
	//<-closeSig

	for sigg := range closeSig {
		fmt.Printf("\nTerminating program: singal %s received\n", sigg)
		os.Exit(1)
	}
}
