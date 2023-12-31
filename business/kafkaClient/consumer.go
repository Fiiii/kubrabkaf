package kafkaClient

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"time"
)

// Consume reads messages from the Kafka broker.
func (e *EthKafka) Consume() {
	run := true
	for run {
		msg, err := e.KConsumer.ReadMessage(time.Second)
		if err == nil {
			fmt.Printf("<--- Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
		} else if !err.(kafka.Error).IsTimeout() {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}

	//e.kConsumer.Close()
}
