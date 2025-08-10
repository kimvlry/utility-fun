package storage_test

import (
    models2 "http_calendar/internal/lib/models"
    "http_calendar/internal/storage"
    "testing"
    "time"
)

func parseDate(dateStr string) models2.Date {
    t, _ := time.Parse("2006-01-02", dateStr)
    return models2.Date{Time: t}
}

func makeEvent(id, userId string, dateStr, desc string) models2.Event {
    return models2.Event{
        Id:     id,
        UserId: userId,
        Date:   parseDate(dateStr),
        Name:   desc,
    }
}

func TestSaveDuplicateEvent(t *testing.T) {
    store := storage.NewInMemoryStorage()
    userId := "1"
    e := makeEvent("10", userId, "2025-07-25", "Test event")

    if err := store.SaveEvent(userId, e); err != nil {
        t.Fatalf("SaveEvent failed: %v", err)
    }
    if err := store.SaveEvent(userId, e); err == nil {
        t.Fatalf("Expected error when saving duplicate, got nil")
    }
}

func TestUpdateNonExistingEvent(t *testing.T) {
    store := storage.NewInMemoryStorage()
    userId := "2"
    e := makeEvent("20", userId, "2025-07-26", "Original")
    if err := store.SaveEvent(userId, e); err != nil {
        t.Fatalf("SaveEvent failed: %v", err)
    }

    e.Name = "Updated"
    if err := store.UpdateEvent(userId, e); err != nil {
        t.Fatalf("UpdateEvent failed: %v", err)
    }

    from := parseDate("2025-07-26").Time
    evs, err := store.GetEvents(userId, from, from)
    if err != nil {
        t.Fatalf("GetEvents failed: %v", err)
    }
    if len(evs) != 1 || evs[0].Name != "Updated" {
        t.Fatalf("Update did not apply, got: %+v", evs)
    }

    eInvalid := makeEvent("999", userId, "2025-07-26", "Nope")
    if err := store.UpdateEvent(userId, eInvalid); err == nil {
        t.Fatalf("Expected error updating nonexistent event, got nil")
    }
}

func TestDeleteNonExistingEvent(t *testing.T) {
    store := storage.NewInMemoryStorage()
    userId := "3"
    date := parseDate("2025-07-27")
    e1 := makeEvent("30", userId, "2025-07-27", "One")
    e2 := makeEvent("31", userId, "2025-07-27", "Two")
    _ = store.SaveEvent(userId, e1)
    _ = store.SaveEvent(userId, e2)

    if err := store.DeleteEvent(userId, date, "30"); err != nil {
        t.Fatalf("DeleteEvent failed: %v", err)
    }
    from := date.Time
    evs, err := store.GetEvents(userId, from, from)
    if err != nil {
        t.Fatalf("GetEvents failed: %v", err)
    }
    if len(evs) != 1 || evs[0].Id != "31" {
        t.Fatalf("Expected only event 31 after deleter, got %v", evs)
    }

    if err := store.DeleteEvent(userId, date, "999"); err == nil {
        t.Fatalf("Expected error deleting nonexistent, got nil")
    }
}

func TestGetEventsRange(t *testing.T) {
    store := storage.NewInMemoryStorage()
    userId := "4"

    err := store.SaveEvent(userId, makeEvent("1", userId, "2025-07-10", "A"))
    if err != nil {
        return
    }
    err = store.SaveEvent(userId, makeEvent("2", userId, "2025-07-15", "B"))
    if err != nil {
        return
    }
    err = store.SaveEvent(userId, makeEvent("3", userId, "2025-07-20", "C"))
    if err != nil {
        return
    }

    from := parseDate("2025-07-11").Time
    to := parseDate("2025-07-18").Time
    evs, err := store.GetEvents(userId, from, to)
    if err != nil {
        t.Fatalf("GetEvents failed: %v", err)
    }
    if len(evs) != 1 || evs[0].Id != "2" {
        t.Fatalf("Expected [2], got ids=%v", extractIds(evs))
    }

    evs, err = store.GetEvents(userId, parseDate("2025-07-10").Time, parseDate("2025-07-20").Time)
    if err != nil {
        t.Fatalf("GetEvents failed: %v", err)
    }
    if len(evs) != 3 {
        t.Fatalf("Expected 3 events, got %d", len(evs))
    }
}

func extractIds(evs []models2.Event) []string {
    ids := make([]string, len(evs))
    for i, e := range evs {
        ids[i] = e.Id
    }
    return ids
}
