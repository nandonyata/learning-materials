package main

import (
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("RabbitMQ in Golang: Getting started tutorial")

	connection, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	defer connection.Close()

	fmt.Println("Successfully connected to RabbitMQ instance")

	// opening a channel over the connection established to interact with RabbitMQ
	channel, err := connection.Channel()
	if err != nil {
		panic(err)
	}
	defer channel.Close()

	// declaring consumer with its properties over channel opened
	msgs, err := channel.Consume(
		"my-queue",           // queue
		"my-consumer-hehehe", // consumer
		true,                 // auto ack
		false,                // exclusive
		false,                // no local
		false,                // no wait
		nil,                  //args
	)
	if err != nil {
		panic(err)
	}

	// print consumed messages from queue
	forever := make(chan bool)
	go func() {
		for msg := range msgs {
			fmt.Printf("Received Message: %s\n", msg.Body)
		}
	}()

	fmt.Println("Waiting for messages...")
	<-forever
}
