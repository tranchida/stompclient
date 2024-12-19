package main

import (
	"fmt"
	"log"

	"github.com/go-stomp/stomp/v3"
)

func main() {

	conn, err := stomp.Dial("tcp", "localhost:61613", stomp.ConnOpt.HeartBeat(0, 0))
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Disconnect()

	for i := 0; i < 10; i++ {
		conn.Send("my-queue", "text/plain", []byte(fmt.Sprintf("hello, world! %d", i)), stomp.SendOpt.Header("persistent", "true"))
	}
	log.Println("Message sent")

}
