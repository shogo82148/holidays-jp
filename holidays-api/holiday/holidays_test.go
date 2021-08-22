package holiday

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
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

func TestFindHolidaysInMonth(t *testing.T) {
	got := findHolidaysInMonth(2000, time.January)
	want := []Holiday{
		{
			Date: "2000-01-01",
			Name: "元日",
		},
		{
			Date: "2000-01-10",
			Name: "成人の日",
		},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("holidays not match: (-want/+got)\n%s", diff)
	}
}
