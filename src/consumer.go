package src

import (
	"log"

	_kafka "github.com/confluentinc/confluent-kafka-go/kafka"
)

// ConsumerConfiguration model
type ConsumerConfiguration struct {
	Broker      string
	GroupID     string
	AutoCommit  bool
	StoreOffset bool
	OffsetReset string
}

// NewConsumerConfiguration factory
func NewConsumerConfiguration() *ConsumerConfiguration {
	return &ConsumerConfiguration{
		AutoCommit:  true,
		StoreOffset: true,
		OffsetReset: "earliest",
	}
}

// ConsumeHandler interface
type ConsumeHandler interface {
	Consume([]byte) error
}

// Consumer model
type Consumer struct {
	config          *_kafka.ConfigMap
	consumeListener []ConsumeHandler
}

// NewConsumer factory
func NewConsumer(config *ConsumerConfiguration) *Consumer {
	consumer := &Consumer{
		consumeListener: []ConsumeHandler{},
	}

	consumer.config = &_kafka.ConfigMap{
		"bootstrap.servers":         config.Broker,
		"auto.offset.reset":         "earliest",
		"broker.address.family":     "v4",
		"session.timeout.ms":        6000,
		"fetch.message.max.bytes":   18000000,
		"receive.message.max.bytes": 1000000000,
	}

	if config.GroupID != "" {
		(*consumer.config)["group.id"] = config.GroupID
	}

	return consumer
}

// Listen to the broker
func (c *Consumer) Listen(h ConsumeHandler, topics []string) error {
	kc, err := _kafka.NewConsumer(c.config)

	if err != nil {
		return err
	}

	kc.SubscribeTopics(topics, nil)
	defer kc.Close()

	for {
		msg, err := kc.ReadMessage(-1)
		if err != nil {
			log.Printf("Consumer error while reading: %v (%v)", err, msg)
			continue
		}

		if err := h.Consume(msg.Value); err != nil {
			log.Printf("Consumer error: %v", err)
		}
	}
}
