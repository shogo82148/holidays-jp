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

func TestFindHolidaysInYear(t *testing.T) {
	got := findHolidaysInYear(2000)
	want := []Holiday{
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
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("holidays not match: (-want/+got)\n%s", diff)
	}
}

func TestCalcHolidaysInMonthWithoutInLieu(t *testing.T) {
	got := calcHolidaysInMonthWithoutInLieu(2022, time.January)
	want := []Holiday{
		{
			Date: "2022-01-01",
			Name: "元日",
		},
		{
			Date: "2022-01-10",
			Name: "成人の日",
		},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("holidays not match: (-want/+got)\n%s", diff)
	}
}
