package email

import (
	"bytes"
	"html/template"
	stdlog "log"
	"os"
	"strings"
	"sync"
	"time"

	"gopkg.in/gomail.v2"

	"github.com/go-playground/log"
)

// FormatFunc is the function that the workers use to create
// a new Formatter per worker allowing reusable go routine safe
// variable to be used within your Formatter function.
type FormatFunc func(email *Email) Formatter

// Formatter is the function used to format the Email entry
type Formatter func(e *log.Entry) *gomail.Message

const (
	gopath          = "GOPATH"
	contentType     = "text/html"
	defaultTemplate = `<!DOCTYPE html>
<html>
    <body>
        <h2>{{ .Message }}</h2>
        {{ if ne .ApplicationID "" }}
            <h4>{{ .ApplicationID }}</h4>
        {{ end }}
        <p>{{ .Level.String }}</p>
        <p>{{ ts . }}</p>
        {{ if ne .Line 0 }}
            {{ display_file . }}:{{ .Line }}
        {{ end }}
        {{ range $f := .Fields }}
            <p><b>{{ $f.Key }}</b>: {{ $f.Value }}</p>
        {{ end }}
    </body>
</html>`
)

// Email is an instance of the email logger
type Email struct {
	buffer          uint // channel buffer
	numWorkers      uint
	formatFunc      FormatFunc
	timestampFormat string
	gopath          string
	fileDisplay     log.FilenameDisplay
	template        *template.Template
	templateHTML    string
	host            string
	port            int
	username        string
	password        string
	from            string
	to              []string
	keepalive       time.Duration
	m               sync.Mutex
}

// New returns a new instance of the email logger
func New(host string, port int, username string, password string, from string, to []string) *Email {

	return &Email{
		buffer:          3,
		numWorkers:      3,
		timestampFormat: log.DefaultTimeFormat,
		fileDisplay:     log.Lshortfile,
		templateHTML:    defaultTemplate,
		host:            host,
		port:            port,
		username:        username,
		password:        password,
		from:            from,
		to:              to,
		keepalive:       time.Second * 30,
		formatFunc:      defaultFormatFunc,
	}
}

// SetKeepAliveTimout tells Email how long to keep the smtp connection
// open when no messsages are being sent; it will automatically reconnect
// on next message that is received.
func (email *Email) SetKeepAliveTimout(keepAlive time.Duration) {
	email.keepalive = keepAlive
}

// SetEmailTemplate sets Email's html template to be used for email body
func (email *Email) SetEmailTemplate(htmlTemplate string) {
	email.templateHTML = htmlTemplate
}

// SetFilenameDisplay tells Email the filename, when present, how to display
func (email *Email) SetFilenameDisplay(fd log.FilenameDisplay) {
	email.fileDisplay = fd
}

// SetBuffersAndWorkers sets the channels buffer size and number of concurrent workers.
// These settings should be thought about together, hence setting both in the same function.
func (email *Email) SetBuffersAndWorkers(size uint, workers uint) {
	email.buffer = size

	if workers == 0 {
		// just in case no log registered yet
		stdlog.Println("Invalid number of workers specified, setting to 1")
		log.Warn("Invalid number of workers specified, setting to 1")

		workers = 1
	}

	email.numWorkers = workers
}

// From returns the Email's From address
func (email *Email) From() string {
	return email.from
}

// To returns the Email's To address
func (email *Email) To() []string {
	return email.to
}

// Template returns the Email's template
func (email *Email) Template() *template.Template {
	return email.template
}

// SetTimestampFormat sets Email's timestamp output format
// Default is : "2006-01-02T15:04:05.000000000Z07:00"
func (email *Email) SetTimestampFormat(format string) {
	email.timestampFormat = format
}

// SetFormatFunc sets FormatFunc each worker will call to get
// a Formatter func
func (email *Email) SetFormatFunc(fn FormatFunc) {
	email.formatFunc = fn
}

// Run starts the logger consuming on the returned channed
func (email *Email) Run() chan<- *log.Entry {

	// pre-setup
	if email.fileDisplay == log.Llongfile {
		// gather $GOPATH for use in stripping off of full name
		// if not found still ok as will be blank
		email.gopath = os.Getenv(gopath)
		if len(email.gopath) != 0 {
			email.gopath += string(os.PathSeparator) + "src" + string(os.PathSeparator)
		}
	}

	// parse email htmlTemplate, will panic if fails
	email.template = template.Must(template.New("email").Funcs(
		template.FuncMap{
			"display_file": func(e *log.Entry) (file string) {

				file = e.File
				if email.fileDisplay == log.Lshortfile {

					for i := len(file) - 1; i > 0; i-- {
						if file[i] == '/' {
							file = file[i+1:]
							break
						}
					}
				} else {

					// additional check, just in case user does
					// have a $GOPATH but code isn't under it.
					if strings.HasPrefix(file, email.gopath) {
						file = file[len(email.gopath):]
					}
				}

				return
			},
			"ts": func(e *log.Entry) (ts string) {
				ts = e.Timestamp.Format(email.timestampFormat)
				return
			},
		},
	).Parse(email.templateHTML))

	ch := make(chan *log.Entry, email.buffer)

	for i := 0; i <= int(email.numWorkers); i++ {
		go email.handleLog(ch)
	}
	return ch
}

func defaultFormatFunc(email *Email) Formatter {
	var err error
	b := new(bytes.Buffer)

	// apparently there is a race condition when I was using
	// email.to... below in the SetHeader for whatever reason
	// so copying the "to" values solves the issue
	// I wonder if it's a flase positive in the race detector.....
	to := make([]string, len(email.to), len(email.to))
	copy(to, email.to)

	template := email.Template()
	message := gomail.NewMessage()

	message.SetHeader("From", email.from)
	message.SetHeader("To", to...)

	return func(e *log.Entry) *gomail.Message {
		b.Reset()
		if err = template.ExecuteTemplate(b, "email", e); err != nil {
			log.WithFields(log.F("error", err)).Error("Error parsing Email handler template")
		}

		message.SetHeader("Subject", e.Message)
		message.SetBody(contentType, b.String())

		return message
	}
}

func (email *Email) handleLog(entries <-chan *log.Entry) {
	var e *log.Entry
	var s gomail.SendCloser
	var err error
	var open bool
	var alreadyTriedSending bool
	var message *gomail.Message
	var count uint8

	formatter := email.formatFunc(email)

	d := gomail.NewDialer(email.host, email.port, email.username, email.password)

	for {
		select {
		case e = <-entries:
			count = 0
			alreadyTriedSending = false
			message = formatter(e)

		REOPEN:
			// check if smtp connection open
			if !open {
				count++
				if s, err = d.Dial(); err != nil {
					log.WithFields(log.F("error", err)).Warn("ERROR connection to smtp server")

					if count == 3 {
						// we tried to reconnect...
						e.Consumed()
						break
					}

					goto REOPEN
				}
				count = 0
				open = true
			}

		RESEND:
			count++
			if err = gomail.Send(s, message); err != nil {

				log.WithFields(log.F("error", err)).Warn("ERROR sending to smtp server, retrying")

				if count == 3 && !alreadyTriedSending {
					// maybe we got disconnected...
					alreadyTriedSending = true
					open = false
					s.Close()
					goto REOPEN
				} else if alreadyTriedSending {
					// we reopened and tried 2 more times, can;t say we didn't try
					log.WithFields(log.F("error", err)).Alert("ERROR sending log via EMAIL, RETRY and REOPEN failed")
					e.Consumed()
					break
				}

				goto RESEND
			}

			e.Consumed()

		case <-time.After(email.keepalive):
			if open {
				s.Close()
				open = false
			}
		}
	}
}
