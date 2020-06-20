package zfs

import (
    "bufio"
    "fmt"
    "github.com/dantheman213/watchdog/pkg/cli"
    libTime "github.com/dantheman213/watchdog/pkg/time"
    "log"
    "math"
    "strings"
    "time"
)

func Start() {
    log.Println("[ZFS] Starting Scheduler Daemon...")

    go startWeekly()
}

func startWeekly() {
    for true {
        log.Println("[ZFS] Weekly Scrub Scheduler has activated...")
        now := time.Now()
        target := libTime.CalculateTimeUntilTargetWeekday(now, time.Wednesday)
        delta := now.Sub(target)

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
