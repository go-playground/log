package log

// F creates a new field key + value entry
func F(key string, value interface{}) Field {
	return Logger.F(key, value)
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

// RegisterDurationFunc registers a custom duration function for Trace events
func RegisterDurationFunc(fn DurationFormatFunc) {
	Logger.RegisterDurationFunc(fn)
}

// SetTimeFormat sets the time format used for Trace events
func SetTimeFormat(format string) {
	Logger.SetTimeFormat(format)
}
