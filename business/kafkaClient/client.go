package kafkaClient

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

var (
	PingTopic = "pingTopic"
	EthTopic  = "ethTopic"
	kafkaUrl  = "localhost:9092"
)

// EthKafka is a wrapper around the Kafka client.
type EthKafka struct {
	KProducer *kafka.Producer
	KConsumer *kafka.Consumer
}

// NewEthKafka creates a new Kafka client.
func NewEthKafka() (*EthKafka, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": kafkaUrl,
	})

	if err != nil {
		return nil, fmt.Errorf("Failed to create producer: %s\n", err)
	}

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": kafkaUrl,
		"group.id":          "myGroup",
		"auto.offset.reset": "smallest"})

	if err != nil {
		return nil, fmt.Errorf("Failed to create consumer: %s\n", err)
	}

	err = consumer.SubscribeTopics([]string{PingTopic, EthTopic}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to topics: %s\n", err)
	}

	return &EthKafka{
		KProducer: producer,
		KConsumer: consumer,
	}, nil
}
