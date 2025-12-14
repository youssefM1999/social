package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
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

func (m *SendGridMailer) Send(templateFile, username, email string, data any, isSandbox bool) error {
	from := mail.NewEmail(FromName, m.fromEmail)
	to := mail.NewEmail(username, email)

	// template parsing and building
	fullPath := "templates/" + templateFile
	tmpl, err := template.ParseFS(FS, fullPath)
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", templateFile, err)
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	body := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(body, "body", data)
	if err != nil {
		return err
	}

	// Extract activation URL from data for plain text version
	var plainTextBody string
	if vars, ok := data.(struct {
		Username      string
		ActivationURL string
	}); ok {
		plainTextBody = fmt.Sprintf("Hi %s,\n\nThanks for signing up for GopherSocial. We're excited to have you on board!\n\nBefore you can start using GopherSocial, you need to confirm your email address. Click the link below to confirm your email address:\n\n%s\n\nIf you want to activate your account manually copy and paste the code from the link above.\n\nIf you didn't sign up for GopherSocial, you can safely ignore this email.\n\nThanks,\nThe GopherSocial Team", vars.Username, vars.ActivationURL)
	} else {
		// Fallback if data structure doesn't match
		plainTextBody = body.String()
	}

	message := mail.NewSingleEmail(from, subject.String(), to, plainTextBody, body.String())

	message.SetMailSettings(&mail.MailSettings{
		SandboxMode: &mail.Setting{
			Enable: &isSandbox,
		},
	})

	if isSandbox {
		log.Printf("WARNING: Sandbox mode is ENABLED - emails will NOT be sent, only validated")
	}

	for i := 0; i < maxRetries; i++ {
		response, err := m.client.Send(message)
		if err != nil {
			log.Printf("Failed to send email: %v, attempt %d of %d", err, i+1, maxRetries)
			log.Printf("Error: %s", err.Error())

			// exponential backoff
			time.Sleep(time.Duration(i+1) * time.Second)
			continue
		}

		// Check response status code
		log.Printf("SendGrid response: status=%d, body=%s", response.StatusCode, response.Body)

		if response.StatusCode >= 200 && response.StatusCode < 300 {
			if isSandbox {
				log.Printf("Email validated successfully (sandbox mode) - status code: %v", response.StatusCode)
			} else {
				log.Printf("Email sent successfully - status code: %v", response.StatusCode)
			}
			return nil
		}

		// Log error details
		log.Printf("SendGrid API error: status %v, body: %s", response.StatusCode, response.Body)

		// Retry on server errors (5xx)
		if response.StatusCode >= 500 {
			log.Printf("Server error from SendGrid, retrying... attempt %d of %d", i+1, maxRetries)
			time.Sleep(time.Duration(i+1) * time.Second)
			continue
		}

		// Don't retry on client errors (4xx) - return immediately
		return fmt.Errorf("SendGrid API error: status code %v, body: %s", response.StatusCode, response.Body)
	}

	return fmt.Errorf("failed to send email after %d attempts", maxRetries)
}
