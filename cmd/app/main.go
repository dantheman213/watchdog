package main

import "github.com/dantheman213/watchdog/pkg/smart"

var quit = make(chan struct{})

func main() {
    smart.Start()

    // blocking
    <-quit
}
