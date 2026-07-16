package models

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	pending    Status = "pending"
	processing Status = "processing"
	failed     Status = "failed"
)

type Notification struct {
	ID        string    `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	Email     string    `db:"email"`
	EventType string    `db:"subject"`
	Status    Status    `db:"status"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
