package notification

import (
	"context"
	"encoding/json"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	emailPkg "gitlab.com/medium-project/medium_notification_service_with_kafka/pkg/email"
)

type Notification struct {
	To      string
	Type    string
	Body    map[string]string
	Subject string
}

func (c *NotificationService) SendEmail(ctx context.Context, event cloudevents.Event) error {
	var (
		req Notification
	)

	err := json.Unmarshal(event.DataEncoded, &req)

	if err != nil {
		return err
	}

	err = emailPkg.SendEmail(&c.cfg, &emailPkg.SendEmailRequest{
		To:      []string{req.To},
		Subject: req.Subject,
		Body:    req.Body,
		Type:    req.Type,
	})

	if err != nil {
		return err
	}

	return nil
}
