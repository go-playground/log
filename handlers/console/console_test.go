package console

import (
	"bytes"
	"io"
	stdlog "log"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/go-playground/log/v7"
)

// NOTES:
// - Run "go test" to run tests
// - Run "gocov test | gocov report" to report on test converage by file
// - Run "gocov test | gocov annotate -" to report on all code and functions, those ,marked with "MISS" were never called
//
// or
// -- may be a good idea to change to output path to somewherelike /tmp
// go test -coverprofile cover.out && go tool cover -html=cover.out -o cover.html

func TestConsoleLogger(t *testing.T) {
	tests := getConsoleLoggerTests()
	buff := new(buffer)

	log.SetExitFunc(func(int) {})

	cLog := New(false)
	cLog.SetWriter(buff)
	cLog.SetDisplayColor(false)
	cLog.SetTimestampFormat("")
	log.AddHandler(cLog, log.AllLevels...)
	defer func() { _ = cLog.Close() }()
	for i, tt := range tests {

		buff.Reset()
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

		if buff.String() != tt.want {
			t.Errorf("test %d: Expected '%s' Got '%s'", i, tt.want, buff.String())
		}
	}
}

func TestConsoleLoggerColor(t *testing.T) {
	log.SetExitFunc(func(int) {})
	tests := getConsoleLoggerColorTests()
	buff := new(buffer)
	cLog := New(false)
	cLog.SetWriter(buff)
	cLog.SetDisplayColor(true)
	cLog.SetTimestampFormat("")

	log.AddHandler(cLog, log.AllLevels...)

	for i, tt := range tests {

		buff.Reset()
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
		case log.AlertLevel:
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
}

func TestConsoleSTDLogCapturing(t *testing.T) {
	buff := new(buffer)
	cLog := New(true)
	cLog.SetDisplayColor(false)
	cLog.SetTimestampFormat("MST")
	cLog.SetWriter(buff)
	log.AddHandler(cLog, log.AllLevels...)

	stdlog.Println("STD LOG message")

	time.Sleep(1000 * time.Millisecond)

	s := buff.String()

	expected := "STD LOG message"

	if !strings.Contains(s, expected) {
		t.Errorf("Expected '%s' Got '%s'", expected, s)
	}
}

type test struct {
	lvl    log.Level
	msg    string
	flds   []log.Field
	want   string
	printf string
}

func getConsoleLoggerTests() []test {
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

func getConsoleLoggerColorTests() []test {
	return []test{
		{
			lvl:    log.DebugLevel,
			msg:    "debugf",
			printf: "%s",
			flds:   nil,
			want: " [32m DEBUG[0m debugf\n",
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

type buffer struct {
	b bytes.Buffer
	m sync.Mutex
}

func (b *buffer) Read(p []byte) (n int, err error) {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.Read(p)
}
func (b *buffer) Write(p []byte) (n int, err error) {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.Write(p)
}
func (b *buffer) String() string {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.String()
}

func (b *buffer) Bytes() []byte {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.Bytes()
}
func (b *buffer) Cap() int {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.Cap()
}
func (b *buffer) Grow(n int) {
	b.m.Lock()
	defer b.m.Unlock()
	b.b.Grow(n)
}
func (b *buffer) Len() int {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.Len()
}
func (b *buffer) Next(n int) []byte {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.Next(n)
}
func (b *buffer) ReadByte() (c byte, err error) {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.ReadByte()
}
func (b *buffer) ReadBytes(delim byte) (line []byte, err error) {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.ReadBytes(delim)
}
func (b *buffer) ReadFrom(r io.Reader) (n int64, err error) {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.ReadFrom(r)
}
func (b *buffer) ReadRune() (r rune, size int, err error) {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.ReadRune()
}
func (b *buffer) ReadString(delim byte) (line string, err error) {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.ReadString(delim)
}
func (b *buffer) Reset() {
	b.m.Lock()
	defer b.m.Unlock()
	b.b.Reset()
}
func (b *buffer) Truncate(n int) {
	b.m.Lock()
	defer b.m.Unlock()
	b.b.Truncate(n)
}
func (b *buffer) UnreadByte() error {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.UnreadByte()
}
func (b *buffer) UnreadRune() error {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.UnreadRune()
}
func (b *buffer) WriteByte(c byte) error {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.WriteByte(c)
}
func (b *buffer) WriteRune(r rune) (n int, err error) {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.WriteRune(r)
}
func (b *buffer) WriteString(s string) (n int, err error) {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.WriteString(s)
}
func (b *buffer) WriteTo(w io.Writer) (n int64, err error) {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.WriteTo(w)
}
