package console

import (
	"bytes"
	"testing"
	"time"

	"github.com/go-playground/log"
	. "gopkg.in/go-playground/assert.v1"
)

// NOTES:
// - Run "go test" to run tests
// - Run "gocov test | gocov report" to report on test converage by file
// - Run "gocov test | gocov annotate -" to report on all code and functions, those ,marked with "MISS" were never called

// or

// -- may be a good idea to change to output path to somewherelike /tmp
// go test -coverprofile cover.out && go tool cover -html=cover.out -o cover.html

func TestConsoleLogger(t *testing.T) {

	buff := new(bytes.Buffer)

	cLog := New()
	cLog.SetWriter(buff)
	cLog.DisplayColor(false)
	cLog.SetChannelBuffer(3)
	cLog.SetTimestampFormat(time.RFC3339)
	cLog.UseMiniTimestamp(true)
	cLog.SetANSIReset(log.Reset)

	log.RegisterHandler(cLog, log.AllLevels...)

	log.Debug("debug")
	Equal(t, buff.String(), "0000  DEBUG debug\n")
	buff.Reset()

	log.Debugf("%s", "debugf")
	Equal(t, buff.String(), "0000  DEBUG debugf\n")
	buff.Reset()

	log.Info("info")
	Equal(t, buff.String(), "0000   INFO info\n")
	buff.Reset()

	log.Infof("%s", "infof")
	Equal(t, buff.String(), "0000   INFO infof\n")
	buff.Reset()

	log.Notice("notice")
	Equal(t, buff.String(), "0000 NOTICE notice\n")
	buff.Reset()

	log.Noticef("%s", "noticef")
	Equal(t, buff.String(), "0000 NOTICE noticef\n")
	buff.Reset()

	log.Warn("warn")
	Equal(t, buff.String(), "0000   WARN warn\n")
	buff.Reset()

	log.Warnf("%s", "warnf")
	Equal(t, buff.String(), "0000   WARN warnf\n")
	buff.Reset()

	log.Error("error")
	Equal(t, buff.String(), "0000  ERROR error\n")
	buff.Reset()

	log.Errorf("%s", "errorf")
	Equal(t, buff.String(), "0000  ERROR errorf\n")
	buff.Reset()

	log.Alert("alert")
	Equal(t, buff.String(), "0000  ALERT alert\n")
	buff.Reset()

	log.Alertf("%s", "alertf")
	Equal(t, buff.String(), "0000  ALERT alertf\n")
	buff.Reset()

	log.Print("print")
	Equal(t, buff.String(), "0000   INFO print\n")
	buff.Reset()

	log.Printf("%s", "printf")
	Equal(t, buff.String(), "0000   INFO printf\n")
	buff.Reset()

	log.Println("println")
	Equal(t, buff.String(), "0000   INFO println\n")
	buff.Reset()

	PanicMatches(t, func() { log.Panic("panic") }, "panic")
	Equal(t, buff.String(), "0000  PANIC panic\n")
	buff.Reset()

	PanicMatches(t, func() { log.Panicf("%s", "panicf") }, "panicf")
	Equal(t, buff.String(), "0000  PANIC panicf\n")
	buff.Reset()

	PanicMatches(t, func() { log.Panicln("panicln") }, "panicln")
	Equal(t, buff.String(), "0000  PANIC panicln\n")
	buff.Reset()

	// WithFields
	log.WithFields(log.F("key", "value")).Debug("debug")
	Equal(t, buff.String(), "0000  DEBUG debug                     key=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Debugf("%s", "debugf")
	Equal(t, buff.String(), "0000  DEBUG debugf                    key=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Info("info")
	Equal(t, buff.String(), "0000   INFO info                      key=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Infof("%s", "infof")
	Equal(t, buff.String(), "0000   INFO infof                     key=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Notice("notice")
	Equal(t, buff.String(), "0000 NOTICE notice                    key=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Noticef("%s", "noticef")
	Equal(t, buff.String(), "0000 NOTICE noticef                   key=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Warn("warn")
	Equal(t, buff.String(), "0000   WARN warn                      key=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Warnf("%s", "warnf")
	Equal(t, buff.String(), "0000   WARN warnf                     key=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Error("error")
	Equal(t, buff.String(), "0000  ERROR error                     key=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Errorf("%s", "errorf")
	Equal(t, buff.String(), "0000  ERROR errorf                    key=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Alert("alert")
	Equal(t, buff.String(), "0000  ALERT alert                     key=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Alertf("%s", "alertf")
	Equal(t, buff.String(), "0000  ALERT alertf                    key=value\n")
	buff.Reset()

	PanicMatches(t, func() { log.WithFields(log.F("key", "value")).Panicf("%s", "panicf") }, "panicf key=value")
	Equal(t, buff.String(), "0000  PANIC panicf                    key=value\n")
	buff.Reset()

	PanicMatches(t, func() { log.WithFields(log.F("key", "value")).Panic("panic") }, "panic key=value")
	Equal(t, buff.String(), "0000  PANIC panic                     key=value\n")
	buff.Reset()

	func() {
		defer log.Trace("trace").End()
	}()

	// TODO: finish up regex
	MatchRegex(t, buff.String(), "^0000\\s\\sTRACE\\strace\\.*")
	buff.Reset()

	func() {
		defer log.Tracef("tracef").End()
	}()

	// TODO: finish up regex
	MatchRegex(t, buff.String(), "^0000\\s\\sTRACE\\stracef\\.*")
	buff.Reset()

	func() {
		defer log.WithFields(log.F("key", "value")).Trace("trace").End()
	}()

	// TODO: finish up regex
	MatchRegex(t, buff.String(), "^0000\\s\\sTRACE\\strace\\.*")
	buff.Reset()

	func() {
		defer log.WithFields(log.F("key", "value")).Tracef("tracef").End()
	}()

	// TODO: finish up regex
	MatchRegex(t, buff.String(), "^0000\\s\\sTRACE\\stracef\\.*")
	buff.Reset()

	// year := time.Now().Format("2006")
	// cLog.UseMiniTimestamp(false)
	// cLog.SetTimestampFormat("2006")

	// log.Info("info")
	// Equal(t, buff.String(), "  INFO["+year+"] info\n")
	// buff.Reset()
}

func TestConsoleLoggerColor(t *testing.T) {

	buff := new(bytes.Buffer)

	cLog := New()
	cLog.SetWriter(buff)
	cLog.DisplayColor(true)
	cLog.SetChannelBuffer(3)
	cLog.SetTimestampFormat(time.RFC3339)
	cLog.UseMiniTimestamp(true)

	log.RegisterHandler(cLog, log.AllLevels...)

	log.Debug("debug")
	Equal(t, buff.String(), "0000 [32m DEBUG[0m debug\n")
	buff.Reset()

	log.Debugf("%s", "debugf")
	Equal(t, buff.String(), "0000 [32m DEBUG[0m debugf\n")
	buff.Reset()

	log.Info("info")
	Equal(t, buff.String(), "0000 [34m  INFO[0m info\n")
	buff.Reset()

	log.Infof("%s", "infof")
	Equal(t, buff.String(), "0000 [34m  INFO[0m infof\n")
	buff.Reset()

	log.Notice("notice")
	Equal(t, buff.String(), "0000 [36;1mNOTICE[0m notice\n")
	buff.Reset()

	log.Notice("%s", "noticef")
	Equal(t, buff.String(), "0000 [36;1mNOTICE[0m %snoticef\n")
	buff.Reset()

	log.Warn("warn")
	Equal(t, buff.String(), "0000 [33;1m  WARN[0m warn\n")
	buff.Reset()

	log.Warnf("%s", "warnf")
	Equal(t, buff.String(), "0000 [33;1m  WARN[0m warnf\n")
	buff.Reset()

	log.Error("error")
	Equal(t, buff.String(), "0000 [31;1m ERROR[0m error\n")
	buff.Reset()

	log.Errorf("%s", "errorf")
	Equal(t, buff.String(), "0000 [31;1m ERROR[0m errorf\n")
	buff.Reset()

	log.Alert("alert")
	Equal(t, buff.String(), "0000 [31m[4m ALERT[0m alert\n")
	buff.Reset()

	log.Alertf("%s", "alertf")
	Equal(t, buff.String(), "0000 [31m[4m ALERT[0m alertf\n")
	buff.Reset()

	log.Print("print")
	Equal(t, buff.String(), "0000 [34m  INFO[0m print\n")
	buff.Reset()

	log.Printf("%s", "printf")
	Equal(t, buff.String(), "0000 [34m  INFO[0m printf\n")
	buff.Reset()

	log.Println("println")
	Equal(t, buff.String(), "0000 [34m  INFO[0m println\n")
	buff.Reset()

	PanicMatches(t, func() { log.Panic("panic") }, "panic")
	Equal(t, buff.String(), "0000 [31m PANIC[0m panic\n")
	buff.Reset()

	PanicMatches(t, func() { log.Panicf("%s", "panicf") }, "panicf")
	Equal(t, buff.String(), "0000 [31m PANIC[0m panicf\n")
	buff.Reset()

	PanicMatches(t, func() { log.Panicln("panicln") }, "panicln")
	Equal(t, buff.String(), "0000 [31m PANIC[0m panicln\n")
	buff.Reset()

	// WithFields
	log.WithFields(log.F("key", "value")).Debug("debug")
	Equal(t, buff.String(), "0000 [32m DEBUG[0m debug                     [32mkey[0m=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Debugf("%s", "debugf")
	Equal(t, buff.String(), "0000 [32m DEBUG[0m debugf                    [32mkey[0m=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Info("info")
	Equal(t, buff.String(), "0000 [34m  INFO[0m info                      [34mkey[0m=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Infof("%s", "infof")
	Equal(t, buff.String(), "0000 [34m  INFO[0m infof                     [34mkey[0m=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Notice("notice")
	Equal(t, buff.String(), "0000 [36;1mNOTICE[0m notice                    [36;1mkey[0m=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Noticef("%s", "noticef")
	Equal(t, buff.String(), "0000 [36;1mNOTICE[0m noticef                   [36;1mkey[0m=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Warn("warn")
	Equal(t, buff.String(), "0000 [33;1m  WARN[0m warn                      [33;1mkey[0m=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Warnf("%s", "warnf")
	Equal(t, buff.String(), "0000 [33;1m  WARN[0m warnf                     [33;1mkey[0m=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Error("error")
	Equal(t, buff.String(), "0000 [31;1m ERROR[0m error                     [31;1mkey[0m=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Errorf("%s", "errorf")
	Equal(t, buff.String(), "0000 [31;1m ERROR[0m errorf                    [31;1mkey[0m=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Alert("alert")
	Equal(t, buff.String(), "0000 [31m[4m ALERT[0m alert                     [31m[4mkey[0m=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Alertf("%s", "alertf")
	Equal(t, buff.String(), "0000 [31m[4m ALERT[0m alertf                    [31m[4mkey[0m=value\n")
	buff.Reset()

	PanicMatches(t, func() { log.WithFields(log.F("key", "value")).Panicf("%s", "panicf") }, "panicf key=value")
	Equal(t, buff.String(), "0000 [31m PANIC[0m panicf                    [31mkey[0m=value\n")
	buff.Reset()

	PanicMatches(t, func() { log.WithFields(log.F("key", "value")).Panic("panic") }, "panic key=value")
	Equal(t, buff.String(), "0000 [31m PANIC[0m panic                     [31mkey[0m=value\n")
	buff.Reset()

	cLog.SetLevelColor(log.DebugLevel, log.LightGreen)

	log.Debug("debug")
	Equal(t, buff.String(), "0000 [32;1m DEBUG[0m debug\n")
	buff.Reset()

	// year := time.Now().Format("2006")
	// cLog.UseMiniTimestamp(false)
	// cLog.SetTimestampFormat("2006")

	// log.Info("info")
	// Equal(t, buff.String(), "[34m  INFO[0m["+year+"] info\n")
	// buff.Reset()

	// log.WithFields(log.F("key", "value")).Info("info")
	// Equal(t, buff.String(), "[34m  INFO[0m["+year+"] info                      [34mkey[0m=value\n")
	// buff.Reset()
}

func TestConsoleLoggerCaller(t *testing.T) {

	buff := new(bytes.Buffer)

	cLog := New()
	cLog.SetWriter(buff)
	cLog.DisplayColor(false)
	cLog.SetChannelBuffer(3)
	cLog.SetTimestampFormat(time.RFC3339)
	cLog.UseMiniTimestamp(true)
	cLog.SetANSIReset(log.Reset)

	log.SetCallerInfo(true)
	log.RegisterHandler(cLog, log.AllLevels...)

	log.Debug("debug")
	Equal(t, buff.String(), "0000  DEBUG console_test.go:382 debug\n")
	buff.Reset()

	log.Debugf("%s", "debugf")
	Equal(t, buff.String(), "0000  DEBUG console_test.go:386 debugf\n")
	buff.Reset()

	log.Info("info")
	Equal(t, buff.String(), "0000   INFO console_test.go:390 info\n")
	buff.Reset()

	log.Infof("%s", "infof")
	Equal(t, buff.String(), "0000   INFO console_test.go:394 infof\n")
	buff.Reset()

	log.Notice("notice")
	Equal(t, buff.String(), "0000 NOTICE console_test.go:398 notice\n")
	buff.Reset()

	log.Noticef("%s", "noticef")
	Equal(t, buff.String(), "0000 NOTICE console_test.go:402 noticef\n")
	buff.Reset()

	log.Warn("warn")
	Equal(t, buff.String(), "0000   WARN console_test.go:406 warn\n")
	buff.Reset()

	log.Warnf("%s", "warnf")
	Equal(t, buff.String(), "0000   WARN console_test.go:410 warnf\n")
	buff.Reset()

	log.Error("error")
	Equal(t, buff.String(), "0000  ERROR console_test.go:414 error\n")
	buff.Reset()

	log.Errorf("%s", "errorf")
	Equal(t, buff.String(), "0000  ERROR console_test.go:418 errorf\n")
	buff.Reset()

	log.Alert("alert")
	Equal(t, buff.String(), "0000  ALERT console_test.go:422 alert\n")
	buff.Reset()

	log.Alertf("%s", "alertf")
	Equal(t, buff.String(), "0000  ALERT console_test.go:426 alertf\n")
	buff.Reset()

	log.Print("print")
	Equal(t, buff.String(), "0000   INFO console_test.go:430 print\n")
	buff.Reset()

	log.Printf("%s", "printf")
	Equal(t, buff.String(), "0000   INFO console_test.go:434 printf\n")
	buff.Reset()

	log.Println("println")
	Equal(t, buff.String(), "0000   INFO console_test.go:438 println\n")
	buff.Reset()

	PanicMatches(t, func() { log.Panic("panic") }, "panic")
	Equal(t, buff.String(), "0000  PANIC console_test.go:442 panic\n")
	buff.Reset()

	PanicMatches(t, func() { log.Panicf("%s", "panicf") }, "panicf")
	Equal(t, buff.String(), "0000  PANIC console_test.go:446 panicf\n")
	buff.Reset()

	PanicMatches(t, func() { log.Panicln("panicln") }, "panicln")
	Equal(t, buff.String(), "0000  PANIC console_test.go:450 panicln\n")
	buff.Reset()

	// WithFields
	log.WithFields(log.F("key", "value")).Debug("debug")
	Equal(t, buff.String(), "0000  DEBUG console_test.go:455 debug                     key=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Debugf("%s", "debugf")
	Equal(t, buff.String(), "0000  DEBUG console_test.go:459 debugf                    key=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Info("info")
	Equal(t, buff.String(), "0000   INFO console_test.go:463 info                      key=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Infof("%s", "infof")
	Equal(t, buff.String(), "0000   INFO console_test.go:467 infof                     key=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Notice("notice")
	Equal(t, buff.String(), "0000 NOTICE console_test.go:471 notice                    key=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Noticef("%s", "noticef")
	Equal(t, buff.String(), "0000 NOTICE console_test.go:475 noticef                   key=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Warn("warn")
	Equal(t, buff.String(), "0000   WARN console_test.go:479 warn                      key=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Warnf("%s", "warnf")
	Equal(t, buff.String(), "0000   WARN console_test.go:483 warnf                     key=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Error("error")
	Equal(t, buff.String(), "0000  ERROR console_test.go:487 error                     key=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Errorf("%s", "errorf")
	Equal(t, buff.String(), "0000  ERROR console_test.go:491 errorf                    key=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Alert("alert")
	Equal(t, buff.String(), "0000  ALERT console_test.go:495 alert                     key=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Alertf("%s", "alertf")
	Equal(t, buff.String(), "0000  ALERT console_test.go:499 alertf                    key=value\n")
	buff.Reset()

	PanicMatches(t, func() { log.WithFields(log.F("key", "value")).Panicf("%s", "panicf") }, "panicf key=value")
	Equal(t, buff.String(), "0000  PANIC console_test.go:503 panicf                    key=value\n")
	buff.Reset()

	PanicMatches(t, func() { log.WithFields(log.F("key", "value")).Panic("panic") }, "panic key=value")
	Equal(t, buff.String(), "0000  PANIC console_test.go:507 panic                     key=value\n")
	buff.Reset()

	func() {
		defer log.Trace("trace").End()
	}()

	// TODO: finish up regex
	MatchRegex(t, buff.String(), "^0000\\s\\sTRACE\\sconsole_test.go:513\\strace\\.*")
	buff.Reset()

	func() {
		defer log.Tracef("tracef").End()
	}()

	// TODO: finish up regex
	MatchRegex(t, buff.String(), "^0000\\s\\sTRACE\\sconsole_test.go:521\\stracef\\.*")
	buff.Reset()

	func() {
		defer log.WithFields(log.F("key", "value")).Trace("trace").End()
	}()

	// TODO: finish up regex
	MatchRegex(t, buff.String(), "^0000\\s\\sTRACE\\sconsole_test.go:529\\strace\\.*")
	buff.Reset()

	func() {
		defer log.WithFields(log.F("key", "value")).Tracef("tracef").End()
	}()

	// TODO: finish up regex
	MatchRegex(t, buff.String(), "^0000\\s\\sTRACE\\sconsole_test.go:537\\stracef\\.*")
	buff.Reset()

	// year := time.Now().Format("2006")
	// cLog.UseMiniTimestamp(false)
	// cLog.SetTimestampFormat("2006")

	// log.Info("info")
	// Equal(t, buff.String(), "  INFO["+year+"] info\n")
	// buff.Reset()
}

func TestConsoleLoggerColorCaller(t *testing.T) {

	buff := new(bytes.Buffer)

	cLog := New()
	cLog.SetWriter(buff)
	cLog.DisplayColor(true)
	cLog.SetChannelBuffer(3)
	cLog.SetTimestampFormat(time.RFC3339)
	cLog.UseMiniTimestamp(true)

	log.SetCallerInfo(true)
	log.RegisterHandler(cLog, log.AllLevels...)

	log.Debug("debug")
	Equal(t, buff.String(), "0000 [32m DEBUG[0m console_test.go:566 debug\n")
	buff.Reset()

	log.Debugf("%s", "debugf")
	Equal(t, buff.String(), "0000 [32m DEBUG[0m console_test.go:570 debugf\n")
	buff.Reset()

	log.Info("info")
	Equal(t, buff.String(), "0000 [34m  INFO[0m console_test.go:574 info\n")
	buff.Reset()

	log.Infof("%s", "infof")
	Equal(t, buff.String(), "0000 [34m  INFO[0m console_test.go:578 infof\n")
	buff.Reset()

	log.Notice("notice")
	Equal(t, buff.String(), "0000 [36;1mNOTICE[0m console_test.go:582 notice\n")
	buff.Reset()

	log.Notice("%s", "noticef")
	Equal(t, buff.String(), "0000 [36;1mNOTICE[0m console_test.go:586 %snoticef\n")
	buff.Reset()

	log.Warn("warn")
	Equal(t, buff.String(), "0000 [33;1m  WARN[0m console_test.go:590 warn\n")
	buff.Reset()

	log.Warnf("%s", "warnf")
	Equal(t, buff.String(), "0000 [33;1m  WARN[0m console_test.go:594 warnf\n")
	buff.Reset()

	log.Error("error")
	Equal(t, buff.String(), "0000 [31;1m ERROR[0m console_test.go:598 error\n")
	buff.Reset()

	log.Errorf("%s", "errorf")
	Equal(t, buff.String(), "0000 [31;1m ERROR[0m console_test.go:602 errorf\n")
	buff.Reset()

	log.Alert("alert")
	Equal(t, buff.String(), "0000 [31m[4m ALERT[0m console_test.go:606 alert\n")
	buff.Reset()

	log.Alertf("%s", "alertf")
	Equal(t, buff.String(), "0000 [31m[4m ALERT[0m console_test.go:610 alertf\n")
	buff.Reset()

	log.Print("print")
	Equal(t, buff.String(), "0000 [34m  INFO[0m console_test.go:614 print\n")
	buff.Reset()

	log.Printf("%s", "printf")
	Equal(t, buff.String(), "0000 [34m  INFO[0m console_test.go:618 printf\n")
	buff.Reset()

	log.Println("println")
	Equal(t, buff.String(), "0000 [34m  INFO[0m console_test.go:622 println\n")
	buff.Reset()

	PanicMatches(t, func() { log.Panic("panic") }, "panic")
	Equal(t, buff.String(), "0000 [31m PANIC[0m console_test.go:626 panic\n")
	buff.Reset()

	PanicMatches(t, func() { log.Panicf("%s", "panicf") }, "panicf")
	Equal(t, buff.String(), "0000 [31m PANIC[0m console_test.go:630 panicf\n")
	buff.Reset()

	PanicMatches(t, func() { log.Panicln("panicln") }, "panicln")
	Equal(t, buff.String(), "0000 [31m PANIC[0m console_test.go:634 panicln\n")
	buff.Reset()

	// WithFields
	log.WithFields(log.F("key", "value")).Debug("debug")
	Equal(t, buff.String(), "0000 [32m DEBUG[0m console_test.go:639 debug                     [32mkey[0m=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Debugf("%s", "debugf")
	Equal(t, buff.String(), "0000 [32m DEBUG[0m console_test.go:643 debugf                    [32mkey[0m=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Info("info")
	Equal(t, buff.String(), "0000 [34m  INFO[0m console_test.go:647 info                      [34mkey[0m=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Infof("%s", "infof")
	Equal(t, buff.String(), "0000 [34m  INFO[0m console_test.go:651 infof                     [34mkey[0m=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Notice("notice")
	Equal(t, buff.String(), "0000 [36;1mNOTICE[0m console_test.go:655 notice                    [36;1mkey[0m=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Noticef("%s", "noticef")
	Equal(t, buff.String(), "0000 [36;1mNOTICE[0m console_test.go:659 noticef                   [36;1mkey[0m=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Warn("warn")
	Equal(t, buff.String(), "0000 [33;1m  WARN[0m console_test.go:663 warn                      [33;1mkey[0m=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Warnf("%s", "warnf")
	Equal(t, buff.String(), "0000 [33;1m  WARN[0m console_test.go:667 warnf                     [33;1mkey[0m=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Error("error")
	Equal(t, buff.String(), "0000 [31;1m ERROR[0m console_test.go:671 error                     [31;1mkey[0m=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Errorf("%s", "errorf")
	Equal(t, buff.String(), "0000 [31;1m ERROR[0m console_test.go:675 errorf                    [31;1mkey[0m=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Alert("alert")
	Equal(t, buff.String(), "0000 [31m[4m ALERT[0m console_test.go:679 alert                     [31m[4mkey[0m=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Alertf("%s", "alertf")
	Equal(t, buff.String(), "0000 [31m[4m ALERT[0m console_test.go:683 alertf                    [31m[4mkey[0m=value\n")
	buff.Reset()

	PanicMatches(t, func() { log.WithFields(log.F("key", "value")).Panicf("%s", "panicf") }, "panicf key=value")
	Equal(t, buff.String(), "0000 [31m PANIC[0m console_test.go:687 panicf                    [31mkey[0m=value\n")
	buff.Reset()

	PanicMatches(t, func() { log.WithFields(log.F("key", "value")).Panic("panic") }, "panic key=value")
	Equal(t, buff.String(), "0000 [31m PANIC[0m console_test.go:691 panic                     [31mkey[0m=value\n")
	buff.Reset()

	cLog.SetLevelColor(log.DebugLevel, log.LightGreen)

	log.Debug("debug")
	Equal(t, buff.String(), "0000 [32;1m DEBUG[0m console_test.go:697 debug\n")
	buff.Reset()
}

func TestConsoleLoggerColorCallerTimeFormat(t *testing.T) {

	year := time.Now().Format("2006")
	buff := new(bytes.Buffer)

	cLog := New()
	cLog.SetWriter(buff)
	cLog.DisplayColor(true)
	cLog.SetChannelBuffer(3)
	cLog.SetTimestampFormat("2006")

	log.SetCallerInfo(true)
	log.RegisterHandler(cLog, log.AllLevels...)

	log.Debug("debug")
	Equal(t, buff.String(), year+" [32m DEBUG[0m console_test.go:716 debug\n")
	buff.Reset()
}
