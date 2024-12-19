package main

import (
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-stomp/stomp/v3"
)

func main() {

	conn, err := stomp.Dial("tcp", "localhost:61613", stomp.ConnOpt.HeartBeat(0, 0))
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Disconnect()

	subcription, err := conn.Subscribe("my-queue", stomp.AckClientIndividual)
	if err != nil {
		log.Fatal(err)
	}

	defer subcription.Unsubscribe()

	messages := make(chan *stomp.Message, 100)

	for i := 0; i < 10; i++ {
		go messageConsumer(i, conn, messages)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

loop:
	for {
		select {
		case <-signalChan:
			log.Println("Received signal, shutting down...")
			break loop

		case message := <-subcription.C:

			if message.Err != nil {
				log.Printf("Error receiving message: %v\n", message.Err)
				continue
			}

			messages <- message

		}
	}

	os.Exit(0)
}

func messageConsumer(id int, conn *stomp.Conn, messages <-chan *stomp.Message) {

	for {
		select {
		case message := <-messages:
			duration := rand.Intn(500) + 1
			log.Printf("WorkerId : %d Message : %s Duration : %d\n", id, message.Body, duration)
			time.Sleep(time.Duration(duration) * time.Millisecond)
			if err := conn.Ack(message); err != nil {
				log.Printf("Error acknowledging message: %v\n", err)
				continue
			}
		}
	}

}
