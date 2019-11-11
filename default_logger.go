package log

import (
	"bufio"
	"fmt"
	"io"
	stdlog "log"
	"os"
	"strconv"

	"github.com/go-playground/ansi/v3"
)

const (
	space   = byte(' ')
	equals  = byte('=')
	newLine = byte('\n')
	base10  = 10
	v       = "%v"
)

var (
	defaultLoggerWriter     io.Writer = os.Stderr                             // here for tests only
	defaultLoggerTimeFormat           = "2006-01-02 15:04:05.000000000Z07:00" // here for tests only
)

// console is an instance of the console logger
type console struct {
	colors          [8]ansi.EscSeq
	writer          io.Writer
	timestampFormat string
	r               *io.PipeReader
}

func newDefaultLogger() *console {
	c := &console{
		colors: [...]ansi.EscSeq{
			DebugLevel:  ansi.Green,
			InfoLevel:   ansi.Blue,
			NoticeLevel: ansi.LightCyan,
			WarnLevel:   ansi.LightYellow,
			ErrorLevel:  ansi.LightRed,
			PanicLevel:  ansi.Red,
			AlertLevel:  ansi.Red + ansi.Underline,
			FatalLevel:  ansi.Red + ansi.Underline + ansi.Blink,
		},
		writer:          defaultLoggerWriter,
		timestampFormat: defaultLoggerTimeFormat,
	}
	done := make(chan struct{})
	go c.handleStdLogger(done)
	<-done // have to wait, it was running too quickly and some messages can be lost
	return c
}

// this will redirect the output of
func (c *console) handleStdLogger(done chan<- struct{}) {
	var w *io.PipeWriter
	c.r, w = io.Pipe()
	stdlog.SetOutput(w)

	scanner := bufio.NewScanner(c.r)
	go func() {
		done <- struct{}{}
	}()

	for scanner.Scan() {
		WithField("stdlog", true).Info(scanner.Text())
	}
	_ = c.r.Close()
	_ = w.Close()
}

// Log handles the log entry
func (c *console) Log(e Entry) {
	var i int

	b := BytePool().Get()
	lvl := e.Level.String()
	color := c.colors[e.Level]
	b = append(b, e.Timestamp.Format(c.timestampFormat)...)
	b = append(b, space)
	b = append(b, color...)

	for i = 0; i < 6-len(lvl); i++ {
		b = append(b, space)
	}
	b = append(b, lvl...)
	b = append(b, ansi.Reset...)
	b = append(b, space)
	b = append(b, e.Message...)

	for _, f := range e.Fields {
		b = append(b, space)
		b = append(b, color...)
		b = append(b, f.Key...)
		b = append(b, ansi.Reset...)
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

// Close cleans up any resources
func (c *console) Close() error {
	// reset the output back to original
	stdlog.SetOutput(os.Stderr)
	// since we reset the output piror to closing we don't have to wait
	if c.r != nil {
		_ = c.r.Close()
	}
	return nil
}
