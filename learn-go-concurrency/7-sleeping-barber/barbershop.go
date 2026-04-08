package main

import (
	"time"

	"github.com/fatih/color"
)

type BarberShop struct {
	ShopCapacity    int
	HairCutDuration time.Duration
	NumberOfBarbers int
	BarbersDoneChan chan bool
	ClientsChan     chan string
	Open            bool
}

func (shop *BarberShop) addBarber(barber string) {
	shop.NumberOfBarbers++

	go func() {
		for {
			isSleeping := false

			// if there are no clients, the barber goes to sleep
			if len(shop.ClientsChan) == 0 {
				isSleeping = true
				color.Blue("Barber %s went to sleep...", barber)
			}

			client, isShopOpen := <-shop.ClientsChan

			// shop is closed, so send the barber home and close this goroutine
			if !isShopOpen {
				color.Yellow("The shop is closed, barber %s preparing to go home...", barber)
				shop.sendBarberHome(barber)
				return
			}

			if isSleeping {
				color.Yellow("%s wakes %s up.", client, barber)
				isSleeping = false
			}

			// cut hair
			shop.cutHair(barber, client)
		}
	}()

}

func (shop *BarberShop) cutHair(barber, client string) {
	color.Green("%s cutting %s hair.", barber, client)
	time.Sleep(shop.HairCutDuration)
	color.Green("%s has finished cutting %s's hair.", barber, client)
}

func (shop *BarberShop) sendBarberHome(barber string) {
	shop.BarbersDoneChan <- true
	color.Blue("%s went home.", barber)
}

func (shop *BarberShop) closeShopForDay() {
	color.Yellow("Closing shop...")
	shop.Open = false
	close(shop.ClientsChan)

	for a := 1; a <= shop.NumberOfBarbers; a++ {
		<-shop.BarbersDoneChan
	}

	close(shop.BarbersDoneChan)

	color.Yellow("-----------------------------------------")
	color.Yellow("The shop is closed, are barbers went home...")
}

func (shop *BarberShop) addClient(client string) {
	color.Green("*** %s arrives!", client)

	if !shop.Open {
		color.Red("Shop is already closed, %s leaves.")
	}

	select {
	case shop.ClientsChan <- client:
		color.Yellow("%s takes a seat in the waiting room.", client)
	default:
		color.Red("The waiting room is full, so %s leaves.", client)
	}
}
