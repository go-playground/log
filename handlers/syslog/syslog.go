package syslog

import (
	"crypto/tls"
	"fmt"
	stdlog "log"
	"os"
	"strconv"
	"strings"

	syslog "github.com/RackSec/srslog"

	"github.com/go-playground/ansi"
	"github.com/go-playground/log"
)

// FormatFunc is the function that the workers use to create
// a new Formatter per worker allowing reusable go routine safe
// variable to be used within your Formatter function.
type FormatFunc func(s *Syslog) Formatter

// Formatter is the function used to format the Redis entry
type Formatter func(e *log.Entry) []byte

const (
	space  = byte(' ')
	equals = byte('=')
	colon  = byte(':')
	base10 = 10
	v      = "%v"
	gopath = "GOPATH"
)

// Syslog is an instance of the syslog logger
type Syslog struct {
	buffer          uint
	numWorkers      uint
	colors          [9]ansi.EscSeq
	writer          *syslog.Writer
	formatFunc      FormatFunc
	timestampFormat string
	gopath          string
	fileDisplay     log.FilenameDisplay
	displayColor    bool
}

var (
	// Colors mapping.
	defaultColors = [...]ansi.EscSeq{
		log.DebugLevel:  ansi.Green,
		log.TraceLevel:  ansi.White,
		log.InfoLevel:   ansi.Blue,
		log.NoticeLevel: ansi.LightCyan,
		log.WarnLevel:   ansi.LightYellow,
		log.ErrorLevel:  ansi.LightRed,
		log.PanicLevel:  ansi.Red,
		log.AlertLevel:  ansi.Red + ansi.Underline,
		log.FatalLevel:  ansi.Red + ansi.Underline + ansi.Blink,
	}
)

// New returns a new instance of the syslog logger
// example: syslog.New("udp", "localhost:514", syslog.LOG_DEBUG, "", nil)
// NOTE: tlsConfig param is optional and only applies when networks in "tcp+tls"
// see TestSyslogTLS func tion int syslog_test.go for an example usage of tlsConfig parameter
func New(network string, raddr string, tag string, tlsConfig *tls.Config) (*Syslog, error) {

	var err error

	s := &Syslog{
		buffer:          3,
		numWorkers:      3,
		colors:          defaultColors,
		displayColor:    false,
		timestampFormat: log.DefaultTimeFormat,
		fileDisplay:     log.Lshortfile,
		formatFunc:      defaultFormatFunc,
	}

	// if non-TLS
	if tlsConfig == nil {
		s.writer, err = syslog.Dial(network, raddr, syslog.LOG_INFO, tag)
	} else {
		s.writer, err = syslog.DialWithTLSConfig(network, raddr, syslog.LOG_INFO, tag, tlsConfig)
	}

	if err != nil {
		return nil, err
	}

	return s, nil
}

// SetFilenameDisplay tells Syslog the filename, when present, how to display
func (s *Syslog) SetFilenameDisplay(fd log.FilenameDisplay) {
	s.fileDisplay = fd
}

// FilenameDisplay returns Syslog's current filename display setting
func (s *Syslog) FilenameDisplay() log.FilenameDisplay {
	return s.fileDisplay
}

// SetDisplayColor tells Syslog to output in color or not
// Default is : true
func (s *Syslog) SetDisplayColor(color bool) {
	s.displayColor = color
}

// DisplayColor returns if logging color or not
func (s *Syslog) DisplayColor() bool {
	return s.displayColor
}

// GetDisplayColor returns the color for the given log level
func (s *Syslog) GetDisplayColor(level log.Level) ansi.EscSeq {
	return s.colors[level]
}

// SetTimestampFormat sets Syslog's timestamp output format
// Default is : 2006-01-02T15:04:05.000000000Z07:00
func (s *Syslog) SetTimestampFormat(format string) {
	s.timestampFormat = format
}

// TimestampFormat returns Syslog's current timestamp output format
func (s *Syslog) TimestampFormat() string {
	return s.timestampFormat
}

// GOPATH returns the GOPATH calculated by Syslog
func (s *Syslog) GOPATH() string {
	return s.gopath
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

	// pre-setup
	if s.fileDisplay == log.Llongfile {
		// gather $GOPATH for use in stripping off of full name
		// if not found still ok as will be blank
		s.gopath = os.Getenv(gopath)
		if len(s.gopath) != 0 {
			s.gopath += string(os.PathSeparator) + "src" + string(os.PathSeparator)
		}
	}

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

	formatter := s.formatFunc(s)

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

func defaultFormatFunc(s *Syslog) Formatter {

	var b []byte
	var file string
	var lvl string
	var i int
	gopath := s.GOPATH()
	tsFormat := s.TimestampFormat()
	fnameDisplay := s.FilenameDisplay()

	if s.DisplayColor() {

		var color ansi.EscSeq

		return func(e *log.Entry) []byte {
			b = b[0:0]
			color = s.GetDisplayColor(e.Level)

			if e.Line == 0 {

				b = append(b, e.Timestamp.Format(tsFormat)...)
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

			} else {
				file = e.File

				if fnameDisplay == log.Lshortfile {

					for i = len(file) - 1; i > 0; i-- {
						if file[i] == '/' {
							file = file[i+1:]
							break
						}
					}
				} else {
					file = file[len(gopath):]
				}

				b = append(b, e.Timestamp.Format(tsFormat)...)
				b = append(b, space)
				b = append(b, color...)

				lvl = e.Level.String()

				for i = 0; i < 6-len(lvl); i++ {
					b = append(b, space)
				}

				b = append(b, lvl...)
				b = append(b, ansi.Reset...)
				b = append(b, space)
				b = append(b, file...)
				b = append(b, colon)
				b = strconv.AppendInt(b, int64(e.Line), base10)
				b = append(b, space)
				b = append(b, e.Message...)
			}

			for _, f := range e.Fields {
				b = append(b, space)
				b = append(b, color...)
				b = append(b, f.Key...)
				b = append(b, ansi.Reset...)
				b = append(b, equals)

				switch f.Value.(type) {
				case string:
					b = append(b, f.Value.(string)...)
				case int:
					b = strconv.AppendInt(b, int64(f.Value.(int)), base10)
				case int8:
					b = strconv.AppendInt(b, int64(f.Value.(int8)), base10)
				case int16:
					b = strconv.AppendInt(b, int64(f.Value.(int16)), base10)
				case int32:
					b = strconv.AppendInt(b, int64(f.Value.(int32)), base10)
				case int64:
					b = strconv.AppendInt(b, f.Value.(int64), base10)
				case uint:
					b = strconv.AppendUint(b, uint64(f.Value.(uint)), base10)
				case uint8:
					b = strconv.AppendUint(b, uint64(f.Value.(uint8)), base10)
				case uint16:
					b = strconv.AppendUint(b, uint64(f.Value.(uint16)), base10)
				case uint32:
					b = strconv.AppendUint(b, uint64(f.Value.(uint32)), base10)
				case uint64:
					b = strconv.AppendUint(b, f.Value.(uint64), base10)
				case float32:
					b = strconv.AppendFloat(b, float64(f.Value.(float32)), 'f', -1, 32)
				case float64:
					b = strconv.AppendFloat(b, f.Value.(float64), 'f', -1, 64)
				case bool:
					b = strconv.AppendBool(b, f.Value.(bool))
				default:
					b = append(b, fmt.Sprintf(v, f.Value)...)
				}
			}

			return b
		}
	}

	return func(e *log.Entry) []byte {
		b = b[0:0]

		if e.Line == 0 {

			b = append(b, e.Timestamp.Format(tsFormat)...)
			b = append(b, space)

			lvl = e.Level.String()

			for i = 0; i < 6-len(lvl); i++ {
				b = append(b, space)
			}

			b = append(b, lvl...)
			b = append(b, space)
			b = append(b, e.Message...)

		} else {
			file = e.File

			if fnameDisplay == log.Lshortfile {

				for i = len(file) - 1; i > 0; i-- {
					if file[i] == '/' {
						file = file[i+1:]
						break
					}
				}
			} else {

				// additional check, just in case user does
				// have a $GOPATH but code isn't under it.
				if strings.HasPrefix(file, gopath) {
					file = file[len(gopath):]
				}
			}

			b = append(b, e.Timestamp.Format(tsFormat)...)
			b = append(b, space)

			lvl = e.Level.String()

			for i = 0; i < 6-len(lvl); i++ {
				b = append(b, space)
			}

			b = append(b, lvl...)
			b = append(b, space)
			b = append(b, file...)
			b = append(b, colon)
			b = strconv.AppendInt(b, int64(e.Line), base10)
			b = append(b, space)
			b = append(b, e.Message...)
		}

		for _, f := range e.Fields {
			b = append(b, space)
			b = append(b, f.Key...)
			b = append(b, equals)

			switch f.Value.(type) {
			case string:
				b = append(b, f.Value.(string)...)
			case int:
				b = strconv.AppendInt(b, int64(f.Value.(int)), base10)
			case int8:
				b = strconv.AppendInt(b, int64(f.Value.(int8)), base10)
			case int16:
				b = strconv.AppendInt(b, int64(f.Value.(int16)), base10)
			case int32:
				b = strconv.AppendInt(b, int64(f.Value.(int32)), base10)
			case int64:
				b = strconv.AppendInt(b, f.Value.(int64), base10)
			case uint:
				b = strconv.AppendUint(b, uint64(f.Value.(uint)), base10)
			case uint8:
				b = strconv.AppendUint(b, uint64(f.Value.(uint8)), base10)
			case uint16:
				b = strconv.AppendUint(b, uint64(f.Value.(uint16)), base10)
			case uint32:
				b = strconv.AppendUint(b, uint64(f.Value.(uint32)), base10)
			case uint64:
				b = strconv.AppendUint(b, f.Value.(uint64), base10)
			case bool:
				b = strconv.AppendBool(b, f.Value.(bool))
			default:
				b = append(b, fmt.Sprintf(v, f.Value)...)
			}
		}

		return b
	}
}
