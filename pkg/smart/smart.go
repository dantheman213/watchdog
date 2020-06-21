package smart

import (
    "fmt"
    "github.com/dantheman213/watchdog/pkg/cli"
    "github.com/dantheman213/watchdog/pkg/common"
    "github.com/dantheman213/watchdog/pkg/config"
    libTime "github.com/dantheman213/watchdog/pkg/time"
    "log"
    "math"
    "time"
)

func Start() {
    if config.Storage.Diagnostics.SMARTTestShort || config.Storage.Diagnostics.SMARTTestLong {
        log.Print("[S.M.A.R.T] Starting Scheduler Daemon...")

        if config.Storage.Diagnostics.SMARTTestShort {
            go startShortTestScheduler()
        }
        if config.Storage.Diagnostics.SMARTTestLong {
            go startLongTestScheduler()
        }
    }
}

func startShortTestScheduler() {
    for true {
        log.Println("[S.M.A.R.T] Short Test Scheduler timer has activated...")

        target := libTime.GetNextScheduleTimeInSeconds(config.Storage.Schedule.SMARTTestShort)
        delta := time.Now().Sub(target)
        sleepSecs := math.Abs(delta.Seconds())
        log.Printf("[S.M.A.R.T] Sleeping short test until %s (%s or %f seconds)\n", target.String(), delta.String(), sleepSecs)
        time.Sleep(time.Duration(sleepSecs) * time.Second)

        testAllDisks("short")
    }
}

func startLongTestScheduler() {
    for true {
        log.Println("[S.M.A.R.T] Long Test Scheduler timer has activated...")

        target := libTime.GetNextScheduleTimeInSeconds(config.Storage.Schedule.SMARTTestLong)
        delta := time.Now().Sub(target)
        sleepSecs := math.Abs(delta.Seconds())
        log.Printf("[S.M.A.R.T] Sleeping long test until %s (%s or %f seconds)\n", target.String(), delta.String(), sleepSecs)
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
