package smart

import (
    "bufio"
    "fmt"
    "github.com/dantheman213/watchdog/pkg/cli"
    "github.com/dantheman213/watchdog/pkg/common"
    "log"
    "math"
    "time"
)

func Start() {
    fmt.Println("Starting S.M.A.R.T Scheduler Daemon...")

    go startDaily()
    go startWeekly()
}

func getDisks() (*[]string, error) {
    disks := make([]string, 0)

    stdout, _, err := cli.RunCommand(`/sbin/fdisk -l | grep Disk | grep -e /dev/nvme -e /dev/sd -e /dev/hd | awk '{print $2}' | sed 's/.$//'`)
    if err != nil {
        fmt.Println("ERROR: could not get disks...")
        return nil, err
    }

    scanner := bufio.NewScanner(&stdout)
    for scanner.Scan() {
        disks = append(disks, scanner.Text())
    }

    return &disks, nil
}

func startDaily() {
    for true {
        fmt.Println("S.M.A.R.T Daily Scan Scheduler activating...")
        //now := time.Now()
        //future := now.AddDate(0, 0, 1)
        //target := time.Date(future.Year(), future.Month(), future.Day(), 0, 0, 0, 0, time.Local)
        //delta := now.Sub(target)
        //
        //sleepSecs := math.Abs(delta.Seconds())
        //fmt.Printf("Sleeping until %s (%s or %f seconds)\n", target.String(), delta.String(), sleepSecs)
        //time.Sleep(time.Duration(sleepSecs) * time.Second)

        startTests("short")
    }
}

func startWeekly() {
    for true {
        fmt.Println("S.M.A.R.T Weekly Scan Scheduler activating...")
        now := time.Now()
        target := common.CalculateEndDate(now, time.Sunday)
        delta := now.Sub(target)

        sleepSecs := math.Abs(delta.Seconds())
        fmt.Printf("Sleeping until %s (%s or %f seconds)\n", target.String(), delta.String(), sleepSecs)
        time.Sleep(time.Duration(sleepSecs) * time.Second)

        startTests("long")
    }
}

// duration = (short, long)
func startTests(duration string) {
    fmt.Printf("Starting '%s' type test for all disks...\n", duration)

    disks, err := getDisks()
    if err != nil {
        log.Fatal(err)
    }

    for _, disk := range *disks {
        _, _, err := cli.RunCommand(fmt.Sprintf("/sbin/smartctl -t %s %s\n", duration, disk))
        if err != nil {
            log.Fatal(err)
        }
    }
}
