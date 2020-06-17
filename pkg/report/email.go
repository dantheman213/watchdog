package report

import (
    "log"
    "net/smtp"
)

func sendEmail(to, body string) {
    from := "...@gmail.com"
    pass := "..."
    //to := "foobarbazz@mailinator.com"

    msg := "From: " + from + "\n" +
        "To: " + to + "\n" +
        "Subject: Hello there\n\n" +
        body

    err := smtp.SendMail("smtp.gmail.com:587",
        smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
        from, []string{to}, []byte(msg))

    if err != nil {
        log.Printf("smtp error: %s", err)
        return
    }

    log.Print("sent, visit http://foobarbazz.mailinator.com")
}
