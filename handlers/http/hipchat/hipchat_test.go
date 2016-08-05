package hipchat

import (
	"io/ioutil"
	stdhttp "net/http"
	"net/http/httptest"
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

const (
	authToken = "FEFEWG45GRT5FRSSEUIOHJEW"
)

func TestHipChat(t *testing.T) {

	tests := []string{
		"\"color\":\"green\"",
		"\"notify\":true",
		"hipchat_test.go:",
	}

	var msg string

	server := httptest.NewServer(stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			msg = err.Error()
			return
		}

		msg = string(b)

		if strings.Contains(msg, "badrequest") {
			w.WriteHeader(stdhttp.StatusBadRequest)
			return
		}

		w.WriteHeader(stdhttp.StatusAccepted)
	}))
	defer server.Close()

	hc, err := New(APIv2, server.URL+"/", "application/json", authToken)
	if err != nil {
		log.Fatalf("Error initializing hipchat received '%s'", err)
	}

	hc.SetBuffersAndWorkers(0, 0)
	hc.SetTimestampFormat("MST")
	hc.SetEmailTemplate(defaultTemplate)
	hc.SetFilenameDisplay(log.Llongfile)
	log.RegisterHandler(hc, log.DebugLevel)

	log.Debug("debug test")

	for i, tt := range tests {
		if !strings.Contains(msg, tt) {
			t.Errorf("Index '%d' Expected '%s' Got '%s'", i, tt, msg)
		}
	}

	hc.GOPATH()

	hc2, err := New(APIv2, server.URL, "application/json", authToken)
	if err != nil {
		log.Fatalf("Error initializing hipchat received '%s'", err)
	}

	hc2.SetBuffersAndWorkers(1, 1)
	hc2.SetTimestampFormat("MST")
	hc2.SetEmailTemplate(defaultTemplate)
	hc2.SetFilenameDisplay(log.Lshortfile)
	log.RegisterHandler(hc2, log.DebugLevel)

	log.Debug("debug test")

	for i, tt := range tests {
		if !strings.Contains(msg, tt) {
			t.Errorf("Index '%d' Expected '%s' Got '%s'", i, tt, msg)
		}
	}
}

func TestBadTemplate(t *testing.T) {

	tests := []string{
		"\"color\":\"green\"",
		"\"notify\":true",
		"\"message\":\"\"",
	}

	var msg string

	server := httptest.NewServer(stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			msg = err.Error()
			return
		}

		msg = string(b)

		if strings.Contains(msg, "badrequest") {
			w.WriteHeader(stdhttp.StatusBadRequest)
			return
		}

		w.WriteHeader(stdhttp.StatusAccepted)
	}))
	defer server.Close()

	hc, err := New(APIv2, server.URL, "application/json", authToken)
	if err != nil {
		log.Fatalf("Error initializing hipchat received '%s'", err)
	}

	hc.SetBuffersAndWorkers(1, 1)
	hc.SetTimestampFormat("MST")
	hc.SetEmailTemplate("{{.NonExistantField}}")
	hc.SetFilenameDisplay(log.Llongfile)
	log.RegisterHandler(hc, log.DebugLevel)

	log.Debug("debug test")

	for i, tt := range tests {
		if !strings.Contains(msg, tt) {
			t.Errorf("Index '%d' Expected '%s' Got '%s'", i, tt, msg)
		}
	}
}

func TestBadURL(t *testing.T) {
	expected := "parse @#$%?auth_test=true: invalid URL escape \"%?a\""
	_, err := New(APIv2, "@#$%", "application/json", authToken)
	if err == nil || err.Error() != expected {
		log.Fatalf("Expected '%s' Got '%s'", err, expected)
	}

	expected = "Get ?auth_test=true: unsupported protocol scheme \"\""
	_, err = New(APIv2, "", "application/json", authToken)
	if err == nil || err.Error() != expected {
		log.Fatalf("Expected '%s' Got '%s'", err, expected)
	}
}

func TestNotAccepted(t *testing.T) {

	server := httptest.NewServer(stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		w.WriteHeader(stdhttp.StatusUnauthorized)
	}))
	defer server.Close()

	expected := "HipChat authorization failed\n "
	_, err := New(APIv2, server.URL, "application/json", authToken)
	if err == nil || err.Error() != expected {
		log.Fatalf("Error Expected '%s' Got '%s'", expected, err.Error())
	}
}
