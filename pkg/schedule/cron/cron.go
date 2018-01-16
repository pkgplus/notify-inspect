package cron

import (
	"encoding/json"
	"errors"
	"fmt"
	// "log"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/xuebing1110/notify-inspect/pkg/holiday"
)

var (
	localZone, _        = time.LoadLocation("Local")
	defaultFirstTime, _ = time.ParseInLocation("20060102150405", "19700101000000", localZone)
)

type CronTaskSetting struct {
	TaskId           string        `json:"taskId"`
	FirstTimeStr     string        `json:"firstTimeStr,omitempty"`
	Interval         string        `json:"interval"` // 1h or 5m
	IntervalDuration time.Duration `json:"intervalDuration"`

	ClockLimitStart string `json:"clockLimitStart,omitempty"`
	ClockLimitEnd   string `json:"clockLimitEnd,omitempty"`
	WeekLimit       string `json:"weekLimit,omitempty"` //

	firstTime time.Time
	loop      int64
}

func newCronTaskSetting(setting_bytes []byte, s *CronTaskSetting) error {
	err := json.Unmarshal(setting_bytes, s)
	if err != nil {
		return err
	}

	return s.Init()
}

func NewCronTaskSetting(setting_bytes []byte) (s *CronTaskSetting, err error) {
	s = new(CronTaskSetting)
	err = newCronTaskSetting(setting_bytes, s)
	return
}

func (s *CronTaskSetting) Init() (err error) {
	// interval
	if s.Interval == "" {
		return errors.New("interval is required")
	}

	// intervalDuration
	if s.IntervalDuration.Seconds() == 0 {
		s.IntervalDuration, err = s.getIntervalDuration()
		if err != nil {
			return
		}
	}

	// firstTimeStr
	if s.FirstTimeStr != "" {
		s.firstTime, err = time.ParseInLocation("20060102150405", s.FirstTimeStr, localZone)
		if err != nil {
			return
		}
	} else {
		s.firstTime = defaultFirstTime
	}

	return
}

func (s *CronTaskSetting) getIntervalDuration() (time.Duration, error) {
	var duration_desc = s.Interval
	if strings.HasSuffix(duration_desc, "d") {
		days, err := strconv.Atoi(duration_desc[:(len(duration_desc) - 1)])
		if err != nil {
			return 0, err
		}
		return time.ParseDuration(fmt.Sprintf("%dh", days*24))
	} else {
		return time.ParseDuration(duration_desc)
	}
}

func (s *CronTaskSetting) NextRunTime(curtime time.Time) (next_run time.Time, err error) {
	err = s.Init()
	if err != nil {
		return
	}

	return s.nextRunTime(curtime)
}

func (s *CronTaskSetting) nextRunTime(curtime time.Time) (next_run time.Time, err error) {
	// loop count too much, return zero
	s.loop++
	if s.loop == 1000 {
		return
	}

	var duration = s.IntervalDuration
	if s.firstTime.After(curtime) {
		return s.firstTime, nil
	}

	var loopCouont = curtime.Sub(s.firstTime) / duration
	next_run = s.firstTime.Add(duration * (loopCouont + 1))

	if len(s.ClockLimitStart) == 5 && strings.Compare(next_run.Format("15:04"), s.ClockLimitStart) < 0 {
		hour, _ := strconv.Atoi(s.ClockLimitStart[:2])
		minute, _ := strconv.Atoi(s.ClockLimitStart[3:5])

		next_run, err = s.nextRunTime(time.Date(
			next_run.Year(), next_run.Month(), next_run.Day(),
			hour, minute, next_run.Second(),
			next_run.Nanosecond(), next_run.Location()).
			Add(0 - time.Nanosecond),
		)
		if err != nil {
			return
		}

		return s.runWeekLimit(next_run), nil
	} else if len(s.ClockLimitEnd) == 5 && strings.Compare(next_run.Format("15:04"), s.ClockLimitEnd) > 0 {
		hour, _ := strconv.Atoi(s.ClockLimitStart[:2])
		minute, _ := strconv.Atoi(s.ClockLimitStart[3:5])
		next_day := next_run.Add(24 * time.Hour)
		next_run, err = s.nextRunTime(
			time.Date(
				next_day.Year(), next_day.Month(), next_day.Day(),
				hour, minute, next_run.Second(),
				next_run.Nanosecond(), next_run.Location()).
				Add(0 - time.Nanosecond),
		)
		if err != nil {
			return
		}

		return s.runWeekLimit(next_run), nil
	} else {
		return s.runWeekLimit(next_run), nil
	}
}

func (s *CronTaskSetting) runWeekLimit(next_run time.Time) time.Time {
	if s.WeekLimit == "weekday" {
		if next_run.Weekday() == time.Sunday {
			next_run = next_run.Add(24 * time.Hour)
		} else if next_run.Weekday() == time.Saturday {
			next_run = next_run.Add(48 * time.Hour)
		}
	} else if s.WeekLimit == "notHoliday" {
		for holiday.IsHoliday(next_run) {
			next_run = next_run.Add(24 * time.Hour)
		}
	} else if s.WeekLimit != "" {
		week_array := strings.Split(s.WeekLimit, ",")
		sort.Strings(week_array)
		var first_week int = -1
		for _, week_str := range week_array {
			week, err := strconv.Atoi(week_str)
			if err != nil {
				continue
			}
			if first_week != -1 {
				first_week = week
			}

			next_run_week := int(next_run.Weekday())
			if next_run_week == week {
				return next_run
			} else if next_run_week < week {
				return next_run.Add(24 * time.Duration(week-next_run_week) * time.Hour)
			} else {

			}
		}

		next_run_week := int(next_run.Weekday())
		return next_run.Add(24 * time.Duration(first_week-next_run_week+7) * time.Hour)
	} else {
		return next_run
	}

	return next_run
}

func (s *CronTaskSetting) String() string {
	bytes, _ := json.Marshal(s)
	return string(bytes)
}

func (s *CronTaskSetting) ParseTaskId() (string, string, string, error) {
	array := strings.Split(s.TaskId, ".")
	if len(array) != 3 {
		return "", "", "", fmt.Errorf("parse %s failed, taskId must be <userid.pluginid.recordid>", s.TaskId)
	}

	return array[0], array[1], array[2], nil
}
