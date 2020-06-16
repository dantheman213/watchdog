package main

import (
    "fmt"
    "github.com/dantheman213/watchdog/pkg/smart"
    "github.com/dantheman213/watchdog/pkg/zfs"
)

var quit = make(chan struct{})

func main() {
    fmt.Println("Starting Watchdog...")

    smart.Start()
    zfs.Start()

    // blocking
    <-quit
}
