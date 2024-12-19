package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-stomp/stomp/v3"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

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

	go func() {

		for {
			select {
			case <-ctx.Done():
				return

			case message := <-subcription.C:

				if message.Err != nil {
					log.Printf("Error receiving message: %v\n", message.Err)
					continue
				}

				ch := make(chan bool)
				go processMessage(message, ch)
				<- ch

				if err := conn.Ack(message); err != nil {
					log.Printf("Error acknowledging message: %v\n", err)
					continue
				}

			}
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	<-signalChan // Block until a signal is received
	log.Println("Shutting down...")
	cancel()

}

func processMessage(message *stomp.Message, ch chan <- bool) {
	log.Println("Message : " + string(message.Body))
	time.Sleep(10 * time.Second)
	ch <- true
}
