package transaction

import (
	"time"
)

func cvtDateToTime(startDate, endDate string) (time.Time, time.Time, error) {
	if startDate == "" || endDate == "" {
		return time.Time{}, time.Time{}, nil
	}

	start, err := time.Parse(time.DateOnly, startDate)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	end, err := time.Parse(time.DateOnly, endDate)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	return start, end, nil
}
