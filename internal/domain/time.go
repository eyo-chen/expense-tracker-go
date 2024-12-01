package domain

import "time"

// IsSameMonth returns true if the two given times are in the same month.
func IsSameMonth(t1, t2 string) bool {
	date1, err := time.Parse(time.DateOnly, t1)
	if err != nil {
		return false
	}

	date2, err := time.Parse(time.DateOnly, t2)
	if err != nil {
		return false
	}

	return date1.Year() == date2.Year() && date1.Month() == date2.Month()
}
