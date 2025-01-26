package mailer

import (
	"bytes"
	"fmt"
	"text/template"

	gomail "gopkg.in/mail.v2"
)

type mailtrapClient struct {
	formEmail string
	apiKey    string
}

func NewMailtrap(apiKey, fromEmail string) (*mailtrapClient, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("mailtrap api key is required")
	}
	return &mailtrapClient{
		formEmail: fromEmail,
		apiKey:    apiKey,
	}, nil
}

func (m *mailtrapClient) Send(templateFile, username, email string, data any, isTest bool) error {
	// template parsing
	tmpl, err := template.ParseFS(EmailTemplate, "templates/"+templateFile)
	if err != nil {
		return fmt.Errorf("failed to parse email template: %v", err)
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return fmt.Errorf("failed to execute email template: %v", err)
	}

	body := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(body, "body", data)
	if err != nil {
		return fmt.Errorf("failed to execute email template: %v", err)
	}

	message := gomail.NewMessage()
	message.SetHeader("From", m.formEmail)
	message.SetHeader("To", email)
	message.SetHeader("Subject", subject.String())

	message.AddAlternative("text/html", body.String())

	dialer := gomail.NewDialer("smtp.mailtrap.io", 587, "api", m.apiKey)

	if err := dialer.DialAndSend(message); err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil

}