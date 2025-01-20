package mailer

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
)

var templateFS embed.FS

type Mailer struct {
	sender string
	auth   smtp.Auth
}

func New(sender, username, password, host string) *Mailer {
	auth := smtp.PlainAuth("", username, password, host)
	return &Mailer{sender: sender, auth: auth}
}

func (m *Mailer) Send(recipient, templateFile string, data any) error {
	// Use the ParseFS() method to parse the required template file from the embedded
	// file system.
	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		log.Printf("failed to parse template: %v", err)
		return err
	}

	// Execute the named template "subject" and store the result in a bytes.Buffer variable.
	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		log.Printf("failed to execute subject template: %v", err)
		return err
	}

	// Follow the same pattern to execute the "plainBody" template and store the result
	// in the plainBody variable.
	plainBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		log.Printf("failed to execute plainBody template: %v", err)
		return err
	}

	// And likewise with the "htmlBody" template.
	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		log.Printf("failed to execute htmlBody template: %v", err)
		return err
	}

	// Create the email message
	msg := []byte("To: " + recipient + "\r\n" +
		"From: " + m.sender + "\r\n" +
		"Subject: " + subject.String() + "\r\n" +
		"MIME-version: 1.0;\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\r\n\r\n" +
		htmlBody.String() + "\r\n")

	// Send the email
	smtpServer := fmt.Sprintf("%s:%d", "smtp.mailtrap.io", 2525)
	err = smtp.SendMail(smtpServer, m.auth, m.sender, []string{recipient}, msg)
	if err != nil {
		log.Printf("failed to send email: %v", err)
		return err
	}

	return nil
}
