package time

import (
    "github.com/dantheman213/watchdog/pkg/config"
    "strconv"
    "strings"
    "time"
)

// Return Time of next target weekday at midnight local time
func CalculateTimeUntilTargetWeekday(t time.Time, wd time.Weekday) time.Time {
    next := int((wd - t.Weekday() + 7) % 7)
    if t.Weekday() == wd {
        next = 7
    }
    y, m, d := t.Date()
    return time.Date(y, m, d+next, 0, 0, 0, 0, t.Location())
}

func GetNextScheduleTimeInSeconds(schedule *[]config.TimeItem) time.Time {
    now := time.Now()
    timeItems := make([]time.Time, 0)

    for _, item := range *schedule {
        t := CalculateTimeUntilTargetWeekday(now, time.Weekday(item.Day))
        parts := strings.Split(item.Time, ":")
        hours, _ := strconv.Atoi(parts[0])
        minutes, _ := strconv.Atoi(parts[1])
        t = t.Add(time.Duration(hours) * time.Hour)
        t = t.Add(time.Duration(minutes) * time.Minute)
        timeItems = append(timeItems, t)
    }

    minSecs := -1.0
    for i, item := range timeItems {
        offset := time.Until(item).Seconds()
        if i == 0 {
            minSecs = offset
            continue
        }
        if minSecs > offset {
            minSecs = offset
        }
    }

    return time.Now().Add(time.Duration(minSecs) * time.Second)
}
