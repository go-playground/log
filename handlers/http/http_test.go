package http

import (
	"fmt"
	"github.com/go-playground/log"
	httplogger "github.com/go-playground/log/handlers/http"
	assert "gopkg.in/go-playground/assert.v1"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var msg string

func TestHttpLogger(t *testing.T) {

	msg = "This is a sample message"

	// Start the test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	hLog, err := httplogger.New(10000, server.URL)
	if err != nil {
		log.Fatal("Could not create new http logger: ", err)
	}
	log.RegisterHandler(hLog, log.AllLevels...)

	log.Info(msg)

}
