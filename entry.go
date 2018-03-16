package log

import (
	"fmt"
	"time"
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
	flds := make([]Field, 0, len(e.Fields))
	e.Fields = append(flds, e.Fields...)
	return e
}

func newEntryWithFields(fields []Field) Entry {
	e := Entry{
		Fields: make([]Field, 0, len(fields)),
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

// WithTrace with add duration of how long the between this function call and
// the subsequent log
func (e Entry) WithTrace() Entry {
	e.start = time.Now()
	return e
}

// WithError add a minimal stack trace to the log Entry
func (e Entry) WithError(err error) Entry {
	return withErrFn(e, err)
}

// Debug logs a debug entry
func (e Entry) Debug(v ...interface{}) {
	e.Message = fmt.Sprint(v...)
	e.Level = DebugLevel
	HandleEntry(e)
}

// Debugf logs a debug entry with formatting
func (e Entry) Debugf(s string, v ...interface{}) {
	e.Message = fmt.Sprintf(s, v...)
	e.Level = DebugLevel
	HandleEntry(e)
}

// Info logs a normal. information, entry
func (e Entry) Info(v ...interface{}) {
	e.Message = fmt.Sprint(v...)
	e.Level = InfoLevel
	HandleEntry(e)
}

// Infof logs a normal. information, entry with formatting
func (e Entry) Infof(s string, v ...interface{}) {
	e.Message = fmt.Sprintf(s, v...)
	e.Level = InfoLevel
	HandleEntry(e)
}

// Notice logs a notice log entry
func (e Entry) Notice(v ...interface{}) {
	e.Message = fmt.Sprint(v...)
	e.Level = NoticeLevel
	HandleEntry(e)
}

// Noticef logs a notice log entry with formatting
func (e Entry) Noticef(s string, v ...interface{}) {
	e.Message = fmt.Sprintf(s, v...)
	e.Level = NoticeLevel
	HandleEntry(e)
}

// Warn logs a warn log entry
func (e Entry) Warn(v ...interface{}) {
	e.Message = fmt.Sprint(v...)
	e.Level = WarnLevel
	HandleEntry(e)
}

// Warnf logs a warn log entry with formatting
func (e Entry) Warnf(s string, v ...interface{}) {
	e.Message = fmt.Sprintf(s, v...)
	e.Level = WarnLevel
	HandleEntry(e)
}

// Panic logs a panic log entry
func (e Entry) Panic(v ...interface{}) {
	e.Message = fmt.Sprint(v...)
	e.Level = PanicLevel
	HandleEntry(e)
	exitFunc(1)
}

// Panicf logs a panic log entry with formatting
func (e Entry) Panicf(s string, v ...interface{}) {
	e.Message = fmt.Sprintf(s, v...)
	e.Level = PanicLevel
	HandleEntry(e)
	exitFunc(1)
}

// Alert logs an alert log entry
func (e Entry) Alert(v ...interface{}) {
	e.Message = fmt.Sprint(v...)
	e.Level = AlertLevel
	HandleEntry(e)
}

// Alertf logs an alert log entry with formatting
func (e Entry) Alertf(s string, v ...interface{}) {
	e.Message = fmt.Sprintf(s, v...)
	e.Level = AlertLevel
	HandleEntry(e)
}

// Fatal logs a fatal log entry
func (e Entry) Fatal(v ...interface{}) {
	e.Message = fmt.Sprint(v...)
	e.Level = FatalLevel
	HandleEntry(e)
	exitFunc(1)
}

// Fatalf logs a fatal log entry with formatting
func (e Entry) Fatalf(s string, v ...interface{}) {
	e.Message = fmt.Sprintf(s, v...)
	e.Level = FatalLevel
	HandleEntry(e)
	exitFunc(1)
}

// Error logs an error log entry
func (e Entry) Error(v ...interface{}) {
	e.Message = fmt.Sprint(v...)
	e.Level = ErrorLevel
	HandleEntry(e)
}

// Errorf logs an error log entry with formatting
func (e Entry) Errorf(s string, v ...interface{}) {
	e.Message = fmt.Sprintf(s, v...)
	e.Level = ErrorLevel
	HandleEntry(e)
}
