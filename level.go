package log

// AllLevels is an array of all log levels, for easier registering of all levels to a handler
var AllLevels = []Level{
	DebugLevel,
	TraceLevel,
	InfoLevel,
	WarnLevel,
	ErrorLevel,
	FatalLevel,
}

// Level of the log
type Level int

// Log levels.
const (
	DebugLevel Level = iota
	TraceLevel
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "DEBUG"
	case TraceLevel:
		return "TRACE"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case FatalLevel:
		return "FATAL"
	default:
		return "Unknow Level"
	}
}
