/*
Package http allows for log messages to be sent via http.

Example

NOTE: you can use the HTTP handler as a base for creating other handlers

    package main

    import (
        stdhttp "net/http"

        "github.com/go-playground/log"
        "github.com/go-playground/log/handlers/http"
    )

    func main() {

        header := make(stdhttp.Header)
        header.Set("Authorization", "Bearer 378978HJJFEWj73JENEWFN3475")

        h, err := http.New("https://logs.logserver.com:4565", "POST", header)
        if err != nil {
            // handle error, most likely URL parsing error
        }

        h.SetFilenameDisplay(log.Llongfile)

        log.RegisterHandler(h, log.AllLevels...)

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
package http
