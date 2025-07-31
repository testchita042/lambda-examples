package adapters

import (
	"time"
)

// SystemClock is an adapter that implements the ClockPort interface
type SystemClock struct{}

// NewSystemClock creates a new SystemClock adapter
func NewSystemClock() *SystemClock {
	return &SystemClock{}
}

// GetCurrentHour returns the current hour in 12-hour format, AM/PM designation, and full time string
func (c *SystemClock) GetCurrentHour() (int, string, string) {
	now := time.Now()

	// Format the full time string
	currentTime := now.Format(time.RFC3339)

	// Get hour in 24-hour format
	hour24 := now.Hour()

	// Convert to 12-hour format
	hour12 := hour24 % 12
	if hour12 == 0 {
		hour12 = 12
	}

	// Determine AM/PM
	amPm := "AM"
	if hour24 >= 12 {
		amPm = "PM"
	}

	return hour12, amPm, currentTime
}
