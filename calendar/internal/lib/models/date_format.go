package models

import "time"

// Date is a wrapper around time.Time that enforces a unified date format.
type Date struct {
	time.Time
}

// String returns the Date formatted as YYYY-MM-DD.
func (d Date) String() string {
	return d.Format(time.DateOnly)
}
