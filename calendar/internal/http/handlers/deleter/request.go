package deleter

import "http_calendar/internal/lib/models"

type Request struct {
	UserId  string      `json:"user_id"    validate:"required"`
	Date    models.Date `json:"date"       validate:"ISO8601date"`
	EventId string      `json:"event_id"   validate:"required"`
}
