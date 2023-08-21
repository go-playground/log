package log

import "fmt"

// ColorizeLevel assigns a specific color to each log level based on its string representation.
func ColorizeLevel(level Level) {
	switch level.String() {
	case "DEBUG":
		PrintColor(DebugColor{})
	case "INFO":
		PrintColor(InfoColor{})
	case "NOTICE":
		PrintColor(NoticeColor{})
	case "WARN":
		PrintColor(WarnColor{})
	case "ERROR":
		PrintColor(ErrorColor{})
	case "PANIC":
		PrintColor(PanicColor{})
	case "ALERT":
		PrintColor(AlertColor{})
	case "FATAL":
		PrintColor(FatalColor{})
	default:
		PrintColor(DefaultColor{})
	}
}

// Custom color types for each log level.
type DefaultColor struct{ value string }
type DebugColor struct{ value string }
type InfoColor struct{ value string }
type NoticeColor struct{ value string }
type WarnColor struct{ value string }
type ErrorColor struct{ value string }
type PanicColor struct{ value string }
type AlertColor struct{ value string }
type FatalColor struct{ value string }

// Methods for each color type that returns the corresponding ANSI escape code.
// These methods implement the Color interface.

// Set Debug Color
func (LevelColor DebugColor) getLevelColor() string {
	LevelColor.value = "\033[36m"
	return LevelColor.value
}

// Set Info Color
func (LevelColor InfoColor) getLevelColor() string {
	LevelColor.value = "\033[32m"
	return LevelColor.value
}

// Set Notice Color
func (LevelColor NoticeColor) getLevelColor() string {
	LevelColor.value = "\033[33m"
	return LevelColor.value
}

// Set Warn Color
func (LevelColor WarnColor) getLevelColor() string {
	LevelColor.value = "\033[35m"
	return LevelColor.value
}

// Set Error Color
func (LevelColor ErrorColor) getLevelColor() string {
	LevelColor.value = "\033[31m"
	return LevelColor.value
}

// Set Panic Color
func (LevelColor PanicColor) getLevelColor() string {
	LevelColor.value = "\033[91m"
	return LevelColor.value
}

// Set Alert Color
func (LevelColor AlertColor) getLevelColor() string {
	LevelColor.value = "\033[93m"
	return LevelColor.value
}

// Set Fatal Color
func (LevelColor FatalColor) getLevelColor() string {
	LevelColor.value = "\033[95m"
	return LevelColor.value
}

// Set Default Color
func (LevelColor DefaultColor) getLevelColor() string {
	LevelColor.value = "\033[0m"
	return LevelColor.value
}

// Color interface defines the getLevelColor method to get the ANSI escape code.
type Color interface {
	getLevelColor() string
}

// PrintColor sets the appropriate color for the log level and prints it.
func PrintColor(C Color) {
	LevelColor := C.getLevelColor()
	fmt.Print(LevelColor)
}
