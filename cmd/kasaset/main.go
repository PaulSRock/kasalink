package main

import (
	"flag"
	"kasalink"
	"log"
)

var host = flag.String("host", "", "Hostname of the strip to talk to")

func main() {
	flag.Parse()
	plug, err := kasalink.NewKasaPowerPlug(*host)
	if err != nil {
		log.Fatal(err)
	}
	defer plug.Close()

	plug.TurnDeviceOff(1)
	//fmt.Println(plug.GetSystemInfo())
}


