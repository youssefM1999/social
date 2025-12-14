package mailer

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
)

const (
	FromName                   = "GopherSocial"
	maxRetries                 = 3
	UserInvitationTemplateFile = "user_invitation.tmpl"
)

//go:embed templates
var FS embed.FS

type Client interface {
	Send(templateFile, username, email string, data any, isSandbox bool) (int, error)
}

type MailData struct {
	Subject string
	Body    string
}

func buildMail(templateFile string, data any) (MailData, error) {
	fullPath := "templates/" + templateFile
	tmpl, err := template.ParseFS(FS, fullPath)
	if err != nil {
		return MailData{}, fmt.Errorf("failed to parse template %s: %w", templateFile, err)
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return MailData{}, err
	}

	body := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(body, "body", data)
	if err != nil {
		return MailData{}, err
	}

	return MailData{Subject: subject.String(), Body: body.String()}, nil
}
