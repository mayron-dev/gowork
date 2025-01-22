package tasks

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/mayron-dev/gowork/config"
)

type EmailDeliveryPayload struct {
	To           string
	Subject      string
	TemplatePath string
	Data         any
	SMTP         struct {
		Host     string
		Port     int
		Username string
		Password string
		From     string
	}
}

func HandleEmailDeliveryTask(ctx context.Context, t *asynq.Task) error {
	var p EmailDeliveryPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	// get template
	env := config.GetEnv()
	body, err := fileService.Download(ctx, env.EMAIL_TEMPLATES_BUCKET, p.TemplatePath)
	if err != nil {
		return fmt.Errorf("fileService.Download failed: %v: %w", err, asynq.SkipRetry)
	}

	// send email
	if err := emailService.SendEmail(ctx, p.To, p.Subject, body, p.Data); err != nil {
		return fmt.Errorf("emailService.SendEmail failed: %v: %w", err, asynq.SkipRetry)
	}

	return nil
}
