package pkg

import (
	"fmt"
	"time"
)

const (
	Second = time.Second
	Minute = time.Minute
	Hour   = time.Hour
	Day    = 24 * time.Hour
	Week   = 7 * Day
	Month  = 4 * Week
	Year   = 12 * Month
)

var timeUnits = []struct {
	size int64
	unit string
}{
	{int64(Year), "year"},
	{int64(Month), "month"},
	{int64(Week), "week"},
	{int64(Day), "day"},
	{int64(Hour), "hour"},
	{int64(Minute), "minute"},
	{int64(Second), "second"},
}

// HumaneDuration is a pretty printed version of a time.Duration value
type HumaneDuration time.Duration

func (d HumaneDuration) String() string {
	return timeFormat(int64(d))
}

// HumaneTime is a pretty printed version of a time.Time value
type HumaneTime time.Time

func (t HumaneTime) String() string {
	diff := time.Since(time.Time(t))
	return timeFormat(int64(diff))
}

func timeFormat(duration int64) string {
	result := ""
	for _, unit := range timeUnits {
		n := duration / unit.size
		if n > 0 {
			if result != "" {
				result += ", "
			}

			result += fmt.Sprintf("%d %s", n, unit.unit)
		}

		duration %= unit.size
	}

	if result == "" {
		return fmt.Sprintf("%d %s", 0, timeUnits[len(timeUnits)-1].unit)
	}

	return result
}
