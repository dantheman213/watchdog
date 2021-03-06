package report

import (
    "bufio"
    "fmt"
    "github.com/dantheman213/watchdog/pkg/cli"
    "github.com/dantheman213/watchdog/pkg/common"
    "github.com/dantheman213/watchdog/pkg/config"
    libTime "github.com/dantheman213/watchdog/pkg/time"
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
    if config.Storage.Diagnostics.EmailReport {
        log.Print("[report] Starting Scheduler Daemon...")
        go startScheduler()
    }
}

func startScheduler() {
    for true {
        log.Println("[report] Scheduler timer has activated...")
        target := libTime.GetNextScheduleTimeInSeconds(config.Storage.Schedule.Report)
        delta := time.Now().Sub(target)
        sleepSecs := math.Abs(delta.Seconds())
        log.Printf("[report] Sleeping until %s (%s or %f seconds)\n", target.String(), delta.String(), sleepSecs)
        time.Sleep(time.Duration(sleepSecs) * time.Second)

        generateReports()
    }
}

func generateReports() {
    log.Print("Generating Reports...")
    header := fmt.Sprintf("<h1>%s</h1><p>Server: <strong>%s</strong></p>", config.Storage.ReportName, config.Storage.ServerName)
    report := ""
    testResultSummary := "<h2>Summary</h2>\n"

    if config.Storage.Diagnostics.SMARTTestShort || config.Storage.Diagnostics.SMARTTestLong {
        disks, err := common.GetDisks()
        if err != nil {
            log.Fatal(err)
        }

        report += "<h2>Details</h2>"
        report += fmt.Sprintf("<h3>S.M.A.R.T Disk Results</h3>")
        for _, disk := range *disks {
            report += "<p>"
            log.Printf("[report] Gathering info on disk %s for report...", disk)

            // S.M.A.R.T reports

            // get a summary PASS/FAIL from each disk
            o, _, err := cli.RunCommand(fmt.Sprintf(`smartctl -a %s | grep -e ": PASSED" -e ": FAILED"`, disk))
            if err != nil {
                log.Println(err)
                testResultSummary += fmt.Sprintf("<strong>Disk %s: ERROR %s</strong><br />", disk, err)
                continue
            }

            scanner := bufio.NewScanner(&o)
            if scanner.Scan() {
                result := scanner.Text()
                if strings.Index(result, "PASSED") > -1 {
                    result = "PASSED"
                } else if strings.Index(result, "PASSED") > -1 {
                    result = "FAILED"
                } else {
                    result = "UNKNOWN"
                }
                testResultSummary += fmt.Sprintf("Disk %s : <strong>%s</strong><br />", disk, result)
            }

            report += fmt.Sprintf("Disk Path: <strong>%s</strong><br />", disk)
            o, _, err = cli.RunCommand(fmt.Sprintf(`/usr/sbin/smartctl -i %s | grep -e SMART -e Available -e "Model Family" -e "Device Model" -e "Serial Number"`, disk))
            if err != nil {
                log.Println(err)
                report += fmt.Sprintf("<strong>%s</strong>", err)
                continue
            }

            scanner = bufio.NewScanner(&o)
            for scanner.Scan() {
                report += fmt.Sprintf("%s<br />", strings.TrimSpace(scanner.Text()))
            }

            o, _, err = cli.RunCommand(fmt.Sprintf(`/usr/sbin/smartctl -a %s | grep -e "test result" -e " PASS" -e " FAIL"`, disk))
            if err != nil {
                log.Println(err)
                report += fmt.Sprintf("<strong>%s</strong><br />", err)
                continue
            }

            scanner = bufio.NewScanner(&o)
            for scanner.Scan() {
                report += fmt.Sprintf("%s<br />", strings.TrimSpace(scanner.Text()))
            }

            o, _, err = cli.RunCommand(fmt.Sprintf(`/usr/sbin/smartctl -l selftest %s | grep -A 10 "=== START OF READ SMART DATA SECTION ==="`, disk))
            if err != nil {
                log.Println(err)
                report += fmt.Sprintf("<strong>%s</strong><br />", err)
                continue
            }

            scanner = bufio.NewScanner(&o)
            for scanner.Scan() {
                report += fmt.Sprintf("%s<br />", strings.TrimSpace(scanner.Text()))
            }

            report += "</p>"
        }
    }

    if config.Storage.Diagnostics.ZFSPoolScrub {
        report += "<p>"
        // ZFS pool report
        o, _, err := cli.RunCommand(`/usr/sbin/zpool status -x`)
        if err != nil {
            log.Println(err)
            testResultSummary += fmt.Sprintf("<strong>%s</strong><br />", err)
        }
        scanner := bufio.NewScanner(&o)
        if scanner.Scan() {
            testResultSummary += fmt.Sprintf("<p>ZFS pool(s): <strong>%s</strong></p>", strings.TrimSpace(scanner.Text()))
        }

        report += fmt.Sprintf("\n\n<h3>ZFS Pool Results</h3>\n\n")
        o, _, err = cli.RunCommand(`/usr/sbin/zpool status -v`)
        if err != nil {
            log.Println(err)
            report += fmt.Sprintf("<strong>%s</strong><br />", err)
        }
        scanner = bufio.NewScanner(&o)
        for scanner.Scan() {
            report += fmt.Sprintf("%s<br />", strings.TrimSpace(scanner.Text()))
        }
        report += "</p>"
    }

    if config.Storage.Diagnostics.UPSTest {
        report += "<p>"
        // UPS status report
        o, _, err := cli.RunCommand(`/usr/sbin/pwrstat -status | grep Result`)
        if err != nil {
            log.Println(err)
            testResultSummary += fmt.Sprintf("<strong>%s</strong><br />", err)
        }
        scanner := bufio.NewScanner(&o)
        if scanner.Scan() {
            testResultSummary += fmt.Sprintf("UPS hardware: <strong>%s</strong><br />", strings.TrimSpace(scanner.Text()))
        }

        report += fmt.Sprintf("<h3>UPS Results</h3>")
        o, _, err = cli.RunCommand(`/usr/sbin/pwrstat -status`)
        if err != nil {
            log.Println(err)
            report += fmt.Sprintf("<strong>%s</strong><br />", err)
        }
        scanner = bufio.NewScanner(&o)
        for scanner.Scan() {
            report += fmt.Sprintf("%s<br />", strings.TrimSpace(scanner.Text()))
        }
        report += "</p>"
    }

    log.Println("[report] preparing to send report email...")
    subject := fmt.Sprintf("%s -- %s", config.Storage.ReportName, config.Storage.ServerName)
    body := header + "\n"
    if config.Storage.Diagnostics.SMARTTestShort || config.Storage.Diagnostics.SMARTTestLong {
        body += testResultSummary + "\n"
    }
    body += report
    payload := fmt.Sprintf("<!DOCTYPE HTML PUBLIC \"-//W3C//DTD HTML 4.01 Transitional//EN\">\n<html>\n<head>\n<meta http-equiv=\"content-type\" content=\"text/html; charset=ISO-8859-1\">\n</head>\n<body bgcolor=\"#ffffff\" text=\"#000000\">\n%s\n</body>\n</html>\n", body)
    sendEmail(config.Storage.EmailAccount.Address, subject, "text/html", payload)
}
