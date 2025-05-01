package kafka

import (
	"errors"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/fentezi/translator/config"
	"github.com/fentezi/translator/internal/entities"
)

const (
	flushTimeout = 5000 // ms
)

var (
	ErrUnknownType = errors.New("unknown event type")
)

type Producer struct {
	producer *kafka.Producer
}

func New(broker config.Kafka) (*Producer, error) {
	address := toAddress(broker)
	conf := &kafka.ConfigMap{
		"bootstrap.servers": address,
	}

	p, err := kafka.NewProducer(conf)
	if err != nil {
		return nil, fmt.Errorf("error creating kafka producer: %w", err)
	}

	return &Producer{producer: p}, nil
}

func (p *Producer) Produce(msg entities.KafkaMessage) error {
	kafkaMsg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &msg.TopicPartition,
			Partition: kafka.PartitionAny,
		},
		Value: msg.Value,
	}

	deliveryChan := make(chan kafka.Event, 1)
	defer close(deliveryChan)

	if err := p.producer.Produce(kafkaMsg, deliveryChan); err != nil {
		return fmt.Errorf("error sending message to kafka producer: %w", err)
	}

	select {
	case e := <-deliveryChan:
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				return fmt.Errorf("delivery failed: %w", ev.TopicPartition.Error)
			}
			return nil
		default:
			return ErrUnknownType
		}
	case <-time.After(5 * time.Second):
		return errors.New("timeout waiting for delivery confirmation")
	}
}

func (p *Producer) Close() {
	p.producer.Flush(flushTimeout)
	p.producer.Close()
}

func toAddress(broker config.Kafka) string {
	address := fmt.Sprintf("%s:%s", broker.Address, broker.Port)
	return address
}
