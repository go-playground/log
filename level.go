package log

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
