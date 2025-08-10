package models

// Event represents a calendar entry.
// Each event is associated with a user, a date, and a text description.
type Event struct {
	Id          int    `json:"id"`      // Id is the unique identifier of the event
	UserId      int    `json:"user_id"` // UserId is the unique identifier of the user who owns the event.
	Date        Date   `json:"date"`    // Date is the date (YYYY-MM-DD) of the event.
	Time        Time   `json:"time"`    // Time is the time (HH:MM) of the event
	Description string `json:"event"`   // Description is a brief description of the event.
}
