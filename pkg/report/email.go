package report

import (
    "fmt"
    "github.com/dantheman213/watchdog/pkg/config"
    "log"
    "net/smtp"
)

func sendEmail(to, subject, body string) {
    from := config.Storage.EmailAccount.Address
    pass := config.Storage.EmailAccount.Password

    msg := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n%s", from, to, subject, body)

    err := smtp.SendMail(fmt.Sprintf("%s:%d", config.Storage.EmailAccount.SMTPHost, config.Storage.EmailAccount.SMTPPort),
        smtp.PlainAuth("", from, pass, config.Storage.EmailAccount.SMTPHost),
        from, []string{to}, []byte(msg))

    if err != nil {
        log.Printf("smtp error: %s", err)
        return
    }

    log.Print("email sent successfully")
}
