package mailer

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/vanng822/go-premailer/premailer"
	"github.com/wneessen/go-mail"
)

type MailTrapSender struct {
	From        string
	FromEmail   string
	Login       string
	Password    string
	SMTPHost    string
	SMTPPort    int
	SMTPAuth    mail.SMTPAuthType
	TemplateDir string
}

func (sender *MailTrapSender) Send(email Message) error {

	if email.Template == "" {
		email.Template = "mail"
	}
	if email.From == "" {
		email.From = sender.From
	}
	if email.FromEmail == "" {
		email.FromEmail = sender.FromEmail
	}

	if email.AttachmentMap == nil {
		email.AttachmentMap = make(map[string]string)
	}

	m := mail.NewMsg()
	if err := m.FromFormat(email.From, email.FromEmail); err != nil {
		return fmt.Errorf("failed to set from address: %s", err)
	}
	if err := m.To(email.To...); err != nil {
		return fmt.Errorf("failed to set To address: %s", err)
	}
	if err := m.Cc(email.CC...); err != nil {
		return fmt.Errorf("failed to set CC address: %s", err)
	}
	if err := m.Bcc(email.BCC...); err != nil {
		return fmt.Errorf("failed to set BCC address: %s", err)
	}
	m.Subject(email.Subject)

	data := map[string]any{
		"message": email.Data,
	}
	email.Message = data

	// generate and set to the message text plain body
	templPlain := fmt.Sprintf("%s/%s.plain.gohtml", sender.TemplateDir, email.Template)
	contentPlain, err := builPlainTextMessage(templPlain, email.Message)
	if err != nil {
		return fmt.Errorf("failed to generate plain text message: %s", err)
	}
	m.SetBodyString(mail.TypeTextPlain, contentPlain)
	// generate and set to the message alternative html formated body
	templFormated := fmt.Sprintf("%s/%s.html.gohtml", sender.TemplateDir, email.Template)
	contentHTML, err := buildHTMLMessage(templFormated, email.Message)
	if err != nil {
		return fmt.Errorf("failed to generate html formated message: %s", err)
	}
	m.AddAlternativeString(mail.TypeTextHTML, contentHTML)

	for _, file := range email.AttachFiles {
		m.AttachFile(file)
	}

	for key, value := range email.AttachmentMap {
		m.AttachFile(value, mail.WithFileName(key))
	}

	c, err := mail.NewClient(
		sender.SMTPHost,
		mail.WithPort(sender.SMTPPort),
		mail.WithSMTPAuth(sender.SMTPAuth),
		mail.WithUsername(sender.Login),
		mail.WithPassword(sender.Password),
	)
	if err != nil {
		return fmt.Errorf("failed to create mail client: %s", err)
	}

	if err = c.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func buildHTMLMessage(templ string, message map[string]any) (string, error) {
	t, err := template.New("email-html").ParseFiles(templ)
	if err != nil {
		return "", fmt.Errorf("failed to create template from %s: %s", templ, err)
	}

	var tpl bytes.Buffer

	if err := t.ExecuteTemplate(&tpl, "body", message); err != nil {
		return "", fmt.Errorf("failed execute template with message %v: %s", message, err)
	}

	formattedMessage, err := inlineCSS(tpl.String())
	if err != nil {
		return "", fmt.Errorf("failed generate inline CSS message from template: %s", err)
	}
	return formattedMessage, nil
}

func inlineCSS(fm string) (string, error) {
	options := premailer.Options{
		RemoveClasses:     false,
		CssToAttributes:   false,
		KeepBangImportant: true,
	}

	prem, err := premailer.NewPremailerFromString(fm, &options)
	if err != nil {
		return "", fmt.Errorf("failed create premailer from string %s: %s", fm, err)
	}

	html, err := prem.Transform()
	if err != nil {
		return "", fmt.Errorf("failed transform premailer to string: %s", err)
	}
	return html, nil
}

func builPlainTextMessage(templ string, message map[string]any) (string, error) {
	t, err := template.New("email-plain").ParseFiles(templ)

	if err != nil {
		return "", fmt.Errorf("failed to create template from %s: %s", templ, err)
	}

	var tpl bytes.Buffer

	if err := t.ExecuteTemplate(&tpl, "body", message); err != nil {
		return "", fmt.Errorf("failed execute template with message %v: %s", message, err)
	}

	return tpl.String(), nil
}
