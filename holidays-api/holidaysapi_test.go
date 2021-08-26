package holidaysapi

import "testing"

func TestParsePath(t *testing.T) {
	tests := []struct {
		path  string
		year  int
		month int
		day   int
		err   bool
	}{
		{
			path: "",
		},
		{
			path: "/",
		},
		{
			path: "2006",
			year: 2006,
		},
		{
			path: "2006/",
			year: 2006,
		},
		{
			path: "2006/",
			year: 2006,
		},
		{
			path: "200/",
			err:  true,
		},
		{
			path: "0200/",
			year: 200,
		},
		{
			path:  "2006/01",
			year:  2006,
			month: 1,
		},
		{
			path:  "2006/01/02",
			year:  2006,
			month: 1,
			day:   2,
		},
		{
			path:  "2006/01/02/",
			year:  2006,
			month: 1,
			day:   2,
		},
		{
			path: "2006/01/02/03",
			err:  true,
		},
	}

	for _, tt := range tests {
		year, month, day, err := parsePath(tt.path)
		if tt.err != (err != nil) {
			t.Errorf("%q: unexpected error: %v", tt.path, err)
		}
		if year != tt.year {
			t.Errorf("%q: unexpected year: want %d, got %d", tt.path, tt.year, year)
		}
		if month != tt.month {
			t.Errorf("%q: unexpected month: want %d, got %d", tt.path, tt.month, month)
		}
		if day != tt.day {
			t.Errorf("%q: unexpected day: want %d, got %d", tt.path, tt.day, day)
		}
	}
}
