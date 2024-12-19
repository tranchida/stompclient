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

}
