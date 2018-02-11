package log

import (
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Fields is the type to send to WithFields
type Fields []Field

// Entry defines a single log entry
type Entry struct {
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Fields    []Field   `json:"fields"`
	Level     Level     `json:"level"`
	start     time.Time
}

func newEntry(e Entry) Entry {
	ne := Entry{
		Fields:    make([]Field, 0, len(e.Fields)),
		Timestamp: time.Now(),
		Message:   e.Message,
		Level:     e.Level,
		start:     e.start,
	}
	ne.Fields = append(ne.Fields, e.Fields...)
	return ne
}

func newEntryWithFields(fields []Field) Entry {
	e := Entry{
		Fields:    make([]Field, 0, len(fields)),
		Timestamp: time.Now(),
	}
	e.Fields = append(e.Fields, fields...)
	return e
}

// WithField returns a new log entry with the supplied field.
func (e Entry) WithField(key string, value interface{}) Entry {
	ne := newEntry(e)
	ne.Fields = append(ne.Fields, Field{Key: key, Value: value})
	return ne
}

// WithFields returns a new log entry with the supplied fields appended
func (e Entry) WithFields(fields ...Field) Entry {
	ne := newEntry(e)
	ne.Fields = append(ne.Fields, fields...)
	return ne
}

// WithTrace withh add duration of how long the between this function call and
// the susequent log
func (e Entry) WithTrace() Entry {
	e.start = time.Now()
	return e
}

// WithError add a minimal stack trace to the log Entry
func (e Entry) WithError(err error) Entry {
	return e.withError(err)
}

// WithError add a minimal stack trace to the log Entry
func (e Entry) withError(err error) Entry {
	ne := newEntry(e)
	ne.Fields = append(ne.Fields, Field{Key: "error", Value: err.Error()})

	var frame errors.Frame

	if s, ok := err.(stackTracer); ok {
		frame = s.StackTrace()[0]
	} else {
		frame = errors.WithStack(err).(stackTracer).StackTrace()[2:][0]
	}

	name := fmt.Sprintf("%n", frame)
	file := fmt.Sprintf("%+s", frame)
	line := fmt.Sprintf("%d", frame)
	parts := strings.Split(file, "\n\t")
	if len(parts) > 1 {
		file = parts[1]
	}
	ne.Fields = append(ne.Fields, Field{Key: "source", Value: fmt.Sprintf("%s: %s:%s", name, file, line)})
	return ne
}

// Debug logs a debug entry
func (e Entry) Debug(v ...interface{}) {
	e.Message = fmt.Sprint(v...)
	e.Level = DebugLevel
	handleEntry(e)
}

// Debugf logs a debug entry with formatting
func (e Entry) Debugf(s string, v ...interface{}) {
	e.Message = fmt.Sprintf(s, v...)
	e.Level = DebugLevel
	handleEntry(e)
}

// Info logs a normal. information, entry
func (e Entry) Info(v ...interface{}) {
	e.Message = fmt.Sprint(v...)
	e.Level = InfoLevel
	handleEntry(e)
}

// Infof logs a normal. information, entry with formatting
func (e Entry) Infof(s string, v ...interface{}) {
	e.Message = fmt.Sprintf(s, v...)
	e.Level = InfoLevel
	handleEntry(e)
}

// Notice logs a notice log entry
func (e Entry) Notice(v ...interface{}) {
	e.Message = fmt.Sprint(v...)
	e.Level = NoticeLevel
	handleEntry(e)
}

// Noticef logs a notice log entry with formatting
func (e Entry) Noticef(s string, v ...interface{}) {
	e.Message = fmt.Sprintf(s, v...)
	e.Level = NoticeLevel
	handleEntry(e)
}

// Warn logs a warn log entry
func (e Entry) Warn(v ...interface{}) {
	e.Message = fmt.Sprint(v...)
	e.Level = WarnLevel
	handleEntry(e)
}

// Warnf logs a warn log entry with formatting
func (e Entry) Warnf(s string, v ...interface{}) {
	e.Message = fmt.Sprintf(s, v...)
	e.Level = WarnLevel
	handleEntry(e)
}

// Panic logs a panic log entry
func (e Entry) Panic(v ...interface{}) {
	e.Message = fmt.Sprint(v...)
	e.Level = PanicLevel
	handleEntry(e)
	exitFunc(1)
}

// Panicf logs a panic log entry with formatting
func (e Entry) Panicf(s string, v ...interface{}) {
	e.Message = fmt.Sprintf(s, v...)
	e.Level = PanicLevel
	handleEntry(e)
	exitFunc(1)
}

// Alert logs an alert log entry
func (e Entry) Alert(v ...interface{}) {
	e.Message = fmt.Sprint(v...)
	e.Level = AlertLevel
	handleEntry(e)
}

// Alertf logs an alert log entry with formatting
func (e Entry) Alertf(s string, v ...interface{}) {
	e.Message = fmt.Sprintf(s, v...)
	e.Level = AlertLevel
	handleEntry(e)
}

// Fatal logs a fatal log entry
func (e Entry) Fatal(v ...interface{}) {
	e.Message = fmt.Sprint(v...)
	e.Level = FatalLevel
	handleEntry(e)
	exitFunc(1)
}

// Fatalf logs a fatal log entry with formatting
func (e Entry) Fatalf(s string, v ...interface{}) {
	e.Message = fmt.Sprintf(s, v...)
	e.Level = FatalLevel
	handleEntry(e)
	exitFunc(1)
}

// Error logs an error log entry
func (e Entry) Error(v ...interface{}) {
	e.Message = fmt.Sprint(v...)
	e.Level = ErrorLevel
	handleEntry(e)
}

// Errorf logs an error log entry with formatting
func (e Entry) Errorf(s string, v ...interface{}) {
	e.Message = fmt.Sprintf(s, v...)
	e.Level = ErrorLevel
	handleEntry(e)
}
