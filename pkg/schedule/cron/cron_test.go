package cron

import (
	"testing"
	"time"
)

func TestNextRunTime(t *testing.T) {
	var cur time.Time
	var next int64

	s, err := newTestSetting()
	if err != nil {
		t.Fatal(err)
	}
	// cur := time.Now().Truncate(time.Minute)
	cur = time.Unix(1505552769, 0).Truncate(time.Minute)
	next_time, _ := s.NextRunTime(cur)
	next = next_time.Unix()
	if next != cur.Add(time.Minute).Unix() {
		t.Fatalf("expect %d but get %d", cur.Add(time.Minute).Unix(), next)
	}

	s2, err := newTestSetting_workDay()
	if err != nil {
		t.Fatal(err)
	}
	next_time, _ = s2.NextRunTime(cur)
	next = next_time.Unix()
	if next != 1505692800 {
		t.Fatalf("expect 1505692800 but get %d", next)
	}

	cur = time.Unix(1506682804, 0).Truncate(time.Minute)
	s3, err := newTestSetting_workDay2()
	if err != nil {
		t.Fatal(err)
	}
	next_time, _ = s3.NextRunTime(cur)
	next = next_time.Unix()
	if next != 1506902400 {
		t.Fatalf("expect 1506902400 but get %d", next)
	}

	// holiday test
	s4, err := newTestSetting_notHoliday()
	if err != nil {
		t.Fatal(err)
	}
	cur = time.Unix(1514761200, 0).Truncate(time.Minute)
	next_time, _ = s4.NextRunTime(cur)
	next = next_time.Unix()
	if next != 1514851200 {
		t.Fatalf("expect 1514851200 but get %d", next)
	}

	// holiday test
	// config: 08:00 ~ 09:00
	// curtime: 2018-02-14 09:00:00
	// holiday: 2018-02-15 ~ 2018-02-21
	// expect: 2018-02-22 08:00:00
	cur = time.Unix(1518570000, 0).Truncate(time.Minute)
	next_time, _ = s4.NextRunTime(cur)
	next = next_time.Unix()
	if next != 1519257600 {
		t.Fatalf("expect 1519257600 but get %d", next)
	}
}

func newTestSetting() (*CronTaskSetting, error) {
	return NewCronTaskSetting([]byte(`
{
"firstTimeStr":"20170916103000",
"interval": "1m"
}
`))
}

func newTestSetting_workDay() (*CronTaskSetting, error) {
	return NewCronTaskSetting([]byte(`
{
"interval": "1m",
"weekLimit":"weekday",
"clockLimitStart":"08:00",
"clockLimitEnd":"12:00"
}
`))
}

func newTestSetting_workDay2() (*CronTaskSetting, error) {
	return NewCronTaskSetting([]byte(`
{
	"interval":"10m",
	"intervalDuration":600000000000,
	"clockLimitStart":"08:00",
	"clockLimitEnd":"09:00",
	"weekLimit":"weekday"
}
`))
}

func newTestSetting_notHoliday() (*CronTaskSetting, error) {
	return NewCronTaskSetting([]byte(`
{
	"interval":"10m",
	"intervalDuration":600000000000,
	"clockLimitStart":"08:00",
	"clockLimitEnd":"09:00",
	"weekLimit":"notHoliday"
}
`))
}
