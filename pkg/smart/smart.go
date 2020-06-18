package smart

import (
    "fmt"
    "github.com/dantheman213/watchdog/pkg/cli"
    "github.com/dantheman213/watchdog/pkg/common"
    "log"
    "math"
    "time"
)

func Start() {
    log.Print("[S.M.A.R.T] Starting Scheduler Daemon...")

    go startDaily()
    go startWeekly()
}

func startDaily() {
    for true {
        log.Println("[S.M.A.R.T] Daily Scan Scheduler timer has activated...")

        now := time.Now()
        future := now.AddDate(0, 0, 1)
        target := time.Date(future.Year(), future.Month(), future.Day(), 0, 0, 0, 0, time.Local)
        delta := now.Sub(target)

        sleepSecs := math.Abs(delta.Seconds())
        log.Printf("[S.M.A.R.T] Sleeping until %s (%s or %f seconds)\n", target.String(), delta.String(), sleepSecs)
        time.Sleep(time.Duration(sleepSecs) * time.Second)

        testAllDisks("short")
    }
}

func startWeekly() {
    for true {
        log.Println("[S.M.A.R.T] Weekly Scan Scheduler timer has activated...")
        now := time.Now()
        target := common.CalculateEndDate(now, time.Sunday)
        delta := now.Sub(target)

        sleepSecs := math.Abs(delta.Seconds())
        log.Printf("[S.M.A.R.T] Sleeping until %s (%s or %f seconds)\n", target.String(), delta.String(), sleepSecs)
        time.Sleep(time.Duration(sleepSecs) * time.Second)

        testAllDisks("long")
    }
}

// duration = (short, long)
func testAllDisks(duration string) {
    log.Printf("[S.M.A.R.T] Starting '%s' type test for all disks...\n", duration)

    disks, err := common.GetDisks()
    if err != nil {
        log.Fatal(err)
    }

    for _, disk := range *disks {
        log.Printf("[S.M.A.R.T] Starting test %s on disk %s\n", duration, disk)
        _, _, err := cli.RunCommand(fmt.Sprintf("/usr/sbin/smartctl -t %s %s", duration, disk))
        if err != nil {
            log.Printf("[S.M.A.R.T] Could not start test on disk %s.\nReceived ERROR: %s\n", disk, err)
        } else {
            log.Printf("[S.M.A.R.T] Started %s disk test successfully on %s", duration, disk)
        }
    }
}
