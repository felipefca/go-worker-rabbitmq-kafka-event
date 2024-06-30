package models

import "time"

type Event struct {
	Message   string
	Date      time.Time
	CreatedBy string
}

func BuildEvent(message string, createdBy string) *Event {
	return &Event{
		Message:   message,
		Date:      time.Now(),
		CreatedBy: createdBy,
	}
}
