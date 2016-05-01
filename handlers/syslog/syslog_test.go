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
	Equal(t, hasString(conn, year+"   WARN warn"), true)

	log.Warnf("%s", "warnf")
	Equal(t, hasString(conn, year+"   WARN warnf"), true)

	log.Error("error")
	Equal(t, hasString(conn, year+"  ERROR error"), true)

	log.Errorf("%s", "errorf")
	Equal(t, hasString(conn, year+"  ERROR errorf"), true)

	log.Alert("alert")
	Equal(t, hasString(conn, year+"  ALERT alert"), true)

	log.Alertf("%s", "alertf")
	Equal(t, hasString(conn, year+"  ALERT alertf"), true)

	log.Print("print")
	Equal(t, hasString(conn, year+"   INFO print"), true)

	log.Printf("%s", "printf")
	Equal(t, hasString(conn, year+"   INFO printf"), true)

	log.Println("println")
	Equal(t, hasString(conn, year+"   INFO println"), true)

	PanicMatches(t, func() { log.Panic("panic") }, "panic")
	Equal(t, hasString(conn, year+"  PANIC panic"), true)

	PanicMatches(t, func() { log.Panicf("%s", "panicf") }, "panicf")
	Equal(t, hasString(conn, year+"  PANIC panicf"), true)

	PanicMatches(t, func() { log.Panicln("panicln") }, "panicln")
	Equal(t, hasString(conn, year+"  PANIC panicln"), true)

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
	Equal(t, hasString(conn, year+"   WARN warn key=value"), true)

	log.WithFields(log.F("key", "value")).Warnf("%s", "warnf")
	Equal(t, hasString(conn, year+"   WARN warnf key=value"), true)

	log.WithFields(log.F("key", "value")).Error("error")
	Equal(t, hasString(conn, year+"  ERROR error key=value"), true)

	log.WithFields(log.F("key", "value")).Errorf("%s", "errorf")
	Equal(t, hasString(conn, year+"  ERROR errorf key=value"), true)

	log.WithFields(log.F("key", "value")).Alert("alert")
	Equal(t, hasString(conn, year+"  ALERT alert key=value"), true)

	log.WithFields(log.F("key", "value")).Alertf("%s", "alertf")
	Equal(t, hasString(conn, year+"  ALERT alertf key=value"), true)

	PanicMatches(t, func() { log.WithFields(log.F("key", "value")).Panicf("%s", "panicf") }, "panicf key=value")
	Equal(t, hasString(conn, year+"  PANIC panicf key=value"), true)

	PanicMatches(t, func() { log.WithFields(log.F("key", "value")).Panic("panic") }, "panic key=value")
	Equal(t, hasString(conn, year+"  PANIC panic key=value"), true)

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
	}

	log.HandleEntry(e)

	Equal(t, hasString(conn, year+"  FATAL fatal"), true)

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
	Equal(t, hasString(conn, year+" [33;1m  WARN[0m warn"), true)

	log.Warnf("%s", "warnf")
	Equal(t, hasString(conn, year+" [33;1m  WARN[0m warnf"), true)

	log.Error("error")
	Equal(t, hasString(conn, year+" [31;1m ERROR[0m error"), true)

	log.Errorf("%s", "errorf")
	Equal(t, hasString(conn, year+" [31;1m ERROR[0m errorf"), true)

	log.Alert("alert")
	Equal(t, hasString(conn, year+" [31m[4m ALERT[0m alert"), true)

	log.Alertf("%s", "alertf")
	Equal(t, hasString(conn, year+" [31m[4m ALERT[0m alertf"), true)

	log.Print("print")
	Equal(t, hasString(conn, year+" [34m  INFO[0m print"), true)

	log.Printf("%s", "printf")
	Equal(t, hasString(conn, year+" [34m  INFO[0m printf"), true)

	log.Println("println")
	Equal(t, hasString(conn, year+" [34m  INFO[0m println"), true)

	PanicMatches(t, func() { log.Panic("panic") }, "panic")
	Equal(t, hasString(conn, year+" [31m PANIC[0m panic"), true)

	PanicMatches(t, func() { log.Panicf("%s", "panicf") }, "panicf")
	Equal(t, hasString(conn, year+" [31m PANIC[0m panicf"), true)

	PanicMatches(t, func() { log.Panicln("panicln") }, "panicln")
	Equal(t, hasString(conn, year+" [31m PANIC[0m panicln"), true)

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
	Equal(t, hasString(conn, year+" [33;1m  WARN[0m warn [33;1mkey[0m=value"), true)

	log.WithFields(log.F("key", "value")).Warnf("%s", "warnf")
	Equal(t, hasString(conn, year+" [33;1m  WARN[0m warnf [33;1mkey[0m=value"), true)

	log.WithFields(log.F("key", "value")).Error("error")
	Equal(t, hasString(conn, year+" [31;1m ERROR[0m error [31;1mkey[0m=value"), true)

	log.WithFields(log.F("key", "value")).Errorf("%s", "errorf")
	Equal(t, hasString(conn, year+" [31;1m ERROR[0m errorf [31;1mkey[0m=value"), true)

	log.WithFields(log.F("key", "value")).Alert("alert")
	Equal(t, hasString(conn, year+" [31m[4m ALERT[0m alert [31m[4mkey[0m=value"), true)

	log.WithFields(log.F("key", "value")).Alertf("%s", "alertf")
	Equal(t, hasString(conn, year+" [31m[4m ALERT[0m alertf [31m[4mkey[0m=value"), true)

	PanicMatches(t, func() { log.WithFields(log.F("key", "value")).Panicf("%s", "panicf") }, "panicf key=value")
	Equal(t, hasString(conn, year+" [31m PANIC[0m panicf [31mkey[0m=value"), true)

	PanicMatches(t, func() { log.WithFields(log.F("key", "value")).Panic("panic") }, "panic key=value")
	Equal(t, hasString(conn, year+" [31m PANIC[0m panic [31mkey[0m=value"), true)

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
	}

	log.HandleEntry(e)

	Equal(t, hasString(conn, year+" [31m[4m[5m FATAL[0m fatal"), true)

	// test changing level color
	sLog.SetLevelColor(log.DebugLevel, log.Red)

	log.Debug("debug")
	Equal(t, hasString(conn, "2016 [31m DEBUG[0m debug"), true)
}

func TestBadAddress(t *testing.T) {
	sLog, err := New("udp", "255.255.255.67", stdsyslog.LOG_DEBUG, "")
	Equal(t, sLog, nil)
	NotEqual(t, err, nil)
}
