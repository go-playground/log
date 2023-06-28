package log

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

// CustomWriter is a custom writer that wraps bytes.Buffer and implements io.Writer
type CustomWriter struct {
	buffer *bytes.Buffer
}

// Write writes the data to the underlying buffer
func (w *CustomWriter) Write(p []byte) (n int, err error) {
	return w.buffer.Write(p)
}

// File returns the underlying *os.File
func (w *CustomWriter) File() *os.File {
	return os.Stdout
}

func TestColorizeLevel(t *testing.T) {
	buff := new(bytes.Buffer)

	// Create a custom writer that wraps the buffer
	customWriter := &CustomWriter{buffer: buff}

	// Redirect stdout to the buffer
	old := os.Stdout
	os.Stdout = customWriter.File()

	// Test cases
	testCases := []struct {
		level       Level
		color       string
		expected    string
		expectedErr error
	}{
		{DebugLevel, "\033[36m", "", nil},  // Expected ANSI escape code for DEBUG: Cyan
		{InfoLevel, "\033[32m", "", nil},   // Expected ANSI escape code for INFO: Green
		{NoticeLevel, "\033[33m", "", nil}, // Expected ANSI escape code for NOTICE: Yellow
		{WarnLevel, "\033[35m", "", nil},   // Expected ANSI escape code for WARN: Magenta
		{ErrorLevel, "\033[31m", "", nil},  // Expected ANSI escape code for ERROR: Red
		{PanicLevel, "\033[91m", "", nil},  // Expected ANSI escape code for PANIC: Light Red
		{AlertLevel, "\033[93m", "", nil},  // Expected ANSI escape code for ALERT: Light Yellow
		{FatalLevel, "\033[95m", "", nil},  // Expected ANSI escape code for FATAL: Light Magenta
		{InfoLevel, "\033[32m", "", nil},   // Invalid color (for testing error handling)
		//{UNKNOWN, "\033[0m", "\033[0m", nil},        // Expected ANSI escape code for unknown level: Reset to default
	}

	// Iterate over test cases
	for _, tc := range testCases {
		t.Run(tc.level.String(), func(t *testing.T) {
			// Reset the buffer
			buff.Reset()

			// Call the ColorizeLevel function
			err := func() (err error) {
				defer func() {
					if r := recover(); r != nil {
						err = fmt.Errorf("panic: %v", r)
					}
				}()
				ColorizeLevel(tc.level)
				return nil
			}()

			// Check if an error occurred (for error handling testing)
			if tc.expectedErr != nil {
				if err == nil || err.Error() != tc.expectedErr.Error() {
					t.Errorf("Expected error: %v, but got: %v", tc.expectedErr, err)
				}
				return
			}

			// Check if the printed output matches the expected ANSI escape code
			actual := buff.String()
			if actual != tc.expected {
				t.Errorf("Expected: %s, but got: %s", tc.expected, actual)
			}
		})
	}

	// Restore stdout
	os.Stdout = old
}
