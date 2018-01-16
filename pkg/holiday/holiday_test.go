package holiday

import (
	"testing"
	"time"
)

func TestIsHoliday(t *testing.T) {
	date, err := time.ParseInLocation("2006-01-02", "2018-01-01", time.Local)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("get date: %d-%d-%d", date.Year(), date.Month(), date.Day())
	if !IsHoliday(date) {
		t.Fatal("expect holiday, but get not")
	}
}
