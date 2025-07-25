package service_test

import (
	"testing"
	"time"

	"http_calendar/internal/models"
	"http_calendar/internal/service"
	"http_calendar/internal/storage"
)

func parseDate(s string) models.Date {
	t, _ := time.Parse("2006-01-02", s)
	return models.Date{Time: t}
}

func makeEvent(id, userId int, dateStr, desc string) models.Event {
	return models.Event{
		Id:          id,
		UserId:      userId,
		Date:        parseDate(dateStr),
		Time:        models.Time{Time: time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC)},
		Description: desc,
	}
}

func TestCreateAndRetrieve(t *testing.T) {
	mem := storage.NewInMemoryStorage()
	svc := service.NewCalendarService(mem, true)
	userId := 1
	e := makeEvent(100, userId, "2025-07-30", "Test Create")

	if err := svc.CreateEvent(userId, e); err != nil {
		t.Fatalf("CreateEvent failed: %v", err)
	}

	evs, err := svc.GetEventsForDay(userId, e.Date)
	if err != nil {
		t.Fatalf("GetEventsForDay failed: %v", err)
	}
	if len(evs) != 1 || evs[0].Id != e.Id {
		t.Fatalf("Expected 1 event id=100, got %v", evs)
	}
}

func TestUpdateNonExisting(t *testing.T) {
	mem := storage.NewInMemoryStorage()
	svc := service.NewCalendarService(mem, true)
	userId := 2
	e := makeEvent(200, userId, "2025-08-01", "Original")
	_ = mem.SaveEvent(userId, e)

	e.Description = "Updated"
	if err := svc.UpdateEvent(userId, e); err != nil {
		t.Fatalf("UpdateEvent failed: %v", err)
	}

	evs, _ := svc.GetEventsForDay(userId, e.Date)
	if evs[0].Description != "Updated" {
		t.Fatalf("Update did not persist, got description=%q", evs[0].Description)
	}

	e2 := makeEvent(999, userId, "2025-08-01", "Nope")
	if err := svc.UpdateEvent(userId, e2); err == nil {
		t.Fatalf("Expected error updating non-existent event, got nil")
	}
}

func TestDeleteNonExisting(t *testing.T) {
	mem := storage.NewInMemoryStorage()
	svc := service.NewCalendarService(mem, true)
	userId := 3
	e1 := makeEvent(300, userId, "2025-08-05", "One")
	e2 := makeEvent(301, userId, "2025-08-05", "Two")
	_ = mem.SaveEvent(userId, e1)
	_ = mem.SaveEvent(userId, e2)

	if err := svc.DeleteEvent(userId, e1.Date, e1.Id); err != nil {
		t.Fatalf("DeleteEvent failed: %v", err)
	}
	evs, _ := svc.GetEventsForDay(userId, e1.Date)
	if len(evs) != 1 || evs[0].Id != e2.Id {
		t.Fatalf("DeleteEvent did not remove correct event, got %v", evs)
	}

	if err := svc.DeleteEvent(userId, e1.Date, 9999); err == nil {
		t.Fatalf("Expected error deleting non-existent event, got nil")
	}
}

func TestGetEventsForDayExplicit(t *testing.T) {
	mem := storage.NewInMemoryStorage()
	svc := service.NewCalendarService(mem, true)
	userId := 6

	e1 := makeEvent(600, userId, "2025-09-01", "DayEvent")
	_ = mem.SaveEvent(userId, e1)

	evs, err := svc.GetEventsForDay(userId, parseDate("2025-09-01"))
	if err != nil {
		t.Fatalf("GetEventsForDay failed: %v", err)
	}
	if len(evs) != 1 || evs[0].Id != e1.Id {
		t.Fatalf("Expected one event on day, got %v", evs)
	}

	evs, err = svc.GetEventsForDay(userId, parseDate("2025-09-02"))
	if err == nil {
		t.Fatalf("Expected error for empty day, got events %v", evs)
	}
}

func TestGetEventsForWeek(t *testing.T) {
	mem := storage.NewInMemoryStorage()
	svc := service.NewCalendarService(mem, true)
	userId := 4
	monday := makeEvent(400, userId, "2025-07-28", "Mon event")
	sunday := makeEvent(401, userId, "2025-08-03", "Sun event")
	_ = mem.SaveEvent(userId, monday)
	_ = mem.SaveEvent(userId, sunday)

	evs, err := svc.GetEventsForWeek(userId, parseDate("2025-07-30"))
	if err != nil {
		t.Fatalf("GetEventsForWeek failed: %v", err)
	}
	if len(evs) != 2 {
		t.Fatalf("Expected 2 events in week, got %d", len(evs))
	}
}

func TestGetEventsForMonth(t *testing.T) {
	mem := storage.NewInMemoryStorage()
	svc := service.NewCalendarService(mem, true)
	userId := 5
	eJuly := makeEvent(500, userId, "2025-07-15", "July")
	eAug := makeEvent(501, userId, "2025-08-01", "Aug")
	_ = mem.SaveEvent(userId, eJuly)
	_ = mem.SaveEvent(userId, eAug)

	evs, err := svc.GetEventsForMonth(userId, parseDate("2025-07-10"))
	if err != nil {
		t.Fatalf("GetEventsForMonth failed: %v", err)
	}
	if len(evs) != 1 || evs[0].Id != eJuly.Id {
		t.Fatalf("Expected only July event, got %v", evs)
	}
}

func TestGetEventsForWeekSundayStart(t *testing.T) {
	mem := storage.NewInMemoryStorage()
	svc := service.NewCalendarService(mem, false)
	userId := 7

	eSun := makeEvent(700, userId, "2025-07-27", "SunEvent")
	eSat := makeEvent(701, userId, "2025-08-02", "SatEvent")
	_ = mem.SaveEvent(userId, eSun)
	_ = mem.SaveEvent(userId, eSat)

	evs, err := svc.GetEventsForWeek(userId, parseDate("2025-07-29"))
	if err != nil {
		t.Fatalf("GetEventsForWeek (Sunday start) failed: %v", err)
	}
	if len(evs) != 2 {
		t.Fatalf("Expected 2 events in sunday-based week, got %d", len(evs))
	}
}
