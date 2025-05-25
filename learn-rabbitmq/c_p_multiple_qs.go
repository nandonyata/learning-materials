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

	// Declare two queues with different routing keys
	queueOne, err := ch.QueueDeclare(
		"queue_one", // Queue name
		true,        // Durable
		false,       // Delete when unused
		false,       // Exclusive
		false,       // No-wait
		nil,         // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare queue_one: %s", err)
	}

	queueTwo, err := ch.QueueDeclare(
		"queue_two", // Queue name
		true,        // Durable
		false,       // Delete when unused
		false,       // Exclusive
		false,       // No-wait
		nil,         // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare queue_two: %s", err)
	}

	// Bind queues to the exchange with different routing keys
	err = ch.QueueBind(
		queueOne.Name,      // Queue name
		"routing_key_one",  // Routing key (specific for this queue)
		"delayed_exchange", // Exchange name
		false,              // No-wait
		nil,                // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to bind queue_one: %s", err)
	}

	err = ch.QueueBind(
		queueTwo.Name,      // Queue name
		"routing_key_two",  // Routing key (specific for this queue)
		"delayed_exchange", // Exchange name
		false,              // No-wait
		nil,                // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to bind queue_two: %s", err)
	}

	// Send a delayed message to queue_one (with routing_key_one)
	message := "This message is for queue_one"
	delayDuration := 10000 // Delay in milliseconds (10 seconds)

	err = ch.Publish(
		"delayed_exchange", // Exchange
		"routing_key_one",  // Routing key
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
		log.Fatalf("Failed to publish a message to queue_one: %s", err)
	}
	fmt.Printf("Sent message to queue_one: %s\n", message)

	// Send a delayed message to queue_two (with routing_key_two)
	message2 := "This message is for queue_two"
	err = ch.Publish(
		"delayed_exchange", // Exchange
		"routing_key_two",  // Routing key
		false,              // Mandatory
		false,              // Immediate
		amqp091.Publishing{
			Headers: amqp091.Table{
				"x-delay": delayDuration, // Delay in ms
			},
			ContentType: "text/plain",
			Body:        []byte(message2),
		},
	)
	if err != nil {
		log.Fatalf("Failed to publish a message to queue_two: %s", err)
	}
	fmt.Printf("Sent message to queue_two: %s\n", message2)

	// Set up consumers for both queues
	msgsOne, err := ch.Consume(
		queueOne.Name, // Queue name
		"",            // Consumer name
		true,          // Auto-acknowledge
		false,         // Exclusive
		false,         // No-local
		false,         // No-wait
		nil,           // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to register consumer for queue_one: %s", err)
	}

	msgsTwo, err := ch.Consume(
		queueTwo.Name, // Queue name
		"",            // Consumer name
		true,          // Auto-acknowledge
		false,         // Exclusive
		false,         // No-local
		false,         // No-wait
		nil,           // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to register consumer for queue_two: %s", err)
	}

	// Wait for messages on both queues
	go func() {
		for msg := range msgsOne {
			fmt.Printf("Received message from queue_one: %s\n", msg.Body)
		}
	}()

	go func() {
		for msg := range msgsTwo {
			fmt.Printf("Received message from queue_two: %s\n", msg.Body)
		}
	}()

	// Block indefinitely to keep the program running
	select {}
}
