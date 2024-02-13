package main

import (
	"fmt"
	"strings"
)

func shout(ping <-chan string, pong chan<- string) {
	for {
		s := <-ping
		pong <- strings.ToUpper(s) + "!!!"
	}
}

func main() {
	ping := make(chan string)
	pong := make(chan string)

	go shout(ping, pong)

	fmt.Println("Type something and press enter (Q to quit)")
	for {
		fmt.Print("--->")

		var userInput string
		_, _ = fmt.Scanln(&userInput)

		if userInput == "Q" {
			break
		}
		ping <- userInput
		response := <-pong
		fmt.Println("Response: " , response)
	}
	fmt.Println("All done, closing channels")
	close(ping)
	close(pong)
}
