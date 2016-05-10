package console

import (
	"fmt"
	"io"
	stdlog "log"
	"os"

	"github.com/go-playground/log"
)

// FormatFunc is the function that the workers use to create
// a new Formatter per worker allowing reusable go routine safe
// variable to be used within your Formatter function.
type FormatFunc func() Formatter

// Formatter is the function used to format the Redis entry
type Formatter func(e *log.Entry) []byte

const (
	defaultTS             = "2006-01-02T15:04:05.000000000Z07:00"
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
	numWorkers      uint
	colors          [9]log.ANSIEscSeq
	writer          io.Writer
	formatFunc      FormatFunc
	timestampFormat string
	displayColor    bool
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
	c := &Console{
		buffer:          0,
		numWorkers:      1,
		colors:          defaultColors,
		writer:          os.Stderr,
		timestampFormat: defaultTS,
		displayColor:    true,
	}

	c.formatFunc = c.defaultFormatFunc

	return c
}

// DisplayColor tells Console to output in color or not
// Default is : true
func (c *Console) DisplayColor(color bool) {
	c.displayColor = color
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

// SetBuffersAndWorkers sets the channels buffer size and number of concurrent workers.
// These settings should be thought about together, hence setting both in the same function.
func (c *Console) SetBuffersAndWorkers(size uint, workers uint) {
	c.buffer = size

	if workers == 0 {
		// just in case no log registered yet
		stdlog.Println("Invalid number of workers specified, setting to 1")
		log.Warn("Invalid number of workers specified, setting to 1")

		workers = 1
	}

	c.numWorkers = workers
}

// SetFormatFunc sets FormatFunc each worker will call to get
// a Formatter func
func (c *Console) SetFormatFunc(fn FormatFunc) {
	c.formatFunc = fn
}

// Run starts the logger consuming on the returned channed
func (c *Console) Run() chan<- *log.Entry {

	// in a big high traffic app, set a higher buffer
	ch := make(chan *log.Entry, c.buffer)

	for i := 0; i <= int(c.numWorkers); i++ {
		go c.handleLog(ch)
	}

	return ch
}

// handleLog consumes and logs any Entry's passed to the channel
func (c *Console) handleLog(entries <-chan *log.Entry) {

	var e *log.Entry
	var b []byte
	formatter := c.formatFunc()

	for e = range entries {

		b = formatter(e)

		fmt.Fprintln(c.writer, string(b))

		e.Consumed()
	}
}

func (c *Console) defaultFormatFunc() Formatter {

	var b []byte
	var file string

	if c.displayColor {

		var color log.ANSIEscSeq

		return func(e *log.Entry) []byte {
			b = b[0:0]
			color = c.colors[e.Level]

			if e.Line == 0 {

				if len(e.Fields) == 0 {
					b = append(b, fmt.Sprintf(colorNoFields, e.Timestamp.Format(c.timestampFormat), color, e.Level, log.Reset, e.Message)...)
				} else {
					b = append(b, fmt.Sprintf(colorFields, e.Timestamp.Format(c.timestampFormat), color, e.Level, log.Reset, e.Message)...)
				}

			} else {
				file = e.File
				for i := len(file) - 1; i > 0; i-- {
					if file[i] == '/' {
						file = file[i+1:]
						break
					}
				}

				if len(e.Fields) == 0 {
					b = append(b, fmt.Sprintf(colorNoFieldsCaller, e.Timestamp.Format(c.timestampFormat), color, e.Level, log.Reset, file, e.Line, e.Message)...)
				} else {
					b = append(b, fmt.Sprintf(colorFieldsCaller, e.Timestamp.Format(c.timestampFormat), color, e.Level, log.Reset, file, e.Line, e.Message)...)
				}
			}

			for _, f := range e.Fields {
				b = append(b, fmt.Sprintf(colorKeyValue, color, f.Key, log.Reset, f.Value)...)
			}

			return b
		}
	}

	return func(e *log.Entry) []byte {
		b = b[0:0]

		if e.Line == 0 {

			if len(e.Fields) == 0 {
				b = append(b, fmt.Sprintf(noColorNoFields, e.Timestamp.Format(c.timestampFormat), e.Level, e.Message)...)
			} else {
				b = append(b, fmt.Sprintf(noColorFields, e.Timestamp.Format(c.timestampFormat), e.Level, e.Message)...)
			}

		} else {
			file = e.File
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					file = file[i+1:]
					break
				}
			}

			if len(e.Fields) == 0 {
				b = append(b, fmt.Sprintf(noColorNoFieldsCaller, e.Timestamp.Format(c.timestampFormat), e.Level, file, e.Line, e.Message)...)
			} else {
				b = append(b, fmt.Sprintf(noColorFieldsCaller, e.Timestamp.Format(c.timestampFormat), e.Level, file, e.Line, e.Message)...)
			}
		}

		for _, f := range e.Fields {
			b = append(b, fmt.Sprintf(noColorKeyValue, f.Key, f.Value)...)
		}

		return b
	}
}
