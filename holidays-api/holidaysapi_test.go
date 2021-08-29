package holidaysapi

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServeHTTP(t *testing.T) {
	h := NewHandler()
	t.Run("not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "http://example.com/", nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		resp := w.Result()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		var v interface{}
		if err := json.Unmarshal(body, &v); err != nil {
			t.Fatal(err)
		}
	})
}

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
