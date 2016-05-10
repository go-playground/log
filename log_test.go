package log

import (
	"bytes"
	"fmt"
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

func TestConsoleLogger(t *testing.T) {

	buff := new(bytes.Buffer)

	th := &testHandler{
		writer: buff,
	}

	Logger.RegisterHandler(th, AllLevels...)

	Logger.Debug("debug")
	Equal(t, buff.String(), "debug")
	buff.Reset()

	Logger.Debugf("%s", "debugf")
	Equal(t, buff.String(), "debugf")
	buff.Reset()

	Logger.Info("info")
	Equal(t, buff.String(), "info")
	buff.Reset()

	Logger.Infof("%s", "infof")
	Equal(t, buff.String(), "infof")
	buff.Reset()

	Logger.Notice("notice")
	Equal(t, buff.String(), "notice")
	buff.Reset()

	Logger.Noticef("%s", "noticef")
	Equal(t, buff.String(), "noticef")
	buff.Reset()

	Logger.Warn("warn")
	Equal(t, buff.String(), "warn")
	buff.Reset()

	Logger.Warnf("%s", "warnf")
	Equal(t, buff.String(), "warnf")
	buff.Reset()

	Logger.Error("error")
	Equal(t, buff.String(), "error")
	buff.Reset()

	Logger.Errorf("%s", "errorf")
	Equal(t, buff.String(), "errorf")
	buff.Reset()

	Logger.Alert("alert")
	Equal(t, buff.String(), "alert")
	buff.Reset()

	Logger.Alertf("%s", "alertf")
	Equal(t, buff.String(), "alertf")
	buff.Reset()

	PanicMatches(t, func() { Logger.Panic("panic") }, "panic")
	Equal(t, buff.String(), "panic")
	buff.Reset()

	PanicMatches(t, func() { Logger.Panicf("%s", "panicf") }, "panicf")
	Equal(t, buff.String(), "panicf")
	buff.Reset()

	// WithFields
	Logger.WithFields(F("key", "value")).Info("info")
	Equal(t, buff.String(), "info key=value")
	buff.Reset()

	Logger.WithFields(F("key", "value")).Infof("%s", "infof")
	Equal(t, buff.String(), "infof key=value")
	buff.Reset()

	Logger.WithFields(F("key", "value")).Notice("notice")
	Equal(t, buff.String(), "notice key=value")
	buff.Reset()

	Logger.WithFields(F("key", "value")).Noticef("%s", "noticef")
	Equal(t, buff.String(), "noticef key=value")
	buff.Reset()

	Logger.WithFields(F("key", "value")).Debug("debug")
	Equal(t, buff.String(), "debug key=value")
	buff.Reset()

	Logger.WithFields(F("key", "value")).Debugf("%s", "debugf")
	Equal(t, buff.String(), "debugf key=value")
	buff.Reset()

	Logger.WithFields(F("key", "value")).Warn("warn")
	Equal(t, buff.String(), "warn key=value")
	buff.Reset()

	Logger.WithFields(F("key", "value")).Warnf("%s", "warnf")
	Equal(t, buff.String(), "warnf key=value")
	buff.Reset()

	Logger.WithFields(F("key", "value")).Error("error")
	Equal(t, buff.String(), "error key=value")
	buff.Reset()

	Logger.WithFields(F("key", "value")).Errorf("%s", "errorf")
	Equal(t, buff.String(), "errorf key=value")
	buff.Reset()

	Logger.WithFields(F("key", "value")).Alert("alert")
	Equal(t, buff.String(), "alert key=value")
	buff.Reset()

	Logger.WithFields(F("key", "value")).Alertf("%s", "alertf")
	Equal(t, buff.String(), "alertf key=value")
	buff.Reset()

	PanicMatches(t, func() { Logger.WithFields(F("key", "value")).Panicf("%s", "panicf") }, "panicf key=value")
	Equal(t, buff.String(), "panicf key=value")
	buff.Reset()

	PanicMatches(t, func() { Logger.WithFields(F("key", "value")).Panic("panic") }, "panic key=value")
	Equal(t, buff.String(), "panic key=value")
	buff.Reset()

	func() {
		defer Logger.Trace("trace").End()
	}()

	// TODO: finish up regex
	MatchRegex(t, buff.String(), "^trace\\s+\\.*")
	buff.Reset()

	func() {
		defer Logger.Tracef("tracef").End()
	}()

	// TODO: finish up regex
	MatchRegex(t, buff.String(), "^tracef\\s+\\.*")
	buff.Reset()

	func() {
		defer Logger.WithFields(F("key", "value")).Trace("trace").End()
	}()

	// TODO: finish up regex
	MatchRegex(t, buff.String(), "^trace\\s+\\.*")
	buff.Reset()

	func() {
		defer Logger.WithFields(F("key", "value")).Tracef("tracef").End()
	}()

	// TODO: finish up regex
	MatchRegex(t, buff.String(), "^tracef\\s+\\.*")
	buff.Reset()

	// Test Custom Entry ( most common case is Unmarshalled from JSON when using centralized logging)
	entry := new(Entry)
	entry.ApplicationID = "APP"
	entry.Level = InfoLevel
	entry.Timestamp = time.Now().UTC()
	entry.Message = "Test Message"
	entry.Fields = make([]Field, 0)
	Logger.HandleEntry(entry)
	Equal(t, buff.String(), "Test Message")
	buff.Reset()
}

func TestConsoleLoggerCaller(t *testing.T) {

	buff := new(bytes.Buffer)

	Logger.SetCallerInfo(true)

	th := &testHandler{
		writer: buff,
	}

	Logger.RegisterHandler(th, AllLevels...)

	Logger.Debug("debug")
	Equal(t, buff.String(), "debug")
	buff.Reset()

	Logger.Debugf("%s", "debugf")
	Equal(t, buff.String(), "debugf")
	buff.Reset()

	Logger.Info("info")
	Equal(t, buff.String(), "info")
	buff.Reset()

	Logger.Infof("%s", "infof")
	Equal(t, buff.String(), "infof")
	buff.Reset()

	Logger.Notice("notice")
	Equal(t, buff.String(), "notice")
	buff.Reset()

	Logger.Noticef("%s", "noticef")
	Equal(t, buff.String(), "noticef")
	buff.Reset()

	Logger.Warn("warn")
	Equal(t, buff.String(), "warn")
	buff.Reset()

	Logger.Warnf("%s", "warnf")
	Equal(t, buff.String(), "warnf")
	buff.Reset()

	Logger.Error("error")
	Equal(t, buff.String(), "error")
	buff.Reset()

	Logger.Errorf("%s", "errorf")
	Equal(t, buff.String(), "errorf")
	buff.Reset()

	Logger.Alert("alert")
	Equal(t, buff.String(), "alert")
	buff.Reset()

	Logger.Alertf("%s", "alertf")
	Equal(t, buff.String(), "alertf")
	buff.Reset()

	PanicMatches(t, func() { Logger.Panic("panic") }, "panic")
	Equal(t, buff.String(), "panic")
	buff.Reset()

	PanicMatches(t, func() { Logger.Panicf("%s", "panicf") }, "panicf")
	Equal(t, buff.String(), "panicf")
	buff.Reset()

	// WithFields
	Logger.WithFields(F("key", "value")).Info("info")
	Equal(t, buff.String(), "info key=value")
	buff.Reset()

	Logger.WithFields(F("key", "value")).Infof("%s", "infof")
	Equal(t, buff.String(), "infof key=value")
	buff.Reset()

	Logger.WithFields(F("key", "value")).Notice("notice")
	Equal(t, buff.String(), "notice key=value")
	buff.Reset()

	Logger.WithFields(F("key", "value")).Noticef("%s", "noticef")
	Equal(t, buff.String(), "noticef key=value")
	buff.Reset()

	Logger.WithFields(F("key", "value")).Debug("debug")
	Equal(t, buff.String(), "debug key=value")
	buff.Reset()

	Logger.WithFields(F("key", "value")).Debugf("%s", "debugf")
	Equal(t, buff.String(), "debugf key=value")
	buff.Reset()

	Logger.WithFields(F("key", "value")).Warn("warn")
	Equal(t, buff.String(), "warn key=value")
	buff.Reset()

	Logger.WithFields(F("key", "value")).Warnf("%s", "warnf")
	Equal(t, buff.String(), "warnf key=value")
	buff.Reset()

	Logger.WithFields(F("key", "value")).Error("error")
	Equal(t, buff.String(), "error key=value")
	buff.Reset()

	Logger.WithFields(F("key", "value")).Errorf("%s", "errorf")
	Equal(t, buff.String(), "errorf key=value")
	buff.Reset()

	Logger.WithFields(F("key", "value")).Alert("alert")
	Equal(t, buff.String(), "alert key=value")
	buff.Reset()

	Logger.WithFields(F("key", "value")).Alertf("%s", "alertf")
	Equal(t, buff.String(), "alertf key=value")
	buff.Reset()

	PanicMatches(t, func() { Logger.WithFields(F("key", "value")).Panicf("%s", "panicf") }, "panicf key=value")
	Equal(t, buff.String(), "panicf key=value")
	buff.Reset()

	PanicMatches(t, func() { Logger.WithFields(F("key", "value")).Panic("panic") }, "panic key=value")
	Equal(t, buff.String(), "panic key=value")
	buff.Reset()

	func() {
		defer Logger.Trace("trace").End()
	}()

	// TODO: finish up regex
	MatchRegex(t, buff.String(), "^trace\\s+\\.*")
	buff.Reset()

	func() {
		defer Logger.Tracef("tracef").End()
	}()

	// TODO: finish up regex
	MatchRegex(t, buff.String(), "^tracef\\s+\\.*")
	buff.Reset()

	func() {
		defer Logger.WithFields(F("key", "value")).Trace("trace").End()
	}()

	// TODO: finish up regex
	MatchRegex(t, buff.String(), "^trace\\s+\\.*")
	buff.Reset()

	func() {
		defer Logger.WithFields(F("key", "value")).Tracef("tracef").End()
	}()

	// TODO: finish up regex
	MatchRegex(t, buff.String(), "^tracef\\s+\\.*")
	buff.Reset()

	// Test Custom Entry ( most common case is Unmarshalled from JSON when using centralized logging)
	entry := new(Entry)
	entry.ApplicationID = "APP"
	entry.Level = InfoLevel
	entry.Timestamp = time.Now().UTC()
	entry.Message = "Test Message"
	entry.Fields = make([]Field, 0)
	Logger.HandleEntry(entry)
	Equal(t, buff.String(), "Test Message")
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
	Logger.RegisterDurationFunc(func(d time.Duration) string {
		return fmt.Sprintf("%gs", d.Seconds())
	})

	Logger.SetTimeFormat(time.RFC1123)
}

func TestEntry(t *testing.T) {

	Logger.SetApplicationID("app-log")

	// Resetting pool to ensure no Entries exist before setting the Application ID
	Logger.entryPool = &sync.Pool{New: func() interface{} {
		return &Entry{
			wg:            new(sync.WaitGroup),
			ApplicationID: Logger.getApplicationID(),
		}
	}}

	e := Logger.entryPool.Get().(*Entry)
	Equal(t, e.ApplicationID, "app-log")
	NotEqual(t, e.wg, nil)

	e = newEntry(InfoLevel, "test", []Field{F("key", "value")}, 0)
	Logger.HandleEntry(e)
}

func TestFatal(t *testing.T) {
	var i int

	exitFunc = func(code int) {
		i = code
	}

	Logger.Fatal("fatal")
	Equal(t, i, 1)

	Logger.Fatalf("fatalf")
	Equal(t, i, 1)

	Logger.WithFields(F("key", "value")).Fatal("fatal")
	Equal(t, i, 1)

	Logger.WithFields(F("key", "value")).Fatalf("fatalf")
	Equal(t, i, 1)
}

func TestColors(t *testing.T) {

	fmt.Printf("%sBlack%s\n", Black, Reset)
	fmt.Printf("%sDarkGray%s\n", DarkGray, Reset)
	fmt.Printf("%sBlue%s\n", Blue, Reset)
	fmt.Printf("%sLightBlue%s\n", LightBlue, Reset)
	fmt.Printf("%sGreen%s\n", Green, Reset)
	fmt.Printf("%sLightGreen%s\n", LightGreen, Reset)
	fmt.Printf("%sCyan%s\n", Cyan, Reset)
	fmt.Printf("%sLightCyan%s\n", LightCyan, Reset)
	fmt.Printf("%sRed%s\n", Red, Reset)
	fmt.Printf("%sLightRed%s\n", LightRed, Reset)
	fmt.Printf("%sMagenta%s\n", Magenta, Reset)
	fmt.Printf("%sLightMagenta%s\n", LightMagenta, Reset)
	fmt.Printf("%sBrown%s\n", Brown, Reset)
	fmt.Printf("%sYellow%s\n", Yellow, Reset)
	fmt.Printf("%sLightGray%s\n", LightGray, Reset)
	fmt.Printf("%sWhite%s\n", White, Reset)

	fmt.Printf("%s%sUnderscoreRed%s\n", Red, Underscore, Reset)
	fmt.Printf("%s%sBlinkRed%s\n", Red, Blink, Reset)
	fmt.Printf("%s%s%sBlinkUnderscoreRed%s\n", Red, Blink, Underscore, Reset)

	fmt.Printf("%s%sRedInverse%s\n", Red, Inverse, Reset)
	fmt.Printf("%sGreenInverse%s\n", Green+Inverse, Reset)
}
