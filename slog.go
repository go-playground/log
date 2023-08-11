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
				g := attr.Key
				if s.group != "" {
					g = s.group + "." + g
				}
				fields = append(fields, s.convertAttrsToFields(g, attr.Value.Group())...)
				return true
			}
			fields = append(fields, s.convertAttrToField(s.group, attr))
			return true
		})
	}
	if record.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{record.PC})
		f, _ := fs.Next()
		sourceBuff := BytePool().Get()
		sourceBuff.B = extractSource(sourceBuff.B, runtimeext.Frame{Frame: f})
		fields = append(fields, Field{Key: "source", Value: string(sourceBuff.B[:len(sourceBuff.B)-1])})
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
		e:     s.e.clone(s.convertAttrsToFields(s.group, attrs)...),
		group: s.group,
	}
}

func (s *slogHandler) convertAttrsToFields(group string, attrs []slog.Attr) []Field {
	fields := make([]Field, 0, len(attrs))

	for _, attr := range attrs {
		switch attr.Value.Kind() {
		case slog.KindGroup:
			g := attr.Key
			if group != "" {
				g = group + "." + g
			}
			fields = append(fields, s.convertAttrsToFields(g, attr.Value.Group())...)
			continue
		default:
			fields = append(fields, s.convertAttrToField(group, attr))
		}
	}
	return fields
}

func (s *slogHandler) convertAttrToField(group string, attr slog.Attr) Field {
	var value any

	switch attr.Value.Kind() {
	case slog.KindLogValuer:
		value = attr.Value.LogValuer().LogValue()
	default:
		value = attr.Value.Any()
	}
	if group == "" {
		return Field{Key: attr.Key, Value: value}
	} else {
		return Field{Key: group + "." + attr.Key, Value: value}
	}
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
	case slog.LevelWarn:
		return WarnLevel
	case slog.LevelError:
		return ErrorLevel
	default:
		switch {
		case level > slog.LevelInfo && level < slog.LevelWarn:
			return NoticeLevel
		case level > slog.LevelError && level < slog.LevelError+4:
			return PanicLevel
		case level > slog.LevelError+4 && level < slog.LevelError+8:
			return AlertLevel
		case level > slog.LevelError+8 && level < slog.LevelError+16:
			return FatalLevel
		}
		return ErrorLevel
	}
}

// RedirectGoStdLog is used to redirect Go's internal std log output to this logger AND registers a handler for slog
// that redirects slog output to this logger.
func RedirectGoStdLog(redirect bool) {
	if redirect {
		slog.SetDefault(slog.New(&slogHandler{e: newEntry()}))
	} else {
		slog.SetDefault(slog.Default())
	}
}
