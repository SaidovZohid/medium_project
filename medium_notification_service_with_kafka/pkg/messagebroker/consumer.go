package event

import (
	"context"
	"errors"

	"github.com/cloudevents/sdk-go/protocol/kafka_sarama/v2"
	cloudevents "github.com/cloudevents/sdk-go/v2"
)

type HandlerFunc func(context.Context, cloudevents.Event) error

type Consumer struct {
	consumerName     string
	topic            string
	handler          HandlerFunc
	cloudEventClient cloudevents.Client
}

// Add Consumer ...
func (kafka *Kafka) AddConsumer(consumerName, topic string, handler HandlerFunc) {
	if kafka.consumers[consumerName] != nil {
		panic(errors.New("consumer with the same name already exists: " + consumerName))
	}

	receiver, err := kafka_sarama.NewConsumer(
		[]string{kafka.cfg.KafkaUrl}, // Kafka connection url
		kafka.saramaConfig,           // Sarama config
		"notification_service",       // Group ID
		topic,                        // Topic
	)

	if err != nil {
		panic(err)
	}

	c, err := cloudevents.NewClient(receiver)

	if err != nil {
		panic(err)
	}

	kafka.consumers[consumerName] = &Consumer{
		consumerName:     consumerName,
		topic:            topic,
		handler:          handler,
		cloudEventClient: c,
	}

}
