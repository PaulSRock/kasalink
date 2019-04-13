package kasalink

import (
	"io"
	"log"
)

func closer(thingToClose io.Closer) {
	var err = thingToClose.Close()
	if err != nil {
		log.Println(err)
	}
}