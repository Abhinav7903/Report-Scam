package mail

import (
	"fmt"
	"log"
	"net/smtp"
	"strings"
)

type Mail struct {
	From     string
	Password string
	AppPass  string
}

// NewMail creates a new Mail instance with the sender's credentials
func NewMail(from, password, appPass string) *Mail {
	return &Mail{
		From:     from,
		Password: password,
		AppPass:  appPass,
	}
}

// SendMail sends an email to the specified recipient(s)
// It prioritizes using the App Password if available.
func (m *Mail) SendMail(to, subject, body string) error {
	// Use AppPass if provided, otherwise fall back to Password
	authPassword := m.AppPass
	if authPassword == "" {
		authPassword = m.Password
	}

	// Log the sender's email for debugging (do not log sensitive credentials)
	fmt.Printf("Attempting to send email from: %s\n", m.From)
	log.Println("Using App Password: ", m.AppPass != "")

	// Create the SMTP authentication
	auth := smtp.PlainAuth("", m.From, authPassword, "smtp.gmail.com")

	// Create the email message
	msg := fmt.Sprintf(
		"From: %s\nTo: %s\nSubject: %s\n\n%s",
		m.From,
		to,
		subject,
		body,
	)

	// Split the recipient addresses if there are multiple
	recipients := strings.Split(to, ",")

	// Use smtp.SendMail to send the email
	err := smtp.SendMail(
		"smtp.gmail.com:587", // SMTP server and port
		auth,                 // Authentication
		m.From,               // Sender email
		recipients,           // Recipient email(s)
		[]byte(msg),          // Email body as bytes
	)

	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	fmt.Println("Email sent successfully")
	return nil
}
