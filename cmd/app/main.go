package main

import (
    "github.com/dantheman213/watchdog/pkg/config"
    "github.com/dantheman213/watchdog/pkg/report"
    "github.com/dantheman213/watchdog/pkg/smart"
    "github.com/dantheman213/watchdog/pkg/zfs"
    "log"
)

var quit = make(chan struct{})

func main() {
    log.Println("Starting Watchdog...")
    config.Load()

    smart.Start()
    zfs.Start()
    report.Start()

    // blocking
    <-quit
}
