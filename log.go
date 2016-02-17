package log

import (
	"fmt"
	"os"
	"sync"
	"time"
)

// HandlerChannels is an array of handler channels
type HandlerChannels []chan<- Entry

// LevelHandlerChannels is a group of Handler channels mapped by Level
type LevelHandlerChannels map[Level]HandlerChannels

type logger struct {
	fieldPool *sync.Pool
	entryPool *sync.Pool
	tracePool *sync.Pool
	channels  LevelHandlerChannels
}

// Logger is the default instance of the log package
var Logger = &logger{
	fieldPool: &sync.Pool{New: func() interface{} {
		return Field{}
	}},
	entryPool: &sync.Pool{New: func() interface{} {
		return new(Entry)
	}},
	tracePool: &sync.Pool{New: func() interface{} {
		return new(TraceEntry)
	}},
	channels: make(LevelHandlerChannels),
}

// // StdLogger interface for being able to replace already built in log instance easily
// type StdLogger interface {
// 	Fatal(v ...interface{})
// 	Fatalf(format string, v ...interface{})
// 	Fatalln(v ...interface{})
// 	Panic(v ...interface{})
// 	Panicf(format string, v ...interface{})
// 	Panicln(v ...interface{})
// 	Print(v ...interface{})
// 	Printf(format string, v ...interface{})
// 	Println(v ...interface{})
// }

// LeveledLogger interface for logging by level
type LeveledLogger interface {
	Debug(v ...interface{})
	Info(v ...interface{})
	Warn(v ...interface{})
	Error(v ...interface{})
	Panic(v ...interface{})
	Fatal(v ...interface{})
	Trace(v ...interface{}) Traceable
	Debugf(msg string, v ...interface{})
	Infof(msg string, v ...interface{})
	Warnf(msg string, v ...interface{})
	Errorf(msg string, v ...interface{})
	Panicf(msg string, v ...interface{})
	Fatalf(msg string, v ...interface{})
	Tracef(msg string, v ...interface{}) Traceable
}

// FieldLeveledLogger interface for logging by level and WithFields
type FieldLeveledLogger interface {
	LeveledLogger
	WithFields(...Field) LeveledLogger
}

var _ FieldLeveledLogger = Logger

// F creates a new field key + value entry
func F(key string, value interface{}) Field {

	fld := Logger.fieldPool.Get().(Field)
	fld.Key = key
	fld.Value = value

	return fld
}

// Debug level formatted message.
func Debug(v ...interface{}) {
	Logger.Debug(v...)
}

// Info level formatted message.
func Info(v ...interface{}) {
	Logger.Info(v...)
}

// Warn level formatted message.
func Warn(v ...interface{}) {
	Logger.Warn(v...)
}

// Error level formatted message.
func Error(v ...interface{}) {
	Logger.Error(v...)
}

// Fatal level formatted message, followed by an exit.
func Fatal(v ...interface{}) {
	Logger.Fatal(v...)
}

// Fatalln level formatted message, followed by an exit.
func Fatalln(v ...interface{}) {
	Logger.Fatal(v...)
}

// Debugf level formatted message.
func Debugf(msg string, v ...interface{}) {
	Logger.Debugf(msg, v...)
}

// Infof level formatted message.
func Infof(msg string, v ...interface{}) {
	Logger.Infof(msg, v...)
}

// Warnf level formatted message.
func Warnf(msg string, v ...interface{}) {
	Logger.Warnf(msg, v...)
}

// Errorf level formatted message.
func Errorf(msg string, v ...interface{}) {
	Logger.Errorf(msg, v...)
}

// Fatalf level formatted message, followed by an exit.
func Fatalf(msg string, v ...interface{}) {
	Logger.Fatalf(msg, v...)
}

// Panic logs an Error level formatted message and then panics
// it is here to let this log package be a drop in replacement
// for the standard logger
func Panic(v ...interface{}) {
	Logger.Panic(v...)
}

// Panicln logs an Error level formatted message and then panics
// it is here to let this log package be a near drop in replacement
// for the standard logger
func Panicln(v ...interface{}) {
	Logger.Panic(v...)
}

// Panicf logs an Error level formatted message and then panics
// it is here to let this log package be a near drop in replacement
// for the standard logger
func Panicf(msg string, v ...interface{}) {
	Logger.Panicf(msg, v...)
}

// Print logs an Info level formatted message
func Print(v ...interface{}) {
	Logger.Info(v...)
}

// Println logs an Info level formatted message
func Println(v ...interface{}) {
	Logger.Info(v...)
}

// Printf logs an Info level formatted message
func Printf(msg string, v ...interface{}) {
	Logger.Infof(msg, v...)
}

// Trace starts a trace & returns Traceable object to End + log
func Trace(v ...interface{}) Traceable {
	return Logger.Trace(v...)
}

// Tracef starts a trace & returns Traceable object to End + log
func Tracef(msg string, v ...interface{}) Traceable {
	return Logger.Tracef(msg, v...)
}

// WithFields returns a log Entry with fields set
func WithFields(fields ...Field) LeveledLogger {
	return Logger.WithFields(fields...)
}

// RegisterHandler adds a new Log Handler and specifies what log levels
// the handler will be passed log entries for
func RegisterHandler(handler Handler, levels ...Level) {
	Logger.RegisterHandler(handler, levels...)
}

// Debug level formatted message.
func (l *logger) Debug(v ...interface{}) {
	e := newEntry(DebugLevel, fmt.Sprint(v...), nil)
	l.HandleEntry(e)
}

// Info level formatted message.
func (l *logger) Info(v ...interface{}) {
	e := newEntry(InfoLevel, fmt.Sprint(v...), nil)
	l.HandleEntry(e)
}

// Warn level formatted message.
func (l *logger) Warn(v ...interface{}) {
	e := newEntry(WarnLevel, fmt.Sprint(v...), nil)
	l.HandleEntry(e)
}

// Error level formatted message.
func (l *logger) Error(v ...interface{}) {
	e := newEntry(ErrorLevel, fmt.Sprint(v...), nil)
	l.HandleEntry(e)
}

// Fatal level formatted message, followed by an exit.
func (l *logger) Fatal(v ...interface{}) {
	e := newEntry(FatalLevel, fmt.Sprint(v...), nil)
	l.HandleEntry(e)
}

// Debugf level formatted message.
func (l *logger) Debugf(msg string, v ...interface{}) {
	e := newEntry(DebugLevel, fmt.Sprintf(msg, v...), nil)
	l.HandleEntry(e)
}

// Infof level formatted message.
func (l *logger) Infof(msg string, v ...interface{}) {
	e := newEntry(InfoLevel, fmt.Sprintf(msg, v...), nil)
	l.HandleEntry(e)
}

// Warnf level formatted message.
func (l *logger) Warnf(msg string, v ...interface{}) {
	e := newEntry(WarnLevel, fmt.Sprintf(msg, v...), nil)
	l.HandleEntry(e)
}

// Errorf level formatted message.
func (l *logger) Errorf(msg string, v ...interface{}) {
	e := newEntry(ErrorLevel, fmt.Sprintf(msg, v...), nil)
	l.HandleEntry(e)
}

// Fatalf level formatted message, followed by an exit.
func (l *logger) Fatalf(msg string, v ...interface{}) {
	e := newEntry(FatalLevel, fmt.Sprintf(msg, v...), nil)
	l.HandleEntry(e)
	os.Exit(1)
}

// Panic logs an Error level formatted message and then panics
func (l *logger) Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	e := newEntry(ErrorLevel, s, nil)
	l.HandleEntry(e)
	panic(s)
}

// Panicf logs an Error level formatted message and then panics
func (l *logger) Panicf(msg string, v ...interface{}) {
	s := fmt.Sprintf(msg, v...)
	e := newEntry(ErrorLevel, s, nil)
	l.HandleEntry(e)
	panic(s)
}

// Trace starts a trace & returns Traceable object to End + log
func (l *logger) Trace(v ...interface{}) Traceable {

	t := l.tracePool.Get().(*TraceEntry)
	t.entry = newEntry(TraceLevel, fmt.Sprint(v...), make([]Field, 0))
	t.start = time.Now().UTC()

	return t
}

// Tracef starts a trace & returns Traceable object to End + log
func (l *logger) Tracef(msg string, v ...interface{}) Traceable {

	t := l.tracePool.Get().(*TraceEntry)
	t.entry = newEntry(TraceLevel, fmt.Sprintf(msg, v...), make([]Field, 0))
	t.start = time.Now().UTC()

	return t
}

// WithFields returns a log Entry with fields set
func (l *logger) WithFields(fields ...Field) LeveledLogger {
	return newEntry(InfoLevel, "", fields)
}

func (l *logger) HandleEntry(e *Entry) {

	// need to dereference as e is put back into the pool
	// and could be reused before the log has been written
	entry := *e

	channels, ok := l.channels[e.Level]
	if !ok {
		fmt.Printf("*********** WARNING no log entry for level %s/n", e.Level)
		goto END
	}

	for _, ch := range channels {
		ch <- entry
	}

END:
	// reclaim entry + fields
	for _, f := range e.Fields {
		l.fieldPool.Put(f)
	}

	l.entryPool.Put(e)
}

// RegisterHandler adds a new Log Handler and specifies what log levels
// the handler will be passed log entries for
func (l *logger) RegisterHandler(handler Handler, levels ...Level) {

	ch := handler.Run()

	for _, level := range levels {

		channels, ok := l.channels[level]
		if !ok {
			channels = make(HandlerChannels, 0)
		}

		l.channels[level] = append(channels, ch)
	}

}
