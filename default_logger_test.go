package log

import (
	"bytes"
	"io"
	stdlog "log"
	"strings"
	"sync"
	"testing"
	"time"
)

// NOTES:
// - Run "go test" to run tests
// - Run "gocov test | gocov report" to report on test converage by file
// - Run "gocov test | gocov annotate -" to report on all code and functions, those ,marked with "MISS" were never called
//
// or
// -- may be a good idea to change to output path to somewherelike /tmp
// go test -coverprofile cover.out && go tool cover -html=cover.out -o cover.html

func TestDefaultLogger(t *testing.T) {
	SetExitFunc(func(int) {})
	tests := getConsoleLoggerColorTests()
	buff := new(buffer)
	defaultLoggerWriter = buff
	defaultLoggerTimeFormat = ""
	cLog := newDefaultLogger()
	defer func() { cLog.Close() }()

	AddHandler(cLog, AllLevels...)

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

		if buff.String() != tt.want {
			t.Errorf("test %d: Expected '%s' Got '%s'", i, tt.want, buff.String())
		}
	}
}

func TestConsoleSTDLogCapturing(t *testing.T) {
	buff := new(buffer)
	defaultLoggerWriter = buff
	defaultLoggerTimeFormat = "MST"
	cLog := newDefaultLogger()
	defer func() { cLog.Close() }()
	AddHandler(cLog, AllLevels...)

	stdlog.Println("STD LOG message")

	time.Sleep(1000 * time.Millisecond)

	s := buff.String()

	expected := "STD LOG message"

	if !strings.Contains(s, expected) {
		t.Errorf("Expected '%s' Got '%s'", expected, s)
	}
}

func getConsoleLoggerColorTests() []test {
	return []test{
		{
			lvl:    DebugLevel,
			msg:    "debugf",
			printf: "%s",
			flds:   nil,
			want: " [32m DEBUG[0m debugf\n",
		},
		{
			lvl:  DebugLevel,
			msg:  "debug",
			flds: nil,
			want: " [32m DEBUG[0m debug\n",
		},
		{
			lvl:    InfoLevel,
			msg:    "infof",
			printf: "%s",
			flds:   nil,
			want: " [34m  INFO[0m infof\n",
		},
		{
			lvl:  InfoLevel,
			msg:  "info",
			flds: nil,
			want: " [34m  INFO[0m info\n",
		},
		{
			lvl:    NoticeLevel,
			msg:    "noticef",
			printf: "%s",
			flds:   nil,
			want: " [36;1mNOTICE[0m noticef\n",
		},
		{
			lvl:  NoticeLevel,
			msg:  "notice",
			flds: nil,
			want: " [36;1mNOTICE[0m notice\n",
		},
		{
			lvl:    WarnLevel,
			msg:    "warnf",
			printf: "%s",
			flds:   nil,
			want: " [33;1m  WARN[0m warnf\n",
		},
		{
			lvl:  WarnLevel,
			msg:  "warn",
			flds: nil,
			want: " [33;1m  WARN[0m warn\n",
		},
		{
			lvl:    ErrorLevel,
			msg:    "errorf",
			printf: "%s",
			flds:   nil,
			want: " [31;1m ERROR[0m errorf\n",
		},
		{
			lvl:  ErrorLevel,
			msg:  "error",
			flds: nil,
			want: " [31;1m ERROR[0m error\n",
		},
		{
			lvl:    AlertLevel,
			msg:    "alertf",
			printf: "%s",
			flds:   nil,
			want: " [31m[4m ALERT[0m alertf\n",
		},
		{
			lvl:  AlertLevel,
			msg:  "alert",
			flds: nil,
			want: " [31m[4m ALERT[0m alert\n",
		},
		{
			lvl:    PanicLevel,
			msg:    "panicf",
			printf: "%s",
			flds:   nil,
			want: " [31m PANIC[0m panicf\n",
		},
		{
			lvl:  PanicLevel,
			msg:  "panic",
			flds: nil,
			want: " [31m PANIC[0m panic\n",
		},
		{
			lvl:    DebugLevel,
			msg:    "debugf",
			printf: "%s",
			flds: []Field{
				F("key", "value"),
			},
			want: " [32m DEBUG[0m debugf [32mkey[0m=value\n",
		},
		{
			lvl: DebugLevel,
			msg: "debug",
			flds: []Field{
				F("key", "value"),
			},
			want: " [32m DEBUG[0m debug [32mkey[0m=value\n",
		},
		{
			lvl:    InfoLevel,
			msg:    "infof",
			printf: "%s",
			flds: []Field{
				F("key", "value"),
			},
			want: " [34m  INFO[0m infof [34mkey[0m=value\n",
		},
		{
			lvl: InfoLevel,
			msg: "info",
			flds: []Field{
				F("key", "value"),
			},
			want: " [34m  INFO[0m info [34mkey[0m=value\n",
		},
		{
			lvl:    NoticeLevel,
			msg:    "noticef",
			printf: "%s",
			flds: []Field{
				F("key", "value"),
			},
			want: " [36;1mNOTICE[0m noticef [36;1mkey[0m=value\n",
		},
		{
			lvl: NoticeLevel,
			msg: "notice",
			flds: []Field{
				F("key", "value"),
			},
			want: " [36;1mNOTICE[0m notice [36;1mkey[0m=value\n",
		},
		{
			lvl:    WarnLevel,
			msg:    "warnf",
			printf: "%s",
			flds: []Field{
				F("key", "value"),
			},
			want: " [33;1m  WARN[0m warnf [33;1mkey[0m=value\n",
		},
		{
			lvl: WarnLevel,
			msg: "warn",
			flds: []Field{
				F("key", "value"),
			},
			want: " [33;1m  WARN[0m warn [33;1mkey[0m=value\n",
		},
		{
			lvl:    ErrorLevel,
			msg:    "errorf",
			printf: "%s",
			flds: []Field{
				F("key", "value"),
			},
			want: " [31;1m ERROR[0m errorf [31;1mkey[0m=value\n",
		},
		{
			lvl: ErrorLevel,
			msg: "error",
			flds: []Field{
				F("key", "value"),
			},
			want: " [31;1m ERROR[0m error [31;1mkey[0m=value\n",
		},
		{
			lvl:    AlertLevel,
			msg:    "alertf",
			printf: "%s",
			flds: []Field{
				F("key", "value"),
			},
			want: " [31m[4m ALERT[0m alertf [31m[4mkey[0m=value\n",
		},
		{
			lvl: AlertLevel,
			msg: "alert",
			flds: []Field{
				F("key", "value"),
			},
			want: " [31m[4m ALERT[0m alert [31m[4mkey[0m=value\n",
		},
		{
			lvl:    PanicLevel,
			msg:    "panicf",
			printf: "%s",
			flds: []Field{
				F("key", "value"),
			},
			want: " [31m PANIC[0m panicf [31mkey[0m=value\n",
		},
		{
			lvl: PanicLevel,
			msg: "panic",
			flds: []Field{
				F("key", "value"),
			},
			want: " [31m PANIC[0m panic [31mkey[0m=value\n",
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
				F("key", float32(5.33)),
				F("key", float64(5.34)),
				F("key", true),
				F("key", struct{ value string }{"struct"}),
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
