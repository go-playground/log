package syslog

import (
	"bytes"
	"fmt"
	"log/syslog"
	"sync"
	"time"

	"github.com/go-playground/log"
)

// Formatter is the function used to format the syslog entry
type Formatter func(e *log.Entry) string

// Syslog is an instance of the syslog logger
type Syslog struct {
	buffer             uint
	colors             [9]log.ANSIEscSeq
	ansiReset          log.ANSIEscSeq
	displayColor       bool
	writer             *syslog.Writer
	hasCustomFormatter bool
	formatter          Formatter
	timestampFormat    string
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

	syslogBuffPool = &sync.Pool{New: func() interface{} {
		return new(bytes.Buffer)
	}}
)

// New returns a new instance of the syslog logger
// example: syslog.New("udp", "localhost:514", syslog.LOG_DEBUG, "")
func New(network string, raddr string, priority syslog.Priority, tag string) (*Syslog, error) {

	var err error

	s := &Syslog{
		buffer:             0,
		colors:             defaultColors,
		ansiReset:          log.Reset,
		displayColor:       false,
		timestampFormat:    time.RFC3339Nano,
		hasCustomFormatter: false,
	}

	s.formatter = s.defaultFormatEntry

	if s.writer, err = syslog.Dial(network, raddr, priority, tag); err != nil {
		return nil, err
	}

	return s, nil
}

// SetLevelColor updates Console's level color values
func (s *Syslog) SetLevelColor(l log.Level, color log.ANSIEscSeq) {
	s.colors[l] = color
}

// SetANSIReset sets the ANSI Reset sequence
func (s *Syslog) SetANSIReset(code log.ANSIEscSeq) {
	s.ansiReset = code
}

// DisplayColor tells Console to output in color or not
// Default is : true
func (s *Syslog) DisplayColor(color bool) {
	s.displayColor = color

	if !s.hasCustomFormatter {
		if color {
			s.formatter = s.defaultFormatEntryColor
		} else {
			s.formatter = s.defaultFormatEntry
		}
	}
}

// SetTimestampFormat sets Console's timestamp output format
// Default is : time.RFC3339Nano
func (s *Syslog) SetTimestampFormat(format string) {
	s.timestampFormat = format
}

// SetFormatter sets the  Syslog entry formatter
// Default is : defaultFormatEntry
func (s *Syslog) SetFormatter(f Formatter) {
	s.formatter = f
	s.hasCustomFormatter = true
}

// SetChannelBuffer tells Syslog what the channel buffer size should be
// Default is : 0
func (s *Syslog) SetChannelBuffer(i uint) {
	s.buffer = i
}

// Run starts the logger consuming on the returned channed
func (s *Syslog) Run() chan<- log.Entry {

	// in a big high traffic app, set a higher buffer
	ch := make(chan log.Entry, s.buffer)

	go s.handleLog(ch)

	return ch
}

// handleLog consumes and logs any Entry's passed to the channel
func (s *Syslog) handleLog(entries <-chan log.Entry) {

	var e log.Entry
	var line string

	for e = range entries {

		line = s.formatter(&e)

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

		e.WG.Done()
	}
}

func (s *Syslog) defaultFormatEntry(e *log.Entry) string {

	buff := syslogBuffPool.Get().(*bytes.Buffer)
	buff.Reset()

	fmt.Fprintf(buff, "%6s[%s] %s", e.Level, e.Timestamp.Format(s.timestampFormat), e.Message)

	for _, f := range e.Fields {
		fmt.Fprintf(buff, " %s=%v", f.Key, f.Value)
	}

	str := buff.String()
	syslogBuffPool.Put(buff)

	return str
}

func (s *Syslog) defaultFormatEntryColor(e *log.Entry) string {

	color := s.colors[e.Level]
	l := len(e.Fields)
	buff := syslogBuffPool.Get().(*bytes.Buffer)
	buff.Reset()

	if l == 0 {
		fmt.Fprintf(buff, "%s%6s%s[%s] %s", color, e.Level, s.ansiReset, e.Timestamp.Format(s.timestampFormat), e.Message)
	} else {
		fmt.Fprintf(buff, "%s%6s%s[%s] %-25s", color, e.Level, s.ansiReset, e.Timestamp.Format(s.timestampFormat), e.Message)
	}

	for _, f := range e.Fields {
		fmt.Fprintf(buff, " %s%s%s=%v", color, f.Key, s.ansiReset, f.Value)
	}

	str := buff.String()
	syslogBuffPool.Put(buff)

	return str
}
