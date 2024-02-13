package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

var seatingCapacity = 5
var arrivalRate = 100 
var cutDuration = 1000 * time.Millisecond
var timeOpen = 10 * time.Second

func main() {

	color.Yellow("The sleeping barber problem")
	color.Yellow("---------------------------")

	clientChan := make(chan string, seatingCapacity)
	doneChan := make(chan bool)

	shop := BarberShop{
		ShopCapacity: seatingCapacity,
        HaircutDuration: cutDuration,
        NumberOfBarbers: 0,
        BarbersDone: doneChan,
        ClientsChan: clientChan,
        Open: true,
	}


	color.Green("Shop is open")
	shop.addBarber("Scary Frankie")
	shop.addBarber("Voodo Karl")
	shop.addBarber("Trogg 'The Log' Frogg")

	shopClosing := make(chan bool)
	closed := make(chan bool)

	go func ()  {
		<- time.After(timeOpen)
		shopClosing <- true
		shop.closeShopForDay()
		closed <- true
	}()

	i:= 1	
	go func (){
		for {
			randomMilliseconds := rand.Int() % (2* arrivalRate)
			select{
				case <-shopClosing:
					return
                case <-time.After(time.Duration(randomMilliseconds) * time.Millisecond):
					shop.addClient(fmt.Sprintf("Client #%d", i))
					i++
			}
		}
	}()
	
	<-closed

}
