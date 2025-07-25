package service

import (
	"time"

	"http_calendar/internal/models"
	"http_calendar/internal/storage"
)

// CalendarService provides business logic for creating, updating,
// deleting, and querying calendar events for users.
// It depends on an abstract storage layer and supports operations
// over day, week, and month periods.
type CalendarService struct {
	// repo is the underlying event storage implementation.
	repo storage.Storage
	// mondayBasedWeek indicates if weeks start on Monday (true) or Sunday (false).
	mondayBasedWeek bool
}

// NewCalendarService constructs a CalendarService.
// repo: implementation of storage.Storage for persisting events.
// isMondayBased: true to treat Monday as the first day of the week.
func NewCalendarService(repo storage.Storage, isMondayBased bool) *CalendarService {
	return &CalendarService{
		repo:            repo,
		mondayBasedWeek: isMondayBased,
	}
}

// CreateEvent creates a new event for the specified user.
// Delegates to repository SaveEvent.
func (s *CalendarService) CreateEvent(userId int, e models.Event) error {
	return s.repo.SaveEvent(userId, e)
}

// UpdateEvent updates an existing event for the specified user.
// Delegates to repository UpdateEvent.
func (s *CalendarService) UpdateEvent(userId int, e models.Event) error {
	return s.repo.UpdateEvent(userId, e)
}

// DeleteEvent removes an event by ID on the given date for the specified user.
// Delegates to repository DeleteEvent.
func (s *CalendarService) DeleteEvent(userId int, date models.Date, eventId int) error {
	return s.repo.DeleteEvent(userId, date, eventId)
}

// GetEventsForDay retrieves all events for a user on the specified date.
// It queries the storage for events in the [date; date] range.
func (s *CalendarService) GetEventsForDay(userId int, date models.Date) ([]models.Event, error) {
	// Use time.Time values for range boundaries
	return s.repo.GetEvents(userId, date.Time, date.Time)
}

// GetEventsForWeek retrieves all events for a user in the week of the given date.
// The week start is determined by mondayBasedWeek setting.
func (s *CalendarService) GetEventsForWeek(userId int, date models.Date) ([]models.Event, error) {
	weekday := int(date.Weekday())

	if s.mondayBasedWeek {
		if weekday == 0 {
			weekday = 7
		}
		weekday--
	}

	start := date.AddDate(0, 0, -weekday)
	end := start.AddDate(0, 0, 6)
	return s.repo.GetEvents(userId, start, end)
}

// GetEventsForMonth retrieves all events for a user in the month of the given date.
// It computes the first and last instants of the month.
func (s *CalendarService) GetEventsForMonth(userId int, date models.Date) ([]models.Event, error) {
	year, mon := date.Year(), date.Month()
	start := time.Date(year, mon, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0).Add(-time.Nanosecond)
	return s.repo.GetEvents(userId, start, end)
}
