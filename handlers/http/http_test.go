package http

import (
	"fmt"
	"github.com/go-playground/log"
	assert "gopkg.in/go-playground/assert.v1"
	"io/ioutil"
	nethttp "net/http"
	"net/http/httptest"
	"testing"
)

var msg string

func TestHttpLogger(t *testing.T) {

	msg = "This is a sample message sent to the http Go log handler"

	// Start the test HTTP server
	server := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		b, err := ioutil.ReadAll(r.Body)
		postData := string(b)
		// Verify there is no error reading the request body
		assert.Equal(t, err, nil)
		// Verify there the data posted matches exactly the message we expect
		assert.Equal(t, postData, msg)
		fmt.Fprintln(w, postData)
	}))
	defer server.Close()

	// Initiate the http logger
	hLog, err := New(10000, server.URL)
	if err != nil {
		log.Fatal("Could not create new http logger: ", err)
	}
	log.RegisterHandler(hLog, log.AllLevels...)

	log.Info(msg)

}
