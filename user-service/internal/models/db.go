package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `db:"id"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	Role         string    `db:"role"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type OutboxEvent struct {
	ID        uuid.UUID `db:"id"`
	Topic     string    `db:"topic"`
	Key       string    `db:"key"`
	Payload   []byte    `db:"payload"`
	CreatedAt time.Time `db:"created_at"`
	Published bool      `db:"published"`
}

func NewOutBoxEvent(topic string, key string, payload []byte) *OutboxEvent {
	return &OutboxEvent{
		ID:        uuid.New(),
		Topic:     topic,
		Key:       key,
		Payload:   payload,
		CreatedAt: time.Now(),
	}
}
