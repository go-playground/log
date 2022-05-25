package http

import (
	"io/ioutil"
	stdhttp "net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/log/v8"
)

// NOTES:
// - Run "go test" to run tests
// - Run "gocov test | gocov report" to report on test converage by file
// - Run "gocov test | gocov annotate -" to report on all code and functions, those ,marked with "MISS" were never called

// or

// -- may be a good idea to change to output path to somewherelike /tmp
// go test -coverprofile cover.out && go tool cover -html=cover.out -o cover.html

func TestHTTPLogger(t *testing.T) {
	log.SetExitFunc(func(int) {})
	tests := getTestHTTPLoggerTests()
	var msg string

	server := httptest.NewServer(stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			msg = err.Error()
			return
		}

		msg = string(b)

		if msg == "  DEBUG badrequest" {
			w.WriteHeader(stdhttp.StatusBadRequest)
			return
		}
	}))
	defer server.Close()

	header := make(stdhttp.Header)
	header.Set("Content-Type", "text/plain")

	hLog, err := New(server.URL, "POST", header)
	if err != nil {
		log.Fatalf("Error initializing HTTP received '%s'", err)
	}
	defer func() { _ = hLog.Close() }()
	hLog.SetTimestampFormat("")
	log.AddHandler(hLog, log.DebugLevel, log.InfoLevel, log.NoticeLevel, log.WarnLevel, log.PanicLevel, log.AlertLevel, log.FatalLevel)
	for i, tt := range tests {

		var l log.Entry

		if tt.flds != nil {
			l = l.WithFields(tt.flds...)
		}

		switch tt.lvl {
		case log.DebugLevel:
			l.Debug(tt.msg)
		case log.InfoLevel:
			l.Info(tt.msg)
		case log.NoticeLevel:
			l.Notice(tt.msg)
		case log.WarnLevel:
			l.Warn(tt.msg)
		case log.ErrorLevel:
			l.Error(tt.msg)
		case log.PanicLevel:
			l.Panic(tt.msg)
		case log.AlertLevel:
			l.Alert(tt.msg)
		}

		if msg != tt.want {
			t.Errorf("test %d: Expected '%s' Got '%s'", i, tt.want, msg)
		}
	}

	log.Debug("badrequest")
}

func TestBadValues(t *testing.T) {

	pErr := `parse "@#$%": invalid URL escape "%"`
	header := make(stdhttp.Header)
	header.Set("Content-Type", "text/plain")

	_, err := New("@#$%", "POST", header)
	if err == nil {
		t.Fatalf("Expected '<nil>' Got '%s'", err)
	}

	if err.Error() != pErr {
		t.Fatalf("Expected '%s' Got '%s'", pErr, err)
	}

	hLog, err := New("http://127.0.0.1:4354", "POST", header)
	if err != nil {
		t.Fatalf("Expected '<nil>' Got '%s'", err)
	}

	hLog.SetFormatFunc(func(h *HTTP) Formatter {
		return func(e log.Entry) []byte {
			return []byte(e.Message)
		}
	})
	log.AddHandler(hLog, log.DebugLevel, log.InfoLevel, log.NoticeLevel, log.WarnLevel, log.PanicLevel, log.AlertLevel, log.FatalLevel)

	log.Debug("debug")
}

type test struct {
	lvl  log.Level
	msg  string
	flds []log.Field
	want string
}

func getTestHTTPLoggerTests() []test {
	return []test{
		{
			lvl:  log.DebugLevel,
			msg:  "debug",
			flds: nil,
			want: "  DEBUG debug",
		},
		{
			lvl:  log.PanicLevel,
			msg:  "panic",
			flds: nil,
			want: "  PANIC panic",
		},
		{
			lvl:  log.InfoLevel,
			msg:  "info",
			flds: nil,
			want: "   INFO info",
		},
		{
			lvl:  log.NoticeLevel,
			msg:  "notice",
			flds: nil,
			want: " NOTICE notice",
		},
		{
			lvl:  log.WarnLevel,
			msg:  "warn",
			flds: nil,
			want: "   WARN warn",
		},
		// {
		// 	lvl:  log.ErrorLevel,
		// 	msg:  "error",
		// 	flds: nil,
		// 	want: "  ERROR http_test.go:78 error",
		// },
		{
			lvl:  log.AlertLevel,
			msg:  "alert",
			flds: nil,
			want: "  ALERT alert",
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
			want: "  DEBUG debug key=string key=1 key=2 key=3 key=4 key=5 key=1 key=2 key=3 key=4 key=5 key=5.33 key=5.34 key=true key={struct}",
		},
	}
}
