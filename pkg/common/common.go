package common

import (
    "bufio"
    "github.com/dantheman213/watchdog/pkg/cli"
    "github.com/dantheman213/watchdog/pkg/config"
    "log"
    "os"
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

func FileExists(filename string) bool {
    info, err := os.Stat(filename)
    if os.IsNotExist(err) {
        return false
    }
    return !info.IsDir()
}

func GetDisks() (*[]string, error) {
    disks := make([]string, 0)

    stdout, _, err := cli.RunCommand(`/sbin/fdisk -l | grep Disk | grep -e /dev/nvme -e /dev/sd -e /dev/hd | awk '{print $2}' | sed 's/.$//'`)
    if err != nil {
        log.Println("ERROR: could not get disks...")
        return nil, err
    }

    scanner := bufio.NewScanner(&stdout)
    for scanner.Scan() {
        disks = append(disks, strings.TrimSpace(scanner.Text()))
    }

    return &disks, nil
}
