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
		fmt.Println("Err in channel")
		panic(err)
	}
	defer channel.Close()

	// declaring queue with its properties over the the channel opened
	queue, err := channel.QueueDeclare(
		"my-queue", // name
		true,       // durable
		false,      // auto delete
		false,      // exclusive
		false,      // no wait
		nil,        // args
	)
	if err != nil {
		fmt.Println("Err in q declare")
		panic(err)
	}

	// publishing a message
	err = channel.Publish(
		"",         // exchange
		"my-queue", // key
		false,      // mandatory
		false,      // immediate
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte("Test Message with delay?"),
			Headers: amqp091.Table{
				"x-delay": 5000,
			},
		},
	)
	if err != nil {
		fmt.Println("Err in publish msg")
		panic(err)
	}

	fmt.Println("Queue status:", queue)
	fmt.Println("Successfully published message")
}
