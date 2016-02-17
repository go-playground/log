package console

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/go-playground/log"
)

// colors.
const (
	none   = 0
	red    = 31
	green  = 32
	yellow = 33
	blue   = 34
	cyan   = 36
	gray   = 37
)

// Console is an instance of the console logger
type Console struct {
	Consumers uint
	Colors    [6]int
	Writer    io.Writer
	mu        sync.Mutex
}

// Colors mapping.
var defaultColors = [...]int{
	log.DebugLevel: gray,
	log.TraceLevel: cyan,
	log.InfoLevel:  blue,
	log.WarnLevel:  yellow,
	log.ErrorLevel: red,
	log.FatalLevel: red,
}

// New returns a new instance of the console logger
func New() *Console {
	return &Console{
		Consumers: 1,
		Colors:    defaultColors,
		Writer:    os.Stderr,
	}
}

// Run starts the logger consuming on the returned channed
func (c *Console) Run() chan<- *log.Entry {

	ch := make(chan *log.Entry)

	// in a big high traffic app, spin up more consumers?
	for i := 0; i < int(c.Consumers); i++ {
		go c.handleLog(ch)
	}

	return ch
}

// handleLog consumes and logs any Entry's passed to the channel
func (c *Console) handleLog(entries <-chan *log.Entry) {

	var e *log.Entry
	var color int
	buff := new(bytes.Buffer)
	// var buff bytes.Buffer

	for e = range entries {

		buff.Reset()
		color = c.Colors[e.Level]

		fmt.Fprintf(buff, "\033[%dm%6s\033[0m[%s] %-25s", color, e.Level, e.Timestamp, e.Message)
		// fmt.Fprintf(c.Writer, "\033[%dm%6s\033[0m[%04d] %-25s", color, e.Level, e.Timestamp, e.Message)

		for _, f := range e.Fields {
			fmt.Fprintf(buff, " \033[%dm%s\033[0m=%v", color, f.Key, f.Value)
		}

		fmt.Fprintln(buff)

		c.mu.Lock()
		c.Writer.Write(buff.Bytes())
		// could have used defer, but overhead really slows it down
		c.mu.Unlock()
	}
}
