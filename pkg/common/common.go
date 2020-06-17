package common

import (
    "os"
    "time"
)

func CalculateEndDate(t time.Time, wd time.Weekday) time.Time {
    next := int((wd - t.Weekday() + 7) % 7)
    y, m, d := t.Date()
    return time.Date(y, m, d+next+1, 0, 0, 0, -1, t.Location())
}

func FileExists(filename string) bool {
    info, err := os.Stat(filename)
    if os.IsNotExist(err) {
        return false
    }
    return !info.IsDir()
}
