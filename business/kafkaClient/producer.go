package kafkaClient

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func (e *EthKafka) Produce(topic string, kafkaMsg []byte) error {
	//defer e.kProducer.Close()

	// Delivery report handler for produced messages
	go func() {
		for ev := range e.kProducer.Events() {
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

	// Produce messages to topic (asynchronously)
	err := e.kProducer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          kafkaMsg,
	}, nil)

	if err != nil {
		return fmt.Errorf("failed to produce message: %v\n", err)
	}

	// Wait for message deliveries before shutting down
	e.kProducer.Flush(15 * 1000)
	return nil
}
