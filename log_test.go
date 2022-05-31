package log

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/go-playground/errors/v5"
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

type testHandler struct {
	writer io.Writer
}

func (th *testHandler) Log(e Entry) {
	s := e.Level.String() + " "
	s += e.Message
	for _, f := range e.Fields {
		s += fmt.Sprintf(" %s=%v", f.Key, f.Value)
	}

	s += "\n"
	if _, err := th.writer.Write([]byte(s)); err != nil {
		panic(err)
	}
}

func TestConsoleLogger1(t *testing.T) {
	SetExitFunc(func(int) {})
	SetWithErrorFn(errorsWithError)

	tests := getLogTests1()
	buff := new(bytes.Buffer)
	th := &testHandler{
		writer: buff,
	}
	logHandlers = map[Level][]Handler{}
	AddHandler(th, AllLevels...)
	for i, tt := range tests {
		buff.Reset()
		var l Entry

		if tt.flds != nil {
			l = l.WithFields(tt.flds...)
		}

		switch tt.lvl {
		case DebugLevel:
			if len(tt.printf) == 0 {
				l.Debug(tt.msg)
			} else {
				l.Debugf(tt.printf, tt.msg)
			}
		case InfoLevel:
			if len(tt.printf) == 0 {
				l.Info(tt.msg)
			} else {
				l.Infof(tt.printf, tt.msg)
			}
		case NoticeLevel:
			if len(tt.printf) == 0 {
				l.Notice(tt.msg)
			} else {
				l.Noticef(tt.printf, tt.msg)
			}
		case WarnLevel:
			if len(tt.printf) == 0 {
				l.Warn(tt.msg)
			} else {
				l.Warnf(tt.printf, tt.msg)
			}
		case ErrorLevel:
			if len(tt.printf) == 0 {
				l.Error(tt.msg)
			} else {
				l.Errorf(tt.printf, tt.msg)
			}
		case PanicLevel:
			if len(tt.printf) == 0 {
				l.Panic(tt.msg)
			} else {
				l.Panicf(tt.printf, tt.msg)
			}
		case AlertLevel:
			if len(tt.printf) == 0 {
				l.Alert(tt.msg)
			} else {
				l.Alertf(tt.printf, tt.msg)
			}
		}

		if buff.String() != tt.want {
			t.Errorf("test %d: Expected '%s' Got '%s'", i, tt.want, buff.String())
		}

	}

	buff.Reset()
	// Test Custom Entry ( most common case is Unmarshalled from JSON when using centralized logging)
	var entry Entry
	entry.Level = InfoLevel
	entry.Timestamp = time.Now().UTC()
	entry.Message = "Test Message"
	entry.Fields = make([]Field, 0)
	HandleEntry(entry)

	if buff.String() != "INFO Test Message\n" {
		t.Errorf("test Custom Entry: Expected '%s' Got '%s'", "INFO Test Message\n", buff.String())
	}
}

func TestConsoleLogger2(t *testing.T) {
	SetExitFunc(func(int) {})
	tests := getLogTests()
	buff := new(bytes.Buffer)
	th := &testHandler{
		writer: buff,
	}
	logHandlers = map[Level][]Handler{}
	AddHandler(th, AllLevels...)

	for i, tt := range tests {
		buff.Reset()
		var l Entry

		if tt.flds != nil {
			l = WithFields(tt.flds...)

			switch tt.lvl {
			case DebugLevel:
				if len(tt.printf) == 0 {
					l.Debug(tt.msg)
				} else {
					l.Debugf(tt.printf, tt.msg)
				}
			case InfoLevel:
				if len(tt.printf) == 0 {
					l.Info(tt.msg)
				} else {
					l.Infof(tt.printf, tt.msg)
				}
			case NoticeLevel:
				if len(tt.printf) == 0 {
					l.Notice(tt.msg)
				} else {
					l.Noticef(tt.printf, tt.msg)
				}
			case WarnLevel:
				if len(tt.printf) == 0 {
					l.Warn(tt.msg)
				} else {
					l.Warnf(tt.printf, tt.msg)
				}
			case ErrorLevel:
				if len(tt.printf) == 0 {
					l.Error(tt.msg)
				} else {
					l.Errorf(tt.printf, tt.msg)
				}
			case PanicLevel:
				func() {
					defer func() {
						_ = recover()
					}()

					if len(tt.printf) == 0 {
						l.Panic(tt.msg)
					} else {
						l.Panicf(tt.printf, tt.msg)
					}
				}()
			case AlertLevel:
				if len(tt.printf) == 0 {
					l.Alert(tt.msg)
				} else {
					l.Alertf(tt.printf, tt.msg)
				}
			}

		} else {

			switch tt.lvl {
			case DebugLevel:
				if len(tt.printf) == 0 {
					Debug(tt.msg)
				} else {
					Debugf(tt.printf, tt.msg)
				}
			case InfoLevel:
				if len(tt.printf) == 0 {
					Info(tt.msg)
				} else {
					Infof(tt.printf, tt.msg)
				}
			case NoticeLevel:
				if len(tt.printf) == 0 {
					Notice(tt.msg)
				} else {
					Noticef(tt.printf, tt.msg)
				}
			case WarnLevel:
				if len(tt.printf) == 0 {
					Warn(tt.msg)
				} else {
					Warnf(tt.printf, tt.msg)
				}
			case ErrorLevel:
				if len(tt.printf) == 0 {
					Error(tt.msg)
				} else {
					Errorf(tt.printf, tt.msg)
				}
			case PanicLevel:
				func() {
					defer func() {
						_ = recover()
					}()

					if len(tt.printf) == 0 {
						Panic(tt.msg)
					} else {
						Panicf(tt.printf, tt.msg)
					}
				}()
			case AlertLevel:
				if len(tt.printf) == 0 {
					Alert(tt.msg)
				} else {
					Alertf(tt.printf, tt.msg)
				}
			}
		}

		if buff.String() != tt.want {
			t.Errorf("test %d: Expected '%s' Got '%s'", i, tt.want, buff.String())
		}
	}

	buff.Reset()
	// Test Custom Entry ( most common case is Unmarshalled from JSON when using centralized logging)
	var entry Entry
	entry.Level = InfoLevel
	entry.Timestamp = time.Now().UTC()
	entry.Message = "Test Message"
	entry.Fields = make([]Field, 0)
	HandleEntry(entry)

	if buff.String() != "INFO Test Message\n" {
		t.Errorf("test Custom Entry: Expected '%s' Got '%s'", "INFO Test Message\n", buff.String())
	}
}

func TestLevel(t *testing.T) {

	tests := []struct {
		value string
		want  string
	}{
		{
			value: Level(255).String(),
			want:  "Unknown Level",
		},
		{
			value: DebugLevel.String(),
			want:  "DEBUG",
		},
		{
			value: InfoLevel.String(),
			want:  "INFO",
		},
		{
			value: NoticeLevel.String(),
			want:  "NOTICE",
		},
		{
			value: WarnLevel.String(),
			want:  "WARN",
		},
		{
			value: ErrorLevel.String(),
			want:  "ERROR",
		},
		{
			value: PanicLevel.String(),
			want:  "PANIC",
		},
		{
			value: AlertLevel.String(),
			want:  "ALERT",
		},
		{
			value: FatalLevel.String(),
			want:  "FATAL",
		},
	}

	for i, tt := range tests {
		if tt.value != tt.want {
			t.Errorf("Test %d: Expected '%s' Got '%s'", i, tt.want, tt.value)
		}
	}
}

func TestFatal(t *testing.T) {
	var i int

	SetExitFunc(func(code int) {
		i = code
	})

	Fatal("fatal")
	if i != 1 {
		t.Errorf("test Fatals: Expected '%d' Got '%d'", 1, i)
	}

	Fatalf("fatalf")
	if i != 1 {
		t.Errorf("test Fatals: Expected '%d' Got '%d'", 1, i)
	}

	WithField("key", "value").Fatal("fatal")
	if i != 1 {
		t.Errorf("test Fatals: Expected '%d' Got '%d'", 1, i)
	}

	WithFields(F("key", "value")).Fatalf("fatalf")
	if i != 1 {
		t.Errorf("test Fatals: Expected '%d' Got '%d'", 1, i)
	}

	Entry{}.WithField("key", "value").Fatalf("fatalf")
	if i != 1 {
		t.Errorf("test Fatals: Expected '%d' Got '%d'", 1, i)
	}
}

func TestWithError(t *testing.T) {
	logHandlers = map[Level][]Handler{}
	buff := new(bytes.Buffer)
	th := &testHandler{
		writer: buff,
	}
	AddHandler(th, AllLevels...)

	terr := fmt.Errorf("this is an %s", "err")
	WithError(terr).Info()
	if !strings.HasSuffix(buff.String(), "log_test.go:361:TestWithError this is an err\n") {
		t.Errorf("Expected '%s' Got '%s'", "log_test.go:361:TestWithError this is an err\n", buff.String())
	}
	buff.Reset()
	Entry{}.WithError(terr).Info()
	if !strings.HasSuffix(buff.String(), "log_test.go:366:TestWithError this is an err\n") {
		t.Errorf("Expected '%s' Got '%s'", "log_test.go:366:TestWithError this is an err\n", buff.String())
	}
	buff.Reset()
	WithError(errors.Wrap(terr, "wrapped error")).Info()
	if !strings.HasSuffix(buff.String(), "log_test.go:371:TestWithError wrapped error: this is an err\n") || !strings.HasPrefix(buff.String(), "INFO  source=") {
		t.Errorf("Expected '%s' Got '%s'", "log_test.go:371:TestWithError wrapped error: this is an err\n", buff.String())
	}
}

func TestWithTrace(t *testing.T) {
	logHandlers = map[Level][]Handler{}
	buff := new(bytes.Buffer)
	th := &testHandler{
		writer: buff,
	}
	AddHandler(th, AllLevels...)

	WithTrace().Info("info")
	if !strings.HasPrefix(buff.String(), "INFO info duration=") {
		t.Errorf("Expected '%s' Got '%s'", "INFO info duration=", buff.String())
	}

	Entry{}.WithTrace().Info("info")
	if !strings.HasPrefix(buff.String(), "INFO info duration=") {
		t.Errorf("Expected '%s' Got '%s'", "INFO info duration=", buff.String())
	}
}

func TestDefaultFields(t *testing.T) {
	logHandlers = map[Level][]Handler{}
	buff := new(bytes.Buffer)
	th := &testHandler{
		writer: buff,
	}
	AddHandler(th, AllLevels...)
	WithDefaultFields(F("key", "value"))
	Info("info")
	if buff.String() != "INFO info key=value\n" {
		t.Errorf("Expected '%s' Got '%s'", "INFO info key=value\n", buff.String())
	}
}

func getLogTests() []test {
	return []test{
		{
			lvl:  PanicLevel,
			msg:  "panicln",
			flds: nil,
			want: "PANIC panicln\n",
		},
		{
			lvl:  DebugLevel,
			msg:  "debug",
			flds: nil,
			want: "DEBUG debug\n",
		},
		{
			lvl:    DebugLevel,
			msg:    "debugf",
			printf: "%s",
			flds:   nil,
			want:   "DEBUG debugf\n",
		},
		{
			lvl:  InfoLevel,
			msg:  "info",
			flds: nil,
			want: "INFO info\n",
		},
		{
			lvl:    InfoLevel,
			msg:    "infof",
			printf: "%s",
			flds:   nil,
			want:   "INFO infof\n",
		},
		{
			lvl:  NoticeLevel,
			msg:  "notice",
			flds: nil,
			want: "NOTICE notice\n",
		},
		{
			lvl:    NoticeLevel,
			msg:    "noticef",
			printf: "%s",
			flds:   nil,
			want:   "NOTICE noticef\n",
		},
		{
			lvl:  WarnLevel,
			msg:  "warn",
			flds: nil,
			want: "WARN warn\n",
		},
		{
			lvl:    WarnLevel,
			msg:    "warnf",
			printf: "%s",
			flds:   nil,
			want:   "WARN warnf\n",
		},
		{
			lvl:  ErrorLevel,
			msg:  "error",
			flds: nil,
			want: "ERROR error\n",
		},
		{
			lvl:    ErrorLevel,
			msg:    "errorf",
			printf: "%s",
			flds:   nil,
			want:   "ERROR errorf\n",
		},
		{
			lvl:  AlertLevel,
			msg:  "alert",
			flds: nil,
			want: "ALERT alert\n",
		},
		{
			lvl:    AlertLevel,
			msg:    "alertf",
			printf: "%s",
			flds:   nil,
			want:   "ALERT alertf\n",
		},
		{
			lvl:  PanicLevel,
			msg:  "panic",
			flds: nil,
			want: "PANIC panic\n",
		},
		{
			lvl:    PanicLevel,
			msg:    "panicf",
			printf: "%s",
			flds:   nil,
			want:   "PANIC panicf\n",
		},
		{
			lvl: DebugLevel,
			msg: "debug",
			flds: []Field{
				F("key", "value"),
			},
			want: "DEBUG debug key=value\n",
		},
		{
			lvl:    DebugLevel,
			msg:    "debugf",
			printf: "%s",
			flds: []Field{
				F("key", "value"),
			},
			want: "DEBUG debugf key=value\n",
		},
		{
			lvl: InfoLevel,
			msg: "info",
			flds: []Field{
				F("key", "value"),
			},
			want: "INFO info key=value\n",
		},
		{
			lvl:    InfoLevel,
			msg:    "infof",
			printf: "%s",
			flds: []Field{
				F("key", "value"),
			},
			want: "INFO infof key=value\n",
		},
		{
			lvl: NoticeLevel,
			msg: "notice",
			flds: []Field{
				F("key", "value"),
			},
			want: "NOTICE notice key=value\n",
		},
		{
			lvl:    NoticeLevel,
			msg:    "noticef",
			printf: "%s",
			flds: []Field{
				F("key", "value"),
			},
			want: "NOTICE noticef key=value\n",
		},
		{
			lvl: WarnLevel,
			msg: "warn",
			flds: []Field{
				F("key", "value"),
			},
			want: "WARN warn key=value\n",
		},
		{
			lvl:    WarnLevel,
			msg:    "warnf",
			printf: "%s",
			flds: []Field{
				F("key", "value"),
			},
			want: "WARN warnf key=value\n",
		},
		{
			lvl: ErrorLevel,
			msg: "error",
			flds: []Field{
				F("key", "value"),
			},
			want: "ERROR error key=value\n",
		},
		{
			lvl:    ErrorLevel,
			msg:    "errorf",
			printf: "%s",
			flds: []Field{
				F("key", "value"),
			},
			want: "ERROR errorf key=value\n",
		},
		{
			lvl: AlertLevel,
			msg: "alert",
			flds: []Field{
				F("key", "value"),
			},
			want: "ALERT alert key=value\n",
		},
		{
			lvl: AlertLevel,
			msg: "alert",
			flds: []Field{
				F("key", "value"),
			},
			want: "ALERT alert key=value\n",
		},
		{
			lvl:    AlertLevel,
			msg:    "alertf",
			printf: "%s",
			flds: []Field{
				F("key", "value"),
			},
			want: "ALERT alertf key=value\n",
		},
		{
			lvl:    PanicLevel,
			msg:    "panicf",
			printf: "%s",
			flds: []Field{
				F("key", "value"),
			},
			want: "PANIC panicf key=value\n",
		},
		{
			lvl: PanicLevel,
			msg: "panic",
			flds: []Field{
				F("key", "value"),
			},
			want: "PANIC panic key=value\n",
		},
		{
			lvl: DebugLevel,
			msg: "debug",
			flds: []Field{
				F("key", "string"),
				F("key", int(1)),
				F("key", int8(2)),
				F("key", int16(3)),
				F("key", int32(4)),
				F("key", int64(5)),
				F("key", uint(1)),
				F("key", uint8(2)),
				F("key", uint16(3)),
				F("key", uint32(4)),
				F("key", uint64(5)),
				F("key", true),
				F("key", struct{ value string }{"struct"}),
			},
			want: "DEBUG debug key=string key=1 key=2 key=3 key=4 key=5 key=1 key=2 key=3 key=4 key=5 key=true key={struct}\n",
		},
	}
}

func getLogTests1() []test {
	return []test{
		{
			lvl:  PanicLevel,
			msg:  "panicln",
			flds: nil,
			want: "PANIC panicln\n",
		},
		{
			lvl:  DebugLevel,
			msg:  "debug",
			flds: nil,
			want: "DEBUG debug\n",
		},
		{
			lvl:    DebugLevel,
			msg:    "debugf",
			printf: "%s",
			flds:   nil,
			want:   "DEBUG debugf\n",
		},
		{
			lvl:  InfoLevel,
			msg:  "info",
			flds: nil,
			want: "INFO info\n",
		},
		{
			lvl:    InfoLevel,
			msg:    "infof",
			printf: "%s",
			flds:   nil,
			want:   "INFO infof\n",
		},
		{
			lvl:  NoticeLevel,
			msg:  "notice",
			flds: nil,
			want: "NOTICE notice\n",
		},
		{
			lvl:    NoticeLevel,
			msg:    "noticef",
			printf: "%s",
			flds:   nil,
			want:   "NOTICE noticef\n",
		},
		{
			lvl:  WarnLevel,
			msg:  "warn",
			flds: nil,
			want: "WARN warn\n",
		},
		{
			lvl:    WarnLevel,
			msg:    "warnf",
			printf: "%s",
			flds:   nil,
			want:   "WARN warnf\n",
		},
		{
			lvl:  ErrorLevel,
			msg:  "error",
			flds: nil,
			want: "ERROR error\n",
		},
		{
			lvl:    ErrorLevel,
			msg:    "errorf",
			printf: "%s",
			flds:   nil,
			want:   "ERROR errorf\n",
		},
		{
			lvl:  AlertLevel,
			msg:  "alert",
			flds: nil,
			want: "ALERT alert\n",
		},
		{
			lvl:    AlertLevel,
			msg:    "alertf",
			printf: "%s",
			flds:   nil,
			want:   "ALERT alertf\n",
		},
		{
			lvl:  PanicLevel,
			msg:  "panic",
			flds: nil,
			want: "PANIC panic\n",
		},
		{
			lvl:    PanicLevel,
			msg:    "panicf",
			printf: "%s",
			flds:   nil,
			want:   "PANIC panicf\n",
		},
		{
			lvl: DebugLevel,
			msg: "debug",
			flds: []Field{
				F("key", "value"),
			},
			want: "DEBUG debug key=value\n",
		},
		{
			lvl:    DebugLevel,
			msg:    "debugf",
			printf: "%s",
			flds: []Field{
				F("key", "value"),
			},
			want: "DEBUG debugf key=value\n",
		},
		{
			lvl: InfoLevel,
			msg: "info",
			flds: []Field{
				F("key", "value"),
			},
			want: "INFO info key=value\n",
		},
		{
			lvl:    InfoLevel,
			msg:    "infof",
			printf: "%s",
			flds: []Field{
				F("key", "value"),
			},
			want: "INFO infof key=value\n",
		},
		{
			lvl: NoticeLevel,
			msg: "notice",
			flds: []Field{
				F("key", "value"),
			},
			want: "NOTICE notice key=value\n",
		},
		{
			lvl:    NoticeLevel,
			msg:    "noticef",
			printf: "%s",
			flds: []Field{
				F("key", "value"),
			},
			want: "NOTICE noticef key=value\n",
		},
		{
			lvl: WarnLevel,
			msg: "warn",
			flds: []Field{
				F("key", "value"),
			},
			want: "WARN warn key=value\n",
		},
		{
			lvl:    WarnLevel,
			msg:    "warnf",
			printf: "%s",
			flds: []Field{
				F("key", "value"),
			},
			want: "WARN warnf key=value\n",
		},
		{
			lvl: ErrorLevel,
			msg: "error",
			flds: []Field{
				F("key", "value"),
			},
			want: "ERROR error key=value\n",
		},
		{
			lvl:    ErrorLevel,
			msg:    "errorf",
			printf: "%s",
			flds: []Field{
				F("key", "value"),
			},
			want: "ERROR errorf key=value\n",
		},
		{
			lvl: AlertLevel,
			msg: "alert",
			flds: []Field{
				F("key", "value"),
			},
			want: "ALERT alert key=value\n",
		},
		{
			lvl: AlertLevel,
			msg: "alert",
			flds: []Field{
				F("key", "value"),
			},
			want: "ALERT alert key=value\n",
		},
		{
			lvl:    AlertLevel,
			msg:    "alertf",
			printf: "%s",
			flds: []Field{
				F("key", "value"),
			},
			want: "ALERT alertf key=value\n",
		},
		{
			lvl:    PanicLevel,
			msg:    "panicf",
			printf: "%s",
			flds: []Field{
				F("key", "value"),
			},
			want: "PANIC panicf key=value\n",
		},
		{
			lvl: PanicLevel,
			msg: "panic",
			flds: []Field{
				F("key", "value"),
			},
			want: "PANIC panic key=value\n",
		},
		{
			lvl: DebugLevel,
			msg: "debug",
			flds: []Field{
				F("key", "string"),
				F("key", int(1)),
				F("key", int8(2)),
				F("key", int16(3)),
				F("key", int32(4)),
				F("key", int64(5)),
				F("key", uint(1)),
				F("key", uint8(2)),
				F("key", uint16(3)),
				F("key", uint32(4)),
				F("key", uint64(5)),
				F("key", true),
				F("key", struct{ value string }{"struct"}),
			},
			want: "DEBUG debug key=string key=1 key=2 key=3 key=4 key=5 key=1 key=2 key=3 key=4 key=5 key=true key={struct}\n",
		},
	}
}

func TestContext(t *testing.T) {

	e := GetContext(context.Background())
	if e.Message != "" {
		t.Errorf("Got '%s' Expected '%s'", e.Message, "")
	}

	l := WithField("key", "value")
	ctx := SetContext(context.Background(), l)
	e = GetContext(ctx)
	if l.Fields[0].Key != e.Fields[0].Key {
		t.Errorf("Got '%s' Expected '%s'", e.Fields[0].Key, "key")
	}
	if l.Fields[0].Value != e.Fields[0].Value {
		t.Errorf("Got '%s' Expected '%s'", e.Fields[0].Value, "value")
	}
}

func TestParseLevel(t *testing.T) {

	tests := []struct {
		value string
		level Level
	}{
		{
			level: Level(255),
			value: "Unknown Level",
		},
		{
			level: DebugLevel,
			value: "DEBUG",
		},
		{
			level: InfoLevel,
			value: "INFO",
		},
		{
			level: NoticeLevel,
			value: "NOTICE",
		},
		{
			level: WarnLevel,
			value: "WARN",
		},
		{
			level: ErrorLevel,
			value: "ERROR",
		},
		{
			level: PanicLevel,
			value: "PANIC",
		},
		{
			level: AlertLevel,
			value: "ALERT",
		},
		{
			level: FatalLevel,
			value: "FATAL",
		},
	}

	for i, tt := range tests {
		entry := Entry{
			Level: tt.level,
		}
		b, _ := json.Marshal(entry)
		_ = json.Unmarshal(b, &entry)
		if entry.Level != tt.level || entry.Level.String() != tt.value {
			t.Errorf("Test %d: Expected '%s' Got '%s'", i, entry.Level, tt.level)
		}
	}
}

func TestWrappedError(t *testing.T) {
	SetExitFunc(func(int) {})
	SetWithErrorFn(errorsWithError)
	logFields = logFields[0:0]
	buff := new(bytes.Buffer)
	th := &testHandler{
		writer: buff,
	}
	logHandlers = map[Level][]Handler{}
	AddHandler(th, AllLevels...)
	err := fmt.Errorf("this is an %s", "error")
	err = errors.Wrap(err, "prefix").AddTypes("Permanent", "Internal").AddTag("key", "value")
	WithError(err).Error("test")
	expected := "log_test.go:993:TestWrappedError prefix: this is an error key=value types=Permanent,Internal\n"
	if !strings.HasSuffix(buff.String(), expected) {
		t.Errorf("got %s Expected %s", buff.String(), expected)
	}
	buff.Reset()
	WithError(err).Error("test")
	if !strings.HasSuffix(buff.String(), expected) {
		t.Errorf("got %s Expected %s", buff.String(), expected)
	}
}

func TestRemoveHandler(t *testing.T) {
	SetExitFunc(func(int) {})
	SetWithErrorFn(errorsWithError)
	logFields = logFields[0:0]
	buff := new(bytes.Buffer)
	th := &testHandler{
		writer: buff,
	}
	logHandlers = map[Level][]Handler{}
	AddHandler(th, InfoLevel)
	RemoveHandler(th)
	if len(logHandlers) != 0 {
		t.Error("expected 0 handlers")
	}

	AddHandler(th, AllLevels...)
	RemoveHandler(th)
	if len(logHandlers) != 0 {
		t.Error("expected 0 handlers")
	}
}

func TestRemoveHandlerLevels(t *testing.T) {
	SetExitFunc(func(int) {})
	SetWithErrorFn(errorsWithError)
	logFields = logFields[0:0]
	buff := new(bytes.Buffer)
	th := &testHandler{
		writer: buff,
	}
	th2 := &testHandler{
		writer: buff,
	}
	logHandlers = map[Level][]Handler{}
	AddHandler(th, InfoLevel)
	RemoveHandlerLevels(th, InfoLevel)
	if len(logHandlers) != 0 {
		t.Error("expected 0 handlers")
	}

	AddHandler(th, InfoLevel)
	AddHandler(th2, InfoLevel)
	RemoveHandlerLevels(th, InfoLevel)
	if len(logHandlers) != 1 {
		t.Error("expected 1 handlers left")
	}
	if len(logHandlers[InfoLevel]) != 1 {
		t.Error("expected 1 handlers with InfoLevel left")
	}
	RemoveHandlerLevels(th2, InfoLevel)
	if len(logHandlers) != 0 {
		t.Error("expected 0 handlers")
	}

	AddHandler(th, AllLevels...)
	RemoveHandlerLevels(th, DebugLevel)
	if len(logHandlers) != 7 {
		t.Error("expected 7 log levels left")
	}

	for _, handlers := range logHandlers {
		if len(handlers) != 1 {
			t.Error("expected 1 handlers for log level")
		}
	}
}
