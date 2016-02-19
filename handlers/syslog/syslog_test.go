package syslog

import (
	stdsyslog "log/syslog"
	"net"
	"strings"
	"sync"
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
	sLog.SetChannelBuffer(3)
	sLog.SetTimestampFormat("2006")
	sLog.SetANSIReset(log.Reset)

	log.RegisterHandler(sLog, log.AllLevels...)

	log.Debug("debug")
	Equal(t, hasString(conn, "DEBUG["+year+"] debug"), true)

	log.Debugf("%s", "debugf")
	Equal(t, hasString(conn, "DEBUG["+year+"] debugf"), true)

	log.Info("info")
	Equal(t, hasString(conn, "INFO["+year+"] info"), true)

	log.Infof("%s", "infof")
	Equal(t, hasString(conn, "INFO["+year+"] infof"), true)

	log.Notice("notice")
	Equal(t, hasString(conn, "NOTICE["+year+"] notice"), true)

	log.Noticef("%s", "noticef")
	Equal(t, hasString(conn, "NOTICE["+year+"] noticef"), true)

	log.Warn("warn")
	Equal(t, hasString(conn, "WARN["+year+"] warn"), true)

	log.Warnf("%s", "warnf")
	Equal(t, hasString(conn, "WARN["+year+"] warnf"), true)

	log.Error("error")
	Equal(t, hasString(conn, "ERROR["+year+"] error"), true)

	log.Errorf("%s", "errorf")
	Equal(t, hasString(conn, "ERROR["+year+"] errorf"), true)

	log.Alert("alert")
	Equal(t, hasString(conn, "ALERT["+year+"] alert"), true)

	log.Alertf("%s", "alertf")
	Equal(t, hasString(conn, "ALERT["+year+"] alertf"), true)

	log.Print("print")
	Equal(t, hasString(conn, "INFO["+year+"] print"), true)

	log.Printf("%s", "printf")
	Equal(t, hasString(conn, "INFO["+year+"] printf"), true)

	log.Println("println")
	Equal(t, hasString(conn, "INFO["+year+"] println"), true)

	PanicMatches(t, func() { log.Panic("panic") }, "panic")
	Equal(t, hasString(conn, "PANIC["+year+"] panic"), true)

	PanicMatches(t, func() { log.Panicf("%s", "panicf") }, "panicf")
	Equal(t, hasString(conn, "PANIC["+year+"] panicf"), true)

	PanicMatches(t, func() { log.Panicln("panicln") }, "panicln")
	Equal(t, hasString(conn, "PANIC["+year+"] panicln"), true)

	// WithFields
	log.WithFields(log.F("key", "value")).Debug("debug")
	Equal(t, hasString(conn, "DEBUG["+year+"] debug key=value"), true)

	log.WithFields(log.F("key", "value")).Debugf("%s", "debugf")
	Equal(t, hasString(conn, "DEBUG["+year+"] debugf key=value"), true)

	log.WithFields(log.F("key", "value")).Info("info")
	Equal(t, hasString(conn, "INFO["+year+"] info key=value"), true)

	log.WithFields(log.F("key", "value")).Infof("%s", "infof")
	Equal(t, hasString(conn, "INFO["+year+"] infof key=value"), true)

	log.WithFields(log.F("key", "value")).Notice("notice")
	Equal(t, hasString(conn, "NOTICE["+year+"] notice key=value"), true)

	log.WithFields(log.F("key", "value")).Noticef("%s", "noticef")
	Equal(t, hasString(conn, "NOTICE["+year+"] noticef key=value"), true)

	log.WithFields(log.F("key", "value")).Warn("warn")
	Equal(t, hasString(conn, "WARN["+year+"] warn key=value"), true)

	log.WithFields(log.F("key", "value")).Warnf("%s", "warnf")
	Equal(t, hasString(conn, "WARN["+year+"] warnf key=value"), true)

	log.WithFields(log.F("key", "value")).Error("error")
	Equal(t, hasString(conn, "ERROR["+year+"] error key=value"), true)

	log.WithFields(log.F("key", "value")).Errorf("%s", "errorf")
	Equal(t, hasString(conn, "ERROR["+year+"] errorf key=value"), true)

	log.WithFields(log.F("key", "value")).Alert("alert")
	Equal(t, hasString(conn, "ALERT["+year+"] alert key=value"), true)

	log.WithFields(log.F("key", "value")).Alertf("%s", "alertf")
	Equal(t, hasString(conn, "ALERT["+year+"] alertf key=value"), true)

	PanicMatches(t, func() { log.WithFields(log.F("key", "value")).Panicf("%s", "panicf") }, "panicf key=value")
	Equal(t, hasString(conn, "PANIC["+year+"] panicf key=value"), true)

	PanicMatches(t, func() { log.WithFields(log.F("key", "value")).Panic("panic") }, "panic key=value")
	Equal(t, hasString(conn, "PANIC["+year+"] panic key=value"), true)

	func() {
		defer log.Trace("trace").End()
	}()

	Equal(t, hasString(conn, "TRACE["+year+"] trace"), true)

	func() {
		defer log.Tracef("tracef").End()
	}()

	Equal(t, hasString(conn, "TRACE["+year+"] tracef"), true)

	func() {
		defer log.WithFields(log.F("key", "value")).Trace("trace").End()
	}()

	Equal(t, hasString(conn, "TRACE["+year+"] trace"), true)

	func() {
		defer log.WithFields(log.F("key", "value")).Tracef("tracef").End()
	}()

	Equal(t, hasString(conn, "TRACE["+year+"] tracef"), true)

	e := &log.Entry{
		WG:        new(sync.WaitGroup),
		Level:     log.FatalLevel,
		Message:   "fatal",
		Timestamp: time.Now(),
	}

	log.HandleEntry(e)

	Equal(t, hasString(conn, "FATAL["+year+"] fatal"), true)

	// Test Custom Formatter
	sLog.SetFormatter(func(e *log.Entry) string {
		return e.Message
	})

	log.Debug("debug")
	Equal(t, hasString(conn, "debug"), true)
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
	sLog.SetChannelBuffer(3)
	sLog.SetTimestampFormat("2006")

	log.RegisterHandler(sLog, log.AllLevels...)

	log.Debug("debug")
	Equal(t, hasString(conn, "[32m DEBUG[0m["+year+"] debug"), true)

	log.Debugf("%s", "debugf")
	Equal(t, hasString(conn, "[32m DEBUG[0m["+year+"] debugf"), true)

	log.Info("info")
	Equal(t, hasString(conn, "[34m  INFO[0m["+year+"] info"), true)

	log.Infof("%s", "infof")
	Equal(t, hasString(conn, "[34m  INFO[0m["+year+"] info"), true)

	log.Notice("notice")
	Equal(t, hasString(conn, "[36;1mNOTICE[0m["+year+"] notice"), true)

	log.Noticef("%s", "noticef")
	Equal(t, hasString(conn, "[36;1mNOTICE[0m["+year+"] noticef"), true)

	log.Warn("warn")
	Equal(t, hasString(conn, "[33;1m  WARN[0m["+year+"] warn"), true)

	log.Warnf("%s", "warnf")
	Equal(t, hasString(conn, "[33;1m  WARN[0m["+year+"] warnf"), true)

	log.Error("error")
	Equal(t, hasString(conn, "[31;1m ERROR[0m["+year+"] error"), true)

	log.Errorf("%s", "errorf")
	Equal(t, hasString(conn, "[31;1m ERROR[0m["+year+"] errorf"), true)

	log.Alert("alert")
	Equal(t, hasString(conn, "[31m[4m ALERT[0m["+year+"] alert"), true)

	log.Alertf("%s", "alertf")
	Equal(t, hasString(conn, "[31m[4m ALERT[0m["+year+"] alertf"), true)

	log.Print("print")
	Equal(t, hasString(conn, "[34m  INFO[0m["+year+"] print"), true)

	log.Printf("%s", "printf")
	Equal(t, hasString(conn, "[34m  INFO[0m["+year+"] printf"), true)

	log.Println("println")
	Equal(t, hasString(conn, "[34m  INFO[0m["+year+"] println"), true)

	PanicMatches(t, func() { log.Panic("panic") }, "panic")
	Equal(t, hasString(conn, "[31m PANIC[0m["+year+"] panic"), true)

	PanicMatches(t, func() { log.Panicf("%s", "panicf") }, "panicf")
	Equal(t, hasString(conn, "[31m PANIC[0m["+year+"] panicf"), true)

	PanicMatches(t, func() { log.Panicln("panicln") }, "panicln")
	Equal(t, hasString(conn, "[31m PANIC[0m["+year+"] panicln"), true)

	// WithFields
	log.WithFields(log.F("key", "value")).Debug("debug")
	Equal(t, hasString(conn, "[32m DEBUG[0m["+year+"] debug                     [32mkey[0m=value"), true)

	log.WithFields(log.F("key", "value")).Debugf("%s", "debugf")
	Equal(t, hasString(conn, "[32m DEBUG[0m["+year+"] debugf                    [32mkey[0m=value"), true)

	log.WithFields(log.F("key", "value")).Info("info")
	Equal(t, hasString(conn, "[34m  INFO[0m["+year+"] info                      [34mkey[0m=value"), true)

	log.WithFields(log.F("key", "value")).Infof("%s", "infof")
	Equal(t, hasString(conn, "[34m  INFO[0m["+year+"] infof                     [34mkey[0m=value"), true)

	log.WithFields(log.F("key", "value")).Notice("notice")
	Equal(t, hasString(conn, "[36;1mNOTICE[0m["+year+"] notice                    [36;1mkey[0m=value"), true)

	log.WithFields(log.F("key", "value")).Noticef("%s", "noticef")
	Equal(t, hasString(conn, "[36;1mNOTICE[0m["+year+"] noticef                   [36;1mkey[0m=value"), true)

	log.WithFields(log.F("key", "value")).Warn("warn")
	Equal(t, hasString(conn, "[33;1m  WARN[0m["+year+"] warn                      [33;1mkey[0m=value"), true)

	log.WithFields(log.F("key", "value")).Warnf("%s", "warnf")
	Equal(t, hasString(conn, "[33;1m  WARN[0m["+year+"] warnf                     [33;1mkey[0m=value"), true)

	log.WithFields(log.F("key", "value")).Error("error")
	Equal(t, hasString(conn, "[31;1m ERROR[0m["+year+"] error                     [31;1mkey[0m=value"), true)

	log.WithFields(log.F("key", "value")).Errorf("%s", "errorf")
	Equal(t, hasString(conn, "[31;1m ERROR[0m["+year+"] errorf                    [31;1mkey[0m=value"), true)

	log.WithFields(log.F("key", "value")).Alert("alert")
	Equal(t, hasString(conn, "[31m[4m ALERT[0m["+year+"] alert                     [31m[4mkey[0m=value"), true)

	log.WithFields(log.F("key", "value")).Alertf("%s", "alertf")
	Equal(t, hasString(conn, "[31m[4m ALERT[0m["+year+"] alertf                    [31m[4mkey[0m=value"), true)

	PanicMatches(t, func() { log.WithFields(log.F("key", "value")).Panicf("%s", "panicf") }, "panicf key=value")
	Equal(t, hasString(conn, "[31m PANIC[0m["+year+"] panicf                    [31mkey[0m=value"), true)

	PanicMatches(t, func() { log.WithFields(log.F("key", "value")).Panic("panic") }, "panic key=value")
	Equal(t, hasString(conn, "[31m PANIC[0m["+year+"] panic                     [31mkey[0m=value"), true)

	func() {
		defer log.Trace("trace").End()
	}()

	Equal(t, hasString(conn, "[37;1m TRACE[0m["+year+"] trace"), true)

	func() {
		defer log.Tracef("tracef").End()
	}()

	Equal(t, hasString(conn, "[37;1m TRACE[0m["+year+"] tracef"), true)

	func() {
		defer log.WithFields(log.F("key", "value")).Trace("trace").End()
	}()

	Equal(t, hasString(conn, "[37;1m TRACE[0m["+year+"] trace"), true)

	func() {
		defer log.WithFields(log.F("key", "value")).Tracef("tracef").End()
	}()

	Equal(t, hasString(conn, "[37;1m TRACE[0m["+year+"] tracef"), true)

	e := &log.Entry{
		WG:        new(sync.WaitGroup),
		Level:     log.FatalLevel,
		Message:   "fatal",
		Timestamp: time.Now(),
	}

	log.HandleEntry(e)

	Equal(t, hasString(conn, "[31m[4m[5m FATAL[0m["+year+"] fatal"), true)

	// test changing level color
	sLog.SetLevelColor(log.DebugLevel, log.Red)

	log.Debug("debug")
	Equal(t, hasString(conn, "[31m DEBUG[0m[2016] debug"), true)
}

func TestBadAddress(t *testing.T) {
	sLog, err := New("udp", "255.255.255.67", stdsyslog.LOG_DEBUG, "")
	Equal(t, sLog, nil)
	NotEqual(t, err, nil)
}
