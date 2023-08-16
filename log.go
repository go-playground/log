package log

import (
	"context"
	"os"
	"sync"
	"time"

	"golang.org/x/term"
)

var (
	bytePool = &byteArrayPool{pool: &sync.Pool{
		New: func() interface{} {
			return &Buffer{
				B: make([]byte, 0, 1024),
			}
		},
	}}
	defaultHandler *Logger
)

func init() {
	if term.IsTerminal(int(os.Stdin.Fd())) {
		h := NewConsoleBuilder().Build()
		AddHandler(h, AllLevels...)
		defaultHandler = h
	}
}

const (
	// DefaultTimeFormat is the default time format when parsing Time values.
	// it is exposed to allow handlers to use and not have to redefine
	DefaultTimeFormat = "2006-01-02T15:04:05.000000000Z07:00" // RFC3339Nano
)

var (
	logFields   []Field
	logHandlers = map[Level][]Handler{}
	exitFunc    = os.Exit
	withErrFn   = errorsWithError
	ctxIdent    = &struct {
		name string
	}{
		name: "log",
	}
	rw = new(sync.RWMutex)
)

// Field is a single Field key and value
type Field struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

// SetExitFunc sets the provided function as the exit function used in Fatal(),
// Fatalf(), Panic() and Panicf(). This is primarily used when wrapping this library,
// you can set this to enable testing (with coverage) of your Fatal() and Fatalf()
// methods.
func SetExitFunc(fn func(code int)) {
	exitFunc = fn
}

// SetWithErrorFn sets a custom WithError function handlers
func SetWithErrorFn(fn func(Entry, error) Entry) {
	withErrFn = fn
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
		return newEntry()
	}
	return v.(Entry)
}

// BytePool returns a sync.Pool of bytes that multiple handlers can use in order to reduce allocation and keep
// a central copy for reuse.
func BytePool() *byteArrayPool {
	return bytePool
}

// HandleEntry handles the log entry and fans out to all handlers with the proper log level
// This is exposed to allow for centralized logging whereby the log entry is marshalled, passed
// to a central logging server, unmarshalled and finally fanned out from there.
func HandleEntry(e Entry) {
	if !e.start.IsZero() {
		e = e.WithField("duration", time.Since(e.start))
	}
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now()
	}

	rw.RLock()
	for _, h := range logHandlers[e.Level] {
		h.Log(e)
	}
	rw.RUnlock()
}

// F creates a new Field using the supplied key + value.
// it is shorthand for defining field manually
func F(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
}

// G creates a new group of fields using the supplied key as the groups name.
func G(key string, fields ...Field) Field {
	return F(key, fields)
}

// AddHandler adds a new log handlers and accepts which log levels that
// handlers will be triggered for
func AddHandler(h Handler, levels ...Level) {
	rw.Lock()
	defer rw.Unlock()
	if defaultHandler != nil {
		removeHandler(defaultHandler)
		defaultHandler = nil
	}
	for _, level := range levels {
		handler := append(logHandlers[level], h)
		logHandlers[level] = handler
	}
}

// RemoveHandler removes an existing handler
func RemoveHandler(h Handler) {
	rw.Lock()
	removeHandler(h)
	rw.Unlock()
}

func removeHandler(h Handler) {
OUTER:
	for lvl, handlers := range logHandlers {
		for i, handler := range handlers {
			if h == handler {
				n := append(handlers[:i], handlers[i+1:]...)
				if len(n) == 0 {
					delete(logHandlers, lvl)
					continue OUTER
				}
				logHandlers[lvl] = n
				continue OUTER
			}
		}
	}
}

// removeHandlerLevels removes the supplied levels, if no more levels exists for the handler
// it will no longer be registered and need to added via AddHandler again.
func removeHandlerLevels(h Handler, levels ...Level) {
	rw.Lock()
	defer rw.Unlock()
OUTER:
	for _, lvl := range levels {
		handlers := logHandlers[lvl]
		for i, handler := range handlers {
			if h == handler {
				n := append(handlers[:i], handlers[i+1:]...)
				if len(n) == 0 {
					delete(logHandlers, lvl)
					continue OUTER
				}
				logHandlers[lvl] = n
				continue OUTER
			}
		}
	}
}

// WithDefaultFields adds fields to the underlying logger instance that will be automatically added to ALL log entries.
func WithDefaultFields(fields ...Field) {
	logFields = append(logFields, fields...)
}

// WithField returns a new log entry with the supplied field.
func WithField(key string, value interface{}) Entry {
	ne := newEntry(Field{Key: key, Value: value})
	return ne
}

// WithFields returns a new log entry with the supplied fields appended
func WithFields(fields ...Field) Entry {
	ne := newEntry(fields...)
	return ne
}

// WithTrace with add duration of how long the between this function call and
// the subsequent log
func WithTrace() Entry {
	ne := newEntry()
	ne.start = time.Now()
	return ne
}

// WithError add a minimal stack trace to the log Entry
func WithError(err error) Entry {
	ne := newEntry()
	return withErrFn(ne, err)
}

// Debug logs a debug entry
func Debug(v ...interface{}) {
	e := newEntry()
	e.Debug(v...)
}

// Debugf logs a debug entry with formatting
func Debugf(s string, v ...interface{}) {
	e := newEntry()
	e.Debugf(s, v...)
}

// Info logs a normal. information, entry
func Info(v ...interface{}) {
	e := newEntry()
	e.Info(v...)
}

// Infof logs a normal. information, entry with formatting
func Infof(s string, v ...interface{}) {
	e := newEntry()
	e.Infof(s, v...)
}

// Notice logs a notice log entry
func Notice(v ...interface{}) {
	e := newEntry()
	e.Notice(v...)
}

// Noticef logs a notice log entry with formatting
func Noticef(s string, v ...interface{}) {
	e := newEntry()
	e.Noticef(s, v...)
}

// Warn logs a warning log entry
func Warn(v ...interface{}) {
	e := newEntry()
	e.Warn(v...)
}

// Warnf logs a warning log entry with formatting
func Warnf(s string, v ...interface{}) {
	e := newEntry()
	e.Warnf(s, v...)
}

// Panic logs a panic log entry
func Panic(v ...interface{}) {
	e := newEntry()
	e.Panic(v...)
}

// Panicf logs a panic log entry with formatting
func Panicf(s string, v ...interface{}) {
	e := newEntry()
	e.Panicf(s, v...)
}

// Alert logs an alert log entry
func Alert(v ...interface{}) {
	e := newEntry()
	e.Alert(v...)
}

// Alertf logs an alert log entry with formatting
func Alertf(s string, v ...interface{}) {
	e := newEntry()
	e.Alertf(s, v...)
}

// Fatal logs a fatal log entry
func Fatal(v ...interface{}) {
	e := newEntry()
	e.Fatal(v...)
}

// Fatalf logs a fatal log entry with formatting
func Fatalf(s string, v ...interface{}) {
	e := newEntry()
	e.Fatalf(s, v...)
}

// Error logs an error log entry
func Error(v ...interface{}) {
	e := newEntry()
	e.Error(v...)
}

// Errorf logs an error log entry with formatting
func Errorf(s string, v ...interface{}) {
	e := newEntry()
	e.Errorf(s, v...)
}

// Handler is an interface that log handlers need to comply with
type Handler interface {
	Log(Entry)
}
