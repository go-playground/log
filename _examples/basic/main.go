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
	defer log.WithTrace().Info("time to run")

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
	log.WithError(err).WithFields(log.Fields{"key": "value"}).Info("test info")

	// predefined global fields
	log.WithDefaultFields(log.Fields{
		"program": "test",
		"version": "0.1.3",
	})

	log.WithField("key", "value").Info("testing default fields")

	// or request scoped default fields
	logger := log.WithFields(log.Fields{
		"request": "req",
		"scoped":  "sco",
	})

	logger.WithField("key", "value").Info("test")
}
