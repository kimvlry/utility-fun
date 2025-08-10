package storage

import (
    "http_calendar/internal/lib/models"
    "time"
)

// Storage is the low‑level data store interface for calendar events.
// It knows nothing about HTTP or business rules — only how to persist and retrieve events.
type Storage interface {
    // SaveEvent stores a new event for userId.
    // Returns an error if an event with the same Id already exists.
    SaveEvent(userId string, e models.Event) (models.Event, error)

    // UpdateEvent updates an existing event for userId.
    // Returns an error if the event does not exist.
    UpdateEvent(userId string, e models.Event) (models.Event, error)

    // DeleteEvent deletes the event with eventId for userId on the specified date.
    // Returns an error if the event or user is not found.
    DeleteEvent(userId string, date models.Date, eventId string) error

    // GetEvents returns all events for userId between from and to inclusive.
    // The `from` and `to` parameters are time.Time values; events whose Date fall within
    // that range (date-only precision) will be returned.
    GetEvents(userId string, from, to time.Time) ([]models.Event, error)
}
