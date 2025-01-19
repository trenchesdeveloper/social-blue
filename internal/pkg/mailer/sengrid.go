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

func NewSendgrid(apikey, fromEmail string) *SendGridMailer {
	client := sendgrid.NewSendClient(apikey)

	return &SendGridMailer{
		fromEmail: fromEmail,
		apiKey:    apikey,
		client:    client,
	}
}

func (m *SendGridMailer) Send(templateFile, username, email string, data any, isTest bool) error {
	from := mail.NewEmail(FromName, m.fromEmail)
	to := mail.NewEmail(username, email)

	// create email template and building
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

	message := mail.NewSingleEmail(from, subject.String(), to, "", body.String())

	message.SetMailSettings(&mail.MailSettings{
		SandboxMode: &mail.Setting{
			Enable: &isTest,
		},
	})

	for i := 0; i < MaxTries; i++ {
		resp, err := m.client.Send(message)
		if err != nil {
			log.Printf("Failed to send to email %v, attempt %d of %d", email, i+1, MaxTries)
			log.Printf("Error: %v", err.Error())
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}

		log.Printf("Email sent to %v, status code: %v", email, resp.StatusCode)
		log.Printf("Response body: %v", resp.Body)

		return nil
	}

	return fmt.Errorf("failed to send email to %v after %d attempts", email, MaxTries)
}
