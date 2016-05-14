package http

import (
	"bytes"
	"fmt"
	"io"
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
type Formatter func(e *log.Entry) io.Reader

const (
	defaultTS       = "2006-01-02T15:04:05.000000000Z07:00"
	format          = "%s %6s %s"
	formatCaller    = "%s %6s %s:%d %s"
	noColorKeyValue = " %s=%v"
)

// HTTP is an instance of the http logger
type HTTP struct {
	buffer          uint // channel buffer
	numWorkers      uint
	remoteHost      string
	formatFunc      FormatFunc
	contentEncoding string
	timestampFormat string
	httpClient      stdhttp.Client
	header          stdhttp.Header
	method          string
}

// New returns a new instance of the http logger
func New(remoteHost string, method string, header stdhttp.Header) (*HTTP, error) {

	if _, err := url.Parse(remoteHost); err != nil {
		return nil, err
	}

	h := &HTTP{
		buffer:          0,
		remoteHost:      remoteHost,
		numWorkers:      1,
		timestampFormat: defaultTS,
		httpClient:      stdhttp.Client{},
		header:          header,
		method:          method,
	}

	h.formatFunc = h.defaultFormatFunc

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

// SetTimestampFormat sets HTTP's timestamp output format
// Default is : "2006-01-02T15:04:05.000000000Z07:00"
func (h *HTTP) SetTimestampFormat(format string) {
	h.timestampFormat = format
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

func (h *HTTP) defaultFormatFunc() Formatter {

	b := new(bytes.Buffer)
	var file string

	return func(e *log.Entry) io.Reader {
		b.Reset()

		if e.Line == 0 {

			b.WriteString(fmt.Sprintf(format, e.Timestamp.Format(h.timestampFormat), e.Level, e.Message))

		} else {
			file = e.File
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					file = file[i+1:]
					break
				}
			}

			b.WriteString(fmt.Sprintf(formatCaller, e.Timestamp.Format(h.timestampFormat), e.Level, file, e.Line, e.Message))
		}

		for _, f := range e.Fields {
			b.WriteString(fmt.Sprintf(noColorKeyValue, f.Key, f.Value))
		}

		return b
	}
}

func (h *HTTP) handleLog(entries <-chan *log.Entry) {
	var e *log.Entry
	var reader io.Reader
	formatter := h.formatFunc()

	for e = range entries {

		reader = formatter(e)

		// TODO: investigate reuse of http.Request... all that changes is the paylod

		// err not gathered as URL parsed during creationg of *HTTP
		req, _ := stdhttp.NewRequest(h.method, h.remoteHost, reader)

		req.Header = h.header
		resp, err := h.httpClient.Do(req)
		if err != nil {
			fmt.Printf("**** WARNING Could not post data to %s: %v\n", h.remoteHost, err)
			goto END
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 299 {
			fmt.Printf("WARNING Received HTTP %d during POST request to %s\n", resp.StatusCode, h.remoteHost)
		}

	END:
		e.Consumed()
	}
}
