package email

import (
	"bytes"
	"html/template"
	"sync"

	"gopkg.in/gomail.v2"

	"github.com/go-playground/log"
)

// FormatFunc is the function that the workers use to create
// a new Formatter per worker allowing reusable go routine safe
// variable to be used within your Formatter function.
type FormatFunc func(email *Email) Formatter

// Formatter is the function used to format the Email entry
type Formatter func(e log.Entry) *gomail.Message

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
	formatter       Formatter
	timestampFormat string
	gopath          string
	template        *template.Template
	host            string
	port            int
	username        string
	password        string
	from            string
	to              []string
	m               sync.Mutex
}

// New returns a new instance of the email logger
func New(host string, port int, username string, password string, from string, to []string) *Email {
	e := &Email{
		timestampFormat: log.DefaultTimeFormat,
		host:            host,
		port:            port,
		username:        username,
		password:        password,
		from:            from,
		to:              to,
	}
	e.SetEmailTemplate(defaultTemplate)
	e.formatter = defaultFormatFunc(e)
	return e
}

// SetTemplate sets Email's html template to be used for email body
func (email *Email) SetTemplate(htmlTemplate string) {
	// parse email htmlTemplate, will panic if fails
	email.template = template.Must(template.New("email").Funcs(
		template.FuncMap{
			"ts": func(e *log.Entry) (ts string) {
				ts = e.Timestamp.Format(email.timestampFormat)
				return
			},
		},
	).Parse(htmlTemplate))
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
	email.formatter = fn(email)
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

	return func(e log.Entry) *gomail.Message {
		b.Reset()
		if err = template.ExecuteTemplate(b, "email", e); err != nil {
			log.WithField("error", err).Error("Error parsing Email handler template")
		}

		message.SetHeader("Subject", e.Message)
		message.SetBody(contentType, b.String())
		return message
	}
}

// Log handles the log entry
func (email *Email) Log(e log.Entry) {
	var s gomail.SendCloser
	var err error
	var open bool
	var alreadyTriedSending bool
	var message *gomail.Message
	var count uint8

	d := gomail.NewDialer(email.host, email.port, email.username, email.password)

	for {
		count = 0
		alreadyTriedSending = false
		message = email.formatter(e)

	REOPEN:
		// check if smtp connection open
		if !open {
			count++
			if s, err = d.Dial(); err != nil {
				log.WithField("error", err).Warn("ERROR connection to smtp server")

				if count == 3 {
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

			log.WithField("error", err).Warn("ERROR sending to smtp server, retrying")

			if count == 3 && !alreadyTriedSending {
				// maybe we got disconnected...
				alreadyTriedSending = true
				open = false
				s.Close()
				goto REOPEN
			} else if alreadyTriedSending {
				// we reopened and tried 2 more times, can't say we didn't try
				log.WithField("error", err).Alert("ERROR sending log via EMAIL, RETRY and REOPEN failed")
				break
			}
			goto RESEND
		}
	}
}
