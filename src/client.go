package src

import (
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// ProducerConfiguration model
type ProducerConfiguration struct {
	Broker string
	Topic  string
}

// Producer model
type Producer struct {
	config   *kafka.ConfigMap
	producer *kafka.Producer
	topic    string
}

// NewProducer factory
func NewProducer(config *ProducerConfiguration) (*Producer, error) {
	res := &Producer{
		config: &kafka.ConfigMap{
			"bootstrap.servers": config.Broker,
		},
		topic: config.Topic,
	}

	p, err := kafka.NewProducer(res.config)

	if err != nil {
		return nil, err
	}

	res.producer = p
	return res, nil
}

// SendJSON transforms to json
func (p *Producer) SendJSON(obj interface{}) error {
	data, _ := json.Marshal(obj)
	return p.Send(data)
}

// Send a message to kafka
func (p *Producer) Send(message []byte) error {
	p.producer.ProduceChannel() <- &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &p.topic,
			Partition: kafka.PartitionAny,
		},
		Value: message,
	}

	return nil
}
