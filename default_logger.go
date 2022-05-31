package log

import (
	"bufio"
	"fmt"
	"io"
	stdlog "log"
	"os"
	"strconv"
	"sync"
)

const (
	space   = byte(' ')
	equals  = byte('=')
	newLine = byte('\n')
	base10  = 10
	v       = "%v"
)

// Builder is used to create a new console logger
type Builder struct {
	writer          io.Writer
	timestampFormat string
	redirect        bool
}

// NewBuilder creates a new Builder for configuring and creating a new console logger
func NewBuilder() *Builder {
	return &Builder{
		writer:          os.Stderr,
		timestampFormat: DefaultTimeFormat,
		redirect:        true,
	}
}

func (b *Builder) WithGoSTDErrLogs(redirect bool) *Builder {
	b.redirect = redirect
	return b
}

func (b *Builder) WithWriter(writer io.Writer) *Builder {
	b.writer = writer
	return b
}

func (b *Builder) WithTimestampFormat(format string) *Builder {
	b.timestampFormat = format
	return b
}

func (b *Builder) Build() *Logger {
	c := &Logger{
		writer:          b.writer,
		timestampFormat: b.timestampFormat,
	}
	if b.redirect {
		ready := make(chan struct{})
		go c.handleStdLogger(ready)
		<-ready // have to wait, it was running too quickly and some messages can be lost
	}
	return c
}

// Logger is an instance of the console logger
type Logger struct {
	m               sync.Mutex
	writer          io.Writer
	r               *io.PipeReader
	timestampFormat string
}

// this will redirect the output of
func (c *Logger) handleStdLogger(ready chan<- struct{}) {
	var w *io.PipeWriter
	c.r, w = io.Pipe()
	stdlog.SetOutput(w)

	scanner := bufio.NewScanner(c.r)
	go func() {
		close(ready)
	}()

	for scanner.Scan() {
		WithField("stdlog", true).Notice(scanner.Text())
	}
	_ = c.r.Close()
	_ = w.Close()
}

// Log handles the log entry
func (c *Logger) Log(e Entry) {
	var lvl string
	var i int
	buff := BytePool().Get()
	buff.B = append(buff.B, e.Timestamp.Format(c.timestampFormat)...)
	buff.B = append(buff.B, space)

	lvl = e.Level.String()

	for i = 0; i < 6-len(lvl); i++ {
		buff.B = append(buff.B, space)
	}

	buff.B = append(buff.B, lvl...)
	buff.B = append(buff.B, space)
	buff.B = append(buff.B, e.Message...)

	for _, f := range e.Fields {
		buff.B = append(buff.B, space)
		buff.B = append(buff.B, f.Key...)
		buff.B = append(buff.B, equals)

		switch t := f.Value.(type) {
		case string:
			buff.B = append(buff.B, t...)
		case int:
			buff.B = strconv.AppendInt(buff.B, int64(t), base10)
		case int8:
			buff.B = strconv.AppendInt(buff.B, int64(t), base10)
		case int16:
			buff.B = strconv.AppendInt(buff.B, int64(t), base10)
		case int32:
			buff.B = strconv.AppendInt(buff.B, int64(t), base10)
		case int64:
			buff.B = strconv.AppendInt(buff.B, t, base10)
		case uint:
			buff.B = strconv.AppendUint(buff.B, uint64(t), base10)
		case uint8:
			buff.B = strconv.AppendUint(buff.B, uint64(t), base10)
		case uint16:
			buff.B = strconv.AppendUint(buff.B, uint64(t), base10)
		case uint32:
			buff.B = strconv.AppendUint(buff.B, uint64(t), base10)
		case uint64:
			buff.B = strconv.AppendUint(buff.B, t, base10)
		case float32:
			buff.B = strconv.AppendFloat(buff.B, float64(t), 'f', -1, 32)
		case float64:
			buff.B = strconv.AppendFloat(buff.B, t, 'f', -1, 64)
		case bool:
			buff.B = strconv.AppendBool(buff.B, t)
		default:
			buff.B = append(buff.B, fmt.Sprintf(v, f.Value)...)
		}
	}
	buff.B = append(buff.B, newLine)

	c.m.Lock()
	_, _ = c.writer.Write(buff.B)
	c.m.Unlock()

	BytePool().Put(buff)
}

// Close cleans up any resources and de-registers the handler with the logger
func (c *Logger) Close() error {
	RemoveHandler(c)
	// reset the output back to original
	// since we reset the output prior to closing we don't have to wait
	stdlog.SetOutput(os.Stderr)
	if c.r != nil {
		_ = c.r.Close()
	}
	return nil
}

func (c *Logger) closeAlreadyLocked() error {
	removeHandler(c)
	// reset the output back to original
	// since we reset the output prior to closing we don't have to wait
	stdlog.SetOutput(os.Stderr)
	if c.r != nil {
		_ = c.r.Close()
	}
	return nil
}
