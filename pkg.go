package log

// F creates a new field key + value entry
func F(key string, value interface{}) Field {
	return Logger.F(key, value)
}

// Debug level formatted message.
func Debug(v ...interface{}) {
	Logger.Debug(v...)
}

// Trace starts a trace & returns Traceable object to End + log
func Trace(v ...interface{}) Traceable {
	return Logger.Trace(v...)
}

// Info level formatted message.
func Info(v ...interface{}) {
	Logger.Info(v...)
}

// Notice level formatted message.
func Notice(v ...interface{}) {
	Logger.Notice(v...)
}

// Warn level formatted message.
func Warn(v ...interface{}) {
	Logger.Warn(v...)
}

// Error level formatted message.
func Error(v ...interface{}) {
	Logger.Error(v...)
}

// Panic logs an Panic level formatted message and then panics
// it is here to let this log package be a drop in replacement
// for the standard logger
func Panic(v ...interface{}) {
	Logger.Panic(v...)
}

// Alert level formatted message.
func Alert(v ...interface{}) {
	Logger.Alert(v...)
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

// Tracef starts a trace & returns Traceable object to End + log
func Tracef(msg string, v ...interface{}) Traceable {
	return Logger.Tracef(msg, v...)
}

// Infof level formatted message.
func Infof(msg string, v ...interface{}) {
	Logger.Infof(msg, v...)
}

// Noticef level formatted message.
func Noticef(msg string, v ...interface{}) {
	Logger.Noticef(msg, v...)
}

// Warnf level formatted message.
func Warnf(msg string, v ...interface{}) {
	Logger.Warnf(msg, v...)
}

// Errorf level formatted message.
func Errorf(msg string, v ...interface{}) {
	Logger.Errorf(msg, v...)
}

// Panicf logs an Panic level formatted message and then panics
// it is here to let this log package be a near drop in replacement
// for the standard logger
func Panicf(msg string, v ...interface{}) {
	Logger.Panicf(msg, v...)
}

// Alertf level formatted message.
func Alertf(msg string, v ...interface{}) {
	Logger.Alertf(msg, v...)
}

// Fatalf level formatted message, followed by an exit.
func Fatalf(msg string, v ...interface{}) {
	Logger.Fatalf(msg, v...)
}

// Panicln logs an Panic level formatted message and then panics
// it is here to let this log package be a near drop in replacement
// for the standard logger
func Panicln(v ...interface{}) {
	Logger.Panic(v...)
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

// WithFields returns a log Entry with fields set
func WithFields(fields ...Field) LeveledLogger {
	return Logger.WithFields(fields...)
}

// HandleEntry send the logs entry out to all the registered handlers
func HandleEntry(e *Entry) {
	Logger.HandleEntry(e)
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

// SetApplicationKey tells the logger to set a constant application key
// that will be set on all log Entry objects. log does not care what it is,
// the application name, app name + hostname.... that's up to you
// it is needed by many logging platforms for separating logs by application
// and even by application server in a distributed app.
func SetApplicationID(id string) {
	Logger.SetApplicationID(id)
}
