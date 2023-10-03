package holidaysapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/shogo82148/holidays-jp/holidays-api/holiday"
)

var jst *time.Location

func init() {
	var err error
	jst, err = time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}
}

var errInvalidDateFormat = errors.New("holidaysapi: invalid date format")

func parseDate(s string) (holiday.Date, error) {
	y, s, ok := strings.Cut(s, "-")
	if !ok {
		return holiday.Date{}, errInvalidDateFormat
	}
	m, s, ok := strings.Cut(s, "-")
	if !ok {
		return holiday.Date{}, errInvalidDateFormat
	}
	d := s

	year, err := parseInt(y, 4)
	if err != nil || year < 1 || year > 9999 {
		return holiday.Date{}, errInvalidDateFormat
	}
	month, err := parseInt(m, 2)
	if err != nil || month < 1 || month > 12 {
		return holiday.Date{}, errInvalidDateFormat
	}
	day, err := parseInt(d, 2)
	if err != nil || day < 1 || day > 31 {
		return holiday.Date{}, errInvalidDateFormat
	}
	return holiday.Date{
		Year:  year,
		Month: time.Month(month),
		Day:   day,
	}, nil
}

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
	if r.Method != http.MethodGet {
		h.responseNotFound(w)
		return
	}

	path := r.URL.Path
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")
	if path == "holidays" {
		if err := h.holidaysInRange(w, r.URL); err != nil {
			h.responseNotFound(w)
		}
		return
	}

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
	now := time.Now().In(jst)
	if year < now.Year() || (year == now.Year() && month < now.Month()) || (year == now.Year() && month == now.Month() && day < now.Day()) {
		w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d", 365*24*60*60))
	} else {
		w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d", 24*60*60))
	}

	d, ok := holiday.FindHoliday(year, month, day)
	if ok {
		h.responseHolidays(w, []holiday.Holiday{d})
	} else {
		h.responseHolidays(w, []holiday.Holiday{})
	}
}

func (h *Handler) holidaysInMonth(w http.ResponseWriter, year int, month time.Month) {
	now := time.Now().In(jst)
	if year < now.Year() || (year == now.Year() && month < now.Month()) {
		w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d", 365*24*60*60))
	} else {
		w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d", 24*60*60))
	}

	holidays := holiday.FindHolidaysInMonth(year, month)
	h.responseHolidays(w, holidays)
}

func (h *Handler) holidaysInYear(w http.ResponseWriter, year int) {
	now := time.Now().In(jst)
	if year < now.Year() {
		w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d", 365*24*60*60))
	} else {
		w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d", 24*60*60))
	}

	holidays := holiday.FindHolidaysInYear(year)
	h.responseHolidays(w, holidays)
}

func (h *Handler) holidaysInRange(w http.ResponseWriter, u *url.URL) error {
	w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d", 24*60*60))

	q := u.Query()
	if !q.Has("from") || !q.Has("to") {
		h.holidaysInYear(w, time.Now().In(jst).Year())
		return nil
	}
	from, err := parseDate(q.Get("from"))
	if err != nil {
		return err
	}
	to, err := parseDate(q.Get("to"))
	if err != nil {
		return err
	}

	holidays := holiday.FindHolidaysInRange(from, to)
	h.responseHolidays(w, holidays)
	return nil
}

func (h *Handler) responseHolidays(w http.ResponseWriter, holidays []holiday.Holiday) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Link", "<https://github.com/sponsors/shogo82148>; rel=\"author\"")

	// ref. https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Strict-Transport-Security#examples
	w.Header().Set("Strict-Transport-Security", "max-age=63072000")

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
	w.Header().Set("Link", "<https://github.com/sponsors/shogo82148>; rel=\"author\"")

	// ref. https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Strict-Transport-Security#examples
	w.Header().Set("Strict-Transport-Security", "max-age=63072000")

	w.WriteHeader(http.StatusNotFound)
	io.WriteString(w, `{"error":"not found","message":"see https://github.com/shogo82148/holidays-jp/ for more information."}`)
}
