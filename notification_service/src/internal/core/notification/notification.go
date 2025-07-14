package notifications

import "time"

type NotificationData struct {
	ID        string    `json:"id"`
	Channel   string    `json:"channel"`
	Payload   string    `json:"payload"`
	Timestamp time.Time `json:"timestamp"`
	Status    string    `json:"status"`
}
