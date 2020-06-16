package cli

import (
    "bytes"
    "os/exec"
)

func RunCommand(command string) (stdout, stderr bytes.Buffer, err error) {
    cmd := exec.Command("bash", "-c", command)

    var outb, errb bytes.Buffer
    cmd.Stdout = &outb
    cmd.Stderr = &errb

    if err := cmd.Run(); err != nil {
        return outb, errb, err
    }

    return outb, errb, nil
}
