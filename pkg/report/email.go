package report

import (
    "fmt"
    "github.com/dantheman213/watchdog/pkg/config"
    "log"
    "net/smtp"
)

func sendEmail(to, subject, mimeType, body string) {
    headerFrom := config.Storage.EmailAccount.Address
    headerTo := fmt.Sprintf("To: %s\n", to)
    headerSubject := fmt.Sprintf("Subject: %s\n", subject)
    headerMIME := fmt.Sprintf(`MIME-version: 1.0;\nContent-Type: %s; charset="UTF-8";\n\n`, mimeType)
    msg := fmt.Sprintf("%s%s%s%s\n%s", headerFrom, headerTo, headerSubject, headerMIME, body)

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
