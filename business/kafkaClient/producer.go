package kafkaClient

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// Produce sends a message to the Kafka broker.
func (e *EthKafka) Produce(topic string, kafkaMsg []byte) error {
	go func() {
		for ev := range e.KProducer.Events() {
			switch ev := ev.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("---> Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	err := e.KProducer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          kafkaMsg,
	}, nil)

	if err != nil {
		return fmt.Errorf("failed to produce message: %v\n", err)
	}

	// Wait for message deliveries before shutting down
	e.KProducer.Flush(15 * 1000)
	return nil
}
