package event

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/Shopify/sarama"
	"gitlab.com/medium-project/medium_notification_service_with_kafka/config"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

type Kafka struct {
	cfg          config.Config
	consumers    map[string]*Consumer
	saramaConfig *sarama.Config
}

func NewKafka(cfg config.Config) (*Kafka, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Version = sarama.V2_0_0_0

	kafka := &Kafka{
		cfg:          cfg,
		consumers:    make(map[string]*Consumer),
		saramaConfig: saramaConfig,
	}

	return kafka, nil
}

func (r *Kafka) RunCustomer(ctx context.Context) {
	var w sync.WaitGroup

	for _, consumer := range r.consumers {
		w.Add(1)
		go func(w *sync.WaitGroup, c *Consumer) {
			defer w.Done()

			err := c.cloudEventClient.StartReceiver(context.Background(), func(ctx context.Context, event cloudevents.Event) {
				err := c.handler(ctx, event)
				if err != nil {
					log.Println(err)
				}
			})

			log.Panic("failed to start customer:", err)
		}(&w, consumer)
		fmt.Println("Key:", consumer.topic, "=>", "customer:", consumer)
	}
	w.Wait()
}
