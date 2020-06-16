package common

import "time"

func CalculateEndDate(t time.Time, wd time.Weekday) time.Time {
    next := int((wd - t.Weekday() + 7) % 7)
    y, m, d := t.Date()
    return time.Date(y, m, d+next+1, 0, 0, 0, -1, t.Location())
}