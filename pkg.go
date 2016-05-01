package log

import (
	"fmt"
	"time"
)

// F creates a new field key + value entry
func F(key string, value interface{}) Field {
	return Logger.F(key, value)
}

// Debug level formatted message.
func Debug(v ...interface{}) {
	e := newEntry(DebugLevel, fmt.Sprint(v...), nil, 2)
	Logger.handleEntry(e)
}

// Trace starts a trace & returns Traceable object to End + log.
// Example defer log.Trace(...).End()
func Trace(v ...interface{}) Traceable {
	t := Logger.tracePool.Get().(*TraceEntry)
	t.entry = newEntry(TraceLevel, fmt.Sprint(v...), make([]Field, 0), 1)
	t.start = time.Now().UTC()

	return t
}

// Info level formatted message.
func Info(v ...interface{}) {
	e := newEntry(InfoLevel, fmt.Sprint(v...), nil, 2)
	Logger.handleEntry(e)
}

// Notice level formatted message.
func Notice(v ...interface{}) {
	e := newEntry(NoticeLevel, fmt.Sprint(v...), nil, 2)
	Logger.handleEntry(e)
}

// Warn level formatted message.
func Warn(v ...interface{}) {
	e := newEntry(WarnLevel, fmt.Sprint(v...), nil, 2)
	Logger.handleEntry(e)
}

// Error level formatted message.
func Error(v ...interface{}) {
	e := newEntry(ErrorLevel, fmt.Sprint(v...), nil, 2)
	Logger.handleEntry(e)
}

// Panic logs an Panic level formatted message and then panics
// it is here to let this log package be a drop in replacement
// for the standard logger
func Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	e := newEntry(PanicLevel, s, nil, 2)
	Logger.handleEntry(e)
	panic(s)
}

// Alert level formatted message.
func Alert(v ...interface{}) {
	s := fmt.Sprint(v...)
	e := newEntry(AlertLevel, s, nil, 2)
	Logger.handleEntry(e)
}

// Fatal level formatted message, followed by an exit.
func Fatal(v ...interface{}) {
	e := newEntry(FatalLevel, fmt.Sprint(v...), nil, 2)
	Logger.handleEntry(e)
	exitFunc(1)
}

// Fatalln level formatted message, followed by an exit.
func Fatalln(v ...interface{}) {
	e := newEntry(FatalLevel, fmt.Sprint(v...), nil, 2)
	Logger.handleEntry(e)
	exitFunc(1)
}

// Debugf level formatted message.
func Debugf(msg string, v ...interface{}) {
	e := newEntry(DebugLevel, fmt.Sprintf(msg, v...), nil, 2)
	Logger.handleEntry(e)
}

// Tracef starts a trace & returns Traceable object to End + log
func Tracef(msg string, v ...interface{}) Traceable {
	t := Logger.tracePool.Get().(*TraceEntry)
	t.entry = newEntry(TraceLevel, fmt.Sprintf(msg, v...), make([]Field, 0), 1)
	t.start = time.Now().UTC()

	return t
}

// Infof level formatted message.
func Infof(msg string, v ...interface{}) {
	e := newEntry(InfoLevel, fmt.Sprintf(msg, v...), nil, 2)
	Logger.handleEntry(e)
}

// Noticef level formatted message.
func Noticef(msg string, v ...interface{}) {
	e := newEntry(NoticeLevel, fmt.Sprintf(msg, v...), nil, 2)
	Logger.handleEntry(e)
}

// Warnf level formatted message.
func Warnf(msg string, v ...interface{}) {
	e := newEntry(WarnLevel, fmt.Sprintf(msg, v...), nil, 2)
	Logger.handleEntry(e)
}

// Errorf level formatted message.
func Errorf(msg string, v ...interface{}) {
	e := newEntry(ErrorLevel, fmt.Sprintf(msg, v...), nil, 2)
	Logger.HandleEntry(e)
}

// Panicf logs an Panic level formatted message and then panics
// it is here to let this log package be a near drop in replacement
// for the standard logger
func Panicf(msg string, v ...interface{}) {
	s := fmt.Sprintf(msg, v...)
	e := newEntry(PanicLevel, s, nil, 2)
	Logger.handleEntry(e)
	panic(s)
}

// Alertf level formatted message.
func Alertf(msg string, v ...interface{}) {
	s := fmt.Sprintf(msg, v...)
	e := newEntry(AlertLevel, s, nil, 2)
	Logger.handleEntry(e)
}

// Fatalf level formatted message, followed by an exit.
func Fatalf(msg string, v ...interface{}) {
	e := newEntry(FatalLevel, fmt.Sprintf(msg, v...), nil, 2)
	Logger.handleEntry(e)
	exitFunc(1)
}

// Panicln logs an Panic level formatted message and then panics
// it is here to let this log package be a near drop in replacement
// for the standard logger
func Panicln(v ...interface{}) {
	s := fmt.Sprint(v...)
	e := newEntry(PanicLevel, s, nil, 2)
	Logger.handleEntry(e)
	panic(s)
}

// Print logs an Info level formatted message
func Print(v ...interface{}) {
	e := newEntry(InfoLevel, fmt.Sprint(v...), nil, 2)
	Logger.handleEntry(e)
}

// Println logs an Info level formatted message
func Println(v ...interface{}) {
	e := newEntry(InfoLevel, fmt.Sprint(v...), nil, 2)
	Logger.handleEntry(e)
}

// Printf logs an Info level formatted message
func Printf(msg string, v ...interface{}) {
	e := newEntry(InfoLevel, fmt.Sprintf(msg, v...), nil, 2)
	Logger.handleEntry(e)
}

// WithFields returns a log Entry with fields set
func WithFields(fields ...Field) LeveledLogger {
	return newEntry(InfoLevel, "", fields, 2)
}

// HandleEntry send the logs entry out to all the registered handlers
func HandleEntry(e *Entry) {
	Logger.handleEntry(e)
}

// RegisterHandler adds a new Log Handler and specifies what log levels
// the handler will be passed log entries for
func RegisterHandler(handler Handler, levels ...Level) {
	Logger.RegisterHandler(handler, levels...)
}

// RegisterDurationFunc registers a custom duration function for Trace events
func RegisterDurationFunc(fn DurationFormatFunc) {
	Logger.RegisterDurationFunc(fn)
}

// SetTimeFormat sets the time format used for Trace events
func SetTimeFormat(format string) {
	Logger.SetTimeFormat(format)
}

// SetCallerInfo tells the logger to gather and set file and line number
// information on the Entry object.
func SetCallerInfo(info bool) {
	Logger.SetCallerInfo(info)
}

// GetCallerInfo returns if the Logger instance is gathering file and
// line number information
func GetCallerInfo() bool {
	return Logger.GetCallerInfo()
}

// SetApplicationID tells the logger to set a constant application key
// that will be set on all log Entry objects. log does not care what it is,
// the application name, app name + hostname.... that's up to you
// it is needed by many logging platforms for separating logs by application
// and even by application server in a distributed app.
func SetApplicationID(id string) {
	Logger.SetApplicationID(id)
}
