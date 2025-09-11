package helpers

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"path/filepath"

	"github.com/mailjet/mailjet-apiv3-go/v4"
)

type EmailData struct {
	Email    string
	Username string
	Link     string
}

func SendEmail(toEmail, username, verifyLink, templateName, subject string) error {
	client := GetMailjetClient()
	sender := GetEmailSender()

	tmpl, err := template.ParseFiles(filepath.Join("shared/templates", templateName+".html"))
	if err != nil {
		return fmt.Errorf("failed to parse email template: %w", err)
	}

	var body bytes.Buffer
	data := EmailData{
		Email:    toEmail,
		Username: username,
		Link:     verifyLink,
	}
	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("failed to render email template: %w", err)
	}

	message := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: sender,
				Name:  "no-reply",
			},
			To: &mailjet.RecipientsV31{
				{Email: toEmail},
			},
			Subject:  fmt.Sprintf("%s - Puspa HIC", subject),
			HTMLPart: body.String(),
		},
	}

	messages := mailjet.MessagesV31{Info: message}
	_, err = client.SendMailV31(&messages)
	if err != nil {
		log.Printf("Failed to send email to %s: %v", toEmail, err)
		return fmt.Errorf("failed to send email to %s: %w", toEmail, err)
	}

	return nil
}
