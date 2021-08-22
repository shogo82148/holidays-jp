package holiday

import (
	"fmt"
	"sort"
	"time"
)

type Holiday struct {
	Date string
	Name string
}

// findHoliday returns whether the specific day is a holiday.
func findHoliday(year int, month time.Month, day int) (Holiday, bool) {
	date := fmt.Sprintf("%04d-%02d-%02d", year, int(month), day)
	idx := sort.Search(len(holidays), func(i int) bool {
		return holidays[i].Date >= date
	})

	if idx < len(holidays) && holidays[idx].Date == date {
		return holidays[idx], true
	}
	return Holiday{}, false
}

// findHolidaysInMonth returns holidays in the specific month.
func findHolidaysInMonth(year int, month time.Month) []Holiday {
	startDate := fmt.Sprintf("%04d-%02d-01", year, int(month))
	endDate := fmt.Sprintf("%04d-%02d-99", year, int(month))

	start := sort.Search(len(holidays), func(i int) bool {
		return holidays[i].Date >= startDate
	})
	end := sort.Search(len(holidays), func(i int) bool {
		return holidays[i].Date >= endDate
	})
	return holidays[start:end]
}
