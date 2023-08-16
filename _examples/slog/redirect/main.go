//go:build go1.21
// +build go1.21

package main

import (
	"github.com/go-playground/log/v8"
	slogredirect "github.com/go-playground/log/v8/handlers/slog"
	stdlog "log"
	"log/slog"
	"os"
)

func main() {

	// This example demonstrates how to redirect the std log and slog to this logger and output back out to any
	// slog.Handler, as well as any other handler(s) registered with this logger.
	log.AddHandler(slogredirect.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: slogredirect.ReplaceAttrFn, // for custom log level output
	})), log.AllLevels...)
	log.WithFields(log.G("grouped", log.F("key", "value"))).Debug("test")
	stdlog.Println("test stdlog")
	slog.Info("test slog", slog.Group("group", "key", "value"))
}
