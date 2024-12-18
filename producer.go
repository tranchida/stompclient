package main

import (
	"log"

	"github.com/go-stomp/stomp/v3"
)

func main() {
	// Create a new producer

	conn, err := stomp.Dial("tcp", "10.108.136.84:61613", stomp.ConnOpt.HeartBeat(0, 0))
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Disconnect()

	conn.Send("my-queue", "text/plain", []byte("hello, world!"), stomp.SendOpt.Header("persistent", "true"))

	log.Println("Message sent")

}