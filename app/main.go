package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

const (
	serverAddress = ":8080"
)

func main() {
	srvCfg := &Config{
		ListenAddr: serverAddress,
	}

	server, err := NewServer(srvCfg)
	if err != nil {
		log.Fatal(err)
	}

	// Handle OS signals to gracefully shut down the consumer
	sig := make(chan os.Signal, 1)
	var wg sync.WaitGroup

	// Start the producer with a 5 seconds interval
	wg.Add(1)
	tn := time.Tick(5 * time.Second)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-tn:
				fmt.Println("Producing message...")
				err = server.EthKafka.Produce("myTopic", []byte("Hello, World!"+time.Now().String()))
				if err != nil {
					fmt.Printf("failed to produce message: %v\n", err)
				}
			}
		}
	}()

	// Start the server
	wg.Add(1)
	go func() {
		defer wg.Done()
		server.Start()
	}()

	// Start the consumer
	wg.Add(1)
	go func() {
		defer wg.Done()
		server.EthKafka.Consume()
	}()

	<-sig
	fmt.Println("Received termination signal. Closing...")

	// Close the consumer and wait for it to finish processing
	wg.Wait()
}
