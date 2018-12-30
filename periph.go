package main

import (
	"log"
	"time"

	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/host"
	"periph.io/x/periph/host/rpi"
)

func main() {
	// Load all the drivers:
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	t := time.NewTicker(500 * time.Millisecond)
	for l := gpio.Low; ; l = !l {
		// Lookup a pin by its location on the board:
		if err := rpi.P1_33.Out(l); err != nil {
			log.Fatal(err)
		}
		<-t.C
	}
}