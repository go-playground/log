package main

import (
	"io"
	stdlog "log"

	"github.com/go-playground/errors/v5"
	"github.com/go-playground/log/v8"
)

func main() {
	log.RedirectGoStdLog(true)

	// Trace
	defer log.WithTrace().Info("time to run")

	log.Debug("debug")
	log.Info("info")
	log.Notice("notice")
	log.Warn("warn")
	log.Error("error")
	// log.Panic("panic") // this will panic
	log.Alert("alert")
	// log.Fatal("fatal") // this will call os.Exit(1)

	err := errors.New("this is the inner error").AddTags(errors.T("inner", "tag"))
	err = errors.Wrap(err, "this is the wrapping error").AddTags(errors.T("outer", "tag"))

	// logging with fields can be used with any of the above
	log.WithError(err).WithFields(log.F("key", "value")).Info("test info")

	// log unwrapped error
	log.WithError(io.EOF).Error("unwrapped error")

	// predefined global fields
	log.WithDefaultFields(log.Fields{
		log.F("program", "test"),
		log.F("version", "0.1.3"),
	}...)

	log.WithField("key", "value").Info("testing default fields")

	// or request scoped default fields
	logger := log.WithFields(
		log.F("request", "req"),
		log.F("scoped", "sco"),
	)

	logger.WithField("key", "value").Info("test")

	stdlog.Println("This was redirected from Go STD output!")
	log.RedirectGoStdLog(false)
	stdlog.Println("This was NOT redirected from Go STD output!")
}
