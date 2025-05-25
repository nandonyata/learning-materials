package main

import (
	"fmt"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

func main() {
	// Connect to RabbitMQ
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	defer conn.Close()

	// Create a channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}
	defer ch.Close()

	// Declare an x-delayed-message exchange
	err = ch.ExchangeDeclare(
		"delayed_exchange",  // Exchange name
		"x-delayed-message", // Type of exchange
		true,                // Durable
		false,               // Auto-deleted
		false,               // Internal
		false,               // No-wait
		amqp091.Table{
			"x-delayed-type": "direct", // The underlying exchange type
		},
	)
	if err != nil {
		log.Fatalf("Failed to declare exchange: %s", err)
	}

	// Declare a queue
	q, err := ch.QueueDeclare(
		"delayed_queue", // Queue name
		true,            // Durable
		false,           // Delete when unused
		false,           // Exclusive
		false,           // No-wait
		nil,             // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare queue: %s", err)
	}

	// Bind the queue to the exchange
	err = ch.QueueBind(
		q.Name,             // Queue name
		"",                 // Routing key
		"delayed_exchange", // Exchange name
		false,              // No-wait
		nil,                // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to bind queue: %s", err)
	}

	// Send a delayed message
	message := "This is a delayed message"
	delayDuration := 10000 // Delay in milliseconds (5 seconds)

	err = ch.Publish(
		"delayed_exchange", // Exchange
		"",                 // Routing key
		false,              // Mandatory
		false,              // Immediate
		amqp091.Publishing{
			Headers: amqp091.Table{
				"x-delay": delayDuration, // Delay in ms
			},
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	if err != nil {
		log.Fatalf("Failed to publish a message: %s", err)
	}
	fmt.Printf("Sent a delayed message: %s\n", message)

	// Wait to receive the delayed message
	msgs, err := ch.Consume(
		q.Name, // Queue name
		"",     // Consumer name
		true,   // Auto-acknowledge
		false,  // Exclusive
		false,  // No-local
		false,  // No-wait
		nil,    // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %s", err)
	}

	// Wait for a message and print it
	for msg := range msgs {
		fmt.Printf("Received message: %s\n", msg.Body)
	}
}
