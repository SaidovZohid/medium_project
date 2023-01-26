package notification

import (
	"gitlab.com/medium-project/medium_notification_service_with_kafka/config"
	event "gitlab.com/medium-project/medium_notification_service_with_kafka/pkg/messagebroker"
)

type NotificationService struct {
	cfg   config.Config
	kafka *event.Kafka
}

func New(cfg config.Config, kafka *event.Kafka) *NotificationService {
	return &NotificationService{
		cfg:   cfg,
		kafka: kafka,
	}
}

func (c *NotificationService) RegisterConsumers() {
	notificationRoute := "v1.notification_service.send_email"

	c.kafka.AddConsumer(
		notificationRoute,
		notificationRoute,
		c.SendEmail,
	)
}
