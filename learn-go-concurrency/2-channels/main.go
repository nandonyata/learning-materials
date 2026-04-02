package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/fatih/color"
)

const NumberOfPizzas = 10

var pizzasMade, pizzasFailed, total int

// Producer is a type for structs that holds two channels: one for pizzas, with all
// information for a given pizza order including whether it was made
// successfully, and another to handle end of processing (when we quit the channel)
type Producer struct {
	data chan PizzaOrder
	quit chan chan error
}

// PizzaOrder is a type for structs that describes a given pizza order. It has the order
// number, a message indicating what happened to the order, and a boolean
// indicating if the order was successfully completed.
type PizzaOrder struct {
	pizzaNumber int
	message     string
	success     bool
}

// Close is simply a method of closing the channel when we are done with it (i.e.
// something is pushed to the quit channel)
func (p *Producer) Close() error {
	ch := make(chan error)
	p.quit <- ch
	return <-ch
}

// makePizza attempts to make a pizza. We generate a random number from 1-12,
// and put in two cases where we can't make the pizza in time. Otherwise,
// we make the pizza without issue. To make things interesting, each pizza
// will take a different length of time to produce (some pizzas are harder than others).
func makePizza(pizzaNumber int) *PizzaOrder {
	if total == NumberOfPizzas {
		return &PizzaOrder{
			pizzaNumber: pizzaNumber,
		}
	}
	pizzaNumber++
	total++
	succes := false
	message := "Pizza made successfully!."

	delay := rand.Intn(3)
	color.Yellow("Start making %dth pizza, and it might take some times..., like %d seconds", pizzaNumber, delay)
	time.Sleep(time.Duration(delay) * time.Second)

	randomNumber := rand.Intn(12)

	if randomNumber < 4 {
		message = "Pizza burned!"
		pizzasFailed++
	} else if randomNumber < 8 {
		message = "Pizza is undercooked!"
		pizzasFailed++
	} else {
		succes = true
		message = "Pizza is cooked perfectly!"
		pizzasMade++
	}

	return &PizzaOrder{
		pizzaNumber: pizzaNumber,
		message:     message,
		success:     succes,
	}
}

// pizzeria is a goroutine that runs in the background and
// calls makePizza to try to make one order each time it iterates through
// the for loop. It executes until it receives something on the quit
// channel. The quit channel does not receive anything until the consumer
// sends it (when the number of orders is greater than or equal to the
// constant NumberOfPizzas).
func pizzeria(pizzaMaker *Producer) {
	var pizzaNumber = 0
	for {
		pizza := makePizza(pizzaNumber)
		if pizza == nil {
			continue
		}

		pizzaNumber = pizza.pizzaNumber
		select {
		case pizzaMaker.data <- *pizza:
		case qc := <-pizzaMaker.quit:
			close(pizzaMaker.data)
			close(qc) // closing a channel will sends default value to the channel (nil)
			// qc <- errors.New("random error") // example of sending an error back to the Close() caller
			return
		}
	}
}

func main() {
	// seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// print out a message
	color.Cyan("The Pizzeria is open for business!")
	color.Cyan("----------------------------------")

	// create a producer
	pizzaJob := &Producer{
		data: make(chan PizzaOrder),
		quit: make(chan chan error),
	}

	// run the producer in the background
	go pizzeria(pizzaJob)

	// create and run consumer
	for pizza := range pizzaJob.data {
		if pizza.success {
			color.Blue("Pizza #%d with message '%s' is ready to delivered.", pizza.pizzaNumber, pizza.message)

		} else {
			color.Red("Pizza #%d with message '%s' is not ready to delivered. The customer seems upset.", pizza.pizzaNumber, pizza.message)

		}

		if pizza.pizzaNumber == NumberOfPizzas {
			color.Green("Time to close the pizzeria :D.")
			if err := pizzaJob.Close(); err != nil {
				color.Red("Error closing the pizzeria: %s", err)
				os.Exit(1)
			}
		}
	}

	// print out the ending message
	color.Cyan("-----------------")
	color.Cyan("Done for the day.")

	color.Cyan("We made %d pizzas, but failed to make %d, with %d attempts in total.", pizzasMade, pizzasFailed, total)

	switch {
	case pizzasFailed > 9:
		color.Red("It was an awful day...")
	case pizzasFailed >= 6:
		color.Red("It was not a very good day...")
	case pizzasFailed >= 4:
		color.Yellow("It was an okay day....")
	case pizzasFailed >= 2:
		color.Yellow("It was a pretty good day!")
	default:
		color.Green("It was a great day!")
	}
}
