package http

import (
	"bytes"
	"fmt"
	"io/ioutil"
	stdlog "log"
	stdhttp "net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/go-playground/log"
)

// FormatFunc is the function that the workers use to create
// a new Formatter per worker allowing reusable go routine safe
// variable to be used within your Formatter function.
type FormatFunc func() Formatter

// Formatter is the function used to format the HTTP entry
type Formatter func(e *log.Entry) []byte

const (
	defaultTS = "2006-01-02T15:04:05.000000000Z07:00"
	space     = byte(' ')
	equals    = byte('=')
	colon     = byte(':')
	base10    = 10
	v         = "%v"
	gopath    = "GOPATH"
)

// HTTP is an instance of the http logger
type HTTP struct {
	buffer          uint // channel buffer
	numWorkers      uint
	remoteHost      string
	formatFunc      FormatFunc
	timestampFormat string
	httpClient      stdhttp.Client
	header          stdhttp.Header
	method          string
	gopath          string
	fileDisplay     log.FilenameDisplay
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
		fileDisplay:     log.Lshortfile,
	}

	h.formatFunc = h.defaultFormatFunc

	return h, nil
}

// SetFilenameDisplay tells HTTP the filename, when present, how to display
func (h *HTTP) SetFilenameDisplay(fd log.FilenameDisplay) {
	h.fileDisplay = fd
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

	// pre-setup
	if h.fileDisplay == log.Llongfile {
		// gather $GOPATH for use in stripping off of full name
		// if not found still ok as will be blank
		h.gopath = os.Getenv(gopath)
		if len(h.gopath) != 0 {
			h.gopath += string(os.PathSeparator) + "src" + string(os.PathSeparator)
		}
	}

	ch := make(chan *log.Entry, h.buffer)

	for i := 0; i <= int(h.numWorkers); i++ {
		go h.handleLog(ch)
	}
	return ch
}

func (h *HTTP) defaultFormatFunc() Formatter {

	var b []byte
	var file string
	var lvl string
	var i int

	return func(e *log.Entry) []byte {
		b = b[0:0]

		if e.Line == 0 {

			b = append(b, e.Timestamp.Format(h.timestampFormat)...)
			b = append(b, space)

			lvl = e.Level.String()

			for i = 0; i < 6-len(lvl); i++ {
				b = append(b, space)
			}

			b = append(b, lvl...)
			b = append(b, space)
			b = append(b, e.Message...)

		} else {
			file = e.File

			if h.fileDisplay == log.Lshortfile {

				for i = len(file) - 1; i > 0; i-- {
					if file[i] == '/' {
						file = file[i+1:]
						break
					}
				}
			} else {
				file = file[len(h.gopath):]
			}

			b = append(b, e.Timestamp.Format(h.timestampFormat)...)
			b = append(b, space)

			lvl = e.Level.String()

			for i = 0; i < 6-len(lvl); i++ {
				b = append(b, space)
			}

			b = append(b, lvl...)
			b = append(b, space)
			b = append(b, file...)
			b = append(b, colon)
			b = strconv.AppendInt(b, int64(e.Line), base10)
			b = append(b, space)
			b = append(b, e.Message...)
		}

		for _, f := range e.Fields {
			b = append(b, space)
			b = append(b, f.Key...)
			b = append(b, equals)

			switch f.Value.(type) {
			case string:
				b = append(b, f.Value.(string)...)
			case int:
				b = strconv.AppendInt(b, int64(f.Value.(int)), base10)
			case int8:
				b = strconv.AppendInt(b, int64(f.Value.(int8)), base10)
			case int16:
				b = strconv.AppendInt(b, int64(f.Value.(int16)), base10)
			case int32:
				b = strconv.AppendInt(b, int64(f.Value.(int32)), base10)
			case int64:
				b = strconv.AppendInt(b, f.Value.(int64), base10)
			case uint:
				b = strconv.AppendUint(b, uint64(f.Value.(uint)), base10)
			case uint8:
				b = strconv.AppendUint(b, uint64(f.Value.(uint8)), base10)
			case uint16:
				b = strconv.AppendUint(b, uint64(f.Value.(uint16)), base10)
			case uint32:
				b = strconv.AppendUint(b, uint64(f.Value.(uint32)), base10)
			case uint64:
				b = strconv.AppendUint(b, f.Value.(uint64), base10)
			case bool:
				b = strconv.AppendBool(b, f.Value.(bool))
			default:
				b = append(b, fmt.Sprintf(v, f.Value)...)
			}
		}

		return b
	}
}

func (h *HTTP) handleLog(entries <-chan *log.Entry) {
	var e *log.Entry
	var b []byte
	var reader *bytes.Reader

	formatter := h.formatFunc()

	req, _ := stdhttp.NewRequest(h.method, h.remoteHost, nil)
	req.Header = h.header

	for e = range entries {

		b = formatter(e)

		reader = bytes.NewReader(b)
		req.Body = ioutil.NopCloser(reader)
		req.ContentLength = int64(reader.Len())

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
