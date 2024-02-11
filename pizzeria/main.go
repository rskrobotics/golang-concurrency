package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

const NumberOfPizzas = 10

var pizzasMade, pizzasFailed, total int

type Producer struct {
	data chan PizzaOrder
	quit chan chan error
}

type PizzaOrder struct {
	pizzaNumber int
	message     string
	success     bool
}

func (p *Producer) Close() error {
	ch := make(chan error)
	p.quit <- ch
	return <-ch
}

func makePizza(pizzaNumber int) *PizzaOrder {
	pizzaNumber++
	if pizzaNumber <= NumberOfPizzas {
		delay := rand.Intn(5) + 1
		fmt.Printf("Received order #%d!\n", pizzaNumber)

		rnd := rand.Intn(12) + 1
		msg := ""
		success := false

		if rnd < 5 {
			pizzasFailed++
		} else {
			pizzasMade++
		}
		total++
		fmt.Printf("Making pizza #%d.... It will take %d seconds\n", pizzaNumber, delay)
		fmt.Println("BLOCK NOW FROM #%d", pizzaNumber)
		time.Sleep(time.Duration(delay) * time.Second)
		fmt.Println("UNBLOCK NOW FROM #%d", pizzaNumber)

		switch {
		case rnd < 2:
			msg = fmt.Sprintf("*** We ran out of ingredients for pizza %d! ***", pizzaNumber)
		case rnd <= 4:
			msg = fmt.Sprintf("*** The cook quit while making pizza %d! ***", pizzaNumber)
		default:
			success = true
			msg = fmt.Sprintf("*** We made pizza %d! ***", pizzaNumber)
		}

		return &PizzaOrder{pizzaNumber, msg, success}

	}
	return &PizzaOrder{pizzaNumber: pizzaNumber}
}

func pizzeria(pizzaMaker *Producer) {
	var i int = 0

	for {
		currentPizza := makePizza(i)
		i = currentPizza.pizzaNumber
		select {
		case pizzaMaker.data <- *currentPizza:
		case quitChan := <-pizzaMaker.quit:
			close(pizzaMaker.data)
			close(quitChan)
			return
		}
	}
}

func main() {

	color.Cyan("The pizzeria is ready!")
	color.Cyan("----------------------")

	pizzaJob := &Producer{
		data: make(chan PizzaOrder),
		quit: make(chan chan error),
	}

	go pizzeria(pizzaJob)

	for i := range pizzaJob.data {
		if i.pizzaNumber <= NumberOfPizzas {
			if i.success {
				color.Green(i.message)
				color.Green("Order %d is out for delivery!\n", i.pizzaNumber)
			} else {
				color.Red(i.message)
				color.Red("Customer is unsatisfied!")
			}
		} else {
			color.Cyan("Done making pizzas!")
			err := pizzaJob.Close()
			if err != nil {
				color.Red("Error closing the channel: ", err)
			}
		}
	}

	color.Cyan("-----------------")
	color.Cyan("Done for the day!")
	color.Cyan("We made %d pizzas and failed to make %d pizzas! Amount of attempts was %d", pizzasMade, pizzasFailed, total)
}
