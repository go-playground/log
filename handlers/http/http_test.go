package http

import (
	"bytes"
	"github.com/go-playground/log"
	httplogger "github.com/go-playground/log/handlers/http"
	assert "gopkg.in/go-playground/assert.v1"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHttpLogger(t *testing.T) {

	hLog, err := httplogger.New(10000, "http://127.0.0.1:8888/")
	if err != nil {
		log.Fatal("Could create new http logger: ", err)
	}
	log.RegisterHandler(hLog, log.AllLevels...)

	msg := "This is a sample message"

	handler := func(w http.ResponseWriter, r *http.Request) {
		postData, err := ioutil.ReadAll(r.Body)
		assert.Equal(t, err, nil)
		io.WriteString(w, string(postData))
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	r := bytes.NewReader([]byte(msg))
	resp, err := http.Post(server.URL, "application/json", r)
	if err != nil {
		t.Fatalf("POST error: %v", err)
	}

	b, err := ioutil.ReadAll(resp.Body)

	assert.Equal(t, strings.Replace(string(b), "\n", "", -1), msg)

}
