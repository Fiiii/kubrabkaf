package main

import (
	"fmt"
	"kubrabkaf/business/kafkaClient"
	"log"
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
	var wg sync.WaitGroup
	wgDone := make(chan bool)

	// Start the producer with a 5 seconds interval
	wg.Add(1)
	tn := time.Tick(5 * time.Second)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-tn:
				fmt.Println("Producing message...")
				err = server.EthKafka.Produce(kafkaClient.PingTopic, []byte("ping! "+time.Now().String()))
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

	// Start the consumers
	wg.Add(1)
	go func() {
		defer wg.Done()
		server.EthKafka.Consume()
	}()

	go func() {
		wg.Wait()
		close(wgDone)
	}()
	//<-sig
	//fmt.Println("Received termination signal. Closing...")

	// Close the consumer and wait for it to finish processing

	select {
	case <-wgDone:
		break
	}

}
