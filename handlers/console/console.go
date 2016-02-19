package console

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/go-playground/log"
)

// Console is an instance of the console logger
type Console struct {
	buffer          uint
	colors          [9]log.ANSIEscSeq
	ansiReset       log.ANSIEscSeq
	writer          io.Writer
	miniTimestamp   bool
	timestampFormat string
	displayColor    bool
	start           time.Time
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
func (c *Console) SetTimestampFormat(format string) {
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
func (c *Console) Run() chan<- log.Entry {

	// in a big high traffic app, set a higher buffer
	ch := make(chan log.Entry, c.buffer)

	if c.displayColor {
		go c.handleLog(ch)
	} else {
		go c.handleLogNoColor(ch)
	}

	return ch
}

func (c *Console) parseMiniTimestamp() int {
	return int(time.Since(c.start) / time.Second)
}

// handleLog consumes and logs any Entry's passed to the channel
func (c *Console) handleLog(entries <-chan log.Entry) {

	var e log.Entry
	var color log.ANSIEscSeq
	var l int

	for e = range entries {

		l = len(e.Fields)
		color = c.colors[e.Level]

		if c.miniTimestamp {
			if l == 0 {
				fmt.Fprintf(c.writer, "%s%6s%s[%04d] %s", color, e.Level, c.ansiReset, c.parseMiniTimestamp(), e.Message)
			} else {
				fmt.Fprintf(c.writer, "%s%6s%s[%04d] %-25s", color, e.Level, c.ansiReset, c.parseMiniTimestamp(), e.Message)
			}
		} else {
			if l == 0 {
				fmt.Fprintf(c.writer, "%s%6s%s[%s] %s", color, e.Level, c.ansiReset, e.Timestamp.Format(c.timestampFormat), e.Message)
			} else {
				fmt.Fprintf(c.writer, "%s%6s%s[%s] %-25s", color, e.Level, c.ansiReset, e.Timestamp.Format(c.timestampFormat), e.Message)
			}
		}

		for _, f := range e.Fields {
			fmt.Fprintf(c.writer, " %s%s%s=%v", color, f.Key, c.ansiReset, f.Value)
		}

		fmt.Fprintln(c.writer)

		e.WG.Done()
	}
}

// handleLogNoColor consumes and logs any Entry's passed to the channel and
// print with no color
func (c *Console) handleLogNoColor(entries <-chan log.Entry) {

	var e log.Entry
	var l int

	for e = range entries {

		l = len(e.Fields)

		if c.miniTimestamp {
			if l == 0 {
				fmt.Fprintf(c.writer, "%6s[%04d] %s", e.Level, c.parseMiniTimestamp(), e.Message)
			} else {
				fmt.Fprintf(c.writer, "%6s[%04d] %-25s", e.Level, c.parseMiniTimestamp(), e.Message)
			}
		} else {
			if l == 0 {
				fmt.Fprintf(c.writer, "%6s[%s] %s", e.Level, e.Timestamp.Format(c.timestampFormat), e.Message)
			} else {
				fmt.Fprintf(c.writer, "%6s[%s] %-25s", e.Level, e.Timestamp.Format(c.timestampFormat), e.Message)
			}
		}

		for _, f := range e.Fields {
			fmt.Fprintf(c.writer, " %s=%v", f.Key, f.Value)
		}

		fmt.Fprintln(c.writer)

		e.WG.Done()
	}
}
