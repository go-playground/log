package log

import (
	"bytes"
	"fmt"
	"io"
	"sync"
	"testing"
	"time"

	. "gopkg.in/go-playground/assert.v1"
)

// NOTES:
// - Run "go test" to run tests
// - Run "gocov test | gocov report" to report on test converage by file
// - Run "gocov test | gocov annotate -" to report on all code and functions, those ,marked with "MISS" were never called
//
// or
//
// -- may be a good idea to change to output path to somewherelike /tmp
// go test -coverprofile cover.out && go tool cover -html=cover.out -o cover.html
//

type testHandler struct {
	writer io.Writer
}

// Run runs handler
func (th *testHandler) Run() chan<- *Entry {
	ch := make(chan *Entry, 0)

	go th.handleLogEntry(ch)

	return ch
}

func (th *testHandler) handleLogEntry(entries <-chan *Entry) {

	var e *Entry

	for e = range entries {
		s := e.Message

		for _, f := range e.Fields {
			s += fmt.Sprintf(" %s=%v", f.Key, f.Value)
		}

		th.writer.Write([]byte(s))

		e.WG.Done()
	}
}

func TestConsoleLogger(t *testing.T) {

	buff := new(bytes.Buffer)

	th := &testHandler{
		writer: buff,
	}

	RegisterHandler(th, AllLevels...)

	Debug("debug")
	Equal(t, buff.String(), "debug")
	buff.Reset()

	Debugf("%s", "debugf")
	Equal(t, buff.String(), "debugf")
	buff.Reset()

	Info("info")
	Equal(t, buff.String(), "info")
	buff.Reset()

	Infof("%s", "infof")
	Equal(t, buff.String(), "infof")
	buff.Reset()

	Notice("notice")
	Equal(t, buff.String(), "notice")
	buff.Reset()

	Noticef("%s", "noticef")
	Equal(t, buff.String(), "noticef")
	buff.Reset()

	Warn("warn")
	Equal(t, buff.String(), "warn")
	buff.Reset()

	Warnf("%s", "warnf")
	Equal(t, buff.String(), "warnf")
	buff.Reset()

	Error("error")
	Equal(t, buff.String(), "error")
	buff.Reset()

	Errorf("%s", "errorf")
	Equal(t, buff.String(), "errorf")
	buff.Reset()

	Alert("alert")
	Equal(t, buff.String(), "alert")
	buff.Reset()

	Alertf("%s", "alertf")
	Equal(t, buff.String(), "alertf")
	buff.Reset()

	Print("print")
	Equal(t, buff.String(), "print")
	buff.Reset()

	Printf("%s", "printf")
	Equal(t, buff.String(), "printf")
	buff.Reset()

	Println("println")
	Equal(t, buff.String(), "println")
	buff.Reset()

	PanicMatches(t, func() { Panic("panic") }, "panic")
	Equal(t, buff.String(), "panic")
	buff.Reset()

	PanicMatches(t, func() { Panicf("%s", "panicf") }, "panicf")
	Equal(t, buff.String(), "panicf")
	buff.Reset()

	PanicMatches(t, func() { Panicln("panicln") }, "panicln")
	Equal(t, buff.String(), "panicln")
	buff.Reset()

	// WithFields
	WithFields(F("key", "value")).Info("info")
	Equal(t, buff.String(), "info key=value")
	buff.Reset()

	WithFields(F("key", "value")).Infof("%s", "infof")
	Equal(t, buff.String(), "infof key=value")
	buff.Reset()

	WithFields(F("key", "value")).Notice("notice")
	Equal(t, buff.String(), "notice key=value")
	buff.Reset()

	WithFields(F("key", "value")).Noticef("%s", "noticef")
	Equal(t, buff.String(), "noticef key=value")
	buff.Reset()

	WithFields(F("key", "value")).Debug("debug")
	Equal(t, buff.String(), "debug key=value")
	buff.Reset()

	WithFields(F("key", "value")).Debugf("%s", "debugf")
	Equal(t, buff.String(), "debugf key=value")
	buff.Reset()

	WithFields(F("key", "value")).Warn("warn")
	Equal(t, buff.String(), "warn key=value")
	buff.Reset()

	WithFields(F("key", "value")).Warnf("%s", "warnf")
	Equal(t, buff.String(), "warnf key=value")
	buff.Reset()

	WithFields(F("key", "value")).Error("error")
	Equal(t, buff.String(), "error key=value")
	buff.Reset()

	WithFields(F("key", "value")).Errorf("%s", "errorf")
	Equal(t, buff.String(), "errorf key=value")
	buff.Reset()

	WithFields(F("key", "value")).Alert("alert")
	Equal(t, buff.String(), "alert key=value")
	buff.Reset()

	WithFields(F("key", "value")).Alertf("%s", "alertf")
	Equal(t, buff.String(), "alertf key=value")
	buff.Reset()

	PanicMatches(t, func() { WithFields(F("key", "value")).Panicf("%s", "panicf") }, "panicf key=value")
	Equal(t, buff.String(), "panicf key=value")
	buff.Reset()

	PanicMatches(t, func() { WithFields(F("key", "value")).Panic("panic") }, "panic key=value")
	Equal(t, buff.String(), "panic key=value")
	buff.Reset()

	func() {
		defer Trace("trace").End()
	}()

	// TODO: finish up regex
	MatchRegex(t, buff.String(), "^trace\\s+\\.*")
	buff.Reset()

	func() {
		defer Tracef("tracef").End()
	}()

	// TODO: finish up regex
	MatchRegex(t, buff.String(), "^tracef\\s+\\.*")
	buff.Reset()

	func() {
		defer WithFields(F("key", "value")).Trace("trace").End()
	}()

	// TODO: finish up regex
	MatchRegex(t, buff.String(), "^trace\\s+\\.*")
	buff.Reset()

	func() {
		defer WithFields(F("key", "value")).Tracef("tracef").End()
	}()

	// TODO: finish up regex
	MatchRegex(t, buff.String(), "^tracef\\s+\\.*")
	buff.Reset()
}

func TestLevel(t *testing.T) {
	l := Level(9999)
	Equal(t, l.String(), "Unknow Level")

	Equal(t, DebugLevel.String(), "DEBUG")
	Equal(t, TraceLevel.String(), "TRACE")
	Equal(t, InfoLevel.String(), "INFO")
	Equal(t, NoticeLevel.String(), "NOTICE")
	Equal(t, WarnLevel.String(), "WARN")
	Equal(t, ErrorLevel.String(), "ERROR")
	Equal(t, PanicLevel.String(), "PANIC")
	Equal(t, AlertLevel.String(), "ALERT")
	Equal(t, FatalLevel.String(), "FATAL")
}

func TestSettings(t *testing.T) {
	RegisterDurationFunc(func(d time.Duration) string {
		return fmt.Sprintf("%gs", d.Seconds())
	})

	SetTimeFormat(time.RFC1123)
}

func TestEntry(t *testing.T) {

	SetApplicationID("app-log")

	// Resetting pool to ensure no Entries exist before setting the Application ID
	Logger.entryPool = &sync.Pool{New: func() interface{} {
		return &Entry{
			WG:            new(sync.WaitGroup),
			ApplicationID: Logger.getApplicationID(),
		}
	}}

	e := Logger.entryPool.Get().(*Entry)
	Equal(t, e.ApplicationID, "app-log")
	NotEqual(t, e.WG, nil)

	e = newEntry(InfoLevel, "test", []Field{F("key", "value")})
	HandleEntry(e)
}
