package time

import (
    "github.com/dantheman213/watchdog/pkg/config"
    "strconv"
    "strings"
    "time"
)

// Return Time of next target weekday at midnight local time
func CalculateTimeUntilTargetWeekday(wd time.Weekday, addHours, addMinutes int) time.Time {
    now := time.Now()
    next := int((wd - now.Weekday() + 7) % 7)
    if now.Weekday() == wd {
        if now.Hour() <= addHours && now.Minute() < addMinutes {
            next = 0
        } else {
            next = 7
        }
    }
    y, m, d := now.Date()
    return time.Date(y, m, d+next, addHours, addMinutes, 0, 0, now.Location())
}

func GetNextScheduleTimeInSeconds(schedule *[]config.TimeItem) time.Time {
    now := time.Now()
    timeItems := make([]time.Time, 0)

    for _, item := range *schedule {
        parts := strings.Split(item.Time, ":")
        hours, _ := strconv.Atoi(parts[0])
        minutes, _ := strconv.Atoi(parts[1])
        t := CalculateTimeUntilTargetWeekday(time.Weekday(item.Day), hours, minutes)
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
