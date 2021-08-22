package holiday

import (
	"testing"
	"time"
)

func TestFindHoliday(t *testing.T) {
	h, ok := findHoliday(2000, time.January, 1)
	if !ok {
		t.Error("want true, but got false")
	}
	if got, want := h.Date, "2000-01-01"; want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}
