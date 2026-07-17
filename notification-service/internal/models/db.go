package models

import (
	"time"

	"github.com/google/uuid"
)

type MsgStatus string

const (
	StatusAccepted      MsgStatus = "accepted"
	StatusDelivered     MsgStatus = "delivered"
	StatusTemporaryFail MsgStatus = "temporary_fail"
	StatusPermanentFail MsgStatus = "permanent_fail"
	StatusOpened        MsgStatus = "opened"
	StatusClicked       MsgStatus = "clicked"
	StatusUnsubscribed  MsgStatus = "unsubscribed"
	StatusComplained    MsgStatus = "complained"
)

var MapMsgStatus = map[string]MsgStatus{
	"accepted":       StatusAccepted,
	"delivered":      StatusDelivered,
	"temporary_fail": StatusTemporaryFail,
	"permanent_fail": StatusPermanentFail,
	"opened":         StatusOpened,
	"clicked":        StatusClicked,
	"unsubscribed":   StatusUnsubscribed,
	"complained":     StatusComplained,
}

type Notification struct {
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	EventID   uuid.UUID `db:"event_id"`
	MessageID string    `db:"message_id"`
	Email     string    `db:"email"`
	EventType string    `db:"subject"`
	Status    MsgStatus `db:"status"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
