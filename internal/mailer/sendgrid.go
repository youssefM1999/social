package mailer

import (
	"fmt"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/youssefM1999/social/internal/retry"
)

type SendGridMailer struct {
	fromEmail string
	apiKey    string
	client    *sendgrid.Client
}

func NewSendgrid(apiKey, fromEmail string) *SendGridMailer {
	client := sendgrid.NewSendClient(apiKey)
	return &SendGridMailer{
		fromEmail: fromEmail,
		apiKey:    apiKey,
		client:    client,
	}
}

func (m *SendGridMailer) Send(templateFile, username, email string, data any, isSandbox bool) (int, error) {
	from := mail.NewEmail(FromName, m.fromEmail)
	to := mail.NewEmail(username, email)

	// template parsing and building
	mailData, err := buildMail(templateFile, data)
	if err != nil {
		return -1, err
	}

	message := mail.NewSingleEmail(from, mailData.Subject, to, "", mailData.Body)

	message.SetMailSettings(&mail.MailSettings{
		SandboxMode: &mail.Setting{
			Enable: &isSandbox,
		},
	})

	var statusCode int
	err = retry.Retry(func() error {
		response, err := m.client.Send(message)
		statusCode = response.StatusCode
		return err
	}, maxRetries)
	if err != nil {
		return statusCode, fmt.Errorf("failed to send email: %w after %d attempts", err, maxRetries)
	}
	return statusCode, nil
}
