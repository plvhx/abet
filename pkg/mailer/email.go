package mailer

import (
    "context"
)

type SendEmail struct {
    Subject     string   `json:"subject"`
    Body        string   `json:"body"`
    ToEmail     string   `json:"toEmail,omitempty"`
    ToListEmail []string `json:"toListEmail,omitempty"`
}

type EmailClient interface {
    SendEmail(ctx context.Context, data SendEmail) error
}
