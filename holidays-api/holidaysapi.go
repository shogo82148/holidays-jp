package holidaysapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/shogo82148/holidays-jp/holidays-api/holiday"
)

// Response is the response of Handler.
type Response struct {
	Holidays []Holiday `json:"holidays"`
}

// Holiday is a holiday.
type Holiday struct {
	Date string `json:"date"`
	Name string `json:"name"`
}

// Handler provides a holiday api.
type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	year, month, day, err := parsePath(r.URL.Path)
	if err != nil {
		h.responseNotFound(w)
		return
	}
	switch {
	case year == 0:
		h.responseNotFound(w)
	case month == 0:
		// 2006
		h.holidaysInYear(w, year)
	case day == 0:
		// 2006/01
		if month < 1 || month > 12 {
			h.responseNotFound(w)
			return
		}
		h.holidaysInMonth(w, year, time.Month(month))
	default:
		// 2006/01/02
		_, err := time.Parse("2006/01/02", fmt.Sprintf("%04d/%02d/%02d", year, month, day))
		if err != nil {
			h.responseNotFound(w)
			return
		}
		h.holiday(w, year, time.Month(month), day)
	}
}

func parsePath(path string) (year, month, day int, err error) {
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")
	if path == "" {
		return
	}

	seg := strings.SplitN(path, "/", 3)
	if len(seg) >= 1 {
		year, err = parseInt(seg[0], 4)
		if err != nil {
			return 0, 0, 0, err
		}
	}
	if len(seg) >= 2 {
		month, err = parseInt(seg[1], 2)
		if err != nil {
			return 0, 0, 0, err
		}
	}
	if len(seg) >= 3 {
		day, err = parseInt(seg[2], 2)
		if err != nil {
			return 0, 0, 0, err
		}
	}
	return
}

func parseInt(s string, digits int) (int, error) {
	if len(s) != digits {
		return 0, errors.New("invalid format")
	}

	var ret int
	for _, ch := range s {
		if '0' <= ch && ch <= '9' {
			ret = ret*10 + int(ch-'0')
		} else {
			return 0, fmt.Errorf("unexpected character: %c", ch)
		}
	}
	return ret, nil
}

func (h *Handler) holiday(w http.ResponseWriter, year int, month time.Month, day int) {
	d, ok := holiday.FindHoliday(year, month, day)
	if ok {
		h.responseHolidays(w, []holiday.Holiday{d})
	} else {
		h.responseHolidays(w, []holiday.Holiday{})
	}
}

func (h *Handler) holidaysInMonth(w http.ResponseWriter, year int, month time.Month) {
	holidays := holiday.FindHolidaysInMonth(year, month)
	h.responseHolidays(w, holidays)
}

func (h *Handler) holidaysInYear(w http.ResponseWriter, year int) {
	holidays := holiday.FindHolidaysInYear(year)
	h.responseHolidays(w, holidays)
}

func (h *Handler) responseHolidays(w http.ResponseWriter, holidays []holiday.Holiday) {
	w.Header().Set("Content-Type", "application/json")

	res := make([]Holiday, 0, len(holidays))
	for _, d := range holidays {
		res = append(res, Holiday{
			Date: d.Date,
			Name: d.Name,
		})
	}
	data, err := json.Marshal(Response{
		Holidays: res,
	})
	if err != nil {
		log.Printf("failed to marshal response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"error":"internal server error"}`)
	}

	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (h *Handler) responseNotFound(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d", 24*60*60))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	io.WriteString(w, `{"error":"not found"}`)
}
