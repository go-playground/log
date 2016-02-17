package console

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/go-playground/log"
)

// colors.
const (
	none     = 0
	red      = 31
	green    = 32
	yellow   = 33
	blue     = 34
	darkGray = 36
	gray     = 37
)

// Console is an instance of the console logger
type Console struct {
	buffer          uint
	colors          [6]int
	writer          io.Writer
	miniTimestamp   bool
	timestampFormat string
	displayColor    bool
	start           time.Time
	mu              sync.Mutex
}

// Colors mapping.
var defaultColors = [...]int{
	log.DebugLevel: green,
	log.TraceLevel: darkGray,
	log.InfoLevel:  blue,
	log.WarnLevel:  yellow,
	log.ErrorLevel: red,
	log.FatalLevel: red,
}

// New returns a new instance of the console logger
func New() *Console {
	return &Console{
		buffer:          0,
		colors:          defaultColors,
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
func (c *Console) SetLevelColor(l log.Level, color int) {
	c.colors[l] = color
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
	var color int

	for e = range entries {
		c.mu.Lock()

		color = c.colors[e.Level]

		if c.miniTimestamp {
			fmt.Fprintf(c.writer, "\033[%dm%6s\033[0m[%04d] %-25s", color, e.Level, c.parseMiniTimestamp(), e.Message)
		} else {
			fmt.Fprintf(c.writer, "\033[%dm%6s\033[0m[%s] %-25s", color, e.Level, e.Timestamp.Format(c.timestampFormat), e.Message)
		}

		for _, f := range e.Fields {
			fmt.Fprintf(c.writer, " \033[%dm%s\033[0m=%v", color, f.Key, f.Value)
		}

		fmt.Fprintln(c.writer)

		c.mu.Unlock()
	}
}

// handleLogNoColor consumes and logs any Entry's passed to the channel and
// print with no color
func (c *Console) handleLogNoColor(entries <-chan log.Entry) {

	var e log.Entry

	for e = range entries {
		c.mu.Lock()

		if c.miniTimestamp {
			fmt.Fprintf(c.writer, "%6s\033[%04d] %-25s", e.Level, c.parseMiniTimestamp(), e.Message)
		} else {
			fmt.Fprintf(c.writer, "%6s\033[%s] %-25s", e.Level, e.Timestamp.Format(c.timestampFormat), e.Message)
		}

		for _, f := range e.Fields {
			fmt.Fprintf(c.writer, " %s=%v", f.Key, f.Value)
		}

		fmt.Fprintln(c.writer)

		c.mu.Unlock()
	}
}
