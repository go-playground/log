package log

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	keyVal = " %s=%v"
	cutset = "\r\n\t "
)

// Entry represents a single log entry.
type Entry struct {
	wg            *sync.WaitGroup
	calldepth     int
	ApplicationID string    `json:"appId"`
	Level         Level     `json:"level"`
	Timestamp     time.Time `json:"timestamp"`
	Message       string    `json:"message"`
	Fields        []Field   `json:"fields"`
	File          string    `json:"file"`
	Line          int       `json:"line"`
}

func newEntry(level Level, message string, fields []Field, calldepth int) *Entry {

	entry := Logger.entryPool.Get().(*Entry)
	entry.calldepth = calldepth
	entry.Level = level
	entry.Message = strings.TrimRight(message, cutset) // need to trim for adding fields later in handlers + why send uneeded whitespace
	entry.Fields = fields
	entry.Timestamp = time.Now().UTC()

	if Logger.logCallerInfo && level != TraceLevel {
		_, entry.File, entry.Line, _ = runtime.Caller(entry.calldepth)
	}

	return entry
}

var _ LeveledLogger = new(Entry)

// Debug level message.
func (e *Entry) Debug(v ...interface{}) {
	e.Level = DebugLevel
	e.Message = fmt.Sprint(v...)
	Logger.handleEntry(e)
}

// Trace starts a trace & returns Traceable object to End + log
func (e *Entry) Trace(v ...interface{}) Traceable {

	e.Level = TraceLevel
	e.Message = fmt.Sprint(v...)

	t := Logger.tracePool.Get().(*TraceEntry)
	t.entry = e
	t.start = time.Now().UTC()

	return t
}

// Info level message.
func (e *Entry) Info(v ...interface{}) {
	e.Level = InfoLevel
	e.Message = fmt.Sprint(v...)
	Logger.handleEntry(e)
}

// Notice level formatted message.
func (e *Entry) Notice(v ...interface{}) {
	e.Level = NoticeLevel
	e.Message = fmt.Sprint(v...)
	Logger.handleEntry(e)
}

// Warn level message.
func (e *Entry) Warn(v ...interface{}) {
	e.Level = WarnLevel
	e.Message = fmt.Sprint(v...)
	Logger.handleEntry(e)
}

// Error level message.
func (e *Entry) Error(v ...interface{}) {
	e.Level = ErrorLevel
	e.Message = fmt.Sprint(v...)
	Logger.handleEntry(e)
}

// Panic logs an Error level formatted message and then panics
func (e *Entry) Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	e.Level = PanicLevel
	e.Message = s
	Logger.handleEntry(e)

	for _, f := range e.Fields {
		s += fmt.Sprintf(keyVal, f.Key, f.Value)
	}

	panic(s)
}

// Alert level message.
func (e *Entry) Alert(v ...interface{}) {
	e.Level = AlertLevel
	e.Message = fmt.Sprint(v...)
	Logger.handleEntry(e)
}

// Fatal level message, followed by an exit.
func (e *Entry) Fatal(v ...interface{}) {
	e.Level = FatalLevel
	e.Message = fmt.Sprint(v...)
	Logger.handleEntry(e)
	exitFunc(1)
}

// Debugf level formatted message.
func (e *Entry) Debugf(msg string, v ...interface{}) {
	e.Level = DebugLevel
	e.Message = fmt.Sprintf(msg, v...)
	Logger.handleEntry(e)
}

// Tracef starts a trace & returns Traceable object to End + log
func (e *Entry) Tracef(msg string, v ...interface{}) Traceable {

	e.Level = TraceLevel
	e.Message = fmt.Sprintf(msg, v...)

	t := Logger.tracePool.Get().(*TraceEntry)
	t.entry = e
	t.start = time.Now().UTC()

	return t
}

// Infof level formatted message.
func (e *Entry) Infof(msg string, v ...interface{}) {
	e.Level = InfoLevel
	e.Message = fmt.Sprintf(msg, v...)
	Logger.handleEntry(e)
}

// Noticef level formatted message.
func (e *Entry) Noticef(msg string, v ...interface{}) {
	e.Level = NoticeLevel
	e.Message = fmt.Sprintf(msg, v...)
	Logger.handleEntry(e)
}

// Warnf level formatted message.
func (e *Entry) Warnf(msg string, v ...interface{}) {
	e.Level = WarnLevel
	e.Message = fmt.Sprintf(msg, v...)
	Logger.handleEntry(e)
}

// Errorf level formatted message.
func (e *Entry) Errorf(msg string, v ...interface{}) {
	e.Level = ErrorLevel
	e.Message = fmt.Sprintf(msg, v...)
	Logger.handleEntry(e)
}

// Panicf logs an Error level formatted message and then panics
func (e *Entry) Panicf(msg string, v ...interface{}) {
	s := fmt.Sprintf(msg, v...)
	e.Level = PanicLevel
	e.Message = s
	Logger.handleEntry(e)

	for _, f := range e.Fields {
		s += fmt.Sprintf(keyVal, f.Key, f.Value)
	}

	panic(s)
}

// Alertf level formatted message.
func (e *Entry) Alertf(msg string, v ...interface{}) {
	e.Level = AlertLevel
	e.Message = fmt.Sprintf(msg, v...)
	Logger.handleEntry(e)
}

// Fatalf level formatted message, followed by an exit.
func (e *Entry) Fatalf(msg string, v ...interface{}) {
	e.Level = FatalLevel
	e.Message = fmt.Sprintf(msg, v...)
	Logger.handleEntry(e)
	exitFunc(1)
}

// Consumed lets the Entry and subsequently the Logger
// instance know that it has been used by a handler
func (e *Entry) Consumed() {
	if e.wg != nil {
		e.wg.Done()
	}
}
