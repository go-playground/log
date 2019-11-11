/*
Package email allows for log messages to be sent via email.

Example

simple email

    package main

    import (
        "github.com/go-playground/log/v7"
        "github.com/go-playground/log/v7/handlers/email"
    )

    func main() {

        email := email.New("smtp.gmail.com", 587, "username", "password", "from@email.com", []string{"to@email.com"})

        log.AddHandler(email, log.WarnLevel, log.AlertLevel, log.PanicLevel)

        log.WithField("omg", "something went wrong").Alert("ALERT!")
    }
*/
package email
