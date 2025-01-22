package email

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"net/smtp"
)

type EmailConfig struct {
	SMTPHost string
	SMTPPort int
	Username string
	Password string
	From     string
}
type IEmailService interface {
	SendEmail(ctx context.Context, to string, subject string, content []byte, data any) error
}

type EmailService struct {
	SMTPHost string
	SMTPPort int
	Username string
	Password string
	From     string
}

func (s *EmailService) SendEmail(ctx context.Context, to string, subject string, content []byte, data interface{}) error {
	tmpl, err := template.New("email-template").Parse(string(content))
	if err != nil {
		return err
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return err
	}

	msg := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\nMIME-Version: 1.0\nContent-Type: text/html; charset=\"UTF-8\"\n\n%s",
		s.From, to, subject, body.String())

	auth := smtp.PlainAuth("", s.Username, s.Password, s.SMTPHost)
	addr := fmt.Sprintf("%s:%d", s.SMTPHost, s.SMTPPort)

	if err := smtp.SendMail(addr, auth, s.From, []string{to}, []byte(msg)); err != nil {
		return err
	}

	return nil
}

func NewEmailService(config *EmailConfig) IEmailService {
	return &EmailService{
		SMTPHost: config.SMTPHost,
		SMTPPort: config.SMTPPort,
		Username: config.Username,
		Password: config.Password,
		From:     config.From,
	}
}
