package mail

import (
	"fmt"
	"html/template"

	"github.com/wneessen/go-mail"
)

const (
	smtpGmailAddress = "smtp.gmail.com"
)

type EmailSender interface {
	SendEmail(
		templateName, subject, content string,
		to, attachFiles []string,
	) error
}

type GmailSender struct {
	name             string
	fromEmailAddress string
	mailClient       *mail.Client
}

func NewGmailSender(name, smtpUser, smtpPass string) (EmailSender, error) {
	client, err := mail.NewClient(smtpGmailAddress,
		mail.WithSMTPAuth(mail.SMTPAuthAutoDiscover),
		mail.WithTLSPortPolicy(mail.TLSMandatory),
		mail.WithUsername(smtpUser),
		mail.WithPassword(smtpPass),
	)
	if err != nil {
		return nil, err
	}

	return &GmailSender{
		name:             name,
		fromEmailAddress: smtpUser,
		mailClient:       client,
	}, nil
}

func (s *GmailSender) SendEmail(
	templateName, subject, content string,
	to, attachFiles []string,
) error {
	msg := mail.NewMsg()

	msg.Subject(subject)
	if err := msg.FromFormat(s.name, s.fromEmailAddress); err != nil {
		return fmt.Errorf("failed to set formatted FROM address: %w", err)
	}
	if err := msg.To(to...); err != nil {
		return fmt.Errorf("failed to add recipient: %w", err)
	}

	htmlTpl, err := template.New(templateName).Parse(content)
	if err != nil {
		return fmt.Errorf("failed to parse text template: %w", err)
	}
	err = msg.SetBodyHTMLTemplate(htmlTpl, nil)
	if err != nil {
		return fmt.Errorf("failed to set html template: %w", err)
	}

	for _, filePath := range attachFiles {
		msg.AttachFile(filePath)
	}

	if err := s.mailClient.DialAndSend(msg); err != nil {
		return fmt.Errorf("failed to deliver mail: %w", err)
	}

	return nil
}
