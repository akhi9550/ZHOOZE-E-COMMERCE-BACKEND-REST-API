package helper

import (
	"time"
)

func GetTimeFromPeriod(timePeriod string) (time.Time, time.Time) {
	endDate := time.Now()

	if timePeriod == "day" {
		startDate := endDate.AddDate(0, 0, -1)
		return startDate, endDate
	}

	if timePeriod == "week" {
		startDate := endDate.AddDate(0, 0, -6)
		return startDate, endDate
	}

	if timePeriod == "year" {
		startDate := endDate.AddDate(-1, 0, 0)
		return startDate, endDate
	}

	return endDate.AddDate(0, 0, -6), endDate
}
