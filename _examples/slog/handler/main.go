//go:build go1.21
// +build go1.21

package main

import (
	"github.com/go-playground/log/v8"
	stdlog "log"
	"log/slog"
)

func main() {

	// This example demonstrates how to redirect the std log and slog to this logger by using it as
	// an slog.Handler.
	log.RedirectGoStdLog(true)
	log.WithFields(log.G("grouped", log.F("key", "value"))).Debug("test")
	stdlog.Println("test stdlog")
	slog.Info("test slog", slog.Group("group", "key", "value"))
}
