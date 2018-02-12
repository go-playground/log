/*
Package console allows for log messages to be sent to a any writer, default os.Stderr.

Example

simple console

    package main

    import (
        "errors"

        "github.com/go-playground/log"
        "github.com/go-playground/log/handlers/console"
    )

    func main() {
        cLog := console.New(true)
        log.AddHandler(cLog, log.AllLevels...)

        // Trace
        defer log.WithTrace().Info("took this long")

        log.Debug("debug")
        log.Info("info")
        log.Notice("notice")
        log.Warn("warn")
        log.Error("error")
        // log.Panic("panic") // this will panic
        log.Alert("alert")
        // log.Fatal("fatal") // this will call os.Exit(1)

        err := errors.New("the is an error")
        // logging with fields can be used with any of the above
        log.WithError(err).WithFields(log.F("key", "value")).Info("test info")

        // predefined global fields
        log.WithDefaultFields(log.Fields{
            {"program", "test"},
            {"version", "0.1.3"},
        }...)

        log.WithField("key", "value").Info("testing default fields")

        // or request scoped default fields
        logger := log.WithFields(
            log.F("request", "req"),
            log.F("scoped", "sco"),
        )

        logger.WithField("key", "value").Info("test")
    }
*/
package console
