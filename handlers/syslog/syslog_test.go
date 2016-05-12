package syslog

import (
	stdsyslog "log/syslog"
	"net"
	"strings"
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

func read(conn *net.UDPConn) (string, error) {

	bytes := make([]byte, 1024)
	read, _, err := conn.ReadFromUDP(bytes)
	if err != nil {
		return "", err
	}

	return string(bytes[0:read]), err
}

func hasString(conn *net.UDPConn, s string) bool {
	read, _ := read(conn)
	return strings.Contains(read, s)
}

func TestBadAddress(t *testing.T) {
	sLog, err := New("udp", "255.255.255.67", stdsyslog.LOG_DEBUG, "")
	Equal(t, sLog, nil)
	NotEqual(t, err, nil)
}

func TestSyslogLogger(t *testing.T) {

	year := time.Now().Format("2006")

	addr, err := net.ResolveUDPAddr("udp", ":2000")
	Equal(t, err, nil)

	conn, err := net.ListenUDP("udp", addr)
	Equal(t, err, nil)
	defer conn.Close()

	sLog, err := New("udp", "127.0.0.1:2000", stdsyslog.LOG_DEBUG, "")
	Equal(t, err, nil)

	sLog.DisplayColor(false)
	sLog.SetBuffersAndWorkers(3, 3)
	sLog.SetTimestampFormat("2006")

	log.RegisterHandler(sLog, log.AllLevels...)

	log.Debug("debug")
	Equal(t, hasString(conn, year+"  DEBUG debug"), true)

	log.Debugf("%s", "debugf")
	Equal(t, hasString(conn, year+"  DEBUG debugf"), true)

	log.Info("info")
	Equal(t, hasString(conn, year+"   INFO info"), true)

	log.Infof("%s", "infof")
	Equal(t, hasString(conn, year+"   INFO infof"), true)

	log.Notice("notice")
	Equal(t, hasString(conn, year+" NOTICE notice"), true)

	log.Noticef("%s", "noticef")
	Equal(t, hasString(conn, year+" NOTICE noticef"), true)

	log.Warn("warn")
	Equal(t, hasString(conn, year+"   WARN syslog_test.go:84 warn"), true)

	log.Warnf("%s", "warnf")
	Equal(t, hasString(conn, year+"   WARN syslog_test.go:87 warnf"), true)

	log.Error("error")
	Equal(t, hasString(conn, year+"  ERROR syslog_test.go:90 error"), true)

	log.Errorf("%s", "errorf")
	Equal(t, hasString(conn, year+"  ERROR syslog_test.go:93 errorf"), true)

	log.Alert("alert")
	Equal(t, hasString(conn, year+"  ALERT syslog_test.go:96 alert"), true)

	log.Alertf("%s", "alertf")
	Equal(t, hasString(conn, year+"  ALERT syslog_test.go:99 alertf"), true)

	log.Print("print")
	Equal(t, hasString(conn, year+"   INFO print"), true)

	log.Printf("%s", "printf")
	Equal(t, hasString(conn, year+"   INFO printf"), true)

	log.Println("println")
	Equal(t, hasString(conn, year+"   INFO println"), true)

	PanicMatches(t, func() { log.Panic("panic") }, "panic")
	Equal(t, hasString(conn, year+"  PANIC syslog_test.go:111 panic"), true)

	PanicMatches(t, func() { log.Panicf("%s", "panicf") }, "panicf")
	Equal(t, hasString(conn, year+"  PANIC syslog_test.go:114 panicf"), true)

	PanicMatches(t, func() { log.Panicln("panicln") }, "panicln")
	Equal(t, hasString(conn, year+"  PANIC syslog_test.go:117 panicln"), true)

	// WithFields
	log.WithFields(log.F("key", "value")).Debug("debug")
	Equal(t, hasString(conn, year+"  DEBUG debug key=value"), true)

	log.WithFields(log.F("key", "value")).Debugf("%s", "debugf")
	Equal(t, hasString(conn, year+"  DEBUG debugf key=value"), true)

	log.WithFields(log.F("key", "value")).Info("info")
	Equal(t, hasString(conn, year+"   INFO info key=value"), true)

	log.WithFields(log.F("key", "value")).Infof("%s", "infof")
	Equal(t, hasString(conn, year+"   INFO infof key=value"), true)

	log.WithFields(log.F("key", "value")).Notice("notice")
	Equal(t, hasString(conn, year+" NOTICE notice key=value"), true)

	log.WithFields(log.F("key", "value")).Noticef("%s", "noticef")
	Equal(t, hasString(conn, year+" NOTICE noticef key=value"), true)

	log.WithFields(log.F("key", "value")).Warn("warn")
	Equal(t, hasString(conn, year+"   WARN syslog_test.go:139 warn key=value"), true)

	log.WithFields(log.F("key", "value")).Warnf("%s", "warnf")
	Equal(t, hasString(conn, year+"   WARN syslog_test.go:142 warnf key=value"), true)

	log.WithFields(log.F("key", "value")).Error("error")
	Equal(t, hasString(conn, year+"  ERROR syslog_test.go:145 error key=value"), true)

	log.WithFields(log.F("key", "value")).Errorf("%s", "errorf")
	Equal(t, hasString(conn, year+"  ERROR syslog_test.go:148 errorf key=value"), true)

	log.WithFields(log.F("key", "value")).Alert("alert")
	Equal(t, hasString(conn, year+"  ALERT syslog_test.go:151 alert key=value"), true)

	log.WithFields(log.F("key", "value")).Alertf("%s", "alertf")
	Equal(t, hasString(conn, year+"  ALERT syslog_test.go:154 alertf key=value"), true)

	PanicMatches(t, func() { log.WithFields(log.F("key", "value")).Panicf("%s", "panicf") }, "panicf key=value")
	Equal(t, hasString(conn, year+"  PANIC syslog_test.go:157 panicf key=value"), true)

	PanicMatches(t, func() { log.WithFields(log.F("key", "value")).Panic("panic") }, "panic key=value")
	Equal(t, hasString(conn, year+"  PANIC syslog_test.go:160 panic key=value"), true)

	func() {
		defer log.Trace("trace").End()
	}()

	Equal(t, hasString(conn, year+"  TRACE trace"), true)

	func() {
		defer log.Tracef("tracef").End()
	}()

	Equal(t, hasString(conn, year+"  TRACE tracef"), true)

	func() {
		defer log.WithFields(log.F("key", "value")).Trace("trace").End()
	}()

	Equal(t, hasString(conn, year+"  TRACE trace"), true)

	func() {
		defer log.WithFields(log.F("key", "value")).Tracef("tracef").End()
	}()

	Equal(t, hasString(conn, year+"  TRACE tracef"), true)

	e := &log.Entry{
		Level:     log.FatalLevel,
		Message:   "fatal",
		Timestamp: time.Now(),
		Line:      3,
		File:      "fake",
	}

	log.HandleEntry(e)

	Equal(t, hasString(conn, year+"  FATAL fake:3 fatal"), true)
}

func TestSyslogLoggerColor(t *testing.T) {

	year := time.Now().Format("2006")

	addr, err := net.ResolveUDPAddr("udp", ":2001")
	Equal(t, err, nil)

	conn, err := net.ListenUDP("udp", addr)
	Equal(t, err, nil)
	defer conn.Close()

	sLog, err := New("udp", "127.0.0.1:2001", stdsyslog.LOG_DEBUG, "")
	Equal(t, err, nil)

	sLog.DisplayColor(true)
	sLog.SetBuffersAndWorkers(3, 3)
	sLog.SetTimestampFormat("2006")

	log.RegisterHandler(sLog, log.AllLevels...)

	log.Debug("debug")
	Equal(t, hasString(conn, year+" [32m DEBUG[0m debug"), true)

	log.Debugf("%s", "debugf")
	Equal(t, hasString(conn, year+" [32m DEBUG[0m debugf"), true)

	log.Info("info")
	Equal(t, hasString(conn, year+" [34m  INFO[0m info"), true)

	log.Infof("%s", "infof")
	Equal(t, hasString(conn, year+" [34m  INFO[0m info"), true)

	log.Notice("notice")
	Equal(t, hasString(conn, year+" [36;1mNOTICE[0m notice"), true)

	log.Noticef("%s", "noticef")
	Equal(t, hasString(conn, year+" [36;1mNOTICE[0m noticef"), true)

	log.Warn("warn")
	Equal(t, hasString(conn, year+" [33;1m  WARN[0m syslog_test.go:238 warn"), true)

	log.Warnf("%s", "warnf")
	Equal(t, hasString(conn, year+" [33;1m  WARN[0m syslog_test.go:241 warnf"), true)

	log.Error("error")
	Equal(t, hasString(conn, year+" [31;1m ERROR[0m syslog_test.go:244 error"), true)

	log.Errorf("%s", "errorf")
	Equal(t, hasString(conn, year+" [31;1m ERROR[0m syslog_test.go:247 errorf"), true)

	log.Alert("alert")
	Equal(t, hasString(conn, year+" [31m[4m ALERT[0m syslog_test.go:250 alert"), true)

	log.Alertf("%s", "alertf")
	Equal(t, hasString(conn, year+" [31m[4m ALERT[0m syslog_test.go:253 alertf"), true)

	log.Print("print")
	Equal(t, hasString(conn, year+" [34m  INFO[0m print"), true)

	log.Printf("%s", "printf")
	Equal(t, hasString(conn, year+" [34m  INFO[0m printf"), true)

	log.Println("println")
	Equal(t, hasString(conn, year+" [34m  INFO[0m println"), true)

	PanicMatches(t, func() { log.Panic("panic") }, "panic")
	Equal(t, hasString(conn, year+" [31m PANIC[0m syslog_test.go:265 panic"), true)

	PanicMatches(t, func() { log.Panicf("%s", "panicf") }, "panicf")
	Equal(t, hasString(conn, year+" [31m PANIC[0m syslog_test.go:268 panicf"), true)

	PanicMatches(t, func() { log.Panicln("panicln") }, "panicln")
	Equal(t, hasString(conn, year+" [31m PANIC[0m syslog_test.go:271 panicln"), true)

	// WithFields
	log.WithFields(log.F("key", "value")).Debug("debug")
	Equal(t, hasString(conn, year+" [32m DEBUG[0m debug [32mkey[0m=value"), true)

	log.WithFields(log.F("key", "value")).Debugf("%s", "debugf")
	Equal(t, hasString(conn, year+" [32m DEBUG[0m debugf [32mkey[0m=value"), true)

	log.WithFields(log.F("key", "value")).Info("info")
	Equal(t, hasString(conn, year+" [34m  INFO[0m info [34mkey[0m=value"), true)

	log.WithFields(log.F("key", "value")).Infof("%s", "infof")
	Equal(t, hasString(conn, year+" [34m  INFO[0m infof [34mkey[0m=value"), true)

	log.WithFields(log.F("key", "value")).Notice("notice")
	Equal(t, hasString(conn, year+" [36;1mNOTICE[0m notice [36;1mkey[0m=value"), true)

	log.WithFields(log.F("key", "value")).Noticef("%s", "noticef")
	Equal(t, hasString(conn, year+" [36;1mNOTICE[0m noticef [36;1mkey[0m=value"), true)

	log.WithFields(log.F("key", "value")).Warn("warn")
	Equal(t, hasString(conn, year+" [33;1m  WARN[0m syslog_test.go:293 warn [33;1mkey[0m=value"), true)

	log.WithFields(log.F("key", "value")).Warnf("%s", "warnf")
	Equal(t, hasString(conn, year+" [33;1m  WARN[0m syslog_test.go:296 warnf [33;1mkey[0m=value"), true)

	log.WithFields(log.F("key", "value")).Error("error")
	Equal(t, hasString(conn, year+" [31;1m ERROR[0m syslog_test.go:299 error [31;1mkey[0m=value"), true)

	log.WithFields(log.F("key", "value")).Errorf("%s", "errorf")
	Equal(t, hasString(conn, year+" [31;1m ERROR[0m syslog_test.go:302 errorf [31;1mkey[0m=value"), true)

	log.WithFields(log.F("key", "value")).Alert("alert")
	Equal(t, hasString(conn, year+" [31m[4m ALERT[0m syslog_test.go:305 alert [31m[4mkey[0m=value"), true)

	log.WithFields(log.F("key", "value")).Alertf("%s", "alertf")
	Equal(t, hasString(conn, year+" [31m[4m ALERT[0m syslog_test.go:308 alertf [31m[4mkey[0m=value"), true)

	PanicMatches(t, func() { log.WithFields(log.F("key", "value")).Panicf("%s", "panicf") }, "panicf key=value")
	Equal(t, hasString(conn, year+" [31m PANIC[0m syslog_test.go:311 panicf [31mkey[0m=value"), true)

	PanicMatches(t, func() { log.WithFields(log.F("key", "value")).Panic("panic") }, "panic key=value")
	Equal(t, hasString(conn, year+" [31m PANIC[0m syslog_test.go:314 panic [31mkey[0m=value"), true)

	func() {
		defer log.Trace("trace").End()
	}()

	Equal(t, hasString(conn, year+" [37;1m TRACE[0m trace"), true)

	func() {
		defer log.Tracef("tracef").End()
	}()

	Equal(t, hasString(conn, year+" [37;1m TRACE[0m tracef"), true)

	func() {
		defer log.WithFields(log.F("key", "value")).Trace("trace").End()
	}()

	Equal(t, hasString(conn, year+" [37;1m TRACE[0m trace"), true)

	func() {
		defer log.WithFields(log.F("key", "value")).Tracef("tracef").End()
	}()

	Equal(t, hasString(conn, year+" [37;1m TRACE[0m tracef"), true)

	e := &log.Entry{
		Level:     log.FatalLevel,
		Message:   "fatal",
		Timestamp: time.Now(),
		Line:      3,
		File:      "fake.txt",
	}

	log.HandleEntry(e)

	Equal(t, hasString(conn, year+" [31m[4m[5m FATAL[0m fake.txt:3 fatal"), true)
}

func TestSyslogLoggerCaller(t *testing.T) {

	year := time.Now().Format("2006")

	addr, err := net.ResolveUDPAddr("udp", ":2002")
	Equal(t, err, nil)

	conn, err := net.ListenUDP("udp", addr)
	Equal(t, err, nil)
	defer conn.Close()

	sLog, err := New("udp", "127.0.0.1:2002", stdsyslog.LOG_DEBUG, "")
	Equal(t, err, nil)

	sLog.DisplayColor(false)
	sLog.SetBuffersAndWorkers(3, 3)
	sLog.SetTimestampFormat("2006")

	log.SetCallerInfo(true)
	log.RegisterHandler(sLog, log.AllLevels...)

	log.Debug("debug")
	Equal(t, hasString(conn, year+"  DEBUG syslog_test.go:375 debug"), true)

	log.Debugf("%s", "debugf")
	Equal(t, hasString(conn, year+"  DEBUG syslog_test.go:378 debugf"), true)

	log.Info("info")
	Equal(t, hasString(conn, year+"   INFO syslog_test.go:381 info"), true)

	log.Infof("%s", "infof")
	Equal(t, hasString(conn, year+"   INFO syslog_test.go:384 infof"), true)

	log.Notice("notice")
	Equal(t, hasString(conn, year+" NOTICE syslog_test.go:387 notice"), true)

	log.Noticef("%s", "noticef")
	Equal(t, hasString(conn, year+" NOTICE syslog_test.go:390 noticef"), true)

	log.Warn("warn")
	Equal(t, hasString(conn, year+"   WARN syslog_test.go:393 warn"), true)

	log.Warnf("%s", "warnf")
	Equal(t, hasString(conn, year+"   WARN syslog_test.go:396 warnf"), true)

	log.Error("error")
	Equal(t, hasString(conn, year+"  ERROR syslog_test.go:399 error"), true)

	log.Errorf("%s", "errorf")
	Equal(t, hasString(conn, year+"  ERROR syslog_test.go:402 errorf"), true)

	log.Alert("alert")
	Equal(t, hasString(conn, year+"  ALERT syslog_test.go:405 alert"), true)

	log.Alertf("%s", "alertf")
	Equal(t, hasString(conn, year+"  ALERT syslog_test.go:408 alertf"), true)

	log.Print("print")
	Equal(t, hasString(conn, year+"   INFO syslog_test.go:411 print"), true)

	log.Printf("%s", "printf")
	Equal(t, hasString(conn, year+"   INFO syslog_test.go:414 printf"), true)

	log.Println("println")
	Equal(t, hasString(conn, year+"   INFO syslog_test.go:417 println"), true)

	PanicMatches(t, func() { log.Panic("panic") }, "panic")
	Equal(t, hasString(conn, year+"  PANIC syslog_test.go:420 panic"), true)

	PanicMatches(t, func() { log.Panicf("%s", "panicf") }, "panicf")
	Equal(t, hasString(conn, year+"  PANIC syslog_test.go:423 panicf"), true)

	PanicMatches(t, func() { log.Panicln("panicln") }, "panicln")
	Equal(t, hasString(conn, year+"  PANIC syslog_test.go:426 panicln"), true)

	// WithFields
	log.WithFields(log.F("key", "value")).Debug("debug")
	Equal(t, hasString(conn, year+"  DEBUG syslog_test.go:430 debug key=value"), true)

	log.WithFields(log.F("key", "value")).Debugf("%s", "debugf")
	Equal(t, hasString(conn, year+"  DEBUG syslog_test.go:433 debugf key=value"), true)

	log.WithFields(log.F("key", "value")).Info("info")
	Equal(t, hasString(conn, year+"   INFO syslog_test.go:436 info key=value"), true)

	log.WithFields(log.F("key", "value")).Infof("%s", "infof")
	Equal(t, hasString(conn, year+"   INFO syslog_test.go:439 infof key=value"), true)

	log.WithFields(log.F("key", "value")).Notice("notice")
	Equal(t, hasString(conn, year+" NOTICE syslog_test.go:442 notice key=value"), true)

	log.WithFields(log.F("key", "value")).Noticef("%s", "noticef")
	Equal(t, hasString(conn, year+" NOTICE syslog_test.go:445 noticef key=value"), true)

	log.WithFields(log.F("key", "value")).Warn("warn")
	Equal(t, hasString(conn, year+"   WARN syslog_test.go:448 warn key=value"), true)

	log.WithFields(log.F("key", "value")).Warnf("%s", "warnf")
	Equal(t, hasString(conn, year+"   WARN syslog_test.go:451 warnf key=value"), true)

	log.WithFields(log.F("key", "value")).Error("error")
	Equal(t, hasString(conn, year+"  ERROR syslog_test.go:454 error key=value"), true)

	log.WithFields(log.F("key", "value")).Errorf("%s", "errorf")
	Equal(t, hasString(conn, year+"  ERROR syslog_test.go:457 errorf key=value"), true)

	log.WithFields(log.F("key", "value")).Alert("alert")
	Equal(t, hasString(conn, year+"  ALERT syslog_test.go:460 alert key=value"), true)

	log.WithFields(log.F("key", "value")).Alertf("%s", "alertf")
	Equal(t, hasString(conn, year+"  ALERT syslog_test.go:463 alertf key=value"), true)

	PanicMatches(t, func() { log.WithFields(log.F("key", "value")).Panicf("%s", "panicf") }, "panicf key=value")
	Equal(t, hasString(conn, year+"  PANIC syslog_test.go:466 panicf key=value"), true)

	PanicMatches(t, func() { log.WithFields(log.F("key", "value")).Panic("panic") }, "panic key=value")
	Equal(t, hasString(conn, year+"  PANIC syslog_test.go:469 panic key=value"), true)

	func() {
		defer log.Trace("trace").End()
	}()

	Equal(t, hasString(conn, year+"  TRACE syslog_test.go:474 trace"), true)

	func() {
		defer log.Tracef("tracef").End()
	}()

	Equal(t, hasString(conn, year+"  TRACE syslog_test.go:480 tracef"), true)

	func() {
		defer log.WithFields(log.F("key", "value")).Trace("trace").End()
	}()

	Equal(t, hasString(conn, year+"  TRACE syslog_test.go:486 trace"), true)

	func() {
		defer log.WithFields(log.F("key", "value")).Tracef("tracef").End()
	}()

	Equal(t, hasString(conn, year+"  TRACE syslog_test.go:492 tracef"), true)

	e := &log.Entry{
		Level:     log.FatalLevel,
		Message:   "fatal",
		Timestamp: time.Now(),
		Line:      5,
		File:      "test.txt",
	}

	log.HandleEntry(e)

	Equal(t, hasString(conn, year+"  FATAL test.txt:5 fatal"), true)
}

func TestSyslogLoggerColorCaller(t *testing.T) {

	year := time.Now().Format("2006")

	addr, err := net.ResolveUDPAddr("udp", ":2003")
	Equal(t, err, nil)

	conn, err := net.ListenUDP("udp", addr)
	Equal(t, err, nil)
	defer conn.Close()

	sLog, err := New("udp", "127.0.0.1:2003", stdsyslog.LOG_DEBUG, "")
	Equal(t, err, nil)

	sLog.DisplayColor(true)
	sLog.SetBuffersAndWorkers(3, 3)
	sLog.SetTimestampFormat("2006")

	log.RegisterHandler(sLog, log.AllLevels...)

	log.Debug("debug")
	Equal(t, hasString(conn, year+" [32m DEBUG[0m syslog_test.go:529 debug"), true)

	log.Debugf("%s", "debugf")
	Equal(t, hasString(conn, year+" [32m DEBUG[0m syslog_test.go:532 debugf"), true)

	log.Info("info")
	Equal(t, hasString(conn, year+" [34m  INFO[0m syslog_test.go:535 info"), true)

	log.Infof("%s", "infof")
	Equal(t, hasString(conn, year+" [34m  INFO[0m syslog_test.go:538 info"), true)

	log.Notice("notice")
	Equal(t, hasString(conn, year+" [36;1mNOTICE[0m syslog_test.go:541 notice"), true)

	log.Noticef("%s", "noticef")
	Equal(t, hasString(conn, year+" [36;1mNOTICE[0m syslog_test.go:544 noticef"), true)

	log.Warn("warn")
	Equal(t, hasString(conn, year+" [33;1m  WARN[0m syslog_test.go:547 warn"), true)

	log.Warnf("%s", "warnf")
	Equal(t, hasString(conn, year+" [33;1m  WARN[0m syslog_test.go:550 warnf"), true)

	log.Error("error")
	Equal(t, hasString(conn, year+" [31;1m ERROR[0m syslog_test.go:553 error"), true)

	log.Errorf("%s", "errorf")
	Equal(t, hasString(conn, year+" [31;1m ERROR[0m syslog_test.go:556 errorf"), true)

	log.Alert("alert")
	Equal(t, hasString(conn, year+" [31m[4m ALERT[0m syslog_test.go:559 alert"), true)

	log.Alertf("%s", "alertf")
	Equal(t, hasString(conn, year+" [31m[4m ALERT[0m syslog_test.go:562 alertf"), true)

	log.Print("print")
	Equal(t, hasString(conn, year+" [34m  INFO[0m syslog_test.go:565 print"), true)

	log.Printf("%s", "printf")
	Equal(t, hasString(conn, year+" [34m  INFO[0m syslog_test.go:568 printf"), true)

	log.Println("println")
	Equal(t, hasString(conn, year+" [34m  INFO[0m syslog_test.go:571 println"), true)

	PanicMatches(t, func() { log.Panic("panic") }, "panic")
	Equal(t, hasString(conn, year+" [31m PANIC[0m syslog_test.go:574 panic"), true)

	PanicMatches(t, func() { log.Panicf("%s", "panicf") }, "panicf")
	Equal(t, hasString(conn, year+" [31m PANIC[0m syslog_test.go:577 panicf"), true)

	PanicMatches(t, func() { log.Panicln("panicln") }, "panicln")
	Equal(t, hasString(conn, year+" [31m PANIC[0m syslog_test.go:580 panicln"), true)

	// WithFields
	log.WithFields(log.F("key", "value")).Debug("debug")
	Equal(t, hasString(conn, year+" [32m DEBUG[0m syslog_test.go:584 debug [32mkey[0m=value"), true)

	log.WithFields(log.F("key", "value")).Debugf("%s", "debugf")
	Equal(t, hasString(conn, year+" [32m DEBUG[0m syslog_test.go:587 debugf [32mkey[0m=value"), true)

	log.WithFields(log.F("key", "value")).Info("info")
	Equal(t, hasString(conn, year+" [34m  INFO[0m syslog_test.go:590 info [34mkey[0m=value"), true)

	log.WithFields(log.F("key", "value")).Infof("%s", "infof")
	Equal(t, hasString(conn, year+" [34m  INFO[0m syslog_test.go:593 infof [34mkey[0m=value"), true)

	log.WithFields(log.F("key", "value")).Notice("notice")
	Equal(t, hasString(conn, year+" [36;1mNOTICE[0m syslog_test.go:596 notice [36;1mkey[0m=value"), true)

	log.WithFields(log.F("key", "value")).Noticef("%s", "noticef")
	Equal(t, hasString(conn, year+" [36;1mNOTICE[0m syslog_test.go:599 noticef [36;1mkey[0m=value"), true)

	log.WithFields(log.F("key", "value")).Warn("warn")
	Equal(t, hasString(conn, year+" [33;1m  WARN[0m syslog_test.go:602 warn [33;1mkey[0m=value"), true)

	log.WithFields(log.F("key", "value")).Warnf("%s", "warnf")
	Equal(t, hasString(conn, year+" [33;1m  WARN[0m syslog_test.go:605 warnf [33;1mkey[0m=value"), true)

	log.WithFields(log.F("key", "value")).Error("error")
	Equal(t, hasString(conn, year+" [31;1m ERROR[0m syslog_test.go:608 error [31;1mkey[0m=value"), true)

	log.WithFields(log.F("key", "value")).Errorf("%s", "errorf")
	Equal(t, hasString(conn, year+" [31;1m ERROR[0m syslog_test.go:611 errorf [31;1mkey[0m=value"), true)

	log.WithFields(log.F("key", "value")).Alert("alert")
	Equal(t, hasString(conn, year+" [31m[4m ALERT[0m syslog_test.go:614 alert [31m[4mkey[0m=value"), true)

	log.WithFields(log.F("key", "value")).Alertf("%s", "alertf")
	Equal(t, hasString(conn, year+" [31m[4m ALERT[0m syslog_test.go:617 alertf [31m[4mkey[0m=value"), true)

	PanicMatches(t, func() { log.WithFields(log.F("key", "value")).Panicf("%s", "panicf") }, "panicf key=value")
	Equal(t, hasString(conn, year+" [31m PANIC[0m syslog_test.go:620 panicf [31mkey[0m=value"), true)

	PanicMatches(t, func() { log.WithFields(log.F("key", "value")).Panic("panic") }, "panic key=value")
	Equal(t, hasString(conn, year+" [31m PANIC[0m syslog_test.go:623 panic [31mkey[0m=value"), true)

	func() {
		defer log.Trace("trace").End()
	}()

	Equal(t, hasString(conn, year+" [37;1m TRACE[0m syslog_test.go:628 trace"), true)

	func() {
		defer log.Tracef("tracef").End()
	}()

	Equal(t, hasString(conn, year+" [37;1m TRACE[0m syslog_test.go:634 tracef"), true)

	func() {
		defer log.WithFields(log.F("key", "value")).Trace("trace").End()
	}()

	Equal(t, hasString(conn, year+" [37;1m TRACE[0m syslog_test.go:640 trace"), true)

	func() {
		defer log.WithFields(log.F("key", "value")).Tracef("tracef").End()
	}()

	Equal(t, hasString(conn, year+" [37;1m TRACE[0m syslog_test.go:646 tracef"), true)

	e := &log.Entry{
		Level:     log.FatalLevel,
		Message:   "fatal",
		Timestamp: time.Now(),
		Line:      54,
		File:      "test.go",
	}

	log.HandleEntry(e)

	Equal(t, hasString(conn, year+" [31m[4m[5m FATAL[0m test.go:54 fatal"), true)
}

func TestBadWorkerCountAndCustomFormatFunc(t *testing.T) {

	addr, err := net.ResolveUDPAddr("udp", ":2004")
	Equal(t, err, nil)

	conn, err := net.ListenUDP("udp", addr)
	Equal(t, err, nil)
	defer conn.Close()

	sLog, err := New("udp", "127.0.0.1:2004", stdsyslog.LOG_DEBUG, "")
	Equal(t, err, nil)

	sLog.DisplayColor(true)
	sLog.SetBuffersAndWorkers(3, 0)
	sLog.SetTimestampFormat("2006")
	sLog.SetFormatFunc(func() Formatter {
		return func(e *log.Entry) string {
			return e.Message
		}
	})

	log.RegisterHandler(sLog, log.AllLevels...)

	log.Debug("debug")
	Equal(t, hasString(conn, "debug"), true)
}
