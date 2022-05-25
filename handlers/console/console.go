package console

import (
	"bufio"
	"fmt"
	"io"
	stdlog "log"
	"os"
	"strconv"

	"github.com/go-playground/ansi/v3"
	"github.com/go-playground/log/v8"
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
	colors          [8]ansi.EscSeq
	writer          io.Writer
	r               *io.PipeReader
	timestampFormat string
	displayColor    bool
	redirectSTDOut  bool
	close           chan struct{}
}

// Colors mapping.
var defaultColors = [...]ansi.EscSeq{
	log.DebugLevel:  ansi.Green,
	log.InfoLevel:   ansi.Blue,
	log.NoticeLevel: ansi.LightCyan,
	log.WarnLevel:   ansi.LightYellow,
	log.ErrorLevel:  ansi.LightRed,
	log.PanicLevel:  ansi.Red,
	log.AlertLevel:  ansi.Red + ansi.Underline,
	log.FatalLevel:  ansi.Red + ansi.Underline + ansi.Blink,
}

// New returns a new instance of the console logger
func New(redirectSTDOut bool) *Console {
	c := &Console{
		colors:          defaultColors,
		writer:          os.Stderr,
		timestampFormat: "2006-01-02 15:04:05.000000000Z07:00",
		displayColor:    true,
		redirectSTDOut:  redirectSTDOut,
		close:           make(chan struct{}),
	}
	if redirectSTDOut {
		done := make(chan struct{})
		go c.handleStdLogger(done)
		<-done // have to wait, it was running too quickly and some messages can be lost
	}
	return c
}

// SetDisplayColor tells Console to output in color or not
// Default is : true
func (c *Console) SetDisplayColor(b bool) {
	c.displayColor = b
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
func (c *Console) handleStdLogger(done chan<- struct{}) {
	var w *io.PipeWriter
	c.r, w = io.Pipe()
	stdlog.SetOutput(w)

	scanner := bufio.NewScanner(c.r)
	go func() {
		done <- struct{}{}
	}()

	for scanner.Scan() {
		log.WithField("stdlog", true).Info(scanner.Text())
	}
	_ = c.r.Close()
	_ = w.Close()
}

// Log handles the log entry
func (c *Console) Log(e log.Entry) {
	var lvl string
	var i int
	b := log.BytePool().Get()
	if c.displayColor {
		color := c.colors[e.Level]

		b = append(b, e.Timestamp.Format(c.timestampFormat)...)
		b = append(b, space)
		b = append(b, color...)

		lvl = e.Level.String()

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
	} else {
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
	}
	_, _ = c.writer.Write(b)
	log.BytePool().Put(b)
}

// Close cleans up any resources and de-registers the handler with the logger
func (c *Console) Close() error {
	log.RemoveHandler(c)
	// reset the output back to original
	// since we reset the output piror to closing we don't have to wait
	stdlog.SetOutput(os.Stderr)
	if c.r != nil {
		_ = c.r.Close()
	}
	return nil
}
