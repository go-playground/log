package log

import (
	"bufio"
	"fmt"
	"io"
	stdlog "log"
	"os"
	"strconv"
)

const (
	space   = byte(' ')
	equals  = byte('=')
	newLine = byte('\n')
	base10  = 10
	v       = "%v"
)

// Console is an instance of the console logger
type Console struct {
	writer          io.Writer
	r               *io.PipeReader
	timestampFormat string
}

// NewDefaultLogger returns a new instance of the console logger
func NewDefaultLogger(redirectGoStdErrLogs bool) *Console {
	c := &Console{
		writer:          os.Stderr,
		timestampFormat: "2006-01-02T15:04:05.000000000Z07:00", // RFC3339Nano
	}
	if redirectGoStdErrLogs {
		ready := make(chan struct{})
		go c.handleStdLogger(ready)
		<-ready // have to wait, it was running too quickly and some messages can be lost
	}
	return c
}

// SetTimestampFormat sets Console's timestamp output format
// Default is : "2006-01-02T15:04:05.000000000Z07:00"
func (c *Console) SetTimestampFormat(format string) {
	c.timestampFormat = format
}

// SetWriter sets Console's wriiter
// Default is : os.Stderr
func (c *Console) SetWriter(w io.Writer) {
	c.writer = w
}

// this will redirect the output of
func (c *Console) handleStdLogger(ready chan<- struct{}) {
	var w *io.PipeWriter
	c.r, w = io.Pipe()
	stdlog.SetOutput(w)

	scanner := bufio.NewScanner(c.r)
	go func() {
		close(ready)
	}()

	for scanner.Scan() {
		WithField("stdlog", true).Info(scanner.Text())
	}
	_ = c.r.Close()
	_ = w.Close()
}

// Log handles the log entry
func (c *Console) Log(e Entry) {
	var lvl string
	var i int
	b := BytePool().Get()
	b = append(b, e.Timestamp.Format(c.timestampFormat)...)
	b = append(b, space)

	lvl = e.Level.String()

	for i = 0; i < 6-len(lvl); i++ {
		b = append(b, space)
	}

	b = append(b, lvl...)
	b = append(b, space)
	b = append(b, e.Message...)

	for _, f := range e.Fields {
		b = append(b, space)
		b = append(b, f.Key...)
		b = append(b, equals)

		switch t := f.Value.(type) {
		case string:
			b = append(b, t...)
		case int:
			b = strconv.AppendInt(b, int64(t), base10)
		case int8:
			b = strconv.AppendInt(b, int64(t), base10)
		case int16:
			b = strconv.AppendInt(b, int64(t), base10)
		case int32:
			b = strconv.AppendInt(b, int64(t), base10)
		case int64:
			b = strconv.AppendInt(b, t, base10)
		case uint:
			b = strconv.AppendUint(b, uint64(t), base10)
		case uint8:
			b = strconv.AppendUint(b, uint64(t), base10)
		case uint16:
			b = strconv.AppendUint(b, uint64(t), base10)
		case uint32:
			b = strconv.AppendUint(b, uint64(t), base10)
		case uint64:
			b = strconv.AppendUint(b, t, base10)
		case float32:
			b = strconv.AppendFloat(b, float64(t), 'f', -1, 32)
		case float64:
			b = strconv.AppendFloat(b, t, 'f', -1, 64)
		case bool:
			b = strconv.AppendBool(b, t)
		default:
			b = append(b, fmt.Sprintf(v, f.Value)...)
		}
	}
	b = append(b, newLine)

	_, _ = c.writer.Write(b)
	BytePool().Put(b)
}

// Close cleans up any resources and de-registers the handler with the logger
func (c *Console) Close() error {
	RemoveHandler(c)
	// reset the output back to original
	// since we reset the output prior to closing we don't have to wait
	stdlog.SetOutput(os.Stderr)
	if c.r != nil {
		_ = c.r.Close()
	}
	return nil
}

func (c *Console) closeAlreadyLocked() error {
	removeHandler(c)
	// reset the output back to original
	// since we reset the output prior to closing we don't have to wait
	stdlog.SetOutput(os.Stderr)
	if c.r != nil {
		_ = c.r.Close()
	}
	return nil
}
