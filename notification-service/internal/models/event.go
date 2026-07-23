package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	SignUpEventType        = "sign_up"
	ResetPasswordEventType = "reset_password"
)

type Event struct {
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	Type      string    `json:"event_type"`
	Timestamp time.Time `json:"timestamp"`
}

type Message struct {
	To      string `json:"email"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}
