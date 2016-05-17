package http

import (
	"fmt"
	"io/ioutil"
	stdhttp "net/http"
	"net/http/httptest"
	"path"
	"runtime"
	"strings"
	"testing"

	"github.com/go-playground/log"
)

// NOTES:
// - Run "go test" to run tests
// - Run "gocov test | gocov report" to report on test converage by file
// - Run "gocov test | gocov annotate -" to report on all code and functions, those ,marked with "MISS" were never called

// or

// -- may be a good idea to change to output path to somewherelike /tmp
// go test -coverprofile cover.out && go tool cover -html=cover.out -o cover.html

func TestHTTPLogger(t *testing.T) {
	var msg string

	server := httptest.NewServer(stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			msg = err.Error()
			return
		}

		msg = string(b)

		if msg == "UTC  DEBUG badrequest" {
			w.WriteHeader(stdhttp.StatusBadRequest)
			return
		}
	}))
	defer server.Close()

	header := make(stdhttp.Header, 0)
	header.Set("Content-Type", "text/plain")

	hLog, err := New(server.URL, "POST", header)
	if err != nil {
		log.Fatalf("Error initializing HTTP recieved '%s'", err)
	}

	hLog.SetBuffersAndWorkers(0, 0)
	hLog.SetTimestampFormat("MST")
	log.RegisterHandler(hLog, log.AllLevels...)

	log.Debug("debug")
	if msg != "UTC  DEBUG debug" {
		t.Errorf("Expected 'UTC  DEBUG' Got '%s'", msg)
	}

	log.Info("info")
	if msg != "UTC   INFO info" {
		t.Errorf("Expected 'UTC   INFO info' Got '%s'", msg)
	}

	log.Notice("notice")
	if msg != "UTC NOTICE notice" {
		t.Errorf("Expected 'UTC NOTICE notice' Got '%s'", msg)
	}

	log.Warn("warn")
	if msg != "UTC   WARN http_test.go:72 warn" {
		t.Errorf("Expected 'UTC   WARN http_test.go:72 warn' Got '%s'", msg)
	}

	log.Error("error")
	if msg != "UTC  ERROR http_test.go:77 error" {
		t.Errorf("Expected 'UTC  ERROR http_test.go:77 error' Got '%s'", msg)
	}

	log.Alert("alert")
	if msg != "UTC  ALERT http_test.go:82 alert" {
		t.Errorf("Expected 'UTC  ALERT http_test.go:82 alert' Got '%s'", msg)
	}

	log.WithFields(
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
	).Debug("debug")
	if msg != "UTC  DEBUG debug key=string key=1 key=2 key=3 key=4 key=5 key=1 key=2 key=3 key=4 key=5 key=true key={struct}" {
		t.Errorf("UTC  DEBUG debug key=string key=1 key=2 key=3 key=4 key=5 key=1 key=2 key=3 key=4 key=5 key=true key={struct}' Got '%s'", msg)
	}

	panicMatchesSkip(t, func() { log.Panic("panic") }, "panic")
	if msg != "UTC  PANIC http_test.go:106 panic" {
		t.Errorf("Expected 'UTC  PANIC http_test.go:106 panic' Got '%s'", msg)
	}

	func() {
		defer log.Trace("trace").End()
	}()

	if !strings.HasPrefix(msg, "UTC  TRACE trace ") {
		t.Errorf("Expected 'UTC  TRACE trace ...' Got '%s'", msg)
	}

	log.Debug("badrequest")
}

func TestBadValues(t *testing.T) {

	pErr := "parse @#$%: invalid URL escape \"%\""
	header := make(stdhttp.Header, 0)
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

	hLog.SetFormatFunc(func() Formatter {
		return func(e *log.Entry) []byte {
			return []byte(e.Message)
		}
	})
	log.RegisterHandler(hLog, log.AllLevels...)

	log.Debug("debug")
}

func panicMatchesSkip(t *testing.T, fn func(), matches string) {

	_, file, line, _ := runtime.Caller(2)

	defer func() {
		if r := recover(); r != nil {
			err := fmt.Sprintf("%s", r)

			if err != matches {
				fmt.Printf("%s:%d Panic...  Expected '%s' Got '%s'", path.Base(file), line, matches, err)
				t.FailNow()
			}
		} else {
			fmt.Printf("%s:%d Panic Expected, none found...  Got '%s'", path.Base(file), line, matches)
			t.FailNow()
		}
	}()

	fn()
}
