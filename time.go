package gotk

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

var (
	_DateRE     = regexp.MustCompile(`^\d+-\d{2}-\d{2}$`)
	_ClockRE    = regexp.MustCompile(`^\d{2}:\d{2}:\d{2}$`)
	_DatetimeRE = regexp.MustCompile(`^\d+-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}$`)
)

func NowMs() string {
	return time.Now().Format("2006-01-02T15:04:05.000Z07:00")
}

// datetime format: "2021-06-24", "09:00:01", "2021-06-24 09:10:11" or "2021-06-24T09:10:11"
func ParseDatetime(value string) (at time.Time, err error) {
	value = strings.TrimSpace(value)
	value = strings.Replace(value, " ", "T", 1)

	switch {
	case _DateRE.Match([]byte(value)):
		return time.ParseInLocation("2006-01-02", value, time.Local)
	case _ClockRE.Match([]byte(value)):
		now := time.Now()
		return time.ParseInLocation(
			"2006-01-02 15:04:05",
			now.Format("2006-01-02")+" "+value,
			time.Local,
		)
	case _DatetimeRE.Match([]byte(value)):
		return time.ParseInLocation("2006-01-02T15:04:05", value, time.Local)
	default:
		return time.Parse(time.RFC3339, value)
	}
}

func GetOClock(shift int) (clock time.Time) {
	clock = time.Now().AddDate(0, 0, shift)
	clock = time.Date(clock.Year(), clock.Month(), clock.Day(), 0, 0, 0, 0, clock.Location())
	return clock
}

/*
ceil time, e.g. TimeCeil('2020-12-01T17:39:07.123+08:00', "M") -> '2020-12-01T17:40:00+08:00'

	valid unit(key or value)
	H: hour, M: minute, S: second
	y: year, s: season, m: month, w: week, d: day
*/
func TimeCeil(at time.Time, tu string) (out time.Time, err error) {
	var (
		location        *time.Location
		year, day, hour int
		minute, second  int
		month           time.Month
	)

	location = at.Location()
	year, month, day = at.Date()
	hour, minute, second = at.Clock()

	switch tu {
	case "second", "S":
		out = time.Date(year, month, day, hour, minute, second, 0, location)
		if at.Sub(out) > time.Millisecond { //!
			out = out.Add(time.Second)
		}
	case "minute", "M":
		out = time.Date(year, month, day, hour, minute, 0, 0, location)
		if at.Sub(out) > time.Second {
			out = out.Add(time.Minute)
		}
	case "hour", "H":
		out = time.Date(year, month, day, hour, 0, 0, 0, location)
		if at.Sub(out) > time.Second {
			out = out.Add(time.Hour)
		}
	case "day", "d":
		out = time.Date(year, month, day, 0, 0, 0, 0, location)
		if at.Sub(out) > time.Second {
			out = out.AddDate(0, 0, 1)
		}
	case "week", "w":
		w := int(at.Weekday()) // 0 for Sunday, the first day of a week is Monday
		if w == 0 {
			w = 7
		}
		year, month, day = at.AddDate(0, 0, 1-w).Date()
		out = time.Date(year, month, day, 0, 0, 0, 0, location)
		if at.Sub(out) > time.Second {
			out = out.AddDate(0, 0, 7)
		}
	case "month", "m":
		out = time.Date(year, month, 1, 0, 0, 0, 0, location)
		if at.Sub(out) > time.Second {
			out = out.AddDate(0, 1, 0)
		}
	case "season", "s":
		month = (month-1)/3*3 + 1
		out = time.Date(year, month, 1, 0, 0, 0, 0, location)
		if at.Sub(out) > time.Second {
			out = out.AddDate(0, 3, 0)
		}
	case "year", "y":
		out = time.Date(year, 1, 1, 0, 0, 0, 0, location)
		if at.Sub(out) > time.Second {
			out = out.AddDate(1, 0, 0)
		}
	default:
		return out, fmt.Errorf("invalid time unit")
	}

	return
}

/*
floor time, e.g. TimeFloor('2020-12-01T17:39:07.123+08:00', "M") -> '2020-12-01T17:39:00+08:00'

	valid unit(key or value)
	H: hour, M: minute, S: second
	y: year, s: season, m: month, w: week, d: day
*/
func TimeFloor(at time.Time, tu string) (out time.Time, err error) {
	var (
		location        *time.Location
		year, day, hour int
		second, minute  int
		month           time.Month
	)

	location = at.Location()
	year, month, day = at.Date()
	hour, minute, second = at.Clock()

	switch tu {
	case "second", "S":
		out = time.Date(year, month, day, hour, minute, second, 0, location)
	case "minute", "M":
		out = time.Date(year, month, day, hour, minute, 0, 0, location)
	case "hour", "H":
		out = time.Date(year, month, day, hour, 0, 0, 0, location)
	case "day", "d":
		out = time.Date(year, month, day, 0, 0, 0, 0, location)
	case "week", "w":
		w := int(at.Weekday()) // 0 for Sunday, the first day of a week is Monday
		if w == 0 {
			w = 7
		}
		year, month, day = at.AddDate(0, 0, 1-w).Date()
		out = time.Date(year, month, day, 0, 0, 0, 0, location)
	case "month", "m":
		out = time.Date(year, month, 1, 0, 0, 0, 0, location)
	case "season", "s":
		month = (month-1)/3*3 + 1
		out = time.Date(year, month, 1, 0, 0, 0, 0, location)
	case "year", "y":
		out = time.Date(year, 1, 1, 0, 0, 0, 0, location)
	default:
		return out, fmt.Errorf("invalid time unit")
	}

	return
}

func TimeTag(at time.Time) string {
	// 2006-01-02T15-03-04MST
	// zone, _ := at.Zone()
	// fmt.Sprintf("%sT%s%s", at.Format(time.DateOnly), strings.Replaat.Format(time.TimeOnly), zone)
	return at.Format("2006-01-02T15-03-04MST")
}

func TruncateHour(t ...time.Time) (start time.Time, end time.Time) {
	var at time.Time

	if len(t) > 0 {
		at = t[0]
	} else {
		at = time.Now()
	}

	year, month, day := at.Date()
	hour, _, _ := at.Clock()
	// fmt.Println(now.Truncate(time.Hour))

	start = time.Date(year, month, day, hour, 0, 0, 0, at.Location())
	if hour == 23 {
		end = time.Date(year, month, day, 0, 0, 0, 0, at.Location()).AddDate(0, 0, 1)
	} else {
		end = time.Date(year, month, day, hour+1, 0, 0, 0, at.Location())
	}

	return start, end
}

func TruncateDate(t ...time.Time) (start time.Time, end time.Time) {
	var at time.Time

	if len(t) > 0 {
		at = t[0]
	} else {
		at = time.Now()
	}

	year, month, day := at.Date()

	start = time.Date(year, month, day, 0, 0, 0, 0, at.Location())
	end = start.AddDate(0, 0, 1)

	return start, end
}

func TruncateMonth(t ...time.Time) (start time.Time, end time.Time) {
	var at time.Time

	if len(t) > 0 {
		at = t[0]
	} else {
		at = time.Now()
	}

	year, month, _ := at.Date()

	start = time.Date(year, month, 1, 0, 0, 0, 0, at.Location())
	end = start.AddDate(0, 1, 0)

	return start, end
}
