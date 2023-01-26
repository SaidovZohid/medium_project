package events

import (
	"context"

	"gitlab.com/medium-project/medium_notification_service_with_kafka/config"
	"gitlab.com/medium-project/medium_notification_service_with_kafka/events/notification"
	messagebroker "gitlab.com/medium-project/medium_notification_service_with_kafka/pkg/messagebroker"
)

// PubsubServer ...
type PubsubServer struct {
	cfg   config.Config
	kafka *messagebroker.Kafka
}

// New ...
func New(cfg config.Config) (*PubsubServer, error) {
	kafka, err := messagebroker.NewKafka(cfg)
	if err != nil {
		return nil, err
	}

	return &PubsubServer{
		cfg:   cfg,
		kafka: kafka,
	}, nil
}

// Run ...
func (s *PubsubServer) Run(ctx context.Context) {
	ntService := notification.New(s.cfg, s.kafka)
	ntService.RegisterConsumers()

	s.kafka.RunCustomer(ctx)
}
