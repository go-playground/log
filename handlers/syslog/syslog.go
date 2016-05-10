package syslog

import (
	"fmt"
	stdlog "log"
	"log/syslog"

	"github.com/go-playground/log"
)

// FormatFunc is the function that the workers use to create
// a new Formatter per worker allowing reusable go routine safe
// variable to be used within your Formatter function.
type FormatFunc func() Formatter

// Formatter is the function used to format the Redis entry
type Formatter func(e *log.Entry) []byte

const (
	defaultTS         = "2006-01-02T15:04:05.000000000Z07:00"
	colorFormat       = "%s %s%6s%s %s"
	colorFormatCaller = "%s %s%6s%s %s:%d %s"
	colorKeyValue     = " %s%s%s=%v"
	format            = "%s %6s %s"
	formatCaller      = "%s %6s %s:%d %s"
	noColorKeyValue   = " %s=%v"
)

// Syslog is an instance of the syslog logger
type Syslog struct {
	buffer          uint
	numWorkers      uint
	colors          [9]log.ANSIEscSeq
	writer          *syslog.Writer
	formatFunc      FormatFunc
	timestampFormat string
	displayColor    bool
}

var (
	// Colors mapping.
	defaultColors = [...]log.ANSIEscSeq{
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
)

// New returns a new instance of the syslog logger
// example: syslog.New("udp", "localhost:514", syslog.LOG_DEBUG, "")
func New(network string, raddr string, priority syslog.Priority, tag string) (*Syslog, error) {

	var err error

	s := &Syslog{
		buffer:          0,
		numWorkers:      1,
		colors:          defaultColors,
		displayColor:    false,
		timestampFormat: defaultTS,
	}

	s.formatFunc = s.defaultFormatFunc

	if s.writer, err = syslog.Dial(network, raddr, priority, tag); err != nil {
		return nil, err
	}

	return s, nil
}

// DisplayColor tells Console to output in color or not
// Default is : true
func (s *Syslog) DisplayColor(color bool) {
	s.displayColor = color
}

// SetTimestampFormat sets Console's timestamp output format
// Default is : 2006-01-02T15:04:05.000000000Z07:00
func (s *Syslog) SetTimestampFormat(format string) {
	s.timestampFormat = format
}

// SetBuffersAndWorkers sets the channels buffer size and number of concurrent workers.
// These settings should be thought about together, hence setting both in the same function.
func (s *Syslog) SetBuffersAndWorkers(size uint, workers uint) {
	s.buffer = size

	if workers == 0 {
		// just in case no log registered yet
		stdlog.Println("Invalid number of workers specified, setting to 1")
		log.Warn("Invalid number of workers specified, setting to 1")

		workers = 1
	}

	s.numWorkers = workers
}

// SetFormatFunc sets FormatFunc each worker will call to get
// a Formatter func
func (s *Syslog) SetFormatFunc(fn FormatFunc) {
	s.formatFunc = fn
}

// Run starts the logger consuming on the returned channed
func (s *Syslog) Run() chan<- *log.Entry {

	// in a big high traffic app, set a higher buffer
	ch := make(chan *log.Entry, s.buffer)

	for i := 0; i <= int(s.numWorkers); i++ {
		go s.handleLog(ch)
	}

	return ch
}

// handleLog consumes and logs any Entry's passed to the channel
func (s *Syslog) handleLog(entries <-chan *log.Entry) {

	var e *log.Entry
	var line string

	formatter := s.formatFunc()

	for e = range entries {

		line = string(formatter(e))

		switch e.Level {
		case log.DebugLevel:
			s.writer.Debug(line)
		case log.TraceLevel, log.InfoLevel:
			s.writer.Info(line)
		case log.NoticeLevel:
			s.writer.Notice(line)
		case log.WarnLevel:
			s.writer.Warning(line)
		case log.ErrorLevel:
			s.writer.Err(line)
		case log.PanicLevel, log.AlertLevel:
			s.writer.Alert(line)
		case log.FatalLevel:
			s.writer.Crit(line)
		}

		e.Consumed()
	}
}

func (s *Syslog) defaultFormatFunc() Formatter {

	var b []byte
	var file string

	if s.displayColor {

		var color log.ANSIEscSeq

		return func(e *log.Entry) []byte {
			b = b[0:0]
			color = s.colors[e.Level]

			if e.Line == 0 {

				b = append(b, fmt.Sprintf(colorFormat, e.Timestamp.Format(s.timestampFormat), color, e.Level, log.Reset, e.Message)...)

			} else {
				file = e.File
				for i := len(file) - 1; i > 0; i-- {
					if file[i] == '/' {
						file = file[i+1:]
						break
					}
				}

				b = append(b, fmt.Sprintf(colorFormatCaller, e.Timestamp.Format(s.timestampFormat), color, e.Level, log.Reset, file, e.Line, e.Message)...)
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

			b = append(b, fmt.Sprintf(format, e.Timestamp.Format(s.timestampFormat), e.Level, e.Message)...)

		} else {
			file = e.File
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					file = file[i+1:]
					break
				}
			}

			b = append(b, fmt.Sprintf(formatCaller, e.Timestamp.Format(s.timestampFormat), e.Level, file, e.Line, e.Message)...)
		}

		for _, f := range e.Fields {
			b = append(b, fmt.Sprintf(noColorKeyValue, f.Key, f.Value)...)
		}

		return b
	}
}
