// Package util provides utility functions and types for common operations.
package util

import (
	"io"
	"net/smtp"
	"strings"
)

// EmailRelay provides a simple interface for sending emails through an SMTP relay server.
// It handles the connection to the SMTP server and the email sending process.
type EmailRelay struct {
	relayHost string // The hostname of the SMTP relay server (without port)
}

// NewEmailRelay creates a new EmailRelay instance with the specified SMTP relay host.
//
// Parameters:
//   - relayHost: The hostname of the SMTP relay server (without port)
//
// Returns:
//   - A pointer to a new EmailRelay instance
func NewEmailRelay(relayHost string) *EmailRelay {
	return &EmailRelay{relayHost: relayHost}
}

// Send sends an email through the configured SMTP relay server.
// It establishes a connection to the server, sets the sender and recipient,
// sends the email body, and properly closes the connection.
//
// Parameters:
//   - from: The email address of the sender
//   - to: The email address of the recipient
//   - body: The full email body, including headers and content
//
// Returns:
//   - An error if any step of the email sending process fails, nil otherwise
func (r *EmailRelay) Send(from, to, body string) error {
	// Connect to the remote SMTP server.
	c, err := smtp.Dial(r.relayHost + ":587")
	if err != nil {
		return err
	}

	// Set the sender and recipient first
	if err = c.Mail(from); err != nil {
		return err
	}
	if err = c.Rcpt(to); err != nil {
		return err
	}

	// Send the email body.
	wc, err := c.Data()
	if err != nil {
		return err
	}

	if _, err = io.Copy(wc, strings.NewReader(body)); err != nil {
		return err
	}
	if err = wc.Close(); err != nil {
		return err
	}

	// Send the QUIT command and close the connection.
	if err = c.Quit(); err != nil {
		return err
	}
	return nil
}
