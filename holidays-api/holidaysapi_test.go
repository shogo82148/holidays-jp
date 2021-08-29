package holidaysapi

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestServeHTTP(t *testing.T) {
	h := NewHandler()
	t.Run("not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "http://example.com/", nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		resp := w.Result()
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("unexpected status code: want %d, got %d", http.StatusNotFound, resp.StatusCode)
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		var v interface{}
		if err := json.Unmarshal(body, &v); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("year", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "http://example.com/2000", nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		resp := w.Result()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("unexpected status code: want %d, got %d", http.StatusOK, resp.StatusCode)
		}
		if resp.Header.Get("Cache-Control") == "" {
			t.Error("Cache-Control is not set")
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		var got Response
		if err := json.Unmarshal(body, &got); err != nil {
			t.Fatal(err)
		}
		want := Response{
			Holidays: []Holiday{
				{
					Date: "2000-01-01",
					Name: "元日",
				},
				{
					Date: "2000-01-10",
					Name: "成人の日",
				},
				{
					Date: "2000-02-11",
					Name: "建国記念の日",
				},
				{
					Date: "2000-03-20",
					Name: "春分の日",
				},
				{
					Date: "2000-04-29",
					Name: "みどりの日",
				},
				{
					Date: "2000-05-03",
					Name: "憲法記念日",
				},
				{
					Date: "2000-05-04",
					Name: "休日",
				},
				{
					Date: "2000-05-05",
					Name: "こどもの日",
				},
				{
					Date: "2000-07-20",
					Name: "海の日",
				},
				{
					Date: "2000-09-15",
					Name: "敬老の日",
				},
				{
					Date: "2000-09-23",
					Name: "秋分の日",
				},
				{
					Date: "2000-10-09",
					Name: "体育の日",
				},
				{
					Date: "2000-11-03",
					Name: "文化の日",
				},
				{
					Date: "2000-11-23",
					Name: "勤労感謝の日",
				},
				{
					Date: "2000-12-23",
					Name: "天皇誕生日",
				},
			},
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("unexpected response: (-want/+got)\n%s", diff)
		}
	})

	t.Run("month", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "http://example.com/2000/01", nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		resp := w.Result()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("unexpected status code: want %d, got %d", http.StatusOK, resp.StatusCode)
		}
		if resp.Header.Get("Cache-Control") == "" {
			t.Error("Cache-Control is not set")
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		var got Response
		if err := json.Unmarshal(body, &got); err != nil {
			t.Fatal(err)
		}
		want := Response{
			Holidays: []Holiday{
				{
					Date: "2000-01-01",
					Name: "元日",
				},
				{
					Date: "2000-01-10",
					Name: "成人の日",
				},
			},
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("unexpected response: (-want/+got)\n%s", diff)
		}
	})

	t.Run("holiday", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "http://example.com/2000/01/01", nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		resp := w.Result()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("unexpected status code: want %d, got %d", http.StatusOK, resp.StatusCode)
		}
		if resp.Header.Get("Cache-Control") == "" {
			t.Error("Cache-Control is not set")
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		var got Response
		if err := json.Unmarshal(body, &got); err != nil {
			t.Fatal(err)
		}
		want := Response{
			Holidays: []Holiday{
				{
					Date: "2000-01-01",
					Name: "元日",
				},
			},
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("unexpected response: (-want/+got)\n%s", diff)
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
