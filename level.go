package log

import (
	"bytes"
	"strings"
)

// AllLevels is an array of all log levels, for easier registering of all levels to a handler
var AllLevels = []Level{
	DebugLevel,
	InfoLevel,
	NoticeLevel,
	WarnLevel,
	ErrorLevel,
	PanicLevel,
	AlertLevel,
	FatalLevel,
}

// Level of the log
type Level uint8

// Log levels.
const (
	DebugLevel Level = iota
	InfoLevel
	NoticeLevel
	WarnLevel
	ErrorLevel
	PanicLevel
	AlertLevel
	FatalLevel // same as syslog CRITICAL
)

func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case NoticeLevel:
		return "NOTICE"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case PanicLevel:
		return "PANIC"
	case AlertLevel:
		return "ALERT"
	case FatalLevel:
		return "FATAL"
	default:
		return "Unknown Level"
	}
}

// ParseLevel parses the provided strings log level or if not supported return 255
func ParseLevel(s string) Level {
	switch strings.ToUpper(s) {
	case "DEBUG":
		return DebugLevel
	case "INFO":
		return InfoLevel
	case "NOTICE":
		return NoticeLevel
	case "WARN":
		return WarnLevel
	case "ERROR":
		return ErrorLevel
	case "PANIC":
		return PanicLevel
	case "ALERT":
		return AlertLevel
	case "FATAL":
		return FatalLevel
	default:
		return 255
	}
}

// MarshalJSON implementation.
func (l Level) MarshalJSON() ([]byte, error) {
	return []byte(`"` + l.String() + `"`), nil
}

// UnmarshalJSON implementation.
func (l *Level) UnmarshalJSON(b []byte) error {
	*l = ParseLevel(string(bytes.Trim(b, `"`)))
	return nil
}
