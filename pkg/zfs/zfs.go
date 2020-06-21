package zfs

import (
    "bufio"
    "fmt"
    "github.com/dantheman213/watchdog/pkg/cli"
    "github.com/dantheman213/watchdog/pkg/config"
    libTime "github.com/dantheman213/watchdog/pkg/time"
    "log"
    "math"
    "strings"
    "time"
)

func Start() {
    if config.Storage.Diagnostics.ZFSPoolScrub {
        log.Println("[ZFS] Starting Scheduler Daemon...")
        go startScheduler()
    }
}

func startScheduler() {
    for true {
        log.Println("[ZFS] Scrub Scheduler has activated...")
        target := libTime.GetNextScheduleTimeInSeconds(config.Storage.Schedule.ZFSTestScrub)
        delta := time.Now().Sub(target)
        sleepSecs := math.Abs(delta.Seconds())
        log.Printf("[ZFS] Sleeping until %s (%s or %f seconds)\n", target.String(), delta.String(), sleepSecs)
        time.Sleep(time.Duration(sleepSecs) * time.Second)

        stdout, _, err := cli.RunCommand(`/usr/sbin/zpool list | sed -n '1d;p' | awk '{print $1}'`)
        if err != nil {
            log.Println(err)
            continue
        }

        pools := make([]string, 0)
        scanner := bufio.NewScanner(&stdout)
        for scanner.Scan() {
            pools = append(pools, strings.TrimSpace(scanner.Text()))
        }

        for _, pool := range pools {
            log.Printf("[ZFS] Starting scrub on pool %s\n", pool)
            _, _, err := cli.RunCommand(fmt.Sprintf("/usr/sbin/zpool scrub %s", pool))
            if err != nil {
                log.Printf("[ZFS] Could not start scrub on pool %s. Received ERROR: %s\n", pool, err)
            } else {
                log.Printf("[ZFS] Started scrub successfully on %s", pool)
            }
        }
    }
}
