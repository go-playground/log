package console

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/go-playground/log"
)

const (
	colorFields           = "%s %s%6s%s %-25s"
	colorNoFields         = "%s %s%6s%s %s"
	colorKeyValue         = " %s%s%s=%v"
	colorFieldsCaller     = "%s %s%6s%s %s:%d %-25s"
	colorNoFieldsCaller   = "%s %s%6s%s %s:%d %s"
	noColorFields         = "%s %6s %-25s"
	noColorNoFields       = "%s %6s %s"
	noColorKeyValue       = " %s=%v"
	noColorFieldsCaller   = "%s %6s %s:%d %-25s"
	noColorNoFieldsCaller = "%s %6s %s:%d %s"
)

// Console is an instance of the console logger
type Console struct {
	buffer          uint
	colors          [9]log.ANSIEscSeq
	ansiReset       log.ANSIEscSeq
	writer          io.Writer
	timestampFormat string
	start           time.Time
	format          string
	formatFields    string
	formatKeyValue  string
	formatTs        func(e *log.Entry) string
	displayColor    bool
	miniTimestamp   bool
}

// Colors mapping.
var defaultColors = [...]log.ANSIEscSeq{
	log.DebugLevel:  log.Green,
	log.TraceLevel:  log.White,
	log.InfoLevel:   log.Blue,
	log.NoticeLevel: log.LightCyan,
	log.WarnLevel:   log.Yellow,
	log.ErrorLevel:  log.LightRed,
	log.PanicLevel:  log.Red,
	log.AlertLevel:  log.Red + log.Underscore,
	log.FatalLevel:  log.Red + log.Underscore + log.Blink,
}

// New returns a new instance of the console logger
func New() *Console {
	return &Console{
		buffer:          0,
		colors:          defaultColors,
		ansiReset:       log.Reset,
		writer:          os.Stderr,
		miniTimestamp:   true,
		timestampFormat: time.RFC3339Nano,
		displayColor:    true,
		start:           time.Now(),
	}
}

// DisplayColor tells Console to output in color or not
// Default is : true
func (c *Console) DisplayColor(color bool) {
	c.displayColor = color
}

// SetTimestampFormat sets Console's timestamp output format
// Default is : time.RFC3339Nano
// automatically calls UseMiniTimestamp(false)
func (c *Console) SetTimestampFormat(format string) {
	c.UseMiniTimestamp(false)
	c.timestampFormat = format
}

// UseMiniTimestamp tells Console to use the mini timestamp(for development mostly)
// Default is : true
func (c *Console) UseMiniTimestamp(mini bool) {
	c.miniTimestamp = mini
}

// SetLevelColor updates Console's level color values
func (c *Console) SetLevelColor(l log.Level, color log.ANSIEscSeq) {
	c.colors[l] = color
}

// SetANSIReset sets the ANSI Reset sequence
func (c *Console) SetANSIReset(code log.ANSIEscSeq) {
	c.ansiReset = code
}

// SetWriter sets Console's wriiter
// Default is : os.Stderr
func (c *Console) SetWriter(w io.Writer) {
	c.writer = w
}

// SetChannelBuffer tells Console what the channel buffer size should be
// Default is : 0
func (c *Console) SetChannelBuffer(i uint) {
	c.buffer = i
}

// Run starts the logger consuming on the returned channed
func (c *Console) Run() chan<- *log.Entry {

	// in a big high traffic app, set a higher buffer
	ch := make(chan *log.Entry, c.buffer)

	if c.miniTimestamp {
		c.formatTs = func(e *log.Entry) string {
			return fmt.Sprintf("%04d", int(time.Since(c.start)/time.Second))
		}
	} else {
		c.formatTs = func(e *log.Entry) string {
			return e.Timestamp.Format(c.timestampFormat)
		}
	}

	// GetCallerInfo may want to add logic to not solely rely upon this
	if log.GetCallerInfo() {
		if c.displayColor {
			c.format = colorNoFieldsCaller
			c.formatFields = colorFieldsCaller
			c.formatKeyValue = colorKeyValue

			go c.handleLogCaller(ch)
		} else {
			c.format = noColorNoFieldsCaller
			c.formatFields = noColorFieldsCaller
			c.formatKeyValue = noColorKeyValue

			go c.handleLogNoColorCaller(ch)
		}
	} else {
		if c.displayColor {
			c.format = colorNoFields
			c.formatFields = colorFields
			c.formatKeyValue = colorKeyValue

			go c.handleLog(ch)
		} else {
			c.format = noColorNoFields
			c.formatFields = noColorFields
			c.formatKeyValue = noColorKeyValue

			go c.handleLogNoColor(ch)
		}
	}

	return ch
}

// handleLog consumes and logs any Entry's passed to the channel
func (c *Console) handleLog(entries <-chan *log.Entry) {

	var e *log.Entry
	var color log.ANSIEscSeq

	for e = range entries {

		color = c.colors[e.Level]

		if len(e.Fields) == 0 {
			fmt.Fprintf(c.writer, c.format, c.formatTs(e), color, e.Level, c.ansiReset, e.Message)
		} else {
			fmt.Fprintf(c.writer, c.formatFields, c.formatTs(e), color, e.Level, c.ansiReset, e.Message)
		}

		for _, f := range e.Fields {
			fmt.Fprintf(c.writer, c.formatKeyValue, color, f.Key, c.ansiReset, f.Value)
		}

		fmt.Fprintln(c.writer)

		e.Consumed()
	}
}

// handleLog consumes and logs any Entry's passed to the channel
func (c *Console) handleLogCaller(entries <-chan *log.Entry) {

	var e *log.Entry
	var color log.ANSIEscSeq
	var file string

	for e = range entries {

		color = c.colors[e.Level]

		file = e.File
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				file = file[i+1:]
				break
			}
		}

		if len(e.Fields) == 0 {
			fmt.Fprintf(c.writer, c.format, c.formatTs(e), color, e.Level, c.ansiReset, file, e.Line, e.Message)
		} else {
			fmt.Fprintf(c.writer, c.formatFields, c.formatTs(e), color, e.Level, c.ansiReset, file, e.Line, e.Message)
		}

		for _, f := range e.Fields {
			fmt.Fprintf(c.writer, c.formatKeyValue, color, f.Key, c.ansiReset, f.Value)
		}

		fmt.Fprintln(c.writer)

		e.Consumed()
	}
}

// handleLogNoColor consumes and logs any Entry's passed to the channel, with no color
func (c *Console) handleLogNoColor(entries <-chan *log.Entry) {

	var e *log.Entry

	for e = range entries {

		if len(e.Fields) == 0 {
			fmt.Fprintf(c.writer, c.format, c.formatTs(e), e.Level, e.Message)
		} else {
			fmt.Fprintf(c.writer, c.formatFields, c.formatTs(e), e.Level, e.Message)
		}

		for _, f := range e.Fields {
			fmt.Fprintf(c.writer, c.formatKeyValue, f.Key, f.Value)
		}

		fmt.Fprintln(c.writer)

		e.Consumed()
	}
}

// handleLogNoColorCaller consumes and logs any Entry's passed to the channel,
// with no color, file and line info
func (c *Console) handleLogNoColorCaller(entries <-chan *log.Entry) {

	var e *log.Entry
	var file string

	for e = range entries {

		file = e.File
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				file = file[i+1:]
				break
			}
		}

		if len(e.Fields) == 0 {
			fmt.Fprintf(c.writer, c.format, c.formatTs(e), e.Level, file, e.Line, e.Message)
		} else {
			fmt.Fprintf(c.writer, c.formatFields, c.formatTs(e), e.Level, file, e.Line, e.Message)
		}

		for _, f := range e.Fields {
			fmt.Fprintf(c.writer, c.formatKeyValue, f.Key, f.Value)
		}

		fmt.Fprintln(c.writer)

		e.Consumed()
	}
}
