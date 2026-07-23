package models

type MailgunWebhook struct {
	Signature Signature `json:"signature"`
	EventData EventData `json:"event-data"`
}

type Signature struct {
	Token     string `json:"token"`
	Signature string `json:"signature"`
	Timestamp string `json:"timestamp"`
}

type EventData struct {
	Event     string         `json:"event"`
	Timestamp float64        `json:"timestamp"`
	Message   MailgunMessage `json:"message"`
}

type MailgunMessage struct {
	Headers Headers `json:"headers"`
}

type Headers struct {
	MessageID string `json:"message-id"`
	To        string `json:"to"`
}
