package config

import (
    "encoding/json"
    "github.com/dantheman213/watchdog/pkg/common"
    "io/ioutil"
    "log"
    "os"
)

type EmailAccount struct {
    Address string
    Username string
    Password string
    SMTPPort int
    SMTPHost string
}

type Schedule struct {
    Report *[]TimeItem
    SMARTTestShort *[]TimeItem
    SMARTTestLong *[]TimeItem
    ZFSTestScrub *[]TimeItem
}

type TimeItem struct {
    Day int
    Time string
}

type Diagnostics struct {
    SMARTTestShort bool
    SMARTTestLong bool
    ZFSPoolScrub bool
    EmailReport bool
}

type Config struct {
    ServerName string
    ReportName string
    Diagnostics *Diagnostics
    EmailAccount *EmailAccount
    Schedule *Schedule
}

var Storage *Config
var filePath string = "/etc/watchdog/config.json"

func Load() {
    if !common.FileExists(filePath) {
        log.Fatal("config file does not exist")
    }

    file, err := os.Open(filePath)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    b, err := ioutil.ReadAll(file)
    if err != nil {
        log.Fatal(err)
    }

    if err := json.Unmarshal([]byte(b), &Storage); err != nil {
        log.Fatal(err)
    }
}
