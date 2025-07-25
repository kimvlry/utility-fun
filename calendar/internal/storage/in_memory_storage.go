package storage

import (
	"sync"
	"time"

	"http_calendar/internal/models"
)

// InMemoryStorage provides a thread-safe, in-memory implementation of the Storage interface.
// It stores events in a nested map structure: userId → date string → slice of Event.
// Date strings use the format YYYY-MM-DD (models.Date.String()).
type InMemoryStorage struct {
	mu      sync.RWMutex                      // protects records for concurrent access
	records map[int]map[string][]models.Event // records[userId][dateKey] = []Event
}

// NewInMemoryStorage initializes and returns a new InMemoryStorage instance.
func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		records: make(map[int]map[string][]models.Event),
	}
}

// SaveEvent adds a new event for the given userId on event.Date.
// Returns an error if an event with the same Id already exists on that date.
// Complexity: O(n) scan of events on that date.
func (c *InMemoryStorage) SaveEvent(userId int, event models.Event) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	dateKey := event.Date.String()

	if _, exists := c.records[userId]; !exists {
		c.records[userId] = make(map[string][]models.Event)
	}

	eventsOnDate := c.records[userId][dateKey]
	for _, e := range eventsOnDate {
		if e.Id == event.Id {
			return NewEventExistsError(event.Id)
		}
	}

	c.records[userId][dateKey] = append(eventsOnDate, event)
	return nil
}

// UpdateEvent modifies an existing event identified by event.Id for the given userId on event.Date.
// Returns an error if the user has no events or the event is not found on that date.
func (c *InMemoryStorage) UpdateEvent(userId int, event models.Event) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	dateKey := event.Date.String()

	userDates, exists := c.records[userId]
	if !exists {
		return NewUserHasNoEventsError(userId)
	}

	eventsOnDate, exists := userDates[dateKey]
	if !exists {
		return NewEventNotFoundError(event.Id)
	}

	for i, e := range eventsOnDate {
		if e.Id == event.Id {
			eventsOnDate[i] = event
			c.records[userId][dateKey] = eventsOnDate
			return nil
		}
	}
	return NewEventNotFoundError(event.Id)
}

// DeleteEvent removes the event with the specified eventId for userId on the given date.
// If this is the last event on that date, the date key is removed. If the user has no more dates, the user is removed.
func (c *InMemoryStorage) DeleteEvent(userId int, date models.Date, eventId int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	userDates, exists := c.records[userId]
	if !exists {
		return NewUserHasNoEventsError(userId)
	}

	dateKey := date.String()
	eventsOnDate, exists := userDates[dateKey]
	if !exists {
		return NewEventNotFoundByDateError(date, userId)
	}

	idx := -1
	for i, e := range eventsOnDate {
		if e.Id == eventId {
			idx = i
			break
		}
	}
	if idx < 0 {
		return NewEventNotFoundError(eventId)
	}

	copy(eventsOnDate[idx:], eventsOnDate[idx+1:])
	eventsOnDate = eventsOnDate[:len(eventsOnDate)-1]

	if len(eventsOnDate) == 0 {
		delete(userDates, dateKey)
		if len(userDates) == 0 {
			delete(c.records, userId)
		}
	} else {
		c.records[userId][dateKey] = eventsOnDate
	}
	return nil
}

// GetEvents returns all events for userId between from and to inclusive.
// Iterates day-by-day, concatenating events for each date key found.
func (c *InMemoryStorage) GetEvents(userId int, from, to time.Time) ([]models.Event, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	userDates, exists := c.records[userId]
	if !exists {
		return nil, NewUserHasNoEventsError(userId)
	}

	var result []models.Event
	start := time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, from.Location())
	end := time.Date(to.Year(), to.Month(), to.Day(), 0, 0, 0, 0, to.Location())

	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		key := models.Date{Time: d}.String()
		if dayEvents, ok := userDates[key]; ok {
			result = append(result, dayEvents...)
		}
	}
	return result, nil
}
