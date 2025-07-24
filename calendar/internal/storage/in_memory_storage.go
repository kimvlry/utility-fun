package storage

import (
	"http_calendar/internal/models"
	"sync"
	"time"
)

// InMemoryStorage is an in-memory implementation of Storage interface.
// Stores events by userId and date string (YYYY-MM-DD).
type InMemoryStorage struct {
	mu              sync.RWMutex
	records         map[int]map[string][]models.Event // string represents models.Date string
	mondayBasedWeek bool                              // true if week starts on Monday
}

func NewInMemoryStorage(mondayBasedWeek bool) *InMemoryStorage {
	return &InMemoryStorage{
		records:         make(map[int]map[string][]models.Event),
		mondayBasedWeek: mondayBasedWeek,
	}
}

// CreateEvent adds a new event if no event with the same Id exists on that date.
func (c *InMemoryStorage) CreateEvent(userId int, event models.Event) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	dateKey := event.Date.String()

	if _, ok := c.records[userId]; !ok {
		c.records[userId] = make(map[string][]models.Event)
	}

	eventsOnDate := c.records[userId][dateKey]
	for _, e := range eventsOnDate {
		if e.Id == event.Id {
			return NewEventExistsError(event.Id)
		}
	}

	c.records[userId][dateKey] = append(c.records[userId][dateKey], event)
	return nil
}

// UpdateEvent updates an existing event identified by Id on the given date.
func (c *InMemoryStorage) UpdateEvent(userId int, event models.Event) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	dateKey := event.Date.String()

	if _, ok := c.records[userId]; !ok {
		return NewUserHasNoEventsError(userId)
	}

	eventsOnDate, ok := c.records[userId][dateKey]
	if !ok {
		return NewEventNotFoundError(event.Id)
	}

	updated := false
	for i, e := range eventsOnDate {
		if e.Id == event.Id {
			eventsOnDate[i] = event
			updated = true
			break
		}
	}
	if !updated {
		return NewEventNotFoundError(event.Id)
	}

	c.records[userId][dateKey] = eventsOnDate
	return nil
}

// DeleteEvent deletes an event by Id on the specified date.
func (c *InMemoryStorage) DeleteEvent(userId int, date models.Date, eventId int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.records[userId]; !ok {
		return NewUserHasNoEventsError(userId)
	}

	eventsOnDate, ok := c.records[userId][date.String()]
	if !ok {
		return NewEventNotFoundByDateError(date, userId)
	}

	index := -1
	for i, e := range eventsOnDate {
		if e.Id == eventId {
			index = i
			break
		}
	}

	if index == -1 {
		return NewEventNotFoundError(eventId)
	}

	copy(eventsOnDate[index:], eventsOnDate[index+1:])
	eventsOnDate = eventsOnDate[:len(eventsOnDate)-1]

	if len(eventsOnDate) == 0 {
		delete(c.records[userId], date.String())
		if len(c.records[userId]) == 0 {
			delete(c.records, userId)
		}
	} else {
		c.records[userId][date.String()] = eventsOnDate
	}

	return nil
}

// GetEventsForDay returns all events for a user on a specific day.
func (c *InMemoryStorage) GetEventsForDay(userId int, date models.Date) ([]models.Event, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if _, ok := c.records[userId]; !ok {
		return nil, NewUserHasNoEventsError(userId)
	}

	eventsOnDate, ok := c.records[userId][date.String()]
	if !ok || len(eventsOnDate) == 0 {
		return nil, NewEventNotFoundByDateError(date, userId)
	}

	eventsCopy := make([]models.Event, len(eventsOnDate))
	copy(eventsCopy, eventsOnDate)
	return eventsCopy, nil
}

// GetEventsForWeek returns all events for a user in the week of the given date.
func (c *InMemoryStorage) GetEventsForWeek(userId int, date models.Date) ([]models.Event, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	userEvents, ok := c.records[userId]
	if !ok {
		return nil, NewUserHasNoEventsError(userId)
	}

	weekday := int(date.Weekday())
	if c.mondayBasedWeek {
		if weekday == 0 {
			weekday = 7
		}
		weekday -= 1
	}

	start := date.AddDate(0, 0, -weekday).Truncate(24 * time.Hour)

	events := make([]models.Event, 0)
	for i := 0; i < 7; i++ {
		day := start.AddDate(0, 0, i)
		key := models.Date{Time: day}.String()

		if dayEvents, ok := userEvents[key]; ok {
			events = append(events, dayEvents...)
		}
	}

	return events, nil
}

// GetEventsForMonth returns all events for a user in the month of the given date.
func (c *InMemoryStorage) GetEventsForMonth(userId int, date models.Date) ([]models.Event, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	userEvents, ok := c.records[userId]
	if !ok {
		return nil, NewUserHasNoEventsError(userId)
	}

	year, month := date.Time.Year(), date.Time.Month()
	start := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)

	events := make([]models.Event, 0)
	for d := start; d.Before(end); d = d.AddDate(0, 0, 1) {
		key := models.Date{Time: d}.String()
		if dayEvents, ok := userEvents[key]; ok {
			events = append(events, dayEvents...)
		}
	}

	return events, nil
}
