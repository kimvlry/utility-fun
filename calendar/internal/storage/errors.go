package storage

import (
    "errors"
    "fmt"
    "http_calendar/internal/lib/models"
)

var (
    errEventExists   = errors.New("event already exists")
    errEventNotFound = errors.New("event not found")
    errNoEvents      = errors.New("user has no events")
)

// NewEventExistsError returns an error indicating that the event already exists
func NewEventExistsError(id string) error {
    return fmt.Errorf("%w: %d", errEventExists, id)
}

// NewEventNotFoundError returns an error indicating that the specified event was not found
func NewEventNotFoundError(id string) error {
    return fmt.Errorf("%w: %d", errEventNotFound, id)
}

// NewEventNotFoundByDateError returns an error indicating that no event was found
// for the given date and user ID.
func NewEventNotFoundByDateError(date models.Date, userId string) error {
    return fmt.Errorf("%w: %+v, user: %d", errEventNotFound, date, userId)
}

// NewUserHasNoEventsError returns an error indicating that no user was found for given ID
func NewUserHasNoEventsError(userId string) error {
    return fmt.Errorf("%w: %d", errNoEvents, userId)
}
