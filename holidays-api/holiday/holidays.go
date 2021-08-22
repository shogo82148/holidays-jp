package holiday

import (
	"fmt"
	"sort"
	"strings"
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

// findHolidaysInYear returns holidays in the specific year.
func findHolidaysInYear(year int) []Holiday {
	startDate := fmt.Sprintf("%04d-01-01", year)
	endDate := fmt.Sprintf("%04d-99-99", year)

	start := sort.Search(len(holidays), func(i int) bool {
		return holidays[i].Date >= startDate
	})
	end := sort.Search(len(holidays), func(i int) bool {
		return holidays[i].Date >= endDate
	})
	return holidays[start:end]
}

type annuallyHolidaysRule struct {
	// BeginYear is a year that the law is enforced
	BeginYear int

	// StaticHolydays are holydays that are on the same date every year
	StaticHolydays []staticHolyday

	// StaticHolydays are holydays that are on the same weekday in the month.
	WeekdayHolydays []weekdayHolyday
}

type staticHolyday struct {
	Date string // MM-DD
	Name string
}

type weekdayHolyday struct {
	Month   time.Month
	Weekday time.Weekday
	Index   int
	Name    string
}

func calcHolidaysInMonthWithoutInLieu(year int, month time.Month) []Holiday {
	// search the rule of this year
	var rule *annuallyHolidaysRule
	for i := len(annuallyHolidaysRules); i > 0; i-- {
		if annuallyHolidaysRules[i-1].BeginYear >= year {
			rule = &annuallyHolidaysRules[i-1]
			break
		}
	}
	if rule == nil {
		return nil
	}

	var holydays []Holiday
	yearPrefix := fmt.Sprintf("%04d-", year)
	monthPrefix := fmt.Sprintf("%02d-", int(month))
	for _, d := range rule.StaticHolydays {
		if strings.HasPrefix(d.Date, monthPrefix) {
			holydays = append(holydays, Holiday{
				Date: yearPrefix + d.Date,
				Name: d.Name,
			})
		}
	}

	weekdayOfFirstDay := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC).Weekday()
	_ = weekdayOfFirstDay
	for _, d := range rule.WeekdayHolydays {
		if d.Month == month {
			day := int(d.Weekday - weekdayOfFirstDay)
			if day < 0 {
				day += 7
			}
			day += d.Index*7 + 1
			holydays = append(holydays, Holiday{
				Date: fmt.Sprintf("%04d-%02d-%02d", year, int(month), day),
				Name: d.Name,
			})
		}
	}

	// TODO: 春分の日, 秋分の日

	return holydays
}
