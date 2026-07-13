package models

import (
	"time"

	"github.com/google/uuid"
)

type UserSignUpEvent struct {
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	TimeStamp time.Time `json:"timestamp"`
}
