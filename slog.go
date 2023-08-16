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
	// List of Groups, each subsequent group belongs to the previous group, except the first
	// which are the top level fields fields before any grouping.
	groups []Field
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

	var current Field
	if len(s.groups) == 0 {
		current = G("")
	} else {
		group := s.groups[len(s.groups)-1]
		last := group.Value.([]Field)
		fields := make([]Field, len(last), len(last)+record.NumAttrs()+1)
		copy(fields, last)

		current = F(group.Key, fields)
	}

	if record.NumAttrs() > 0 {
		record.Attrs(func(attr slog.Attr) bool {
			current.Value = s.convertAttrToField(current.Value.([]Field), attr)
			return true
		})
	}
	if record.Level >= slog.LevelError && record.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{record.PC})
		f, _ := fs.Next()
		sourceBuff := BytePool().Get()
		sourceBuff.B = extractSource(sourceBuff.B, runtimeext.Frame{Frame: f})
		current.Value = append(current.Value.([]Field), F(slog.SourceKey, string(sourceBuff.B[:len(sourceBuff.B)-1])))
		BytePool().Put(sourceBuff)
	}

	for i := len(s.groups) - 2; i >= 0; i-- {
		group := s.groups[i]
		gf := group.Value.([]Field)
		copied := make([]Field, len(gf), len(gf)+1)
		copy(copied, gf)
		current = G(group.Key, append(copied, current)...)
	}

	var e Entry
	if current.Key == "" {
		e = Entry{Fields: current.Value.([]Field)}
	} else {
		e = Entry{Fields: []Field{current}}
	}
	e.Message = record.Message
	e.Level = convertSlogLevel(record.Level)
	e.Timestamp = record.Time

	HandleEntry(e)
	return nil
}

func (s *slogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	var groups []Field
	if len(s.groups) == 0 {
		groups = []Field{G("", s.convertAttrsToFields(nil, attrs)...)}
	} else {
		groups = make([]Field, len(s.groups))
		copy(groups, s.groups)

		l := len(groups) - 1
		current := groups[l]
		currentFields := current.Value.([]Field)
		copiedFields := make([]Field, len(currentFields), len(currentFields)+len(attrs))
		copy(copiedFields, currentFields)
		groups[l].Value = s.convertAttrsToFields(copiedFields, attrs)
	}

	return &slogHandler{
		groups: groups,
	}
}

func (s *slogHandler) convertAttrsToFields(fields []Field, attrs []slog.Attr) []Field {
	for _, attr := range attrs {
		if attr.Key == "" {
			continue
		}
		if attr.Key == slog.TimeKey && attr.Value.Time().IsZero() {
			continue
		}
		fields = s.convertAttrToField(fields, attr)
	}
	return fields
}

func (s *slogHandler) convertAttrToField(fields []Field, attr slog.Attr) []Field {
	var value any

	switch attr.Value.Kind() {
	case slog.KindLogValuer:
		return s.convertAttrToField(fields, slog.Attr{Key: attr.Key, Value: attr.Value.LogValuer().LogValue()})

	case slog.KindGroup:
		attrs := attr.Value.Group()
		groupedFields := make([]Field, 0, len(attrs))
		value = s.convertAttrsToFields(groupedFields, attrs)

	default:
		value = attr.Value.Any()
	}
	return append(fields, F(attr.Key, value))
}

func (s *slogHandler) WithGroup(name string) slog.Handler {
	groups := make([]Field, len(s.groups), len(s.groups)+1)
	copy(groups, s.groups)

	return &slogHandler{
		groups: append(groups, G(name)),
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
		slog.SetDefault(slog.New(&slogHandler{}))
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
