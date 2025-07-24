package models

import "time"

// Time is a wrapper around time.Time that enforces a unified time format.
type Time struct {
	time.Time
}

// String returns the Time formatted as HH:MM.
func (d Time) String() string {
	return d.Time.Format("15:04")
}
