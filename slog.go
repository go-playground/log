//go:build go1.21
// +build go1.21

package log

import (
	"context"
	runtimeext "github.com/go-playground/pkg/v5/runtime"
	"log/slog"
	"runtime"
)

var _ slog.Handler = (*slogHandler)(nil)

type slogHandler struct {
	e     Entry
	group string
}

// Enabled returns if the current logging level is enabled. In the case of this log package in this Level has a
// handler registered.
func (s *slogHandler) Enabled(_ context.Context, level slog.Level) bool {
	rw.RLock()
	_, enabled := logHandlers[convertSlogLevel(level)]
	rw.RUnlock()
	return enabled
}

func (s *slogHandler) Handle(ctx context.Context, record slog.Record) error {

	var fields []Field
	if record.NumAttrs() > 0 {
		fields = make([]Field, 0, record.NumAttrs())
		record.Attrs(func(attr slog.Attr) bool {
			if attr.Value.Kind() == slog.KindGroup {
				fields = append(fields, Field{Key: attr.Key, Value: s.convertAttrsToFields(attr.Value.Group())})
				return true
			}
			fields = append(fields, s.convertAttrToField(attr))
			return true
		})
	}
	if record.Level >= slog.LevelError && record.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{record.PC})
		f, _ := fs.Next()
		sourceBuff := BytePool().Get()
		sourceBuff.B = extractSource(sourceBuff.B, runtimeext.Frame{Frame: f})
		fields = append(fields, Field{Key: slog.SourceKey, Value: string(sourceBuff.B[:len(sourceBuff.B)-1])})
		BytePool().Put(sourceBuff)
	}
	e := s.e.clone(fields...)
	e.Message = record.Message
	e.Level = convertSlogLevel(record.Level)
	HandleEntry(e)
	return nil
}

func (s *slogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &slogHandler{
		e:     s.e.clone(s.convertAttrsToFields(attrs)...),
		group: s.group,
	}
}

func (s *slogHandler) convertAttrsToFields(attrs []slog.Attr) []Field {
	fields := make([]Field, 0, len(attrs))

	for _, attr := range attrs {
		switch attr.Value.Kind() {
		case slog.KindGroup:
			fields = append(fields, Field{Key: attr.Key, Value: s.convertAttrsToFields(attr.Value.Group())})
			continue
		default:
			fields = append(fields, s.convertAttrToField(attr))
		}
	}
	return fields
}

func (s *slogHandler) convertAttrToField(attr slog.Attr) Field {
	var value any

	switch attr.Value.Kind() {
	case slog.KindLogValuer:
		value = attr.Value.LogValuer().LogValue()
	default:
		value = attr.Value.Any()
	}
	return Field{Key: attr.Key, Value: value}
}

func (s *slogHandler) WithGroup(name string) slog.Handler {
	return &slogHandler{
		e:     s.e.clone(),
		group: name,
	}
}

func convertSlogLevel(level slog.Level) Level {
	switch level {
	case slog.LevelDebug:
		return DebugLevel
	case slog.LevelInfo:
		return InfoLevel
	case SlogNoticeLevel:
		return NoticeLevel
	case slog.LevelWarn:
		return WarnLevel
	case slog.LevelError:
		return ErrorLevel
	case SlogPanicLevel:
		return PanicLevel
	case SlogAlertLevel:
		return AlertLevel
	case SlogFatalLevel:
		return FatalLevel
	default:
		switch {
		case level > slog.LevelInfo && level < slog.LevelWarn:
			return NoticeLevel
		case level > slog.LevelError && level <= SlogPanicLevel:
			return PanicLevel
		case level > SlogPanicLevel && level <= SlogAlertLevel:
			return AlertLevel
		case level > SlogAlertLevel && level <= SlogFatalLevel:
			return FatalLevel
		}
		return InfoLevel
	}
}

var (
	prevSlogLogger *slog.Logger
)

// RedirectGoStdLog is used to redirect Go's internal std log output to this logger AND registers a handler for slog
// that redirects slog output to this logger.
//
// If you intend to use this log interface with another slog handler then you should not use this function and instead
// register a handler with slog directly and register the slog redirect, found under the handlers package or other
// custom redirect handler with this logger.
func RedirectGoStdLog(redirect bool) {
	if redirect {
		prevSlogLogger = slog.Default()
		slog.SetDefault(slog.New(&slogHandler{e: newEntry()}))
	} else if prevSlogLogger != nil {
		slog.SetDefault(prevSlogLogger)
		prevSlogLogger = nil
	}
}

// slog log levels.
const (
	SlogDebugLevel  slog.Level = slog.LevelDebug
	SlogInfoLevel   slog.Level = slog.LevelInfo
	SlogWarnLevel   slog.Level = slog.LevelWarn
	SlogErrorLevel  slog.Level = slog.LevelError
	SlogNoticeLevel slog.Level = slog.LevelInfo + 2
	SlogPanicLevel  slog.Level = slog.LevelError + 4
	SlogAlertLevel  slog.Level = SlogPanicLevel + 4
	SlogFatalLevel  slog.Level = SlogAlertLevel + 4 // same as syslog CRITICAL
)
