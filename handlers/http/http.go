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
	"strings"

	"github.com/go-playground/log"
)

// FormatFunc is the function that the workers use to create
// a new Formatter per worker allowing reusable go routine safe
// variable to be used within your Formatter function.
type FormatFunc func(h HTTP) Formatter

// Formatter is the function used to format the HTTP entry
type Formatter func(e *log.Entry) []byte

const (
	space  = byte(' ')
	equals = byte('=')
	colon  = byte(':')
	base10 = 10
	v      = "%v"
	gopath = "GOPATH"
)

// HTTP interface to allow for defining handlers based upon this one.
type HTTP interface {
	SetFilenameDisplay(fd log.FilenameDisplay)
	FilenameDisplay() log.FilenameDisplay
	SetBuffersAndWorkers(size uint, workers uint)
	SetTimestampFormat(format string)
	TimestampFormat() string
	GOPATH() string
	SetFormatFunc(fn FormatFunc)
	Run() chan<- *log.Entry
	FormatFunc() FormatFunc
	Method() string
	RemoteHost() string
	Headers() stdhttp.Header
	Buffers() uint
	Workers() uint
}

// HTTP is an instance of the http logger
type internalHTTP struct {
	buffer          uint // channel buffer
	numWorkers      uint
	remoteHost      string
	formatFunc      FormatFunc
	timestampFormat string
	header          stdhttp.Header
	method          string
	gopath          string
	fileDisplay     log.FilenameDisplay
}

var _ HTTP = new(internalHTTP)

// New returns a new instance of the http logger
func New(remoteHost string, method string, header stdhttp.Header) (HTTP, error) {

	if _, err := url.Parse(remoteHost); err != nil {
		return nil, err
	}

	return &internalHTTP{
		buffer:          3,
		numWorkers:      3,
		remoteHost:      remoteHost,
		timestampFormat: log.DefaultTimeFormat,
		header:          header,
		method:          method,
		fileDisplay:     log.Lshortfile,
		formatFunc:      defaultFormatFunc,
	}, nil
}

// SetFilenameDisplay tells HTTP the filename, when present, how to display
func (h *internalHTTP) SetFilenameDisplay(fd log.FilenameDisplay) {
	h.fileDisplay = fd
}

// FilenameDisplay returns Console's current filename display setting
func (h *internalHTTP) FilenameDisplay() log.FilenameDisplay {
	return h.fileDisplay
}

// SetBuffersAndWorkers sets the channels buffer size and number of concurrent workers.
// These settings should be thought about together, hence setting both in the same function.
func (h *internalHTTP) SetBuffersAndWorkers(size uint, workers uint) {
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
func (h *internalHTTP) SetTimestampFormat(format string) {
	h.timestampFormat = format
}

// TimestampFormat returns HTTP's current timestamp output format
func (h *internalHTTP) TimestampFormat() string {
	return h.timestampFormat
}

// GOPATH returns the GOPATH calculated by HTTP
func (h *internalHTTP) GOPATH() string {
	return h.gopath
}

// SetFormatFunc sets FormatFunc each worker will call to get
// a Formatter func
func (h *internalHTTP) SetFormatFunc(fn FormatFunc) {
	h.formatFunc = fn
}

// FormatFunc returns FormatFunc registered with HTTP
func (h *internalHTTP) FormatFunc() FormatFunc {
	return h.formatFunc
}

// Method returns http method registered with HTTP
func (h *internalHTTP) Method() string {
	return h.method
}

// RemoteHost returns the remote host registered to HTTP
func (h *internalHTTP) RemoteHost() string {
	return h.remoteHost
}

// Headers returns the http headers registered to HTTP
func (h *internalHTTP) Headers() stdhttp.Header {
	return h.header
}

// Buffers returns the http buffer count registered to HTTP
func (h *internalHTTP) Buffers() uint {
	return h.buffer
}

// Workers returns the http worker count registered to HTTP
func (h *internalHTTP) Workers() uint {
	return h.numWorkers
}

// Run starts the logger consuming on the returned channed
func (h *internalHTTP) Run() chan<- *log.Entry {

	// pre-setup
	if h.fileDisplay == log.Llongfile {
		// gather $GOPATH for use in stripping off of full name
		// if not found still ok as will be blank
		h.gopath = os.Getenv(gopath)
		if len(h.gopath) != 0 {
			h.gopath += string(os.PathSeparator) + "src" + string(os.PathSeparator)
		}
	}

	ch := make(chan *log.Entry, h.Buffers())

	for i := 0; i <= int(h.Workers()); i++ {
		go HandleLog(h, ch)
	}
	return ch
}

func defaultFormatFunc(h HTTP) Formatter {

	var b []byte
	var file string
	var lvl string
	var i int
	gopath := h.GOPATH()
	tsFormat := h.TimestampFormat()
	fnameDisplay := h.FilenameDisplay()

	return func(e *log.Entry) []byte {
		b = b[0:0]

		if e.Line == 0 {

			b = append(b, e.Timestamp.Format(tsFormat)...)
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

			if fnameDisplay == log.Lshortfile {

				for i = len(file) - 1; i > 0; i-- {
					if file[i] == '/' {
						file = file[i+1:]
						break
					}
				}
			} else {

				// additional check, just in case user does
				// have a $GOPATH but code isn't under it.
				if strings.HasPrefix(file, gopath) {
					file = file[len(gopath):]
				}
			}

			b = append(b, e.Timestamp.Format(tsFormat)...)
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
			case float32:
				b = strconv.AppendFloat(b, float64(f.Value.(float32)), 'f', -1, 32)
			case float64:
				b = strconv.AppendFloat(b, f.Value.(float64), 'f', -1, 64)
			case bool:
				b = strconv.AppendBool(b, f.Value.(bool))
			default:
				b = append(b, fmt.Sprintf(v, f.Value)...)
			}
		}

		return b
	}
}

// HandleLog is the default http log handler
func HandleLog(h HTTP, entries <-chan *log.Entry) {
	var e *log.Entry
	var b []byte
	var reader *bytes.Reader

	formatter := h.FormatFunc()(h)
	remoteHost := h.RemoteHost()
	httpClient := stdhttp.Client{}

	req, _ := stdhttp.NewRequest(h.Method(), remoteHost, nil)
	req.Header = h.Headers()

	for e = range entries {

		b = formatter(e)

		reader = bytes.NewReader(b)
		req.Body = ioutil.NopCloser(reader)
		req.ContentLength = int64(reader.Len())

		resp, err := httpClient.Do(req)
		if err != nil {
			log.Error("Could not post data to %s: %v\n", remoteHost, err)
			goto END
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 299 {
			bt, _ := ioutil.ReadAll(resp.Body)
			log.Error("Received HTTP %d during POST request to %s body: %s\n", resp.StatusCode, remoteHost, string(bt))
		}

	END:
		e.Consumed()
	}
}
