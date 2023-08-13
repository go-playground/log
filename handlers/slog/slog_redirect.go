//go:build go1.21
// +build go1.21

package slog

import (
	"context"
	"github.com/go-playground/log/v8"
	"log/slog"
)

// Handler implementation.
type Handler struct {
	handler slog.Handler
}

// New handler wraps an slog.Handler for log output.
//
// Calling this function automatically calls the slog.RedirectGoStdLog function in order to intercept and forward
// the Go standard library log output to this handler.
func New(handler slog.Handler) *Handler {
	log.RedirectGoStdLog(true)
	return &Handler{handler: handler}
}

// Log handles the log entry
func (h *Handler) Log(e log.Entry) {
	attrs := make([]slog.Attr, 0, len(e.Fields))
	for _, f := range e.Fields {
		attrs = append(attrs, slog.Any(f.Key, f.Value))
	}

	r := slog.NewRecord(e.Timestamp, slog.Level(e.Level), e.Message, 0)
	r.AddAttrs(h.convertFields(e.Fields)...)
	h.handler.Handle(context.Background(), r)
}

func (h *Handler) convertFields(fields []log.Field) []slog.Attr {
	attrs := make([]slog.Attr, 0, len(fields))
	for _, f := range fields {
		switch t := f.Value.(type) {
		case []log.Field:
			a := h.convertFields(t)
			arr := make([]any, 0, len(a))
			for _, v := range a {
				arr = append(arr, v)
			}
			attrs = append(attrs, slog.Group(f.Key, arr...))
		default:
			attrs = append(attrs, slog.Any(f.Key, f.Value))
		}
	}
	return attrs
}

// ReplaceAttrFn can be used with slog.HandlerOptions to replace attributes.
// This function replaces the "level" attribute to get the custom log levels of this package.
var ReplaceAttrFn = func(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.LevelKey {
		level := log.Level(a.Value.Any().(slog.Level))
		a.Value = slog.StringValue(level.String())
		//switch {
		//case level > slog.LevelInfo && level < slog.LevelWarn:
		//	a.Value = slog.StringValue(log.NoticeLevel.String())
		//case level > slog.LevelError && level <= log.SlogPanicLevel:
		//	a.Value = slog.StringValue(log.PanicLevel.String())
		//case level > log.SlogPanicLevel && level <= log.SlogAlertLevel:
		//	a.Value = slog.StringValue(log.AlertLevel.String())
		//case level > log.SlogAlertLevel && level <= log.SlogFatalLevel:
		//	a.Value = slog.StringValue(log.FatalLevel.String())
		//}
	}
	return a
}
