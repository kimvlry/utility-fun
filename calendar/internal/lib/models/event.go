package models

// Event represents a calendar entry.
// Each event is associated with a user, a date, and a text description.
type Event struct {
	Id     string `json:"id"`      // Id is the unique identifier of the event
	UserId string `json:"user_id"` // UserId is the unique identifier of the user who owns the event.
	Date   Date   `json:"date"`    // Date is the date (YYYY-MM-DD) of the event.
	// TODO: time of the day
	Name string `json:"event"` // Name is a name or brief description of the event.
}
