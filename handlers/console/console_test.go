package console

import (
	"bytes"
	"io"
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
	cLog.SetBuffersAndWorkers(3, 3)
	cLog.SetTimestampFormat("MST")

	log.RegisterHandler(cLog, log.AllLevels...)

	log.Debug("debug")
	Equal(t, buff.String(), "UTC  DEBUG debug\n")
	buff.Reset()

	log.Debugf("%s", "debugf")
	Equal(t, buff.String(), "UTC  DEBUG debugf\n")
	buff.Reset()

	log.Info("info")
	Equal(t, buff.String(), "UTC   INFO info\n")
	buff.Reset()

	log.Infof("%s", "infof")
	Equal(t, buff.String(), "UTC   INFO infof\n")
	buff.Reset()

	log.Notice("notice")
	Equal(t, buff.String(), "UTC NOTICE notice\n")
	buff.Reset()

	log.Noticef("%s", "noticef")
	Equal(t, buff.String(), "UTC NOTICE noticef\n")
	buff.Reset()

	log.Warn("warn")
	Equal(t, buff.String(), "UTC   WARN console_test.go:58 warn\n")
	buff.Reset()

	log.Warnf("%s", "warnf")
	Equal(t, buff.String(), "UTC   WARN console_test.go:62 warnf\n")
	buff.Reset()

	log.Error("error")
	Equal(t, buff.String(), "UTC  ERROR console_test.go:66 error\n")
	buff.Reset()

	log.Errorf("%s", "errorf")
	Equal(t, buff.String(), "UTC  ERROR console_test.go:70 errorf\n")
	buff.Reset()

	log.Alert("alert")
	Equal(t, buff.String(), "UTC  ALERT console_test.go:74 alert\n")
	buff.Reset()

	log.Alertf("%s", "alertf")
	Equal(t, buff.String(), "UTC  ALERT console_test.go:78 alertf\n")
	buff.Reset()

	log.Print("print")
	Equal(t, buff.String(), "UTC   INFO print\n")
	buff.Reset()

	log.Printf("%s", "printf")
	Equal(t, buff.String(), "UTC   INFO printf\n")
	buff.Reset()

	log.Println("println")
	Equal(t, buff.String(), "UTC   INFO println\n")
	buff.Reset()

	PanicMatches(t, func() { log.Panic("panic") }, "panic")
	Equal(t, buff.String(), "UTC  PANIC console_test.go:94 panic\n")
	buff.Reset()

	PanicMatches(t, func() { log.Panicf("%s", "panicf") }, "panicf")
	Equal(t, buff.String(), "UTC  PANIC console_test.go:98 panicf\n")
	buff.Reset()

	PanicMatches(t, func() { log.Panicln("panicln") }, "panicln")
	Equal(t, buff.String(), "UTC  PANIC console_test.go:102 panicln\n")
	buff.Reset()

	// WithFields
	log.WithFields(log.F("key", "value")).Debug("debug")
	Equal(t, buff.String(), "UTC  DEBUG debug                     key=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Debugf("%s", "debugf")
	Equal(t, buff.String(), "UTC  DEBUG debugf                    key=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Info("info")
	Equal(t, buff.String(), "UTC   INFO info                      key=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Infof("%s", "infof")
	Equal(t, buff.String(), "UTC   INFO infof                     key=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Notice("notice")
	Equal(t, buff.String(), "UTC NOTICE notice                    key=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Noticef("%s", "noticef")
	Equal(t, buff.String(), "UTC NOTICE noticef                   key=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Warn("warn")
	Equal(t, buff.String(), "UTC   WARN console_test.go:131 warn                      key=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Warnf("%s", "warnf")
	Equal(t, buff.String(), "UTC   WARN console_test.go:135 warnf                     key=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Error("error")
	Equal(t, buff.String(), "UTC  ERROR console_test.go:139 error                     key=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Errorf("%s", "errorf")
	Equal(t, buff.String(), "UTC  ERROR console_test.go:143 errorf                    key=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Alert("alert")
	Equal(t, buff.String(), "UTC  ALERT console_test.go:147 alert                     key=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Alertf("%s", "alertf")
	Equal(t, buff.String(), "UTC  ALERT console_test.go:151 alertf                    key=value\n")
	buff.Reset()

	PanicMatches(t, func() { log.WithFields(log.F("key", "value")).Panicf("%s", "panicf") }, "panicf key=value")
	Equal(t, buff.String(), "UTC  PANIC console_test.go:155 panicf                    key=value\n")
	buff.Reset()

	PanicMatches(t, func() { log.WithFields(log.F("key", "value")).Panic("panic") }, "panic key=value")
	Equal(t, buff.String(), "UTC  PANIC console_test.go:159 panic                     key=value\n")
	buff.Reset()

	func() {
		defer log.Trace("trace").End()
	}()

	MatchRegex(t, buff.String(), "^UTC\\s\\sTRACE\\strace\\.*")
	buff.Reset()

	func() {
		defer log.Tracef("tracef").End()
	}()

	MatchRegex(t, buff.String(), "^UTC\\s\\sTRACE\\stracef\\.*")
	buff.Reset()

	func() {
		defer log.WithFields(log.F("key", "value")).Trace("trace").End()
	}()

	MatchRegex(t, buff.String(), "^UTC\\s\\sTRACE\\strace\\.*")
	buff.Reset()

	func() {
		defer log.WithFields(log.F("key", "value")).Tracef("tracef").End()
	}()

	MatchRegex(t, buff.String(), "^UTC\\s\\sTRACE\\stracef\\.*")
	buff.Reset()
}

func TestConsoleLoggerColor(t *testing.T) {

	buff := new(bytes.Buffer)

	cLog := New()
	cLog.SetWriter(buff)
	cLog.DisplayColor(true)
	cLog.SetBuffersAndWorkers(3, 3)
	cLog.SetTimestampFormat("MST")

	log.RegisterHandler(cLog, log.AllLevels...)

	log.Debug("debug")
	Equal(t, buff.String(), "UTC [32m DEBUG[0m debug\n")
	buff.Reset()

	log.Debugf("%s", "debugf")
	Equal(t, buff.String(), "UTC [32m DEBUG[0m debugf\n")
	buff.Reset()

	log.Info("info")
	Equal(t, buff.String(), "UTC [34m  INFO[0m info\n")
	buff.Reset()

	log.Infof("%s", "infof")
	Equal(t, buff.String(), "UTC [34m  INFO[0m infof\n")
	buff.Reset()

	log.Notice("notice")
	Equal(t, buff.String(), "UTC [36;1mNOTICE[0m notice\n")
	buff.Reset()

	log.Notice("%s", "noticef")
	Equal(t, buff.String(), "UTC [36;1mNOTICE[0m %snoticef\n")
	buff.Reset()

	log.Warn("warn")
	Equal(t, buff.String(), "UTC [33;1m  WARN[0m console_test.go:228 warn\n")
	buff.Reset()

	log.Warnf("%s", "warnf")
	Equal(t, buff.String(), "UTC [33;1m  WARN[0m console_test.go:232 warnf\n")
	buff.Reset()

	log.Error("error")
	Equal(t, buff.String(), "UTC [31;1m ERROR[0m console_test.go:236 error\n")
	buff.Reset()

	log.Errorf("%s", "errorf")
	Equal(t, buff.String(), "UTC [31;1m ERROR[0m console_test.go:240 errorf\n")
	buff.Reset()

	log.Alert("alert")
	Equal(t, buff.String(), "UTC [31m[4m ALERT[0m console_test.go:244 alert\n")
	buff.Reset()

	log.Alertf("%s", "alertf")
	Equal(t, buff.String(), "UTC [31m[4m ALERT[0m console_test.go:248 alertf\n")
	buff.Reset()

	log.Print("print")
	Equal(t, buff.String(), "UTC [34m  INFO[0m print\n")
	buff.Reset()

	log.Printf("%s", "printf")
	Equal(t, buff.String(), "UTC [34m  INFO[0m printf\n")
	buff.Reset()

	log.Println("println")
	Equal(t, buff.String(), "UTC [34m  INFO[0m println\n")
	buff.Reset()

	PanicMatches(t, func() { log.Panic("panic") }, "panic")
	Equal(t, buff.String(), "UTC [31m PANIC[0m console_test.go:264 panic\n")
	buff.Reset()

	PanicMatches(t, func() { log.Panicf("%s", "panicf") }, "panicf")
	Equal(t, buff.String(), "UTC [31m PANIC[0m console_test.go:268 panicf\n")
	buff.Reset()

	PanicMatches(t, func() { log.Panicln("panicln") }, "panicln")
	Equal(t, buff.String(), "UTC [31m PANIC[0m console_test.go:272 panicln\n")
	buff.Reset()

	// WithFields
	log.WithFields(log.F("key", "value")).Debug("debug")
	Equal(t, buff.String(), "UTC [32m DEBUG[0m debug                     [32mkey[0m=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Debugf("%s", "debugf")
	Equal(t, buff.String(), "UTC [32m DEBUG[0m debugf                    [32mkey[0m=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Info("info")
	Equal(t, buff.String(), "UTC [34m  INFO[0m info                      [34mkey[0m=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Infof("%s", "infof")
	Equal(t, buff.String(), "UTC [34m  INFO[0m infof                     [34mkey[0m=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Notice("notice")
	Equal(t, buff.String(), "UTC [36;1mNOTICE[0m notice                    [36;1mkey[0m=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Noticef("%s", "noticef")
	Equal(t, buff.String(), "UTC [36;1mNOTICE[0m noticef                   [36;1mkey[0m=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Warn("warn")
	Equal(t, buff.String(), "UTC [33;1m  WARN[0m console_test.go:301 warn                      [33;1mkey[0m=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Warnf("%s", "warnf")
	Equal(t, buff.String(), "UTC [33;1m  WARN[0m console_test.go:305 warnf                     [33;1mkey[0m=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Error("error")
	Equal(t, buff.String(), "UTC [31;1m ERROR[0m console_test.go:309 error                     [31;1mkey[0m=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Errorf("%s", "errorf")
	Equal(t, buff.String(), "UTC [31;1m ERROR[0m console_test.go:313 errorf                    [31;1mkey[0m=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Alert("alert")
	Equal(t, buff.String(), "UTC [31m[4m ALERT[0m console_test.go:317 alert                     [31m[4mkey[0m=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Alertf("%s", "alertf")
	Equal(t, buff.String(), "UTC [31m[4m ALERT[0m console_test.go:321 alertf                    [31m[4mkey[0m=value\n")
	buff.Reset()

	PanicMatches(t, func() { log.WithFields(log.F("key", "value")).Panicf("%s", "panicf") }, "panicf key=value")
	Equal(t, buff.String(), "UTC [31m PANIC[0m console_test.go:325 panicf                    [31mkey[0m=value\n")
	buff.Reset()

	PanicMatches(t, func() { log.WithFields(log.F("key", "value")).Panic("panic") }, "panic key=value")
	Equal(t, buff.String(), "UTC [31m PANIC[0m console_test.go:329 panic                     [31mkey[0m=value\n")
	buff.Reset()
}

func TestConsoleLoggerCaller(t *testing.T) {

	buff := new(bytes.Buffer)

	cLog := New()
	cLog.SetWriter(buff)
	cLog.DisplayColor(false)
	cLog.SetBuffersAndWorkers(3, 3)
	cLog.SetTimestampFormat("MST")

	log.SetCallerInfo(true)
	log.RegisterHandler(cLog, log.AllLevels...)

	log.Debug("debug")
	Equal(t, buff.String(), "UTC  DEBUG console_test.go:347 debug\n")
	buff.Reset()

	log.Debugf("%s", "debugf")
	Equal(t, buff.String(), "UTC  DEBUG console_test.go:351 debugf\n")
	buff.Reset()

	log.Info("info")
	Equal(t, buff.String(), "UTC   INFO console_test.go:355 info\n")
	buff.Reset()

	log.Infof("%s", "infof")
	Equal(t, buff.String(), "UTC   INFO console_test.go:359 infof\n")
	buff.Reset()

	log.Notice("notice")
	Equal(t, buff.String(), "UTC NOTICE console_test.go:363 notice\n")
	buff.Reset()

	log.Noticef("%s", "noticef")
	Equal(t, buff.String(), "UTC NOTICE console_test.go:367 noticef\n")
	buff.Reset()

	log.Warn("warn")
	Equal(t, buff.String(), "UTC   WARN console_test.go:371 warn\n")
	buff.Reset()

	log.Warnf("%s", "warnf")
	Equal(t, buff.String(), "UTC   WARN console_test.go:375 warnf\n")
	buff.Reset()

	log.Error("error")
	Equal(t, buff.String(), "UTC  ERROR console_test.go:379 error\n")
	buff.Reset()

	log.Errorf("%s", "errorf")
	Equal(t, buff.String(), "UTC  ERROR console_test.go:383 errorf\n")
	buff.Reset()

	log.Alert("alert")
	Equal(t, buff.String(), "UTC  ALERT console_test.go:387 alert\n")
	buff.Reset()

	log.Alertf("%s", "alertf")
	Equal(t, buff.String(), "UTC  ALERT console_test.go:391 alertf\n")
	buff.Reset()

	log.Print("print")
	Equal(t, buff.String(), "UTC   INFO console_test.go:395 print\n")
	buff.Reset()

	log.Printf("%s", "printf")
	Equal(t, buff.String(), "UTC   INFO console_test.go:399 printf\n")
	buff.Reset()

	log.Println("println")
	Equal(t, buff.String(), "UTC   INFO console_test.go:403 println\n")
	buff.Reset()

	PanicMatches(t, func() { log.Panic("panic") }, "panic")
	Equal(t, buff.String(), "UTC  PANIC console_test.go:407 panic\n")
	buff.Reset()

	PanicMatches(t, func() { log.Panicf("%s", "panicf") }, "panicf")
	Equal(t, buff.String(), "UTC  PANIC console_test.go:411 panicf\n")
	buff.Reset()

	PanicMatches(t, func() { log.Panicln("panicln") }, "panicln")
	Equal(t, buff.String(), "UTC  PANIC console_test.go:415 panicln\n")
	buff.Reset()

	// WithFields
	log.WithFields(log.F("key", "value")).Debug("debug")
	Equal(t, buff.String(), "UTC  DEBUG console_test.go:420 debug                     key=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Debugf("%s", "debugf")
	Equal(t, buff.String(), "UTC  DEBUG console_test.go:424 debugf                    key=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Info("info")
	Equal(t, buff.String(), "UTC   INFO console_test.go:428 info                      key=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Infof("%s", "infof")
	Equal(t, buff.String(), "UTC   INFO console_test.go:432 infof                     key=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Notice("notice")
	Equal(t, buff.String(), "UTC NOTICE console_test.go:436 notice                    key=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Noticef("%s", "noticef")
	Equal(t, buff.String(), "UTC NOTICE console_test.go:440 noticef                   key=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Warn("warn")
	Equal(t, buff.String(), "UTC   WARN console_test.go:444 warn                      key=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Warnf("%s", "warnf")
	Equal(t, buff.String(), "UTC   WARN console_test.go:448 warnf                     key=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Error("error")
	Equal(t, buff.String(), "UTC  ERROR console_test.go:452 error                     key=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Errorf("%s", "errorf")
	Equal(t, buff.String(), "UTC  ERROR console_test.go:456 errorf                    key=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Alert("alert")
	Equal(t, buff.String(), "UTC  ALERT console_test.go:460 alert                     key=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Alertf("%s", "alertf")
	Equal(t, buff.String(), "UTC  ALERT console_test.go:464 alertf                    key=value\n")
	buff.Reset()

	PanicMatches(t, func() { log.WithFields(log.F("key", "value")).Panicf("%s", "panicf") }, "panicf key=value")
	Equal(t, buff.String(), "UTC  PANIC console_test.go:468 panicf                    key=value\n")
	buff.Reset()

	PanicMatches(t, func() { log.WithFields(log.F("key", "value")).Panic("panic") }, "panic key=value")
	Equal(t, buff.String(), "UTC  PANIC console_test.go:472 panic                     key=value\n")
	buff.Reset()

	func() {
		defer log.Trace("trace").End()
	}()

	// TODO: finish up regex
	MatchRegex(t, buff.String(), "^UTC\\s\\sTRACE\\sconsole_test.go:478\\strace\\.*")
	buff.Reset()

	func() {
		defer log.Tracef("tracef").End()
	}()

	// TODO: finish up regex
	MatchRegex(t, buff.String(), "^UTC\\s\\sTRACE\\sconsole_test.go:486\\stracef\\.*")
	buff.Reset()

	func() {
		defer log.WithFields(log.F("key", "value")).Trace("trace").End()
	}()

	// TODO: finish up regex
	MatchRegex(t, buff.String(), "^UTC\\s\\sTRACE\\sconsole_test.go:494\\strace\\.*")
	buff.Reset()

	func() {
		defer log.WithFields(log.F("key", "value")).Tracef("tracef").End()
	}()

	// TODO: finish up regex
	MatchRegex(t, buff.String(), "^UTC\\s\\sTRACE\\sconsole_test.go:502\\stracef\\.*")
	buff.Reset()
}

func TestConsoleLoggerColorCaller(t *testing.T) {

	buff := new(bytes.Buffer)

	cLog := New()
	cLog.SetWriter(buff)
	cLog.DisplayColor(true)
	cLog.SetBuffersAndWorkers(3, 3)
	cLog.SetTimestampFormat("MST")

	log.SetCallerInfo(true)
	log.RegisterHandler(cLog, log.AllLevels...)

	log.Debug("debug")
	Equal(t, buff.String(), "UTC [32m DEBUG[0m console_test.go:522 debug\n")
	buff.Reset()

	log.Debugf("%s", "debugf")
	Equal(t, buff.String(), "UTC [32m DEBUG[0m console_test.go:526 debugf\n")
	buff.Reset()

	log.Info("info")
	Equal(t, buff.String(), "UTC [34m  INFO[0m console_test.go:530 info\n")
	buff.Reset()

	log.Infof("%s", "infof")
	Equal(t, buff.String(), "UTC [34m  INFO[0m console_test.go:534 infof\n")
	buff.Reset()

	log.Notice("notice")
	Equal(t, buff.String(), "UTC [36;1mNOTICE[0m console_test.go:538 notice\n")
	buff.Reset()

	log.Notice("%s", "noticef")
	Equal(t, buff.String(), "UTC [36;1mNOTICE[0m console_test.go:542 %snoticef\n")
	buff.Reset()

	log.Warn("warn")
	Equal(t, buff.String(), "UTC [33;1m  WARN[0m console_test.go:546 warn\n")
	buff.Reset()

	log.Warnf("%s", "warnf")
	Equal(t, buff.String(), "UTC [33;1m  WARN[0m console_test.go:550 warnf\n")
	buff.Reset()

	log.Error("error")
	Equal(t, buff.String(), "UTC [31;1m ERROR[0m console_test.go:554 error\n")
	buff.Reset()

	log.Errorf("%s", "errorf")
	Equal(t, buff.String(), "UTC [31;1m ERROR[0m console_test.go:558 errorf\n")
	buff.Reset()

	log.Alert("alert")
	Equal(t, buff.String(), "UTC [31m[4m ALERT[0m console_test.go:562 alert\n")
	buff.Reset()

	log.Alertf("%s", "alertf")
	Equal(t, buff.String(), "UTC [31m[4m ALERT[0m console_test.go:566 alertf\n")
	buff.Reset()

	log.Print("print")
	Equal(t, buff.String(), "UTC [34m  INFO[0m console_test.go:570 print\n")
	buff.Reset()

	log.Printf("%s", "printf")
	Equal(t, buff.String(), "UTC [34m  INFO[0m console_test.go:574 printf\n")
	buff.Reset()

	log.Println("println")
	Equal(t, buff.String(), "UTC [34m  INFO[0m console_test.go:578 println\n")
	buff.Reset()

	PanicMatches(t, func() { log.Panic("panic") }, "panic")
	Equal(t, buff.String(), "UTC [31m PANIC[0m console_test.go:582 panic\n")
	buff.Reset()

	PanicMatches(t, func() { log.Panicf("%s", "panicf") }, "panicf")
	Equal(t, buff.String(), "UTC [31m PANIC[0m console_test.go:586 panicf\n")
	buff.Reset()

	PanicMatches(t, func() { log.Panicln("panicln") }, "panicln")
	Equal(t, buff.String(), "UTC [31m PANIC[0m console_test.go:590 panicln\n")
	buff.Reset()

	// WithFields
	log.WithFields(log.F("key", "value")).Debug("debug")
	Equal(t, buff.String(), "UTC [32m DEBUG[0m console_test.go:595 debug                     [32mkey[0m=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Debugf("%s", "debugf")
	Equal(t, buff.String(), "UTC [32m DEBUG[0m console_test.go:599 debugf                    [32mkey[0m=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Info("info")
	Equal(t, buff.String(), "UTC [34m  INFO[0m console_test.go:603 info                      [34mkey[0m=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Infof("%s", "infof")
	Equal(t, buff.String(), "UTC [34m  INFO[0m console_test.go:607 infof                     [34mkey[0m=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Notice("notice")
	Equal(t, buff.String(), "UTC [36;1mNOTICE[0m console_test.go:611 notice                    [36;1mkey[0m=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Noticef("%s", "noticef")
	Equal(t, buff.String(), "UTC [36;1mNOTICE[0m console_test.go:615 noticef                   [36;1mkey[0m=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Warn("warn")
	Equal(t, buff.String(), "UTC [33;1m  WARN[0m console_test.go:619 warn                      [33;1mkey[0m=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Warnf("%s", "warnf")
	Equal(t, buff.String(), "UTC [33;1m  WARN[0m console_test.go:623 warnf                     [33;1mkey[0m=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Error("error")
	Equal(t, buff.String(), "UTC [31;1m ERROR[0m console_test.go:627 error                     [31;1mkey[0m=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Errorf("%s", "errorf")
	Equal(t, buff.String(), "UTC [31;1m ERROR[0m console_test.go:631 errorf                    [31;1mkey[0m=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Alert("alert")
	Equal(t, buff.String(), "UTC [31m[4m ALERT[0m console_test.go:635 alert                     [31m[4mkey[0m=value\n")
	buff.Reset()

	log.WithFields(log.F("key", "value")).Alertf("%s", "alertf")
	Equal(t, buff.String(), "UTC [31m[4m ALERT[0m console_test.go:639 alertf                    [31m[4mkey[0m=value\n")
	buff.Reset()

	PanicMatches(t, func() { log.WithFields(log.F("key", "value")).Panicf("%s", "panicf") }, "panicf key=value")
	Equal(t, buff.String(), "UTC [31m PANIC[0m console_test.go:643 panicf                    [31mkey[0m=value\n")
	buff.Reset()

	PanicMatches(t, func() { log.WithFields(log.F("key", "value")).Panic("panic") }, "panic key=value")
	Equal(t, buff.String(), "UTC [31m PANIC[0m console_test.go:647 panic                     [31mkey[0m=value\n")
	buff.Reset()
}

func TestConsoleLoggerColorCallerTimeFormat(t *testing.T) {

	year := time.Now().Format("2006")
	buff := new(bytes.Buffer)

	cLog := New()
	cLog.SetWriter(buff)
	cLog.DisplayColor(true)
	cLog.SetBuffersAndWorkers(3, 3)
	cLog.SetTimestampFormat("2006")

	log.SetCallerInfo(true)
	log.RegisterHandler(cLog, log.AllLevels...)

	log.Debug("debug")
	Equal(t, buff.String(), year+" [32m DEBUG[0m console_test.go:666 debug\n")
	buff.Reset()
}

func TestBadWorkerCount(t *testing.T) {

	year := time.Now().Format("2006")
	buff := new(bytes.Buffer)
	cLog := New()
	cLog.SetWriter(buff)
	cLog.SetTimestampFormat("2006")
	cLog.SetBuffersAndWorkers(3, 0)

	log.RegisterHandler(cLog, log.AllLevels...)

	log.Debug("debug")
	Equal(t, buff.String(), year+" [32m DEBUG[0m console_test.go:682 debug\n")
	buff.Reset()
}

func TestCustomFormatFunc(t *testing.T) {

	buff := new(bytes.Buffer)
	cLog := New()
	cLog.SetWriter(buff)
	cLog.SetTimestampFormat("2006")
	cLog.SetBuffersAndWorkers(3, 2)
	cLog.SetFormatFunc(func() Formatter {
		return func(e *log.Entry) io.WriterTo {
			b := new(bytes.Buffer)
			b.WriteString(e.Message)
			return b
		}
	})

	log.RegisterHandler(cLog, log.AllLevels...)

	log.Debug("debug")
	Equal(t, buff.String(), "debug")
	buff.Reset()
}
