package http

import (
	"bytes"
	"fmt"
	stdlog "log"
	stdhttp "net/http"
	"net/url"

	"github.com/go-playground/log"
)

// FormatFunc is the function that the workers use to create
// a new Formatter per worker allowing reusable go routine safe
// variable to be used within your Formatter function.
type FormatFunc func() Formatter

// Formatter is the function used to format the HTTP entry
type Formatter func(e *log.Entry) []byte

// HTTP is an instance of the http logger
type HTTP struct {
	buffer          uint // channel buffer
	remoteHost      string
	formatFunc      FormatFunc
	contentEncoding string
	httpClient      stdhttp.Client
	numWorkers      uint
}

// New returns a new instance of the http logger
func New(remoteHost string, contentEncoding string) (*HTTP, error) {

	if _, err := url.Parse(remoteHost); err != nil {
		return nil, err
	}

	h := &HTTP{
		buffer:          0,
		remoteHost:      remoteHost,
		contentEncoding: contentEncoding,
		formatFunc:      formatFunc,
		numWorkers:      1,
		httpClient:      stdhttp.Client{},
	}

	return h, nil
}

// SetBuffersAndWorkers sets the channels buffer size and number of concurrent workers.
// These settings should be thought about together, hence setting both in the same function.
func (h *HTTP) SetBuffersAndWorkers(size uint, workers uint) {
	h.buffer = size

	if workers == 0 {
		// just in case no log registered yet
		stdlog.Println("Invalid number of workers specified, setting to 1")
		log.Warn("Invalid number of workers specified, setting to 1")

		workers = 1
	}

	h.numWorkers = workers
}

// SetFormatFunc sets FormatFunc each worker will call to get
// a Formatter func
func (h *HTTP) SetFormatFunc(fn FormatFunc) {
	h.formatFunc = fn
}

// Run starts the logger consuming on the returned channed
func (h *HTTP) Run() chan<- *log.Entry {

	ch := make(chan *log.Entry, h.buffer)

	for i := 0; i <= int(h.numWorkers); i++ {
		go h.handleLog(ch)
	}
	return ch
}

func formatFunc() Formatter {

	var b []byte

	return func(e *log.Entry) []byte {
		b = b[0:0]

		b = append(b, "TEST"...)
		return b
	}
}

func (h *HTTP) handleLog(entries <-chan *log.Entry) {
	var e *log.Entry
	var payload []byte
	formatter := h.formatFunc()

	for e = range entries {

		payload = formatter(e)
		b := bytes.NewBuffer(payload)

		// TODO: investigate reuse of http.Request... al that changes is the paylod
		// maybe b.Reset()

		// Issue POST request to send off data
		req, err := stdhttp.NewRequest("POST", h.remoteHost, b)
		if err != nil {
			log.Info(fmt.Sprintf("[Error] Could not initialize new request: %v\n", err))
		}

		req.Header.Add("Content-Type", h.contentEncoding)
		resp, err := h.httpClient.Do(req)
		if err != nil {
			log.Info(fmt.Sprintf("[Error] Could not post data to %s: %v\n", h.remoteHost, err))
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 299 {
			log.Info(fmt.Sprintf("[Error] Received HTTP %d during POST request to %s\n", resp.StatusCode, h.remoteHost))
		}

		e.Consumed()
	}
}
