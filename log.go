package log

import (
	"context"
	"os"
	"time"
)

const (
	// DefaultTimeFormat is the default time format when parsing Time values.
	// it is exposed to allow handlers to use and not have to redefine
	DefaultTimeFormat = "2006-01-02T15:04:05.000000000Z07:00"
)

var (
	logFields   []Field
	logHandlers = map[Level][]Handler{}
	exitFunc    = os.Exit
	ctxIdent    = &struct {
		name string
	}{
		name: "log",
	}
)

// Field is a single Field key and value
type Field struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

// SetExitFunc sets the provided function as the exit function used in Fatal(),
// Fatalf(), Panic() and Panicf(). This is primarily used when wrapping this library,
// you can set this to to enable testing (with coverage) of your Fatal() and Fatalf()
// methods.
func SetExitFunc(fn func(code int)) {
	exitFunc = fn
}

// SetContext sets a log entry into the provided context
func SetContext(ctx context.Context, e Entry) context.Context {
	return context.WithValue(ctx, ctxIdent, e)
}

// GetContext returns the log Entry found in the context,
// or a new Default log Entry if none is found
func GetContext(ctx context.Context) Entry {
	v := ctx.Value(ctxIdent)
	if v == nil {
		return newEntryWithFields(nil)
	}
	return v.(Entry)
}

func handleEntry(e Entry) {
	if !e.start.IsZero() {
		e = e.WithField("duration", time.Since(e.start))
	}
	for _, h := range logHandlers[e.Level] {
		h.Log(e)
	}
}

// F creates a new Field using the supplied key + value.
// it is shorthand for defining field manually
func F(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
}

// AddHandler adds a new log handler and accepts which log levels that
// handler will be triggered for
func AddHandler(h Handler, levels ...Level) {
	for _, level := range levels {
		handler := append(logHandlers[level], h)
		logHandlers[level] = handler
	}
}

// WithDefaultFields adds fields to the underlying logger instance
func WithDefaultFields(fields ...Field) {
	logFields = append(logFields, fields...)
}

// WithField returns a new log entry with the supplied field.
func WithField(key string, value interface{}) Entry {
	ne := newEntryWithFields(logFields)
	ne.Fields = append(ne.Fields, Field{Key: key, Value: value})
	return ne
}

// WithFields returns a new log entry with the supplied fields appended
func WithFields(fields ...Field) Entry {
	ne := newEntryWithFields(logFields)
	ne.Fields = append(ne.Fields, fields...)
	return ne
}

// WithTrace withh add duration of how long the between this function call and
// the susequent log
func WithTrace() Entry {
	ne := newEntryWithFields(logFields)
	ne.start = time.Now()
	return ne
}

// WithError add a minimal stack trace to the log Entry
func WithError(err error) Entry {
	ne := newEntryWithFields(logFields)
	return ne.withError(err)
}

// Debug logs a debug entry
func Debug(v ...interface{}) {
	e := newEntryWithFields(logFields)
	e.Debug(v...)
}

// Debugf logs a debug entry with formatting
func Debugf(s string, v ...interface{}) {
	e := newEntryWithFields(logFields)
	e.Debugf(s, v...)
}

// Info logs a normal. information, entry
func Info(v ...interface{}) {
	e := newEntryWithFields(logFields)
	e.Info(v...)
}

// Infof logs a normal. information, entry with formatiing
func Infof(s string, v ...interface{}) {
	e := newEntryWithFields(logFields)
	e.Infof(s, v...)
}

// Notice logs a notice log entry
func Notice(v ...interface{}) {
	e := newEntryWithFields(logFields)
	e.Notice(v...)
}

// Noticef logs a notice log entry with formatting
func Noticef(s string, v ...interface{}) {
	e := newEntryWithFields(logFields)
	e.Noticef(s, v...)
}

// Warn logs a warn log entry
func Warn(v ...interface{}) {
	e := newEntryWithFields(logFields)
	e.Warn(v...)
}

// Warnf logs a warn log entry with formatting
func Warnf(s string, v ...interface{}) {
	e := newEntryWithFields(logFields)
	e.Warnf(s, v...)
}

// Panic logs a panic log entry
func Panic(v ...interface{}) {
	e := newEntryWithFields(logFields)
	e.Panic(v...)
}

// Panicf logs a panic log entry with formatting
func Panicf(s string, v ...interface{}) {
	e := newEntryWithFields(logFields)
	e.Panicf(s, v...)
}

// Alert logs an alert log entry
func Alert(v ...interface{}) {
	e := newEntryWithFields(logFields)
	e.Alert(v...)
}

// Alertf logs an alert log entry with formatting
func Alertf(s string, v ...interface{}) {
	e := newEntryWithFields(logFields)
	e.Alertf(s, v...)
}

// Fatal logs a fatal log entry
func Fatal(v ...interface{}) {
	e := newEntryWithFields(logFields)
	e.Fatal(v...)
}

// Fatalf logs a fatal log entry with formatting
func Fatalf(s string, v ...interface{}) {
	e := newEntryWithFields(logFields)
	e.Fatalf(s, v...)
}

// Error logs an error log entry
func Error(v ...interface{}) {
	e := newEntryWithFields(logFields)
	e.Error(v...)
}

// Errorf logs an error log entry with formatting
func Errorf(s string, v ...interface{}) {
	e := newEntryWithFields(logFields)
	e.Errorf(s, v...)
}

// Handler is an interface that log handlers need to comply with
type Handler interface {
	Log(Entry)
}
