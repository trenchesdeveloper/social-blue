package mailer

import "embed"

const (
	FromName = "SocialBlue"
	MaxTries = 3
	UserWelcomeTemplate = "user_invitation"
)

//go:embed templates
var EmailTemplate embed.FS

type Client interface {
	Send(templateFile, username, email string, data any, isTest bool) error
}