package zfs

import (
    "bufio"
    "fmt"
    "github.com/dantheman213/watchdog/pkg/cli"
    "github.com/dantheman213/watchdog/pkg/common"
    "math"
    "strings"
    "time"
)

func Start() {
    log.Println("Starting ZFS Scheduler Daemon...")

    go startWeekly()
}

func startWeekly() {
    for true {
        log.Println("ZFS Weekly Scrub Scheduler has activated...")
        now := time.Now()
        target := common.CalculateEndDate(now, time.Wednesday)
        delta := now.Sub(target)

        sleepSecs := math.Abs(delta.Seconds())
        log.Printf("Sleeping until %s (%s or %f seconds)\n", target.String(), delta.String(), sleepSecs)
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
            log.Printf("Starting ZFS scrub on pool %s\n", pool)
            _, _, err := cli.RunCommand(fmt.Sprintf("/usr/sbin/zpool scrub %s", pool))
            if err != nil {
                log.Printf("Could not start ZFS scrub on pool %s. Received ERROR: %s\n", pool, err)
            } else {
                log.Printf("Started ZFS scrub successfully on %s", pool)
            }
        }
    }
}
