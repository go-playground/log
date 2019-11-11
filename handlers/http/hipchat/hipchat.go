package hipchat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	stdhttp "net/http"
	"strings"

	"github.com/go-playground/log/v7"
	"github.com/go-playground/log/v7/handlers/http"
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
	URL       string `json:"url"`              // min 1 max 250
	URLRetina string `json:"url@2x"`           // min 1 max 250, the thumbnail url in retina
	Width     uint   `json:"width,omitempty"`  // width of image
	Height    uint   `json:"height,omitempty"` // height of image
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
	log.InfoLevel:   "purple",
	log.NoticeLevel: "purple",
	log.WarnLevel:   "yellow",
	log.ErrorLevel:  "red",
	log.PanicLevel:  "red",
	log.AlertLevel:  "red",
	log.FatalLevel:  "red",
}

const (
	method          = "POST"
	defaultTemplate = `<p><b>{{ .Level.String }}</b></p>
        <p>{{ ts . }}</p>
        <p><b>{{ .Message }}</b></p>
        {{ range $f := .Fields }}
            <p><b>{{ $f.Key }}</b>: {{ $f.Value }}</p>
        {{ end }}`
)

// HipChat object
type HipChat struct {
	*http.HTTP
	colors      [8]string
	template    *template.Template
	api         APIVersion
	application string
}

// New returns a new instance of the HipChat logger
func New(api APIVersion, remoteHost string, contentType string, authToken string, application string) (*HipChat, error) {
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

	hc := &HipChat{
		colors:      defaultColors,
		api:         api,
		application: application,
	}

	// not checking error because url.Parse() is the only thin that can fail,
	// and we've already checked above that it was OK sending the test request
	hc.HTTP, _ = http.New(strings.TrimRight(remoteHost, "/")+"/notification", method, header)
	hc.HTTP.SetFormatFunc(formatFunc(hc))
	hc.SetTemplate(defaultTemplate)
	return hc, nil
}

// SetTemplate sets Hipchats html template to be used for email body
func (hc *HipChat) SetTemplate(htmlTemplate string) {
	hc.template = template.Must(template.New("hipchat").Funcs(
		template.FuncMap{
			"ts": func(e log.Entry) (ts string) {
				ts = e.Timestamp.Format(hc.TimestampFormat())
				return
			},
		},
	).Parse(htmlTemplate))
}

func formatFunc(hc *HipChat) http.FormatFunc {
	return func(h *http.HTTP) http.Formatter {
		var bt []byte

		b := new(bytes.Buffer)
		body := new(Body)
		body.Notify = true

		return func(e log.Entry) []byte {

			bt = bt[0:0]
			b.Reset()
			body.From = hc.application
			body.Color = hc.colors[e.Level]

			_ = hc.template.ExecuteTemplate(b, "hipchat", e)
			body.Message = b.String()

			// shouldn't be possible to fail here
			// at least with the default handler...
			bt, _ = json.Marshal(body)
			return bt
		}
	}
}

// Close cleans up any resources and de-registers the handler with the logger
func (hc *HipChat) Close() error {
	log.RemoveHandler(hc)
	return nil
}
