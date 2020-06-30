package report

import (
    "fmt"
    "github.com/dantheman213/watchdog/pkg/config"
    "log"
    "net/smtp"
)

func sendEmail(to, subject, mimeType, body string) {
    headerFrom := fmt.Sprintf("From: %s", config.Storage.EmailAccount.Address)
    headerSubject := fmt.Sprintf("Subject: %s\n", subject)
    headerTo := fmt.Sprintf("To: %s\n", to)
    headerMIME := "MIME-version: 1.0;\n"
    headerContentType := fmt.Sprintf("Content-Type: %s; charset=\"UTF-8\";\n", mimeType)
    msg := fmt.Sprintf("%s%s%s%s%s\n%s", headerFrom, headerSubject, headerTo, headerMIME, headerContentType, body)

    err := smtp.SendMail(
                fmt.Sprintf("%s:%d", config.Storage.EmailAccount.SMTPHost, config.Storage.EmailAccount.SMTPPort),
                smtp.PlainAuth("", config.Storage.EmailAccount.Address, config.Storage.EmailAccount.Password, config.Storage.EmailAccount.SMTPHost),
                config.Storage.EmailAccount.Address, []string{to}, []byte(msg),
            )

    if err != nil {
        log.Printf("smtp error: %s", err)
        return
    }

    log.Print("email sent successfully")
}
