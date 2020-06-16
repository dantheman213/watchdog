package main

import (
    "github.com/dantheman213/watchdog/pkg/smart"
    "github.com/dantheman213/watchdog/pkg/zfs"
    "log"
)

var quit = make(chan struct{})

func main() {
    log.Println("Starting Watchdog...")

    smart.Start()
    zfs.Start()

    // blocking
    <-quit
}
