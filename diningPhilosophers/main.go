package main

import (
	"fmt"
	"sync"

	"time"
)

type Philosopher struct {
	name      string
	rightFork int
	leftFork  int
}

type Order struct {
	order []string
	mtx   sync.Mutex
}



var philosophers = []Philosopher{
	{name: "Plato", leftFork: 4, rightFork: 0},
	{name: "Socrates", leftFork: 0, rightFork: 1},
	{name: "Aristotle", leftFork: 1, rightFork: 2},
	{name: "Pascal", leftFork: 2, rightFork: 3},
	{name: "Locke", leftFork: 3, rightFork: 4},
}

// how many times eating should be performed per philosopher
var hunger = 3
var eatTime = time.Second * 0
var thinkTime = time.Second * 0
var order *Order = &Order{[]string{}, sync.Mutex{}}


func main() {
	//print out a welcome message
	fmt.Println("Dining Philosophers Problem")
	fmt.Println("---------------------------")
	fmt.Println("The table is empty.")

	//start the eating
	dine()
	//print out finishing message
	fmt.Println("The table is empty.")
}

func dine() {
	// eatTime = 0 * time.Second
	// sleepTime = 0 * time.Second
	// thinkTime = 0 * time.Second

	wg := &sync.WaitGroup{}
	wg.Add(len(philosophers))
	seated := &sync.WaitGroup{}
	seated.Add(len(philosophers))

	forks := make(map[int]*sync.Mutex)

	for i := 0; i < len(philosophers); i++ {
		forks[i] = &sync.Mutex{}
	}

	for i := 0; i < len(philosophers); i++ {
		//fire off a goroutine for the current philosopher
		go diningProblem(philosophers[i], wg, forks, seated, order)
	}

	wg.Wait()
	for i, v := range order.order{
		fmt.Printf("%v has finished in order %v\n", v, i)
	}
}

func diningProblem(philosopher Philosopher, wg *sync.WaitGroup, forks map[int]*sync.Mutex, seated *sync.WaitGroup, order *Order) {
	defer wg.Done()

	// seat
	fmt.Printf("%s is seated at the table.\n", philosopher.name)
	seated.Done()
	seated.Wait()

	for i := hunger; i > 0; i-- {
		fmt.Printf("%s is hungry.\n", philosopher.name)

		if philosopher.leftFork > philosopher.rightFork {
			fmt.Printf("%s is hungry.\n", philosopher.name)
			forks[philosopher.rightFork].Lock()
			fmt.Printf("%s takes the right fork.\n", philosopher.name)
			forks[philosopher.leftFork].Lock()
			fmt.Printf("%s takes the left fork.\n", philosopher.name)
		} else {
			forks[philosopher.leftFork].Lock()
			fmt.Printf("%s takes the left fork.\n", philosopher.name)
			forks[philosopher.rightFork].Lock()
			fmt.Printf("%s takes the right fork.\n", philosopher.name)
		}
		fmt.Printf("%s has both forks and is eating.\n", philosopher.name)
		time.Sleep(eatTime)
		fmt.Printf("%s has both forks and is thinking.\n", philosopher.name)
		time.Sleep(thinkTime)
		forks[philosopher.leftFork].Unlock()
		forks[philosopher.rightFork].Unlock()
		fmt.Printf("%s put down the forks.\n", philosopher.name)
	}
	fmt.Println(philosopher.name, "is done eating.")
	order.mtx.Lock()
	order.order = append(order.order, philosopher.name)
	order.mtx.Unlock()
	fmt.Println(philosopher.name, "left the table.")
}
