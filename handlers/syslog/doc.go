/*
Package syslog allows for log messages to be sent via syslog.

Example

NOTE: syslog uses github.com/RackSec/srslog as the stdlib syslog
      is no longer being maintained or added to as of this discussion
      https://github.com/golang/go/issues/13449#issuecomment-161204716

    package main

    import (
        stdsyslog "log/syslog"

        "github.com/go-playground/log/v8"
        "github.com/go-playground/log/v8/handlers/syslog"
    )

    func main() {

        sysLog, err := syslog.New("udp", "log.logs.com:4863", stdsyslog.LOG_DEBUG, "")
        if err != nil {
            // handle error
        }

        log.AddHandler(sysLog, log.AllLevels...)

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

        // logging with fields can be used with any of the above
        log.WithField("key", "value").Info("test info")
    }
*/
package syslog
