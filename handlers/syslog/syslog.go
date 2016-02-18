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
	buffer    uint
	writer    *syslog.Writer
	formatter Formatter
}

var syslogBuffPool = &sync.Pool{New: func() interface{} {
	return new(bytes.Buffer)
}}

// New returns a new instance of the syslog logger
// example: syslog.New("udp", "localhost:514", syslog.LOG_DEBUG, "")
func New(network string, raddr string, priority syslog.Priority, tag string) (*Syslog, error) {

	var err error

	s := &Syslog{
		buffer:    0,
		formatter: defaultFormatEntry,
	}

	if s.writer, err = syslog.Dial(network, raddr, priority, tag); err != nil {
		return nil, err
	}

	return s, nil
}

// SetFormatter sets the  Syslog entry formatter
// Default is : defaultFormatEntry
func (s *Syslog) SetFormatter(f Formatter) {
	s.formatter = f
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

func defaultFormatEntry(e *log.Entry) string {

	buff := syslogBuffPool.Get().(*bytes.Buffer)
	buff.Reset()

	fmt.Fprintf(buff, "%6s[%s] %s", e.Level, e.Timestamp.Format(time.RFC3339Nano), e.Message)

	for _, f := range e.Fields {
		fmt.Fprintf(buff, " %s=%v", f.Key, f.Value)
	}

	s := buff.String()
	syslogBuffPool.Put(buff)

	return s
}
