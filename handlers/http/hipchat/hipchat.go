package hipchat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	stdhttp "net/http"
	"os"
	"strings"

	"github.com/go-playground/log"
	"github.com/go-playground/log/handlers/http"
)

// APIVersion specifies the HipChat API version to use
type APIVersion uint8

// Supported API Versions
const (
	APIv2 APIVersion = iota
)

// Value is a HipChat value object
type Value struct {
	URL   string `json:"url,omitempty"`   // min 1 max unlimited, Url to be opened when a user clicks on the label
	Style string `json:"style,omitempty"` // AUI Integrations for now supporting only lozenges, lozenge-success, lozenge-error, lozenge-current, lozenge-complete, lozenge-moved, lozenge.
	Label string `json:"label"`           // min 1 max unlimited
	Icon  *Icon  `json:"icon,omitempty"`  // icon to display
}

// Attribute is a HipChat attribute object
type Attribute struct {
	Value Value  `json:"value"`           // attribute value
	Label string `json:"label,omitempty"` // min 1 max 50
}

// Icon is a HipChat icon object
type Icon struct {
	URL       string `json:"url"`    // min 1 max unlimited
	URLRetina string `json:"url@2x"` // min 1 max unlimited, the icon url in retina
}

// Activity is a HipChat activity object
type Activity struct {
	HTML string `json:"html"`           // 1 - unlimited, Html for the activity to show in one line a summary of the action that happened
	Icon *Icon  `json:"icon,omitempty"` // icon to display
}

// Thumbnail is a HipChat thumbnail object
type Thumbnail struct {
	URL       string `json:"url"`               // min 1 max 250
	URLRetina string `json:"url@2x"`            // min 1 max 250, the thumbnail url in retina
	Width     uint   `json:"width,omitempty"`   // width of image
	Height    uint   `json:"height, omitempty"` // height of image
}

// Description is a HipChat description object
type Description struct {
	Value  string `json:"value"`  // min 1 max 1000
	Format string `json:"format"` // html or text
}

// Card is a custom Hipchat
type Card struct {
	Style            string       `json:"style"`                 // min 1 max 16
	Description      *Description `json:"description,omitempty"` // description object
	Format           string       `json:"format,omitempty"`      // compact or medium
	URL              string       `json:"url,omitempty"`         // 1 - unlimited
	Title            string       `json:"title"`                 // min 1 max 500
	HipChatThumbnail *Thumbnail   `json:"thumbnail,omitempty"`   // thumbnail object
	Attributes       []Attribute  `json:"attributes,omitempty"`  // List of attributes to show below the card. Sample {label}:{value.icon} {value.label}
	ID               string       `json:"id"`                    // min 1 max unlimited, An id that will help HipChat recognise the same card when it is sent multiple times
	Icon             *Icon        `json:"icon,omitempty"`        // icon to display
}

// Body encompases the structure needed to post
// data to a specific room
type Body struct {
	From          string `json:"from,omitempty"`           // min 0 max 64
	MessageFormat string `json:"message_format,omitempty"` // html or text
	Color         string `json:"color,omitempty"`          // yellow, green, red, purple, gray, random
	AttachTo      string `json:"attach_to,omitempty"`      // min 0 max 36
	Notify        bool   `json:"notify,omitempty"`         // Default false
	Message       string `json:"message"`                  // min 0 max 10,000
	Card          *Card  `json:"card,omitempty"`
}

// Colors mapping.
var defaultColors = [...]string{
	log.DebugLevel:  "green",
	log.TraceLevel:  "gray",
	log.InfoLevel:   "purple",
	log.NoticeLevel: "purple",
	log.WarnLevel:   "yellow",
	log.ErrorLevel:  "red",
	log.PanicLevel:  "red",
	log.AlertLevel:  "red",
	log.FatalLevel:  "red",
}

const (
	gopath          = "GOPATH"
	method          = "POST"
	defaultTemplate = `<p><b>{{ .Level.String }}</b></p>
        <p>{{ ts . }}</p>
        {{ if ne .Line 0 }}
            {{ display_file . }}:{{ .Line }}
        {{ end }}
        <p><b>{{ .Message }}</b></p>
        {{ range $f := .Fields }}
            <p><b>{{ $f.Key }}</b>: {{ $f.Value }}</p>
        {{ end }}`
)

// HipChat object
type HipChat struct {
	http.HTTP
	colors       [9]string
	template     *template.Template
	templateHTML string
	gopath       string
	api          APIVersion
}

// New returns a new instance of the HipChat logger
func New(api APIVersion, remoteHost string, contentType string, authToken string) (*HipChat, error) {

	// test here https://developer.atlassian.com/hipchat/guide/hipchat-rest-api that api token has access

	authToken = "Bearer " + authToken

	client := &stdhttp.Client{}

	req, err := stdhttp.NewRequest("GET", remoteHost+"?auth_test=true", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", authToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != stdhttp.StatusAccepted {
		bt, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("HipChat authorization failed\n %s", string(bt))
	}

	header := make(stdhttp.Header)
	header.Set("Content-Type", contentType)
	header.Set("Authorization", authToken)

	// not checking error because url.Parse() is the only thin that can fail,
	// and we've already checked above that it was OK sending the test request
	h, _ := http.New(strings.TrimRight(remoteHost, "/")+"/notification", method, header)

	h.SetFormatFunc(defaultFormatFunc)

	return &HipChat{
		HTTP:         h,
		colors:       defaultColors,
		templateHTML: defaultTemplate,
		api:          api,
	}, nil
}

// GetDisplayColor returns the color for the given log level
func (hc *HipChat) GetDisplayColor(level log.Level) string {
	return hc.colors[level]
}

// SetEmailTemplate sets Email's html template to be used for email body
func (hc *HipChat) SetEmailTemplate(htmlTemplate string) {
	hc.templateHTML = htmlTemplate
}

// Template returns the HipChats's template
func (hc *HipChat) Template() *template.Template {
	return hc.template
}

// GOPATH returns the GOPATH calculated by HTTP
func (hc *HipChat) GOPATH() string {
	return hc.gopath
}

// Run starts the logger consuming on the returned channed
func (hc *HipChat) Run() chan<- *log.Entry {

	fileDisplay := hc.FilenameDisplay()
	tsFormat := hc.TimestampFormat()

	// parse HipChat htmlTemplate, will panic if fails
	hc.template = template.Must(template.New("hipchat").Funcs(
		template.FuncMap{
			"display_file": func(e *log.Entry) (file string) {

				file = e.File
				if fileDisplay == log.Lshortfile {

					for i := len(file) - 1; i > 0; i-- {
						if file[i] == '/' {
							file = file[i+1:]
							break
						}
					}
				} else {

					// additional check, just in case user does
					// have a $GOPATH but code isn't under it.
					if strings.HasPrefix(file, hc.GOPATH()) {
						file = file[len(hc.GOPATH()):]
					}
				}

				return
			},
			"ts": func(e *log.Entry) (ts string) {
				ts = e.Timestamp.Format(tsFormat)
				return
			},
		},
	).Parse(hc.templateHTML))

	// pre-setup
	if fileDisplay == log.Llongfile {
		// gather $GOPATH for use in stripping off of full name
		// if not found still ok as will be blank
		hc.gopath = os.Getenv(gopath)
		if len(hc.gopath) != 0 {
			hc.gopath += string(os.PathSeparator) + "src" + string(os.PathSeparator)
		}
	}

	ch := make(chan *log.Entry, hc.Buffers())

	for i := 0; i <= int(hc.Workers()); i++ {
		go http.HandleLog(hc, ch)
	}
	return ch
}

func defaultFormatFunc(h http.HTTP) http.Formatter {

	var bt []byte
	var err error

	hc := h.(*HipChat)
	b := new(bytes.Buffer)
	template := hc.Template()

	body := new(Body)
	body.Notify = true

	return func(e *log.Entry) []byte {

		bt = bt[0:0]
		b.Reset()
		body.From = e.ApplicationID
		body.Color = hc.GetDisplayColor(e.Level)

		if err = template.ExecuteTemplate(b, "hipchat", e); err != nil {
			log.WithFields(log.F("error", err)).Error("Error parsing HipChat handler template")
		}

		body.Message = b.String()

		// shouldn't be possible to fail here
		// at least with the default handler...
		bt, _ = json.Marshal(body)

		return bt
	}
}
