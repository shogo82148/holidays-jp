package holiday

import (
	"cmp"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
)

// Date represents a date.
type Date struct {
	Year  int
	Month time.Month
	Day   int
}

func (a Date) cmp(b Date) int {
	if a.Year != b.Year {
		return cmp.Compare(a.Year, b.Year)
	}
	if a.Month != b.Month {
		return cmp.Compare(a.Month, b.Month)
	}
	return cmp.Compare(a.Day, b.Day)
}

func (d Date) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", d.Year, int(d.Month), d.Day)
}

// firstDay returns the first day of the month.
func (d Date) firstDay() Date {
	return Date{d.Year, d.Month, 1}
}

// nextMonth returns the first day of the next month.
func (d Date) nextMonth() Date {
	if d.Month == time.December {
		return Date{d.Year + 1, time.January, 1}
	}
	return Date{d.Year, d.Month + 1, 1}
}

// FindHoliday returns whether the specific day is a holiday.
func FindHoliday(year int, month time.Month, day int) (Holiday, bool) {
	if holidaysStartYear <= year && year <= holidaysEndYear {
		// return from pre-calculated holidays
		return findHoliday(year, month, day)
	}

	// calculate holidays based on the law
	date := fmt.Sprintf("%04d-%02d-%02d", year, int(month), day)
	holidays := calcHolidaysInMonth(year, month)
	for _, d := range holidays {
		if d.Date == date {
			return d, true
		}
	}
	return Holiday{}, false
}

// FindHolidaysInMonth returns holidays in the month.
func FindHolidaysInMonth(year int, month time.Month) []Holiday {
	if holidaysStartYear <= year && year <= holidaysEndYear {
		// return from pre-calculated holidays
		return findHolidaysInMonth(year, month)
	}

	// calculate holidays based on the law
	return calcHolidaysInMonth(year, month)
}

// FindHolidaysInYear returns holidays in the year.
func FindHolidaysInYear(year int) []Holiday {
	if holidaysStartYear <= year && year <= holidaysEndYear {
		// return from pre-calculated holidays
		return findHolidaysInYear(year)
	}

	// calculate holidays based on the law
	return calcHolidaysInYear(year)
}

func FindHolidaysInRange(from, to Date) []Holiday {
	if from.cmp(to) > 0 {
		from, to = to, from
	}
	if holidaysStartYear <= from.Year && to.Year <= holidaysEndYear {
		// return from pre-calculated holidays
		return findHolidaysInRange(from, to)
	}

	// calculate holidays based on the law
	return calcHolidaysInRange(from, to)
}

const dateLayout = "2006-01-02"

func mustParseDate(date string) time.Time {
	d, err := time.Parse(dateLayout, date)
	if err != nil {
		panic(err)
	}
	return d
}

type Holiday struct {
	Date string
	Name string
}

type withDate []Holiday

func (s withDate) Len() int           { return len(s) }
func (s withDate) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s withDate) Less(i, j int) bool { return s[i].Date < s[j].Date }

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
	startDate := Date{year, month, 1}
	endDate := Date{year, month, 31}
	return findHolidaysInRange(startDate, endDate)
}

// findHolidaysInYear returns holidays in the specific year.
func findHolidaysInYear(year int) []Holiday {
	startDate := Date{year, time.January, 1}
	endDate := Date{year, time.December, 31}
	return findHolidaysInRange(startDate, endDate)
}

// findHolidaysInRange returns holidays in the specific range.
func findHolidaysInRange(from, to Date) []Holiday {
	startDate := from.String()
	endDate := to.String()

	start := sort.Search(len(holidays), func(i int) bool {
		return holidays[i].Date >= startDate
	})
	end := sort.Search(len(holidays), func(i int) bool {
		return holidays[i].Date >= endDate
	})
	if end < len(holidays) && holidays[end].Date == endDate {
		end++
	}
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
	for i := 0; i < len(annuallyHolidaysRules); i++ {
		if year >= annuallyHolidaysRules[i].BeginYear {
			rule = &annuallyHolidaysRules[i]
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

	// Vernal Equinox Day
	if month == time.March {
		holydays = append(holydays, Holiday{
			Date: fmt.Sprintf("%04d-%02d-%02d", year, int(month), vernalEquinoxDay(year)),
			Name: "春分の日",
		})
	}

	// Autumnal Equinox Day
	if month == time.September {
		holydays = append(holydays, Holiday{
			Date: fmt.Sprintf("%04d-%02d-%02d", year, int(month), autumnalEquinoxDay(year)),
			Name: "秋分の日",
		})
	}

	yearMonthPrefix := yearPrefix + monthPrefix
	for _, d := range specialHolidays {
		if strings.HasPrefix(d.Date, yearMonthPrefix) {
			holydays = append(holydays, d)
		}
	}

	sort.Sort(withDate(holydays))
	return holydays
}

func calcHolidaysInMonth(year int, month time.Month) []Holiday {
	holidays := calcHolidaysInMonthWithoutInLieu(year, month)

	// 昭和六十年法律第百三号
	// 国民の祝日に関する法律の一部を改正する法律
	// 衆議院制定法律: https://www.shugiin.go.jp/internet/itdb_housei.nsf/html/houritsu/10319851227103.htm
	if year >= 1986 {
		var extraHolidays []Holiday
		for i := 0; i < len(holidays)-1; i++ {
			holidayA := mustParseDate(holidays[i].Date)
			holidayB := mustParseDate(holidays[i+1].Date)

			// > 第三条に次の一項を加える。
			// > ３　その前日及び翌日が「国民の祝日」である日（日曜日にあたる日及び前項に規定する休日にあたる日を除く。）は、休日とする。
			if holidayB.Sub(holidayA) == 2*24*time.Hour {
				d := holidayA.Add(24 * time.Hour)
				if d.Weekday() != time.Sunday {
					extraHolidays = append(extraHolidays, Holiday{
						Date: d.Format(dateLayout),
						Name: "休日",
					})
				}
			}
		}

		// Handle edge cases that span months
		if len(holidays) > 0 {
			firstHolidayInMonth := mustParseDate(holidays[0].Date)
			beforeTwoDays := firstHolidayInMonth.Add(-2 * 24 * time.Hour)
			if firstHolidayInMonth.Month() != beforeTwoDays.Month() && firstHolidayInMonth.Weekday() != time.Monday {
				// the first day in the month might be a holiday
				previousHolidays := calcHolidaysInMonthWithoutInLieu(
					beforeTwoDays.Year(), beforeTwoDays.Month(),
				)
				if len(previousHolidays) > 0 && previousHolidays[len(previousHolidays)-1].Date == beforeTwoDays.Format(dateLayout) {
					extraHolidays = append(extraHolidays, Holiday{
						Date: firstHolidayInMonth.Add(-24 * time.Hour).Format(dateLayout),
						Name: "休日",
					})
				}
			}

			lastHolidayInMonth := mustParseDate(holidays[len(holidays)-1].Date)
			afterTwoDays := lastHolidayInMonth.Add(2 * 24 * time.Hour)
			if lastHolidayInMonth.Month() != afterTwoDays.Month() && lastHolidayInMonth.Weekday() != time.Monday {
				// the last day in the month might be a holiday
				nextHolidays := calcHolidaysInMonthWithoutInLieu(
					afterTwoDays.Year(), afterTwoDays.Month(),
				)
				if len(nextHolidays) > 0 && nextHolidays[0].Date == afterTwoDays.Format(dateLayout) {
					extraHolidays = append(extraHolidays, Holiday{
						Date: lastHolidayInMonth.Add(24 * time.Hour).Format(dateLayout),
						Name: "休日",
					})
				}
			}
		}

		holidays = append(holidays, extraHolidays...)
		sort.Sort(withDate(holidays))
	}

	// 昭和四十八年法律第十号
	// 国民の祝日に関する法律の一部を改正する法律
	// 衆議院制定法律: https://www.shugiin.go.jp/internet/itdb_housei.nsf/html/houritsu/07119730412010.htm
	//
	// > 第三条に次の一項を加える。
	// > ２　「国民の祝日」が日曜日にあたるときは、その翌日を休日とする。
	if 1973 <= year && year < 2007 {
		var holidaysInLieu []Holiday
		for _, holiday := range holidays {

			// This law was enacted on April 12, 1973,
			// so it did not apply to holidays before that date.
			if holiday.Date <= "1973-04-12" {
				continue
			}

			d, err := time.Parse(dateLayout, holiday.Date)
			if err != nil {
				panic(err)
			}
			if d.Weekday() != time.Sunday {
				continue
			}
			d = d.Add(24 * time.Hour)
			if !contains(holidays, d.Format(dateLayout)) {
				holidaysInLieu = append(holidaysInLieu, Holiday{
					Date: d.Format(dateLayout),
					Name: "休日",
				})
			}
		}
		holidays = append(holidays, holidaysInLieu...)
		sort.Sort(withDate(holidays))
	}

	// 平成十七年法律第四十三号
	// 国民の祝日に関する法律の一部を改正する法律
	// 衆議院制定法律: https://www.shugiin.go.jp/internet/itdb_housei.nsf/html/housei/16220050520043.htm
	// 官報: https://kanpou.npb.go.jp/old/20050520/20050520g00109/20050520g001090005f.html
	//
	// > 第三条第二項中「あたるときは、その翌日」を「当たるときは、その日後においてその日に最も近い「国民の祝日」でない日」に改め、
	// > 同条第三項中「日曜日にあたる日及び前項に規定する休日にあたる日を除く。」を「「国民の祝日」でない日に限る。」に改める。
	if year >= 2007 {
		var holidaysInLieu []Holiday
		for _, holiday := range holidays {
			d, err := time.Parse(dateLayout, holiday.Date)
			if err != nil {
				panic(err)
			}
			if d.Weekday() != time.Sunday {
				continue
			}
			d = d.Add(24 * time.Hour)
			for contains(holidays, d.Format(dateLayout)) {
				d = d.Add(24 * time.Hour)
			}
			holidaysInLieu = append(holidaysInLieu, Holiday{
				Date: d.Format(dateLayout),
				Name: "休日",
			})
		}
		holidays = append(holidays, holidaysInLieu...)
		sort.Sort(withDate(holidays))
	}

	return holidays
}

func contains(holidays []Holiday, date string) bool {
	for _, d := range holidays {
		if d.Date == date {
			return true
		}
	}
	return false
}

func calcHolidaysInYear(year int) []Holiday {
	var result []Holiday
	for month := time.January; month <= time.December; month++ {
		holidays := calcHolidaysInMonth(year, month)
		result = append(result, holidays...)
	}
	return result
}

func calcHolidaysInRange(from, to Date) []Holiday {
	if from.cmp(to) > 0 {
		from, to = to, from
	}

	firstDay := to.firstDay()

	startDate := from.String()
	endDate := to.String()
	var result []Holiday
	for d := from.firstDay(); d.cmp(firstDay) <= 0; d = d.nextMonth() {
		holidays := calcHolidaysInMonth(d.Year, d.Month)
		for _, h := range holidays {
			if startDate <= h.Date && h.Date <= endDate {
				result = append(result, h)
			}
		}
	}
	return result
}

// from 長沢 工(1999) "日の出・日の入りの計算 天体の出没時刻の求め方" 株式会社地人書館
var sunLongitudeTable = [...][3]float64{
	{0.0200, 355.05, 719.981},
	{0.0048, 234.95, 19.341},
	{0.0020, 247.1, 329.64},
	{0.0018, 297.8, 4452.67},
	{0.0018, 251.3, 0.20},
	{0.0015, 343.2, 450.37},
	{0.0013, 81.4, 225.18},
	{0.0008, 132.5, 659.29},
	{0.0007, 153.3, 90.38},
	{0.0007, 206.8, 30.35},
	{0.0006, 29.8, 337.18},
	{0.0005, 207.4, 1.50},
	{0.0005, 291.2, 22.81},
	{0.0004, 234.9, 315.56},
	{0.0004, 157.3, 299.30},
	{0.0004, 21.1, 720.02},
	{0.0003, 352.5, 1079.97},
	{0.0003, 329.7, 44.43},
}

// julianYear is a number of julian years from J2000.0(2000/01/01 12:00 Terrestrial Time)
type julianYear float64

var j2000 = time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC).Unix()

func time2JulianYear(t time.Time) julianYear {
	d := t.Unix() - j2000

	// convert UTC(Coordinated Universal Time) into TAI(International Atomic Time)
	d += 36 // TAI - UTC = 36seconds (at 2015/08)

	// convert TAI into TT(Terrestrial Time)
	d += 32
	return julianYear(float64(d) / ((365*24 + 6) * 60 * 60))
}

func sunLongitude(jy julianYear) float64 {
	t := float64(jy)
	l := normalizeDegree(360.00769 * t)
	l = normalizeDegree(l + 280.4603)
	l = normalizeDegree(l + (1.9146-0.00005*t)*sin(357.538+359.991*t))
	for _, b := range sunLongitudeTable {
		l = normalizeDegree(l + b[0]*sin(b[1]+b[2]*t))
	}
	return l
}

func sin(x float64) float64 {
	return math.Sin(x / 180 * math.Pi)
}

func normalizeDegree(x float64) float64 {
	x = math.Mod(x, 360)
	if x < 0 {
		x += 360
	}
	return x
}

var jst *time.Location

func init() {
	var err error
	jst, err = time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}
}

func vernalEquinoxDay(year int) int {
	for i := 10; i <= 31; i++ {
		t := time.Date(year, time.March, i, 0, 0, 0, 0, jst)
		l := sunLongitude(time2JulianYear(t))
		if l < 180 {
			return i - 1
		}
	}
	return 0
}

func autumnalEquinoxDay(year int) int {
	for i := 10; i <= 30; i++ {
		t := time.Date(year, time.September, i, 0, 0, 0, 0, jst)
		l := sunLongitude(time2JulianYear(t))
		if l >= 180 {
			return i - 1
		}
	}
	return 0
}
