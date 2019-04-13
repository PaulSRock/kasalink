package main

import (
	"flag"
	"log"

	"github.com/PaulSRock/kasalink"
)

var host = flag.String("host", "", "Hostname of the strip to talk to")
var plugNumber = flag.Int("plug", 0, "Plug to switch")
var on = flag.Bool("on", false, "Plug state")

func main() {
	flag.Parse()
	plug, err := kasalink.NewKasaPowerPlug(*host)
	if err != nil {
		log.Fatal(err)
	}

	defer plug.Close()

	if *on {
		_, err := plug.TurnDeviceOn(*plugNumber)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		_, err := plug.TurnDeviceOff(*plugNumber)
		if err != nil {
			log.Fatal(err)
		}
	}
}
