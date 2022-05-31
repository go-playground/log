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

//NOTES:
//- Run "go test" to run tests
//- Run "gocov test | gocov report" to report on test converage by file
//- Run "gocov test | gocov annotate -" to report on all code and functions, those ,marked with "MISS" were never called
//
//or
//-- may be a good idea to change to output path to somewherelike /tmp
//go test -coverprofile cover.out && go tool cover -html=cover.out -o cover.html

func TestConsoleLogger(t *testing.T) {
	tests := getConsoleLoggerTests()
	buff := new(buffer)

	SetExitFunc(func(int) {})

	cLog := NewBuilder().WithWriter(buff).WithTimestampFormat("").Build()
	AddHandler(cLog, AllLevels...)
	defer func() { _ = cLog.Close() }()
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
}

func TestConsoleSTDLogCapturing(t *testing.T) {
	buff := new(buffer)
	cLog := NewBuilder().WithWriter(buff).WithTimestampFormat("MST").Build()
	AddHandler(cLog, AllLevels...)

	stdlog.Println("STD LOG message")

	time.Sleep(1000 * time.Millisecond)

	s := buff.String()

	expected := "STD LOG message"

	if !strings.Contains(s, expected) {
		t.Errorf("Expected '%s' Got '%s'", expected, s)
	}
}

type test struct {
	lvl    Level
	msg    string
	flds   []Field
	want   string
	printf string
}

func getConsoleLoggerTests() []test {
	return []test{
		{
			lvl:  DebugLevel,
			msg:  "debug",
			flds: nil,
			want: "  DEBUG debug\n",
		},
		{
			lvl:    DebugLevel,
			msg:    "debugf",
			printf: "%s",
			flds:   nil,
			want:   "  DEBUG debugf\n",
		},
		{
			lvl:  InfoLevel,
			msg:  "info",
			flds: nil,
			want: "   INFO info\n",
		},
		{
			lvl:    InfoLevel,
			msg:    "infof",
			printf: "%s",
			flds:   nil,
			want:   "   INFO infof\n",
		},
		{
			lvl:  NoticeLevel,
			msg:  "notice",
			flds: nil,
			want: " NOTICE notice\n",
		},
		{
			lvl:    NoticeLevel,
			msg:    "noticef",
			printf: "%s",
			flds:   nil,
			want:   " NOTICE noticef\n",
		},
		{
			lvl:  WarnLevel,
			msg:  "warn",
			flds: nil,
			want: "   WARN warn\n",
		},
		{
			lvl:    WarnLevel,
			msg:    "warnf",
			printf: "%s",
			flds:   nil,
			want:   "   WARN warnf\n",
		},
		{
			lvl:  ErrorLevel,
			msg:  "error",
			flds: nil,
			want: "  ERROR error\n",
		},
		{
			lvl:    ErrorLevel,
			msg:    "errorf",
			printf: "%s",
			flds:   nil,
			want:   "  ERROR errorf\n",
		},
		{
			lvl:  AlertLevel,
			msg:  "alert",
			flds: nil,
			want: "  ALERT alert\n",
		},
		{
			lvl:    AlertLevel,
			msg:    "alertf",
			printf: "%s",
			flds:   nil,
			want:   "  ALERT alertf\n",
		},
		{
			lvl:  PanicLevel,
			msg:  "panic",
			flds: nil,
			want: "  PANIC panic\n",
		},
		{
			lvl:    PanicLevel,
			msg:    "panicf",
			printf: "%s",
			flds:   nil,
			want:   "  PANIC panicf\n",
		},
		{
			lvl: DebugLevel,
			msg: "debug",
			flds: []Field{
				F("key", "value"),
			},
			want: "  DEBUG debug key=value\n",
		},
		{
			lvl:    DebugLevel,
			msg:    "debugf",
			printf: "%s",
			flds: []Field{
				F("key", "value"),
			},
			want: "  DEBUG debugf key=value\n",
		},
		{
			lvl: InfoLevel,
			msg: "info",
			flds: []Field{
				F("key", "value"),
			},
			want: "   INFO info key=value\n",
		},
		{
			lvl:    InfoLevel,
			msg:    "infof",
			printf: "%s",
			flds: []Field{
				F("key", "value"),
			},
			want: "   INFO infof key=value\n",
		},
		{
			lvl: NoticeLevel,
			msg: "notice",
			flds: []Field{
				F("key", "value"),
			},
			want: " NOTICE notice key=value\n",
		},
		{
			lvl:    NoticeLevel,
			msg:    "noticef",
			printf: "%s",
			flds: []Field{
				F("key", "value"),
			},
			want: " NOTICE noticef key=value\n",
		},
		{
			lvl: WarnLevel,
			msg: "warn",
			flds: []Field{
				F("key", "value"),
			},
			want: "   WARN warn key=value\n",
		},
		{
			lvl:    WarnLevel,
			msg:    "warnf",
			printf: "%s",
			flds: []Field{
				F("key", "value"),
			},
			want: "   WARN warnf key=value\n",
		},
		{
			lvl: ErrorLevel,
			msg: "error",
			flds: []Field{
				F("key", "value"),
			},
			want: "  ERROR error key=value\n",
		},
		{
			lvl:    ErrorLevel,
			msg:    "errorf",
			printf: "%s",
			flds: []Field{
				F("key", "value"),
			},
			want: "  ERROR errorf key=value\n",
		},
		{
			lvl: AlertLevel,
			msg: "alert",
			flds: []Field{
				F("key", "value"),
			},
			want: "  ALERT alert key=value\n",
		},
		{
			lvl: AlertLevel,
			msg: "alert",
			flds: []Field{
				F("key", "value"),
			},
			want: "  ALERT alert key=value\n",
		},
		{
			lvl:    AlertLevel,
			msg:    "alertf",
			printf: "%s",
			flds: []Field{
				F("key", "value"),
			},
			want: "  ALERT alertf key=value\n",
		},
		{
			lvl:    PanicLevel,
			msg:    "panicf",
			printf: "%s",
			flds: []Field{
				F("key", "value"),
			},
			want: "  PANIC panicf key=value\n",
		},
		{
			lvl: PanicLevel,
			msg: "panic",
			flds: []Field{
				F("key", "value"),
			},
			want: "  PANIC panic key=value\n",
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
			want: "  DEBUG debug key=string key=1 key=2 key=3 key=4 key=5 key=1 key=2 key=3 key=4 key=5 key=true key={struct}\n",
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
