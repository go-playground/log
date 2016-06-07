/*
Package email allows for log messages to be sent via email.

Example

simple email

    package main

    import (
        "github.com/go-playground/log"
        "github.com/go-playground/log/handlers/email"
    )

    func main() {

        email := email.New("smtp.gmail.com", 587, "username", "password", "from@email.com", []string{"to@email.com"})
        email.SetFilenameDisplay(log.Llongfile)

        log.RegisterHandler(email, log.WarnLevel, log.AlertLevel, log.PanicLevel)

        log.WithFields(log.F("error", "something went wrong")).StackTrace().Alert("ALERT!")
    }
*/
package email
