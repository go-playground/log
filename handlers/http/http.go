package http

import (
	"bytes"
	"fmt"
	"io/ioutil"
	stdhttp "net/http"
	"net/url"
	"strconv"
	"sync"

	"github.com/go-playground/log"
)

// FormatFunc is the function that the workers use to create
// a new Formatter per worker allowing reusable go routine safe
// variable to be used within your Formatter function.
type FormatFunc func(h *HTTP) Formatter

// Formatter is the function used to format the HTTP entry
type Formatter func(e log.Entry) []byte

const (
	space  = byte(' ')
	equals = byte('=')
	base10 = 10
	v      = "%v"
)

// HTTP is an instance of the http logger
type HTTP struct {
	remoteHost      string
	formatter       Formatter
	formatFunc      FormatFunc
	timestampFormat string
	header          stdhttp.Header
	method          string
	client          stdhttp.Client
	once            sync.Once
}

// New returns a new instance of the http logger
func New(remoteHost string, method string, header stdhttp.Header) (*HTTP, error) {
	if _, err := url.Parse(remoteHost); err != nil {
		return nil, err
	}

	h := &HTTP{
		remoteHost:      remoteHost,
		timestampFormat: log.DefaultTimeFormat,
		header:          header,
		method:          method,
		client:          stdhttp.Client{},
		formatFunc:      defaultFormatFunc,
	}
	return h, nil
}

// SetTimestampFormat sets HTTP's timestamp output format
// Default is : "2006-01-02T15:04:05.000000000Z07:00"
func (h *HTTP) SetTimestampFormat(format string) {
	h.timestampFormat = format
}

// TimestampFormat returns HTTP's current timestamp output format
func (h *HTTP) TimestampFormat() string {
	return h.timestampFormat
}

// SetFormatFunc sets FormatFunc each worker will call to get
// a Formatter func
func (h *HTTP) SetFormatFunc(fn FormatFunc) {
	h.formatFunc = fn
}

// Method returns http method registered with HTTP
func (h *HTTP) Method() string {
	return h.method
}

// RemoteHost returns the remote host registered to HTTP
func (h *HTTP) RemoteHost() string {
	return h.remoteHost
}

// Headers returns the http headers registered to HTTP
func (h *HTTP) Headers() stdhttp.Header {
	return h.header
}

func defaultFormatFunc(h *HTTP) Formatter {
	var b []byte
	var lvl string
	var i int
	tsFormat := h.TimestampFormat()

	return func(e log.Entry) []byte {
		b = b[0:0]

		b = append(b, e.Timestamp.Format(tsFormat)...)
		b = append(b, space)

		lvl = e.Level.String()

		for i = 0; i < 6-len(lvl); i++ {
			b = append(b, space)
		}

		b = append(b, lvl...)
		b = append(b, space)
		b = append(b, e.Message...)

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

// Log handles the log entry
func (h *HTTP) Log(e log.Entry) {
	h.once.Do(func() {
		h.formatter = h.formatFunc(h)
	})
	remoteHost := h.RemoteHost()
	req, _ := stdhttp.NewRequest(h.Method(), remoteHost, nil)
	req.Header = h.Headers()

	b := h.formatter(e)

	reader := bytes.NewReader(b)
	req.Body = ioutil.NopCloser(reader)
	req.ContentLength = int64(reader.Len())

	resp, err := h.client.Do(req)
	if err != nil {
		log.Errorf("Could not post data to %s: %v\n", remoteHost, err)
		return
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 299 {
		bt, _ := ioutil.ReadAll(resp.Body)
		log.Errorf("Received HTTP %d during POST request to %s body: %s\n", resp.StatusCode, remoteHost, string(bt))
	}
}
