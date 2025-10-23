// Package mailer provides functionality for sending emails.
package mailer

import (
	"fmt"

	"github.com/wneessen/go-mail"
)

const (
	maxRetries = 3
)

type Mail struct {
	Sender EmailSender
}

type EmailSender interface {
	Send(email Message) error
}

type Message struct {
	From          string
	FromEmail     string
	Subject       string
	To            []string
	CC            []string
	BCC           []string
	Data          any
	Message       map[string]any
	AttachFiles   []string
	AttachmentMap map[string]string
	Template      string
}

type MailConf struct {
	EmailService   string `json:"email_service" validate:"required,min=2,max=100"`
	SenderName     string `json:"sender_name" validate:"required,min=2,max=100"`
	SenderEmail    string `json:"sender_email" validate:"required,email"`
	EmailLogin     string `json:"email_login" validate:"required,min=2,max=100"`
	EmailPassword  string `json:"email_password" validate:"required,min=2,max=100"`
	PathToTemplate string `json:"path_to_template" validate:"required,min=2,max=100"`
}

func NewMailSender(conf MailConf) (EmailSender, error) {
	switch conf.EmailService {
	case "mailtrap":
		return &MailTrapSender{
			From:        conf.SenderName,
			FromEmail:   conf.SenderEmail,
			Login:       conf.EmailLogin,
			Password:    conf.EmailPassword,
			SMTPHost:    "sandbox.smtp.mailtrap.io",
			SMTPPort:    2525,
			SMTPAuth:    mail.SMTPAuthPlain,
			TemplateDir: conf.PathToTemplate,
		}, nil
	default:
		return nil, fmt.Errorf("not implemented any other mail service except mailtrap")
	}
}
