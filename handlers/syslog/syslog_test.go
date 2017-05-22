package syslog

import (
	"crypto/tls"
	"crypto/x509"
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
	tests := getSyslogLoggerTests()

	addr, err := net.ResolveUDPAddr("udp", ":2000")
	if err != nil {
		t.Errorf("Expected '%v' Got '%v'", nil, err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		t.Errorf("Expected '%v' Got '%v'", nil, err)
	}
	defer conn.Close()

	sLog, err := New("udp", "127.0.0.1:2000", "", nil)
	if err != nil {
		t.Errorf("Expected '%v' Got '%v'", nil, err)
	}

	sLog.SetDisplayColor(false)
	sLog.SetBuffersAndWorkers(3, 3)
	sLog.SetTimestampFormat("MST")
	log.SetCallerInfoLevels(log.WarnLevel, log.ErrorLevel, log.PanicLevel, log.AlertLevel, log.FatalLevel)
	log.RegisterHandler(sLog, log.AllLevels...)

	for i, tt := range tests {

		var l log.LeveledLogger

		if tt.flds != nil {
			l = log.WithFields(tt.flds...)
		} else {
			l = log.Logger
		}

		switch tt.lvl {
		case log.DebugLevel:
			if len(tt.printf) == 0 {
				l.Debug(tt.msg)
			} else {
				l.Debugf(tt.printf, tt.msg)
			}
		case log.TraceLevel:
			if len(tt.printf) == 0 {
				l.Trace(tt.msg).End()
			} else {
				l.Tracef(tt.printf, tt.msg).End()
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
			func() {
				defer func() {
					recover()
				}()

				if len(tt.printf) == 0 {
					l.Panic(tt.msg)
				} else {
					l.Panicf(tt.printf, tt.msg)
				}
			}()
		case log.AlertLevel:
			if len(tt.printf) == 0 {
				l.Alert(tt.msg)
			} else {
				l.Alertf(tt.printf, tt.msg)
			}
		}

		if s := hasString(conn); !strings.HasSuffix(s, tt.want) {

			if tt.lvl == log.TraceLevel {
				if !strings.Contains(s, tt.want) {
					t.Errorf("test %d: Contains Suffix '%s' Got '%s'", i, tt.want, s)
				}
				continue
			}

			t.Errorf("test %d: Expected Suffix '%s' Got '%s'", i, tt.want, s)
		}
	}
}

func TestSyslogLoggerColor(t *testing.T) {

	tests := getSyslogLoggerColorTests()

	addr, err := net.ResolveUDPAddr("udp", ":2001")
	if err != nil {
		t.Errorf("Expected '%v' Got '%v'", nil, err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		t.Errorf("Expected '%v' Got '%s'", nil, err)
	}
	defer conn.Close()

	sLog, err := New("udp", "127.0.0.1:2001", "", nil)
	if err != nil {
		t.Errorf("Expected '%v' Got '%s'", nil, err)
	}

	sLog.SetDisplayColor(true)
	sLog.SetBuffersAndWorkers(3, 3)
	sLog.SetTimestampFormat("MST")

	log.RegisterHandler(sLog, log.AllLevels...)

	for i, tt := range tests {

		var l log.LeveledLogger

		if tt.flds != nil {
			l = log.WithFields(tt.flds...)
		} else {
			l = log.Logger
		}

		switch tt.lvl {
		case log.DebugLevel:
			if len(tt.printf) == 0 {
				l.Debug(tt.msg)
			} else {
				l.Debugf(tt.printf, tt.msg)
			}
		case log.TraceLevel:
			if len(tt.printf) == 0 {
				l.Trace(tt.msg).End()
			} else {
				l.Tracef(tt.printf, tt.msg).End()
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
			func() {
				defer func() {
					recover()
				}()

				if len(tt.printf) == 0 {
					l.Panic(tt.msg)
				} else {
					l.Panicf(tt.printf, tt.msg)
				}
			}()
		case log.AlertLevel:
			if len(tt.printf) == 0 {
				l.Alert(tt.msg)
			} else {
				l.Alertf(tt.printf, tt.msg)
			}
		}

		if s := hasString(conn); !strings.HasSuffix(s, tt.want) {

			if tt.lvl == log.TraceLevel {
				if !strings.Contains(s, tt.want) {
					t.Errorf("test %d: Expected Contains '%s' Got '%s'", i, tt.want, s)
				}
				continue
			}

			t.Errorf("test %d: Expected Suffix '%s' Got '%s'", i, tt.want, s)
		}
	}

	e := &log.Entry{
		Level:     log.FatalLevel,
		Message:   "fatal",
		Timestamp: time.Now().UTC(),
		Line:      259,
		File:      "syslog_test.go",
	}

	log.HandleEntry(e)

	if s := hasString(conn); !strings.Contains(s, "UTC [31m[4m[5m FATAL[0m syslog_test.go:259 fatal\n") {
		t.Errorf("test fatal: Expected Contains '%s' Got '%s'", "UTC [31m[4m[5m FATAL[0m syslog_test.go:259 fatal\n", s)
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
	defer conn.Close()

	sLog, err := New("udp", "127.0.0.1:2004", "", nil)
	if err != nil {
		log.Errorf("Expected '%v' Got '%v'", nil, err)
	}

	sLog.SetDisplayColor(true)
	sLog.SetBuffersAndWorkers(3, 0)
	sLog.SetTimestampFormat("2006")
	sLog.SetFormatFunc(func(s *Syslog) Formatter {
		return func(e *log.Entry) []byte {
			return []byte(e.Message)
		}
	})

	log.RegisterHandler(sLog, log.AllLevels...)

	log.Debug("debug")
	if s := hasString(conn); s != "debug" {
		log.Errorf("Expected '%s' Got '%s'", "debug", s)
	}
}

func TestSetFilename(t *testing.T) {

	addr, err := net.ResolveUDPAddr("udp", ":2005")
	if err != nil {
		log.Errorf("Expected '%v' Got '%v'", nil, err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Errorf("Expected '%v' Got '%v'", nil, err)
	}
	defer conn.Close()

	sLog, err := New("udp", "127.0.0.1:2005", "", nil)
	if err != nil {
		log.Errorf("Expected '%v' Got '%v'", nil, err)
	}

	sLog.SetDisplayColor(false)
	sLog.SetBuffersAndWorkers(3, 1)
	sLog.SetTimestampFormat("MST")
	sLog.SetFilenameDisplay(log.Llongfile)

	log.RegisterHandler(sLog, log.AllLevels...)

	log.Error("error")
	if s := hasString(conn); !strings.Contains(s, "log/handlers/syslog/syslog_test.go:337 error") {
		t.Errorf("Expected '%s' Got '%s'", "log/handlers/syslog/syslog_test.go:337 error", s)
	}
}

func TestSetFilenameColor(t *testing.T) {
	addr, err := net.ResolveUDPAddr("udp", ":2006")
	if err != nil {
		log.Errorf("Expected '%v' Got '%v'", nil, err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Errorf("Expected '%v' Got '%v'", nil, err)
	}
	defer conn.Close()

	sLog, err := New("udp", "127.0.0.1:2006", "", nil)
	if err != nil {
		log.Errorf("Expected '%v' Got '%v'", nil, err)
	}

	sLog.SetDisplayColor(true)
	sLog.SetBuffersAndWorkers(3, 1)
	sLog.SetTimestampFormat("MST")
	sLog.SetFilenameDisplay(log.Llongfile)

	log.RegisterHandler(sLog, log.AllLevels...)

	log.Error("error")
	if s := hasString(conn); !strings.Contains(s, "log/handlers/syslog/syslog_test.go:367 error") {
		t.Errorf("Expected '%s' Got '%s'", "log/handlers/syslog/syslog_test.go:367 error", s)
	}
}

func TestSyslogTLS(t *testing.T) {

	// setup server

	addr, err := net.ResolveTCPAddr("tcp", ":2022")
	if err != nil {
		log.Errorf("Expected '%v' Got '%v'", nil, err)
	}

	cnn, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Errorf("Expected '%v' Got '%v'", nil, err)
	}
	defer cnn.Close()

	tlsConfig := &tls.Config{
		Certificates: make([]tls.Certificate, 1),
	}

	tlsConfig.Certificates[0], err = tls.X509KeyPair([]byte(publicKey), []byte(privateKey))
	if err != nil {
		log.Errorf("Expected '%v' Got '%v'", nil, err)
	}

	conn := tls.NewListener(cnn, tlsConfig)

	var msg string
	var m sync.Mutex

	go func() {
		client, err := conn.Accept()
		if err != nil {
			log.Errorf("Expected '%v' Got '%v'", nil, err)
		}

		b := make([]byte, 1024)

		read, err := client.Read(b)
		if err != nil {
			log.Errorf("Expected '%v' Got '%v'", nil, err)
		}

		m.Lock()
		defer m.Unlock()
		msg = string(b[0:read])
	}()

	// setup client log

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM([]byte(publicKey))

	config := &tls.Config{
		RootCAs: pool,
	}

	sLog, err := New("tcp+tls", "127.0.0.1:2022", "", config)
	if err != nil {
		log.Fatalf("Expected '%v' Got '%v'", nil, err)
	}

	//sLog.SetDisplayColor(true)
	// sLog.SetBuffersAndWorkers(3, 0)
	// sLog.SetTimestampFormat("2006")
	sLog.SetFormatFunc(func(s *Syslog) Formatter {
		return func(e *log.Entry) []byte {
			return []byte(e.Message)
		}
	})

	log.RegisterHandler(sLog, log.AllLevels...)

	log.Debug("debug")
	time.Sleep(500 * time.Millisecond)

	m.Lock()
	defer m.Unlock()

	if msg != "debug" {
		log.Errorf("Expected '%s' Got '%s'", "debug", msg)
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
			want: "UTC  DEBUG debug\n",
		},
		{
			lvl:    log.DebugLevel,
			msg:    "debugf",
			printf: "%s",
			flds:   nil,
			want:   "UTC  DEBUG debugf\n",
		},
		{
			lvl:  log.InfoLevel,
			msg:  "info",
			flds: nil,
			want: "UTC   INFO info\n",
		},
		{
			lvl:    log.InfoLevel,
			msg:    "infof",
			printf: "%s",
			flds:   nil,
			want:   "UTC   INFO infof\n",
		},
		{
			lvl:  log.NoticeLevel,
			msg:  "notice",
			flds: nil,
			want: "UTC NOTICE notice\n",
		},
		{
			lvl:    log.NoticeLevel,
			msg:    "noticef",
			printf: "%s",
			flds:   nil,
			want:   "UTC NOTICE noticef\n",
		},
		{
			lvl:  log.WarnLevel,
			msg:  "warn",
			flds: nil,
			want: "UTC   WARN syslog_test.go:101 warn\n",
		},
		{
			lvl:    log.WarnLevel,
			msg:    "warnf",
			printf: "%s",
			flds:   nil,
			want:   "UTC   WARN syslog_test.go:103 warnf\n",
		},
		{
			lvl:  log.ErrorLevel,
			msg:  "error",
			flds: nil,
			want: "UTC  ERROR syslog_test.go:107 error\n",
		},
		{
			lvl:    log.ErrorLevel,
			msg:    "errorf",
			printf: "%s",
			flds:   nil,
			want:   "UTC  ERROR syslog_test.go:109 errorf\n",
		},
		{
			lvl:  log.AlertLevel,
			msg:  "alert",
			flds: nil,
			want: "UTC  ALERT syslog_test.go:125 alert\n",
		},
		{
			lvl:    log.AlertLevel,
			msg:    "alertf",
			printf: "%s",
			flds:   nil,
			want:   "UTC  ALERT syslog_test.go:127 alertf\n",
		},
		{
			lvl:  log.PanicLevel,
			msg:  "panic",
			flds: nil,
			want: "UTC  PANIC syslog_test.go:118 panic\n",
		},
		{
			lvl:    log.PanicLevel,
			msg:    "panicf",
			printf: "%s",
			flds:   nil,
			want:   "UTC  PANIC syslog_test.go:120 panicf\n",
		},
		{
			lvl: log.DebugLevel,
			msg: "debug",
			flds: []log.Field{
				log.F("key", "value"),
			},
			want: "UTC  DEBUG debug key=value\n",
		},
		{
			lvl:    log.DebugLevel,
			msg:    "debugf",
			printf: "%s",
			flds: []log.Field{
				log.F("key", "value"),
			},
			want: "UTC  DEBUG debugf key=value\n",
		},
		{
			lvl: log.InfoLevel,
			msg: "info",
			flds: []log.Field{
				log.F("key", "value"),
			},
			want: "UTC   INFO info key=value\n",
		},
		{
			lvl:    log.InfoLevel,
			msg:    "infof",
			printf: "%s",
			flds: []log.Field{
				log.F("key", "value"),
			},
			want: "UTC   INFO infof key=value\n",
		},
		{
			lvl: log.NoticeLevel,
			msg: "notice",
			flds: []log.Field{
				log.F("key", "value"),
			},
			want: "UTC NOTICE notice key=value\n",
		},
		{
			lvl:    log.NoticeLevel,
			msg:    "noticef",
			printf: "%s",
			flds: []log.Field{
				log.F("key", "value"),
			},
			want: "UTC NOTICE noticef key=value\n",
		},
		{
			lvl: log.WarnLevel,
			msg: "warn",
			flds: []log.Field{
				log.F("key", "value"),
			},
			want: "UTC   WARN syslog_test.go:101 warn key=value\n",
		},
		{
			lvl:    log.WarnLevel,
			msg:    "warnf",
			printf: "%s",
			flds: []log.Field{
				log.F("key", "value"),
			},
			want: "UTC   WARN syslog_test.go:103 warnf key=value\n",
		},
		{
			lvl: log.ErrorLevel,
			msg: "error",
			flds: []log.Field{
				log.F("key", "value"),
			},
			want: "UTC  ERROR syslog_test.go:107 error key=value\n",
		},
		{
			lvl:    log.ErrorLevel,
			msg:    "errorf",
			printf: "%s",
			flds: []log.Field{
				log.F("key", "value"),
			},
			want: "UTC  ERROR syslog_test.go:109 errorf key=value\n",
		},
		{
			lvl: log.AlertLevel,
			msg: "alert",
			flds: []log.Field{
				log.F("key", "value"),
			},
			want: "UTC  ALERT syslog_test.go:125 alert key=value\n",
		},
		{
			lvl: log.AlertLevel,
			msg: "alert",
			flds: []log.Field{
				log.F("key", "value"),
			},
			want: "UTC  ALERT syslog_test.go:125 alert key=value\n",
		},
		{
			lvl:    log.AlertLevel,
			msg:    "alertf",
			printf: "%s",
			flds: []log.Field{
				log.F("key", "value"),
			},
			want: "UTC  ALERT syslog_test.go:127 alertf key=value\n",
		},
		{
			lvl:    log.PanicLevel,
			msg:    "panicf",
			printf: "%s",
			flds: []log.Field{
				log.F("key", "value"),
			},
			want: "UTC  PANIC syslog_test.go:120 panicf key=value\n",
		},
		{
			lvl: log.PanicLevel,
			msg: "panic",
			flds: []log.Field{
				log.F("key", "value"),
			},
			want: "UTC  PANIC syslog_test.go:118 panic key=value\n",
		},
		{
			lvl:  log.TraceLevel,
			msg:  "trace",
			flds: nil,
			want: "UTC  TRACE trace",
		},
		{
			lvl:    log.TraceLevel,
			msg:    "tracef",
			printf: "%s",
			flds:   nil,
			want:   "UTC  TRACE tracef",
		},
		{
			lvl: log.TraceLevel,
			msg: "trace",
			flds: []log.Field{
				log.F("key", "value"),
			},
			want: "UTC  TRACE trace key=value",
		},
		{
			lvl:    log.TraceLevel,
			msg:    "tracef",
			printf: "%s",
			flds: []log.Field{
				log.F("key", "value"),
			},
			want: "UTC  TRACE tracef key=value",
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
			want: "UTC  DEBUG debug key=string key=1 key=2 key=3 key=4 key=5 key=1 key=2 key=3 key=4 key=5 key=true key={struct}\n",
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
			want: "UTC [32m DEBUG[0m debugf\n",
		},
		{
			lvl:  log.DebugLevel,
			msg:  "debug",
			flds: nil,
			want: "UTC [32m DEBUG[0m debug\n",
		},
		{
			lvl:    log.InfoLevel,
			msg:    "infof",
			printf: "%s",
			flds:   nil,
			want: "UTC [34m  INFO[0m infof\n",
		},
		{
			lvl:  log.InfoLevel,
			msg:  "info",
			flds: nil,
			want: "UTC [34m  INFO[0m info\n",
		},
		{
			lvl:    log.NoticeLevel,
			msg:    "noticef",
			printf: "%s",
			flds:   nil,
			want: "UTC [36;1mNOTICE[0m noticef\n",
		},
		{
			lvl:  log.NoticeLevel,
			msg:  "notice",
			flds: nil,
			want: "UTC [36;1mNOTICE[0m notice\n",
		},
		{
			lvl:    log.WarnLevel,
			msg:    "warnf",
			printf: "%s",
			flds:   nil,
			want: "UTC [33;1m  WARN[0m syslog_test.go:210 warnf\n",
		},
		{
			lvl:  log.WarnLevel,
			msg:  "warn",
			flds: nil,
			want: "UTC [33;1m  WARN[0m syslog_test.go:208 warn\n",
		},
		{
			lvl:    log.ErrorLevel,
			msg:    "errorf",
			printf: "%s",
			flds:   nil,
			want: "UTC [31;1m ERROR[0m syslog_test.go:216 errorf\n",
		},
		{
			lvl:  log.ErrorLevel,
			msg:  "error",
			flds: nil,
			want: "UTC [31;1m ERROR[0m syslog_test.go:214 error\n",
		},
		{
			lvl:    log.AlertLevel,
			msg:    "alertf",
			printf: "%s",
			flds:   nil,
			want: "UTC [31m[4m ALERT[0m syslog_test.go:234 alertf\n",
		},
		{
			lvl:  log.AlertLevel,
			msg:  "alert",
			flds: nil,
			want: "UTC [31m[4m ALERT[0m syslog_test.go:232 alert\n",
		},
		{
			lvl:    log.PanicLevel,
			msg:    "panicf",
			printf: "%s",
			flds:   nil,
			want: "UTC [31m PANIC[0m syslog_test.go:227 panicf\n",
		},
		{
			lvl:  log.PanicLevel,
			msg:  "panic",
			flds: nil,
			want: "UTC [31m PANIC[0m syslog_test.go:225 panic\n",
		},
		{
			lvl:    log.DebugLevel,
			msg:    "debugf",
			printf: "%s",
			flds: []log.Field{
				log.F("key", "value"),
			},
			want: "UTC [32m DEBUG[0m debugf [32mkey[0m=value\n",
		},
		{
			lvl: log.DebugLevel,
			msg: "debug",
			flds: []log.Field{
				log.F("key", "value"),
			},
			want: "UTC [32m DEBUG[0m debug [32mkey[0m=value\n",
		},
		{
			lvl:    log.InfoLevel,
			msg:    "infof",
			printf: "%s",
			flds: []log.Field{
				log.F("key", "value"),
			},
			want: "UTC [34m  INFO[0m infof [34mkey[0m=value\n",
		},
		{
			lvl: log.InfoLevel,
			msg: "info",
			flds: []log.Field{
				log.F("key", "value"),
			},
			want: "UTC [34m  INFO[0m info [34mkey[0m=value\n",
		},
		{
			lvl:    log.NoticeLevel,
			msg:    "noticef",
			printf: "%s",
			flds: []log.Field{
				log.F("key", "value"),
			},
			want: "UTC [36;1mNOTICE[0m noticef [36;1mkey[0m=value\n",
		},
		{
			lvl: log.NoticeLevel,
			msg: "notice",
			flds: []log.Field{
				log.F("key", "value"),
			},
			want: "UTC [36;1mNOTICE[0m notice [36;1mkey[0m=value\n",
		},
		{
			lvl:    log.WarnLevel,
			msg:    "warnf",
			printf: "%s",
			flds: []log.Field{
				log.F("key", "value"),
			},
			want: "UTC [33;1m  WARN[0m syslog_test.go:210 warnf [33;1mkey[0m=value\n",
		},
		{
			lvl: log.WarnLevel,
			msg: "warn",
			flds: []log.Field{
				log.F("key", "value"),
			},
			want: "UTC [33;1m  WARN[0m syslog_test.go:208 warn [33;1mkey[0m=value\n",
		},
		{
			lvl:    log.ErrorLevel,
			msg:    "errorf",
			printf: "%s",
			flds: []log.Field{
				log.F("key", "value"),
			},
			want: "UTC [31;1m ERROR[0m syslog_test.go:216 errorf [31;1mkey[0m=value\n",
		},
		{
			lvl: log.ErrorLevel,
			msg: "error",
			flds: []log.Field{
				log.F("key", "value"),
			},
			want: "UTC [31;1m ERROR[0m syslog_test.go:214 error [31;1mkey[0m=value\n",
		},
		{
			lvl:    log.AlertLevel,
			msg:    "alertf",
			printf: "%s",
			flds: []log.Field{
				log.F("key", "value"),
			},
			want: "UTC [31m[4m ALERT[0m syslog_test.go:234 alertf [31m[4mkey[0m=value\n",
		},
		{
			lvl: log.AlertLevel,
			msg: "alert",
			flds: []log.Field{
				log.F("key", "value"),
			},
			want: "UTC [31m[4m ALERT[0m syslog_test.go:232 alert [31m[4mkey[0m=value\n",
		},
		{
			lvl:    log.PanicLevel,
			msg:    "panicf",
			printf: "%s",
			flds: []log.Field{
				log.F("key", "value"),
			},
			want: "UTC [31m PANIC[0m syslog_test.go:227 panicf [31mkey[0m=value\n",
		},
		{
			lvl: log.PanicLevel,
			msg: "panic",
			flds: []log.Field{
				log.F("key", "value"),
			},
			want: "UTC [31m PANIC[0m syslog_test.go:225 panic [31mkey[0m=value\n",
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
			want: "UTC [32m DEBUG[0m debug [32mkey[0m=string [32mkey[0m=1 [32mkey[0m=2 [32mkey[0m=3 [32mkey[0m=4 [32mkey[0m=5 [32mkey[0m=1 [32mkey[0m=2 [32mkey[0m=3 [32mkey[0m=4 [32mkey[0m=5 [32mkey[0m=5.33 [32mkey[0m=5.34 [32mkey[0m=true [32mkey[0m={struct}\n",
		},
	}
}

const (
	publicKey = `-----BEGIN CERTIFICATE-----
MIIC4DCCAcigAwIBAgIQQiBLO9euTB+dokvZeAoNQjANBgkqhkiG9w0BAQsFADAU
MRIwEAYDVQQDDAkxMjcuMC4wLjEwHhcNMTYwNzIwMTA0MjA4WhcNMTgxMjMxMDAw
MDAwWjAUMRIwEAYDVQQDDAkxMjcuMC4wLjEwggEiMA0GCSqGSIb3DQEBAQUAA4IB
DwAwggEKAoIBAQDbrVoGAj4NplB4qzHa7GgYCk2kbsLrb15POIRVmsgpfsLumMsJ
xjI+/wW+ff6SLvmEoONhiQyhUzzqUm9XlC+CWm0rAPZabT0P11WFxL0fFW9lANSa
rvmcIxAo54Vx0VugXntFSV7MTPhu08Xy6kU6d0D2+t6eY8qNsaP0CSGLXHiDMQH4
g5mtL5RUtDb0JRBCX96BaMj3Y9T9oYPQf++tlsv4QbdvgvUVkRz9T66hmLOSpN3I
VlWKbSUtjysFR1Q5TokJmQy4VkRIjKUsrhoOayFlbzXAz40qumvOUKdB2H2YNbY3
yXCTD0O2Em/Hwav1XDueRX3IEKdWDauPEBpfAgMBAAGjLjAsMA8GA1UdEQQIMAaH
BH8AAAEwCwYDVR0PBAQDAgGuMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQELBQAD
ggEBAGr/EQC2O/nPLLTQaZ9sio1S7aT0wsRGvqReeyMBeRPB1OlFa1EYtHbs9tAz
Hsh6jUqTEyxQ3aorx1bIwuKcbKmw+G4zLf0r7v0M59xIJzECaV9zHqNxenYw4i2L
53R3yu7Q9QvWqFgpHRErk2+F7CoUGsCsmI+J90Pk1KwPcEieoRJJc9di9z7wUORa
QciCwomYnFl/5dQLepV2oIEWCEFwm1oFzSuIboffK3T3vzIYDe424b2mo+/O2J4l
lQFRwOT4kY5o246rn/sm4ag/hD5sQxu1tairnyLeSBSYLq0qef0VGIF78L5QonFR
YwYOENQ14l1n9/lbTd4WBsFqHNA=
-----END CERTIFICATE-----`
	privateKey = `-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDbrVoGAj4NplB4
qzHa7GgYCk2kbsLrb15POIRVmsgpfsLumMsJxjI+/wW+ff6SLvmEoONhiQyhUzzq
Um9XlC+CWm0rAPZabT0P11WFxL0fFW9lANSarvmcIxAo54Vx0VugXntFSV7MTPhu
08Xy6kU6d0D2+t6eY8qNsaP0CSGLXHiDMQH4g5mtL5RUtDb0JRBCX96BaMj3Y9T9
oYPQf++tlsv4QbdvgvUVkRz9T66hmLOSpN3IVlWKbSUtjysFR1Q5TokJmQy4VkRI
jKUsrhoOayFlbzXAz40qumvOUKdB2H2YNbY3yXCTD0O2Em/Hwav1XDueRX3IEKdW
DauPEBpfAgMBAAECggEBAK5fQuckHoeNLbErCs7g+puiihDszpI9e5ncncaprxqp
ASiNZhVjGn1Axxl3P3xgBzXM09CXDcx8mwzQ1IqrGK8bAi6xe9s5fM+3OK6PBSPI
SvzclOYX4BCdEHW3mQhIi7eXZ7gOzk3TBxxJw4XXiY4oHQwvBEiro5unlyHdoZ/R
FUZhXvy5E5a0/bkcCgCzYmzQPqCfB/+xCCC1P3TBNt7x6qOmTVJ7Z5TFdF47R9Ld
/G72+p96x50YJpOhfyvxqSIi/79NQVQsf8kLGo62PuCsJ2pT6V1Hb2LTDGLbUhRX
/45ycTIE9YzTso1mWIWInl4wN6NsUPyGeY/kWc9cTdkCgYEA85pPN7FPOtvj/6gl
iOL8YutkplLfhmf4B9tqSIWMsOaEK69KSJC7BTLEEZ/EMypba50fJ/HL5CfYAYz3
7YgIoresUd4gDCmA4TaaCOC/WXeuRBzyNYn2A66yfufv3T+cnjER+q6l8x2mf5Pc
Y+L3kvAUvjCcrgbhDvcrY3eZW3sCgYEA5ttWAnyCbp2Fw/6oNICsfifHpO8AbRHp
2dndMc4N9rKY3t0dwPXCSEwgedIgj76vSBUUR/39vltBmG90Diz7lYful2Uzs0Zp
JJWY5CU6rosohIiq4vef5J8uHzGkhfYoOmlMyTqT3YJvJFrv4IUS1PkF2nLU6tEF
DAn+3B6URW0CgYBsMuzeqsWrOgHiCxho3ZEGitFQwtx/gWx8aOujPJZJ+IlaMeiH
pKk83NiTj2gA5d5nRQmSn2ZVd5EM10VD3rkfNP+3+TY40LJq1erC6Lh1D6B6pnS6
bQW1iwHDNlem6NsytE7tDmetPU03u0AXqbcXL8W22DavYWTTVduSuYuHQwKBgG1o
F1v4TAxGRQW841R2gskK6zfEOOx3597hvE2FPOLkg0RjgF1ZWyjOQznYlqvpD8LW
kpUHz0BumSi38UVilhyoni9Lu/PDc8LtztaYujXMJ3igGHSWLEW6Fq6b5T/DiA8e
plBbnYYF8cxF+JbsGh+qoNaFQ1jBlGW/OvRw3Y4FAoGBAKMeQFBTPw6/JN0yAJld
rkahdsRBeRYQQPGv4O2fHN50E9IKZrcF7wFOs+xK6KYJVl0UexxRYz/Z7d2zWHFC
8GDWydoyRAicwuyldIbL+uA517KCzreYWnVumo5gpkkJeb2L0Kr3lk3/pz4Kfq1b
iJo4hgFbcvqHpSPfZDTLK9RX
-----END PRIVATE KEY-----`
)
