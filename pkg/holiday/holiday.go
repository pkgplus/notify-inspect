package holiday

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

const (
	URL_HOLIDAY = "http://v.juhe.cn/calendar/month"
)

var (
	mutex        sync.RWMutex
	load_year    int
	calendar_key = "b08417aa8d791abd716bfd90c5f45068"
	holiday_map  = make(map[string]string)
	loaded_years = make(map[int]bool)
)

type Holiday struct {
	Name     string             `json:"name"`
	Festival string             `json:"festival"`
	Desc     string             `json:"desc"`
	Rest     string             `json:"rest"`
	List     []DayHolidayDetail `json:"list"`
}

type DayHolidayDetail struct {
	Date   string `json:"date"`
	Status string `json:"status"`
}

func init() {
	key := os.Getenv("JUHE_CALENDAR_KEY")
	if key != "" {
		calendar_key = key
	}

	// loadHoliday(time.Now().Year())
}

type HolidayYearMonthResp struct {
	Reason string `json:"reason"`
	Result struct {
		Data struct {
			Year        string `json:"year"`
			YearMonth   string `json:"year-month"`
			HolidayData string `json:"holiday"`
		} `json:"data"`
	} `json:"result"`
}

func loadHoliday(year int) {
	if loaded_years[year] {
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	loaded_years[year] = true
	params := url.Values{}
	params.Set("key", calendar_key)

	for i := 1; i <= 12; i++ {
		year_month := fmt.Sprintf("%d-%d", year, i)
		params.Set("year-month", year_month)
		req_url := URL_HOLIDAY + "?" + params.Encode()
		// log.Printf("get %s", req_url)

		resp, err := http.Get(req_url)
		if err != nil {
			log.Printf("ERROR load %s holiday failed: %v", year_month, err)
			continue
		}
		defer resp.Body.Close()
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("ERROR load %s holiday failed: %v", year_month, err)
			continue
		}

		hs, err := decodeResult(data)
		if err != nil {
			log.Printf("ERROR load %s holiday failed: %v", year_month, err)
			continue
		}

		for _, h := range hs {
			for _, hday := range h.List {
				holiday_map[hday.Date] = hday.Status
			}
		}
	}
}

func decodeResult(data []byte) ([]*Holiday, error) {
	hs := make([]*Holiday, 0)

	resp := new(HolidayYearMonthResp)
	err := json.Unmarshal(data, resp)
	if err != nil {
		return hs, err
	}

	if resp.Reason != "Success" {
		return hs, fmt.Errorf(resp.Reason)
	}

	err = json.Unmarshal([]byte(resp.Result.Data.HolidayData), &hs)
	if err != nil {
		return hs, err
	}

	return hs, nil
}

func IsHoliday(date time.Time) bool {
	loadHoliday(date.Year())

	mutex.RLock()
	defer mutex.RUnlock()

	day_str := fmt.Sprintf("%d-%d-%d", date.Year(), date.Month(), date.Day())
	status, found := holiday_map[day_str]
	if !found {
		if date.Weekday() == time.Saturday || date.Weekday() == time.Sunday {
			return true
		} else {
			return false
		}
	}

	if status == "1" {
		return true
	} else {
		return false
	}
}
