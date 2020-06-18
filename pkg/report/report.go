package report

import (
    "bufio"
    "fmt"
    "github.com/dantheman213/watchdog/pkg/cli"
    "github.com/dantheman213/watchdog/pkg/common"
    "github.com/dantheman213/watchdog/pkg/config"
    "log"
    "math"
    "strings"
    "time"
)

// Example SMART report
//
// Model Family:     Seagate Barracuda Compute
// Device Model:     ST8000DM004-2CX188
// Serial Number:    ABCD1234
// SMART support is: Available - device has SMART capability.
// SMART support is: Enabled
// SMART overall-health self-assessment test result: PASSED
// === START OF READ SMART DATA SECTION ===
// SMART Self-test log structure revision number 1
// Num  Test_Description    Status                  Remaining  LifeTime(hours)  LBA_of_first_error
// # 1  Short offline       Completed without error       00%      500         -
// # 2  Short offline       Completed without error       00%      480         -
// # 3  Short offline       Completed without error       00%      460         -

func Start() {
    log.Print("[report] Starting Scheduler Daemon...")

    go startWeekly()
}

func startWeekly() {
    for true {
        log.Println("Report Weekly Scheduler timer has activated...")

        now := time.Now()
        target := common.CalculateEndDate(now, time.Saturday)
        delta := now.Sub(target)

        sleepSecs := math.Abs(delta.Seconds())
        log.Printf("[report] Sleeping until %s (%s or %f seconds)\n", target.String(), delta.String(), sleepSecs)
        time.Sleep(time.Duration(sleepSecs) * time.Second)

        generateReports()
    }
}

func generateReports() {
    log.Print("Generating Reports...")
    report := ""

    disks, err := common.GetDisks()
    if err != nil {
        log.Fatal(err)
    }

    for _, disk := range *disks {
        log.Printf("[report] Gathering info on disk %s for report...", disk)

        report += fmt.Sprintf("\n\nDisk Path: %s\n", disk)
        o, _, err := cli.RunCommand(fmt.Sprintf(`/usr/sbin/smartctl -i %s | grep -e SMART -e Available -e "Model Family" -e "Device Model" -e "Serial Number"`, disk))
        if err != nil {
            log.Println(err)
            report += fmt.Sprintf("\n%s\n", err)
            continue
        }

        scanner := bufio.NewScanner(&o)
        for scanner.Scan() {
            report += fmt.Sprintf("%s\n", strings.TrimSpace(scanner.Text()))
        }

        o, _, err = cli.RunCommand(fmt.Sprintf(`/usr/sbin/smartctl -a %s | grep -e "test result" -e " PASS" -e " FAIL"`, disk))
        if err != nil {
            log.Println(err)
            report += fmt.Sprintf("\n%s\n", err)
            continue
        }

        scanner = bufio.NewScanner(&o)
        for scanner.Scan() {
            report += fmt.Sprintf("%s\n", strings.TrimSpace(scanner.Text()))
        }

        o, _, err = cli.RunCommand(fmt.Sprintf(`/usr/sbin/smartctl -l selftest %s | grep -A 5 "=== START OF READ SMART DATA SECTION ==="`, disk))
        if err != nil {
            log.Println(err)
            report += fmt.Sprintf("\n%s\n", err)
            continue
        }

        scanner = bufio.NewScanner(&o)
        for scanner.Scan() {
            report += fmt.Sprintf("%s\n", strings.TrimSpace(scanner.Text()))
        }
    }

    log.Println("[report] preparing to send report email")
    sendEmail(config.Storage.EmailAccount.Address, "Watchdog Diagnostics Server Results", report)
}
