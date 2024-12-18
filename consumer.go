package main

import (
	"log"

	"github.com/go-stomp/stomp/v3"
)

func main() {

	conn, err := stomp.Dial("tcp", "localhost:61613", stomp.ConnOpt.HeartBeat(0, 0))
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Disconnect()

	subcription, err := conn.Subscribe("my-queue", stomp.AckAuto)
	if err != nil {
		log.Fatal(err)
	}

	defer subcription.Unsubscribe()

	for {
		message := <-subcription.C
		if message.Err != nil {
			log.Println("timeout")
			continue
		}
		log.Println("Message : " + string(message.Body))

	}

}
