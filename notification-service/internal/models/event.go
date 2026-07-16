package models

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID        uuid.UUID `json:"id"`
	Topic     string    `json:"topic"`
	Key       string    `json:"key"`
	Type      string    `json:"type"`
	Payload   UserData  `json:"payload"`
	CreatedAt time.Time `json:"created_at"`
}

type UserData struct {
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type Message struct {
	To   string `json:"email"`
	Name string `json:"name"`
	Body string `json:"body"`
}
