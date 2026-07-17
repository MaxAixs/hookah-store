package models

import (
	"time"

	"github.com/google/uuid"
)

type NotificationResponse struct {
	ID        string    `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	EventType string    `json:"event_type"`
	Status    MsgStatus `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
