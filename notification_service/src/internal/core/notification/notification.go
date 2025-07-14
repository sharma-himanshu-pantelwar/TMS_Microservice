package notifications

import "time"

type Notification struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Message   string    `json:"message"`
	Channel   string    `json:"channel"`
	CreatedAt time.Time `json:"created_at"`
}

type NotificationServiceImpl interface {
	ProcessNotification(notification Notification) error
}
