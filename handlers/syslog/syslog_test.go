package syslog

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/go-playground/log"
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

func hasString(conn *net.UDPConn) string {
	read, _ := read(conn)
	return read
}

func TestSyslogLogger(t *testing.T) {
	log.SetExitFunc(func(int) {})
	tests := getSyslogLoggerTests()

	addr, err := net.ResolveUDPAddr("udp", ":2000")
	if err != nil {
		t.Errorf("Expected '%v' Got '%v'", nil, err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		t.Errorf("Expected '%v' Got '%v'", nil, err)
	}
	defer func() {
		_ = conn.Close()
	}()

	sLog, err := New("udp", "127.0.0.1:2000", "", nil)
	if err != nil {
		t.Errorf("Expected '%v' Got '%v'", nil, err)
	}

	sLog.SetDisplayColor(false)
	sLog.SetTimestampFormat("")
	log.AddHandler(sLog, log.AllLevels...)

	for i, tt := range tests {

		var l log.Entry

		if tt.flds != nil {
			l = l.WithFields(tt.flds...)
		}

		switch tt.lvl {
		case log.DebugLevel:
			if len(tt.printf) == 0 {
				l.Debug(tt.msg)
			} else {
				l.Debugf(tt.printf, tt.msg)
			}
		case log.InfoLevel:
			if len(tt.printf) == 0 {
				l.Info(tt.msg)
			} else {
				l.Infof(tt.printf, tt.msg)
			}
		case log.NoticeLevel:
			if len(tt.printf) == 0 {
				l.Notice(tt.msg)
			} else {
				l.Noticef(tt.printf, tt.msg)
			}
		case log.WarnLevel:
			if len(tt.printf) == 0 {
				l.Warn(tt.msg)
			} else {
				l.Warnf(tt.printf, tt.msg)
			}
		case log.ErrorLevel:
			if len(tt.printf) == 0 {
				l.Error(tt.msg)
			} else {
				l.Errorf(tt.printf, tt.msg)
			}
		case log.PanicLevel:
			if len(tt.printf) == 0 {
				l.Panic(tt.msg)
			} else {
				l.Panicf(tt.printf, tt.msg)
			}
		case log.AlertLevel:
			if len(tt.printf) == 0 {
				l.Alert(tt.msg)
			} else {
				l.Alertf(tt.printf, tt.msg)
			}
		}

		if s := hasString(conn); !strings.HasSuffix(s, tt.want) {
			t.Errorf("test %d: Expected Suffix '%s' Got '%s'", i, tt.want, s)
		}
	}
}

func TestSyslogLoggerColor(t *testing.T) {
	log.SetExitFunc(func(int) {})
	tests := getSyslogLoggerColorTests()

	addr, err := net.ResolveUDPAddr("udp", ":2001")
	if err != nil {
		t.Errorf("Expected '%v' Got '%v'", nil, err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		t.Errorf("Expected '%v' Got '%s'", nil, err)
	}
	defer func() {
		_ = conn.Close()
	}()

	sLog, err := New("udp", "127.0.0.1:2001", "", nil)
	if err != nil {
		t.Errorf("Expected '%v' Got '%s'", nil, err)
	}

	sLog.SetDisplayColor(true)
	sLog.SetTimestampFormat("")

	log.AddHandler(sLog, log.AllLevels...)

	for i, tt := range tests {

		var l log.Entry

		if tt.flds != nil {
			l = l.WithFields(tt.flds...)
		}

		switch tt.lvl {
		case log.DebugLevel:
			if len(tt.printf) == 0 {
				l.Debug(tt.msg)
			} else {
				l.Debugf(tt.printf, tt.msg)
			}

		case log.InfoLevel:
			if len(tt.printf) == 0 {
				l.Info(tt.msg)
			} else {
				l.Infof(tt.printf, tt.msg)
			}
		case log.NoticeLevel:
			if len(tt.printf) == 0 {
				l.Notice(tt.msg)
			} else {
				l.Noticef(tt.printf, tt.msg)
			}
		case log.WarnLevel:
			if len(tt.printf) == 0 {
				l.Warn(tt.msg)
			} else {
				l.Warnf(tt.printf, tt.msg)
			}
		case log.ErrorLevel:
			if len(tt.printf) == 0 {
				l.Error(tt.msg)
			} else {
				l.Errorf(tt.printf, tt.msg)
			}
		case log.PanicLevel:
			if len(tt.printf) == 0 {
				l.Panic(tt.msg)
			} else {
				l.Panicf(tt.printf, tt.msg)
			}
		case log.AlertLevel:
			if len(tt.printf) == 0 {
				l.Alert(tt.msg)
			} else {
				l.Alertf(tt.printf, tt.msg)
			}
		case log.FatalLevel:
			if len(tt.printf) == 0 {
				l.Fatal(tt.msg)
			} else {
				l.Fatalf(tt.printf, tt.msg)
			}
		}

		if s := hasString(conn); !strings.HasSuffix(s, tt.want) {
			t.Errorf("test %d: Expected Suffix '%s' Got '%s'", i, tt.want, s)
		}
	}
}

func TestBadAddress(t *testing.T) {
	sLog, err := New("udp", "255.255.255.67", "", nil)
	if err == nil {
		log.Errorf("Expected '%v' Got '%v'", "not nil", err)
	}

	if sLog != nil {
		log.Errorf("Expected '%v' Got '%v'", nil, sLog)
	}
}

func TestBadWorkerCountAndCustomFormatFunc(t *testing.T) {

	addr, err := net.ResolveUDPAddr("udp", ":2004")
	if err != nil {
		log.Errorf("Expected '%v' Got '%v'", nil, err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Errorf("Expected '%v' Got '%v'", nil, err)
	}
	defer func() {
		_ = conn.Close()
	}()

	sLog, err := New("udp", "127.0.0.1:2004", "", nil)
	if err != nil {
		log.Errorf("Expected '%v' Got '%v'", nil, err)
	}

	sLog.SetDisplayColor(true)
	sLog.SetTimestampFormat("2006")
	sLog.SetFormatFunc(func(s *Syslog) Formatter {
		return func(e log.Entry) []byte {
			return []byte(e.Message)
		}
	})

	log.AddHandler(sLog, log.AllLevels...)

	log.Debug("debug")
	if s := hasString(conn); s != "debug" {
		log.Errorf("Expected '%s' Got '%s'", "debug", s)
	}
}

func TestSyslogTLS(t *testing.T) {

	// setup server

	addr, err := net.ResolveTCPAddr("tcp", ":2022")
	if err != nil {
		log.Errorf("Expected '%v' Got '%v'", nil, err)
	}

	conn, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Errorf("Expected '%v' Got '%v'", nil, err)
	}
	defer func() {
		_ = conn.Close()
	}()

	var msg string
	var m sync.Mutex

	go func() {
		client, err := conn.Accept()
		if err != nil {
			panic(fmt.Sprintf("Expected '%v' Got '%v'", nil, err))
		}

		b := make([]byte, 1024)

		m.Lock()
		defer m.Unlock()
		read, err := client.Read(b)
		if err != nil {
			panic(fmt.Sprintf("Expected '%v' Got '%v'", nil, err))
		}
		msg = string(b[0:read])
	}()

	sLog, err := New("tcp", "127.0.0.1:2022", "", nil)
	if err != nil {
		panic(fmt.Sprintf("Expected '%v' Got '%v'", nil, err))
	}

	sLog.SetFormatFunc(func(s *Syslog) Formatter {
		return func(e log.Entry) []byte {
			return []byte(e.Message)
		}
	})

	log.AddHandler(sLog, log.AllLevels...)

	log.Debug("debug")
	time.Sleep(500 * time.Millisecond)

	m.Lock()
	defer m.Unlock()

	if strings.HasSuffix(msg, "debug") {
		t.Fatalf("Expected '%s' Got '%s'", "debug", msg)
	}
}

type test struct {
	lvl    log.Level
	msg    string
	flds   []log.Field
	want   string
	printf string
}

func getSyslogLoggerTests() []test {
	return []test{
		{
			lvl:  log.DebugLevel,
			msg:  "debug",
			flds: nil,
			want: "  DEBUG debug\n",
		},
		{
			lvl:    log.DebugLevel,
			msg:    "debugf",
			printf: "%s",
			flds:   nil,
			want:   "  DEBUG debugf\n",
		},
		{
			lvl:  log.InfoLevel,
			msg:  "info",
			flds: nil,
			want: "   INFO info\n",
		},
		{
			lvl:    log.InfoLevel,
			msg:    "infof",
			printf: "%s",
			flds:   nil,
			want:   "   INFO infof\n",
		},
		{
			lvl:  log.NoticeLevel,
			msg:  "notice",
			flds: nil,
			want: " NOTICE notice\n",
		},
		{
			lvl:    log.NoticeLevel,
			msg:    "noticef",
			printf: "%s",
			flds:   nil,
			want:   " NOTICE noticef\n",
		},
		{
			lvl:  log.WarnLevel,
			msg:  "warn",
			flds: nil,
			want: "   WARN warn\n",
		},
		{
			lvl:    log.WarnLevel,
			msg:    "warnf",
			printf: "%s",
			flds:   nil,
			want:   "   WARN warnf\n",
		},
		{
			lvl:  log.ErrorLevel,
			msg:  "error",
			flds: nil,
			want: "  ERROR error\n",
		},
		{
			lvl:    log.ErrorLevel,
			msg:    "errorf",
			printf: "%s",
			flds:   nil,
			want:   "  ERROR errorf\n",
		},
		{
			lvl:  log.AlertLevel,
			msg:  "alert",
			flds: nil,
			want: "  ALERT alert\n",
		},
		{
			lvl:    log.AlertLevel,
			msg:    "alertf",
			printf: "%s",
			flds:   nil,
			want:   "  ALERT alertf\n",
		},
		{
			lvl:  log.PanicLevel,
			msg:  "panic",
			flds: nil,
			want: "  PANIC panic\n",
		},
		{
			lvl:    log.PanicLevel,
			msg:    "panicf",
			printf: "%s",
			flds:   nil,
			want:   "  PANIC panicf\n",
		},
		{
			lvl: log.DebugLevel,
			msg: "debug",
			flds: []log.Field{
				log.F("key", "value"),
			},
			want: "  DEBUG debug key=value\n",
		},
		{
			lvl:    log.DebugLevel,
			msg:    "debugf",
			printf: "%s",
			flds: []log.Field{
				log.F("key", "value"),
			},
			want: "  DEBUG debugf key=value\n",
		},
		{
			lvl: log.InfoLevel,
			msg: "info",
			flds: []log.Field{
				log.F("key", "value"),
			},
			want: "   INFO info key=value\n",
		},
		{
			lvl:    log.InfoLevel,
			msg:    "infof",
			printf: "%s",
			flds: []log.Field{
				log.F("key", "value"),
			},
			want: "   INFO infof key=value\n",
		},
		{
			lvl: log.NoticeLevel,
			msg: "notice",
			flds: []log.Field{
				log.F("key", "value"),
			},
			want: " NOTICE notice key=value\n",
		},
		{
			lvl:    log.NoticeLevel,
			msg:    "noticef",
			printf: "%s",
			flds: []log.Field{
				log.F("key", "value"),
			},
			want: " NOTICE noticef key=value\n",
		},
		{
			lvl: log.WarnLevel,
			msg: "warn",
			flds: []log.Field{
				log.F("key", "value"),
			},
			want: "   WARN warn key=value\n",
		},
		{
			lvl:    log.WarnLevel,
			msg:    "warnf",
			printf: "%s",
			flds: []log.Field{
				log.F("key", "value"),
			},
			want: "   WARN warnf key=value\n",
		},
		{
			lvl: log.ErrorLevel,
			msg: "error",
			flds: []log.Field{
				log.F("key", "value"),
			},
			want: "  ERROR error key=value\n",
		},
		{
			lvl:    log.ErrorLevel,
			msg:    "errorf",
			printf: "%s",
			flds: []log.Field{
				log.F("key", "value"),
			},
			want: "  ERROR errorf key=value\n",
		},
		{
			lvl: log.AlertLevel,
			msg: "alert",
			flds: []log.Field{
				log.F("key", "value"),
			},
			want: "  ALERT alert key=value\n",
		},
		{
			lvl: log.AlertLevel,
			msg: "alert",
			flds: []log.Field{
				log.F("key", "value"),
			},
			want: "  ALERT alert key=value\n",
		},
		{
			lvl:    log.AlertLevel,
			msg:    "alertf",
			printf: "%s",
			flds: []log.Field{
				log.F("key", "value"),
			},
			want: "  ALERT alertf key=value\n",
		},
		{
			lvl:    log.PanicLevel,
			msg:    "panicf",
			printf: "%s",
			flds: []log.Field{
				log.F("key", "value"),
			},
			want: "  PANIC panicf key=value\n",
		},
		{
			lvl: log.PanicLevel,
			msg: "panic",
			flds: []log.Field{
				log.F("key", "value"),
			},
			want: "  PANIC panic key=value\n",
		},
		{
			lvl: log.DebugLevel,
			msg: "debug",
			flds: []log.Field{
				log.F("key", "string"),
				log.F("key", int(1)),
				log.F("key", int8(2)),
				log.F("key", int16(3)),
				log.F("key", int32(4)),
				log.F("key", int64(5)),
				log.F("key", uint(1)),
				log.F("key", uint8(2)),
				log.F("key", uint16(3)),
				log.F("key", uint32(4)),
				log.F("key", uint64(5)),
				log.F("key", true),
				log.F("key", struct{ value string }{"struct"}),
			},
			want: "  DEBUG debug key=string key=1 key=2 key=3 key=4 key=5 key=1 key=2 key=3 key=4 key=5 key=true key={struct}\n",
		},
	}
}

func getSyslogLoggerColorTests() []test {
	return []test{
		{
			lvl:    log.DebugLevel,
			msg:    "debugf",
			printf: "%s",
			flds:   nil,
			want: " [32m DEBUG[0m debugf\n",
		},
		{
			lvl:    log.FatalLevel,
			msg:    "fatal",
			printf: "%s",
			flds:   nil,
			want: " [31m[4m[5m FATAL[0m fatal\n",
		},
		{
			lvl:  log.DebugLevel,
			msg:  "debug",
			flds: nil,
			want: " [32m DEBUG[0m debug\n",
		},
		{
			lvl:    log.InfoLevel,
			msg:    "infof",
			printf: "%s",
			flds:   nil,
			want: " [34m  INFO[0m infof\n",
		},
		{
			lvl:  log.InfoLevel,
			msg:  "info",
			flds: nil,
			want: " [34m  INFO[0m info\n",
		},
		{
			lvl:    log.NoticeLevel,
			msg:    "noticef",
			printf: "%s",
			flds:   nil,
			want: " [36;1mNOTICE[0m noticef\n",
		},
		{
			lvl:  log.NoticeLevel,
			msg:  "notice",
			flds: nil,
			want: " [36;1mNOTICE[0m notice\n",
		},
		{
			lvl:    log.WarnLevel,
			msg:    "warnf",
			printf: "%s",
			flds:   nil,
			want: " [33;1m  WARN[0m warnf\n",
		},
		{
			lvl:  log.WarnLevel,
			msg:  "warn",
			flds: nil,
			want: " [33;1m  WARN[0m warn\n",
		},
		{
			lvl:    log.ErrorLevel,
			msg:    "errorf",
			printf: "%s",
			flds:   nil,
			want: " [31;1m ERROR[0m errorf\n",
		},
		{
			lvl:  log.ErrorLevel,
			msg:  "error",
			flds: nil,
			want: " [31;1m ERROR[0m error\n",
		},
		{
			lvl:    log.AlertLevel,
			msg:    "alertf",
			printf: "%s",
			flds:   nil,
			want: " [31m[4m ALERT[0m alertf\n",
		},
		{
			lvl:  log.AlertLevel,
			msg:  "alert",
			flds: nil,
			want: " [31m[4m ALERT[0m alert\n",
		},
		{
			lvl:    log.PanicLevel,
			msg:    "panicf",
			printf: "%s",
			flds:   nil,
			want: " [31m PANIC[0m panicf\n",
		},
		{
			lvl:  log.PanicLevel,
			msg:  "panic",
			flds: nil,
			want: " [31m PANIC[0m panic\n",
		},
		{
			lvl:    log.DebugLevel,
			msg:    "debugf",
			printf: "%s",
			flds: []log.Field{
				log.F("key", "value"),
			},
			want: " [32m DEBUG[0m debugf [32mkey[0m=value\n",
		},
		{
			lvl: log.DebugLevel,
			msg: "debug",
			flds: []log.Field{
				log.F("key", "value"),
			},
			want: " [32m DEBUG[0m debug [32mkey[0m=value\n",
		},
		{
			lvl:    log.InfoLevel,
			msg:    "infof",
			printf: "%s",
			flds: []log.Field{
				log.F("key", "value"),
			},
			want: " [34m  INFO[0m infof [34mkey[0m=value\n",
		},
		{
			lvl: log.InfoLevel,
			msg: "info",
			flds: []log.Field{
				log.F("key", "value"),
			},
			want: " [34m  INFO[0m info [34mkey[0m=value\n",
		},
		{
			lvl:    log.NoticeLevel,
			msg:    "noticef",
			printf: "%s",
			flds: []log.Field{
				log.F("key", "value"),
			},
			want: " [36;1mNOTICE[0m noticef [36;1mkey[0m=value\n",
		},
		{
			lvl: log.NoticeLevel,
			msg: "notice",
			flds: []log.Field{
				log.F("key", "value"),
			},
			want: " [36;1mNOTICE[0m notice [36;1mkey[0m=value\n",
		},
		{
			lvl:    log.WarnLevel,
			msg:    "warnf",
			printf: "%s",
			flds: []log.Field{
				log.F("key", "value"),
			},
			want: " [33;1m  WARN[0m warnf [33;1mkey[0m=value\n",
		},
		{
			lvl: log.WarnLevel,
			msg: "warn",
			flds: []log.Field{
				log.F("key", "value"),
			},
			want: " [33;1m  WARN[0m warn [33;1mkey[0m=value\n",
		},
		{
			lvl:    log.ErrorLevel,
			msg:    "errorf",
			printf: "%s",
			flds: []log.Field{
				log.F("key", "value"),
			},
			want: " [31;1m ERROR[0m errorf [31;1mkey[0m=value\n",
		},
		{
			lvl: log.ErrorLevel,
			msg: "error",
			flds: []log.Field{
				log.F("key", "value"),
			},
			want: " [31;1m ERROR[0m error [31;1mkey[0m=value\n",
		},
		{
			lvl:    log.AlertLevel,
			msg:    "alertf",
			printf: "%s",
			flds: []log.Field{
				log.F("key", "value"),
			},
			want: " [31m[4m ALERT[0m alertf [31m[4mkey[0m=value\n",
		},
		{
			lvl: log.AlertLevel,
			msg: "alert",
			flds: []log.Field{
				log.F("key", "value"),
			},
			want: " [31m[4m ALERT[0m alert [31m[4mkey[0m=value\n",
		},
		{
			lvl:    log.PanicLevel,
			msg:    "panicf",
			printf: "%s",
			flds: []log.Field{
				log.F("key", "value"),
			},
			want: " [31m PANIC[0m panicf [31mkey[0m=value\n",
		},
		{
			lvl: log.PanicLevel,
			msg: "panic",
			flds: []log.Field{
				log.F("key", "value"),
			},
			want: " [31m PANIC[0m panic [31mkey[0m=value\n",
		},
		{
			lvl: log.DebugLevel,
			msg: "debug",
			flds: []log.Field{
				log.F("key", "string"),
				log.F("key", int(1)),
				log.F("key", int8(2)),
				log.F("key", int16(3)),
				log.F("key", int32(4)),
				log.F("key", int64(5)),
				log.F("key", uint(1)),
				log.F("key", uint8(2)),
				log.F("key", uint16(3)),
				log.F("key", uint32(4)),
				log.F("key", uint64(5)),
				log.F("key", float32(5.33)),
				log.F("key", float64(5.34)),
				log.F("key", true),
				log.F("key", struct{ value string }{"struct"}),
			},
			want: " [32m DEBUG[0m debug [32mkey[0m=string [32mkey[0m=1 [32mkey[0m=2 [32mkey[0m=3 [32mkey[0m=4 [32mkey[0m=5 [32mkey[0m=1 [32mkey[0m=2 [32mkey[0m=3 [32mkey[0m=4 [32mkey[0m=5 [32mkey[0m=5.33 [32mkey[0m=5.34 [32mkey[0m=true [32mkey[0m={struct}\n",
		},
	}
}
