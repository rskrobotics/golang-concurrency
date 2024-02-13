package main

import (
	"time"

	"github.com/fatih/color"
)

type BarberShop struct {
	ShopCapacity    int
	HaircutDuration time.Duration
	NumberOfBarbers int
	BarbersDone     chan bool
	ClientsChan     chan string
	Open            bool
}

func (shop *BarberShop) addBarber(barber string) {
	shop.NumberOfBarbers++

	go func() {
		isSleeping := false
		color.Yellow("%s goes to work, and checks for clients", barber)

		for {
			if len(shop.ClientsChan) == 0 {
				color.Yellow("There is nothing to do, so %s takes a nap", barber)
				isSleeping = true
			}
			client, shopOpen := <-shop.ClientsChan

			if shopOpen {
				if isSleeping {
					color.Yellow("%s wakes %s up.", client, barber)
					isSleeping = false
				}
				shop.cutHair(barber, client)
			} else {
				shop.sendBarberHome(barber)
				return
			}

		}
	}()
}

func (shop *BarberShop) cutHair(barber, client string) {
	color.Green("%s cuts %s's hair", barber, client)
	time.Sleep(shop.HaircutDuration)
	color.Green("%s is finished with %s's haircut", barber, client)

}

func (shop *BarberShop) sendBarberHome(barber string) {
	color.Cyan("%s is going home", barber)
	shop.BarbersDone <- true
}

func (shop *BarberShop) closeShopForDay() {
	color.Cyan("Shop is closing for the day.")
	close(shop.ClientsChan)
	shop.Open = false

	for a := 01; a <= shop.NumberOfBarbers; a++ {
		<-shop.BarbersDone
	}
	close(shop.BarbersDone)

	color.Green("The barbershop is now closed for the day, everyone has gone home.")
	color.Green("-----------------------------------------------------------------")
}

func (shop *BarberShop) addClient(client string) {
	color.Green("*** %s arrives!", client)
	if shop.Open{
		select {
		case shop.ClientsChan <- client:
			color.Yellow("%s takes a seat in the waiting room.", client)
		default:
			color.Red("The waiting room is full, so %s leaves.", client)
		}
	} else {
		color.Red("The shop is closed for the day, so %s is not allowed to enter.", client)
	}
}
