package main

import (
	"context"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-stomp/stomp/v3"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

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
		wg.Add(1)
		go processMessageWorker(i, ctx, &wg, conn, messages)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

loop:
	for {
		select {
		case <-signalChan:
			log.Println("Received signal, shutting down...")
			cancel()
			break loop

		case message := <-subcription.C:

			if message.Err != nil {
				log.Printf("Error receiving message: %v\n", message.Err)
				continue
			}

			messages <- message

		}
	}

	log.Println("Wait go routine shutdown")
	wg.Wait()
	log.Println("all go routine shutdown")

	os.Exit(0)
}

func processMessageWorker(id int, ctx context.Context, wg *sync.WaitGroup, conn *stomp.Conn, messages <-chan *stomp.Message) {

	defer wg.Done()

workerLoop:
	for {
		select {
		case <-ctx.Done():
			log.Printf("WorkerId : %d shutting down...\n", id)
			break workerLoop
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

	log.Printf("WorkerId : %d Shutdown complete\n", id)
}
