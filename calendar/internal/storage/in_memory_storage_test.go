package storage_test

import (
	"http_calendar/internal/models"
	"http_calendar/internal/storage"
	"testing"
	"time"
)

func parseDate(dateStr string) models.Date {
	t, _ := time.Parse("2006-01-02", dateStr)
	return models.Date{Time: t}
}

func parseTime(timeStr string) models.Time {
	t, _ := time.Parse("15:04", timeStr)
	return models.Time{Time: t}
}

func makeEvent(id, userId int, dateStr, timeStr, desc string) models.Event {
	return models.Event{
		Id:          id,
		UserId:      userId,
		Date:        parseDate(dateStr),
		Time:        parseTime(timeStr),
		Description: desc,
	}
}

func TestCreateDuplicateEvent(t *testing.T) {
	store := storage.NewInMemoryStorage(true)
	userId := 1
	event := makeEvent(10, userId, "2025-07-25", "14:30", "Test event")

	err := store.CreateEvent(userId, event)
	if err != nil {
		t.Fatalf("CreateEvent failed: %v", err)
	}

	err = store.CreateEvent(userId, event)
	if err == nil {
		t.Fatalf("Expected error when creating duplicate event, got nil")
	}
}

func TestUpdateEvent(t *testing.T) {
	store := storage.NewInMemoryStorage(true)
	userId := 2
	event := makeEvent(20, userId, "2025-07-26", "09:00", "Original event")

	err := store.CreateEvent(userId, event)
	if err != nil {
		t.Fatalf("CreateEvent failed: %v", err)
	}

	event.Description = "Updated event"
	err = store.UpdateEvent(userId, event)
	if err != nil {
		t.Fatalf("UpdateEvent failed: %v", err)
	}

	events, err := store.GetEventsForDay(userId, parseDate("2025-07-26"))
	if err != nil {
		t.Fatalf("GetEventsForDay failed: %v", err)
	}
	if len(events) != 1 || events[0].Description != "Updated event" {
		t.Fatalf("Event not updated correctly, got: %+v", events)
	}

	fakeEvent := makeEvent(999, userId, "2025-07-26", "10:00", "Fake event")
	err = store.UpdateEvent(userId, fakeEvent)
	if err == nil {
		t.Fatalf("Expected error when updating non-existent event, got nil")
	}
}

func TestDeleteEvent(t *testing.T) {
	store := storage.NewInMemoryStorage(true)
	userId := 3
	date := parseDate("2025-07-27")

	event1 := makeEvent(30, userId, "2025-07-27", "08:00", "Event 1")
	event2 := makeEvent(31, userId, "2025-07-27", "10:00", "Event 2")

	err := store.CreateEvent(userId, event1)
	if err != nil {
		t.Fatalf("CreateEvent failed: %v", err)
	}
	err = store.CreateEvent(userId, event2)
	if err != nil {
		t.Fatalf("CreateEvent failed: %v", err)
	}

	err = store.DeleteEvent(userId, date, 30)
	if err != nil {
		t.Fatalf("DeleteEvent failed: %v", err)
	}

	events, err := store.GetEventsForDay(userId, date)
	if err != nil {
		t.Fatalf("GetEventsForDay failed: %v", err)
	}
	if len(events) != 1 || events[0].Id != 31 {
		t.Fatalf("DeleteEvent did not remove correct event, got: %+v", events)
	}

	err = store.DeleteEvent(userId, date, 999)
	if err == nil {
		t.Fatalf("Expected error when deleting non-existent event, got nil")
	}
}

func TestGetEventsForDay(t *testing.T) {
	store := storage.NewInMemoryStorage(true)
	userId := 4
	date := parseDate("2025-07-28")

	event := makeEvent(40, userId, "2025-07-28", "12:00", "Day event")
	err := store.CreateEvent(userId, event)
	if err != nil {
		t.Fatalf("CreateEvent failed: %v", err)
	}

	events, err := store.GetEventsForDay(userId, date)
	if err != nil {
		t.Fatalf("GetEventsForDay failed: %v", err)
	}
	if len(events) != 1 || events[0].Id != 40 {
		t.Fatalf("GetEventsForDay returned wrong events, got: %+v", events)
	}

	_, err = store.GetEventsForDay(userId, parseDate("2025-07-29"))
	if err == nil {
		t.Fatalf("Expected error for no events on date, got nil")
	}
}

func TestGetEventsForWeek(t *testing.T) {
	store := storage.NewInMemoryStorage(true)
	userId := 5

	eventMon := makeEvent(50, userId, "2025-07-21", "09:00", "Monday event")
	err := store.CreateEvent(userId, eventMon)
	if err != nil {
		t.Fatalf("CreateEvent failed: %v", err)
	}

	eventSun := makeEvent(51, userId, "2025-07-27", "18:00", "Sunday event")
	err = store.CreateEvent(userId, eventSun)
	if err != nil {
		t.Fatalf("CreateEvent failed: %v", err)
	}

	events, err := store.GetEventsForWeek(userId, parseDate("2025-07-23"))
	if err != nil {
		t.Fatalf("GetEventsForWeek failed: %v", err)
	}

	if len(events) != 2 {
		t.Fatalf("Expected 2 events in week, got %d", len(events))
	}
}

func TestGetEventsForMonth(t *testing.T) {
	store := storage.NewInMemoryStorage(true)
	userId := 6

	event1 := makeEvent(60, userId, "2025-07-10", "10:00", "Event July 10")
	event2 := makeEvent(61, userId, "2025-07-15", "15:00", "Event July 15")
	event3 := makeEvent(62, userId, "2025-08-01", "09:00", "Event August 1")

	err := store.CreateEvent(userId, event1)
	if err != nil {
		t.Fatalf("CreateEvent failed: %v", err)
	}
	err = store.CreateEvent(userId, event2)
	if err != nil {
		t.Fatalf("CreateEvent failed: %v", err)
	}
	err = store.CreateEvent(userId, event3)
	if err != nil {
		t.Fatalf("CreateEvent failed: %v", err)
	}

	events, err := store.GetEventsForMonth(userId, parseDate("2025-07-20"))
	if err != nil {
		t.Fatalf("GetEventsForMonth failed: %v", err)
	}

	if len(events) != 2 {
		t.Fatalf("Expected 2 events in July, got %d", len(events))
	}
}
