package kafkaClient

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type EthKafka struct {
	kProducer *kafka.Producer
	kConsumer *kafka.Consumer
}

func NewEthKafka() (*EthKafka, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
	})

	if err != nil {
		return nil, fmt.Errorf("Failed to create producer: %s\n", err)
	}

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "myGroup",
		"auto.offset.reset": "smallest"})

	if err != nil {
		return nil, fmt.Errorf("Failed to create consumer: %s\n", err)
	}

	err = consumer.SubscribeTopics([]string{"myTopic"}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %s\n", err)
	}

	if err != nil {
		return nil, fmt.Errorf("Failed to create consumer: %s\n", err)
	}

	return &EthKafka{
		kProducer: producer,
		kConsumer: consumer,
	}, nil
}
