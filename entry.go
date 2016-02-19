package log

import (
	"fmt"
	"sync"
	"time"
)

const (
	keyVal = " %s=%v"
)

// Entry represents a single log entry.
type Entry struct {
	WG            *sync.WaitGroup `json:"-"`
	ApplicationID string          `json:"appId"`
	Level         Level           `json:"level"`
	Timestamp     time.Time       `json:"timestamp"`
	Message       string          `json:"message"`
	Fields        []Field         `json:"fields"`
}

func newEntry(level Level, message string, fields []Field) *Entry {

	entry := Logger.entryPool.Get().(*Entry)
	entry.Level = level
	entry.Message = message
	entry.Fields = fields
	entry.Timestamp = time.Now().UTC()

	return entry
}

var _ LeveledLogger = new(Entry)

// Debug level message.
func (e *Entry) Debug(v ...interface{}) {
	e.Level = DebugLevel
	e.Message = fmt.Sprint(v...)
	Logger.HandleEntry(e)
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
	Logger.HandleEntry(e)
}

// Notice level formatted message.
func (e *Entry) Notice(v ...interface{}) {
	e.Level = NoticeLevel
	e.Message = fmt.Sprint(v...)
	Logger.HandleEntry(e)
}

// Warn level message.
func (e *Entry) Warn(v ...interface{}) {
	e.Level = WarnLevel
	e.Message = fmt.Sprint(v...)
	Logger.HandleEntry(e)
}

// Error level message.
func (e *Entry) Error(v ...interface{}) {
	e.Level = ErrorLevel
	e.Message = fmt.Sprint(v...)
	Logger.HandleEntry(e)
}

// Panic logs an Error level formatted message and then panics
func (e *Entry) Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	e.Level = PanicLevel
	e.Message = s
	Logger.HandleEntry(e)

	for _, f := range e.Fields {
		s += fmt.Sprintf(keyVal, f.Key, f.Value)
	}

	panic(s)
}

// Alert level message.
func (e *Entry) Alert(v ...interface{}) {
	e.Level = AlertLevel
	e.Message = fmt.Sprint(v...)
	Logger.HandleEntry(e)
}

// Fatal level message, followed by an exit.
func (e *Entry) Fatal(v ...interface{}) {
	e.Level = FatalLevel
	e.Message = fmt.Sprint(v...)
	Logger.HandleEntry(e)
	exitFunc(1)
}

// Debugf level formatted message.
func (e *Entry) Debugf(msg string, v ...interface{}) {
	e.Level = DebugLevel
	e.Message = fmt.Sprintf(msg, v...)
	Logger.HandleEntry(e)
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
	Logger.HandleEntry(e)
}

// Noticef level formatted message.
func (e *Entry) Noticef(msg string, v ...interface{}) {
	e.Level = NoticeLevel
	e.Message = fmt.Sprintf(msg, v...)
	Logger.HandleEntry(e)
}

// Warnf level formatted message.
func (e *Entry) Warnf(msg string, v ...interface{}) {
	e.Level = WarnLevel
	e.Message = fmt.Sprintf(msg, v...)
	Logger.HandleEntry(e)
}

// Errorf level formatted message.
func (e *Entry) Errorf(msg string, v ...interface{}) {
	e.Level = ErrorLevel
	e.Message = fmt.Sprintf(msg, v...)
	Logger.HandleEntry(e)
}

// Panicf logs an Error level formatted message and then panics
func (e *Entry) Panicf(msg string, v ...interface{}) {
	s := fmt.Sprintf(msg, v...)
	e.Level = PanicLevel
	e.Message = s
	Logger.HandleEntry(e)

	for _, f := range e.Fields {
		s += fmt.Sprintf(keyVal, f.Key, f.Value)
	}

	panic(s)
}

// Alertf level formatted message.
func (e *Entry) Alertf(msg string, v ...interface{}) {
	e.Level = AlertLevel
	e.Message = fmt.Sprintf(msg, v...)
	Logger.HandleEntry(e)
}

// Fatalf level formatted message, followed by an exit.
func (e *Entry) Fatalf(msg string, v ...interface{}) {
	e.Level = FatalLevel
	e.Message = fmt.Sprintf(msg, v...)
	Logger.HandleEntry(e)
	exitFunc(1)
}
