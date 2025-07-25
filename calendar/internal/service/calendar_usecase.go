package service

import "http_calendar/internal/models"

// CalendarUsecase defines the business‑logic operations on calendar events.
// It abstracts over any storage or transport layer.
type CalendarUsecase interface {
	// CreateEvent creates a new event for userId.
	// Returns an error if saving to storage fails.
	CreateEvent(userId int, e models.Event) error

	// UpdateEvent updates an existing event for userId.
	// Returns an error if the event does not exist or update fails.
	UpdateEvent(userId int, e models.Event) error

	// DeleteEvent deletes the event with eventId on the given date for userId.
	// Returns an error if the event is not found.
	DeleteEvent(userId int, date models.Date) error

	// GetEventsForDay returns all events for userId on the specified date.
	// If no events exist, returns an empty slice and no error.
	GetEventsForDay(userId int, date models.Date) ([]models.Event, error)

	// GetEventsForWeek returns all events for userId in the week containing date.
	// Week start (Sunday or Monday) is determined by service configuration.
	GetEventsForWeek(userId int, date models.Date) ([]models.Event, error)

	// GetEventsForMonth returns all events for userId in the month containing date.
	GetEventsForMonth(userId int, date models.Date) ([]models.Event, error)
}
