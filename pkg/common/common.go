package common

import (
    "bufio"
    "github.com/dantheman213/watchdog/pkg/cli"
    "log"
    "os"
    "strings"
)

func FileExists(filename string) bool {
    info, err := os.Stat(filename)
    if os.IsNotExist(err) {
        return false
    }
    return !info.IsDir()
}

func GetDisks() (*[]string, error) {
    disks := make([]string, 0)

    stdout, _, err := cli.RunCommand(`/sbin/fdisk -l | grep Disk | grep -e /dev/nvme -e /dev/sd -e /dev/hd | awk '{print $2}' | sed 's/.$//'`)
    if err != nil {
        log.Println("ERROR: could not get disks...")
        return nil, err
    }

    scanner := bufio.NewScanner(&stdout)
    for scanner.Scan() {
        disks = append(disks, strings.TrimSpace(scanner.Text()))
    }

    return &disks, nil
}
