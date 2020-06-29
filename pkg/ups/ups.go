package ups

import (
    "github.com/dantheman213/watchdog/pkg/cli"
    "github.com/dantheman213/watchdog/pkg/config"
    libTime "github.com/dantheman213/watchdog/pkg/time"
    "log"
    "math"
    "time"
)

func Start() {
    if config.Storage.Diagnostics.UPSTest {
        log.Println("[UPS] Starting Scheduler Daemon...")
        go startScheduler()
    }
}

func startScheduler() {
    for true {
        log.Println("[UPS] Scrub Scheduler has activated...")
        target := libTime.GetNextScheduleTimeInSeconds(config.Storage.Schedule.UPSTest)
        delta := time.Now().Sub(target)
        sleepSecs := math.Abs(delta.Seconds())
        log.Printf("[UPS] Sleeping until %s (%s or %f seconds)\n", target.String(), delta.String(), sleepSecs)
        time.Sleep(time.Duration(sleepSecs) * time.Second)

        log.Println("[UPS] Starting battery and power test...")
        _, _, err := cli.RunCommand(`/usr/sbin/pwrstat -test`)
        if err != nil {
            log.Println(err)
            continue
        }
    }
}
