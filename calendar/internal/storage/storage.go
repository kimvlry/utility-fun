package storage

import (
	"http_calendar/internal/models"
)

// Storage is an interface defining operations on events storage
type Storage interface {
	CreateEvent(userId int, e models.Event) error
	UpdateEvent(userId int, e models.Event) error
	DeleteEvent(userId int, date models.Date) error
	GetEventsForDay(userId int, date models.Date) ([]models.Event, error)
	GetEventsForWeek(userId int, date models.Date) ([]models.Event, error)
	GetEventsForMonth(userId int, date models.Date) ([]models.Event, error)
}
