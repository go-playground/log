/*
Package syslog allows for log messages to be sent via syslog.

Example

NOTE: syslog uses github.com/RackSec/srslog as the stdlib syslog
      is no longer being maintained or added to as of this discussion
      https://github.com/golang/go/issues/13449#issuecomment-161204716

    package main

    import (
        stdsyslog "log/syslog"

        "github.com/go-playground/log"
        "github.com/go-playground/log/handlers/syslog"
    )

    func main() {

        sysLog, err := syslog.New("udp", "log.logs.com:4863", stdsyslog.LOG_DEBUG, "")
        if err != nil {
            // handle error
        }

        sysLog.SetFilenameDisplay(log.Llongfile)

        log.RegisterHandler(sysLog, log.AllLevels...)

        // Trace
        defer log.Trace("trace").End()

        log.Debug("debug")
        log.Info("info")
        log.Notice("notice")
        log.Warn("warn")
        log.Error("error")
        // log.Panic("panic") // this will panic
        log.Alert("alert")
        // log.Fatal("fatal") // this will call os.Exit(1)

        // logging with fields can be used with any of the above
        log.WithFields(log.F("key", "value")).Info("test info")
    }
*/
package syslog
