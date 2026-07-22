package models

import (
	"time"

	"github.com/google/uuid"
)

type AuthEventType string

const (
	AuthEventSignUp        AuthEventType = "sign_up"
	AuthEventResetPassword AuthEventType = "reset_password"
)

type UserAuthEvent struct {
	UserID    uuid.UUID      `json:"user_id"`
	Email     string         `json:"email"`
	EventType AuthEventType  `json:"event_type"`
	TimeStamp time.Time      `json:"timestamp"`
}
