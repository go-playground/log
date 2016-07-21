package email

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"

	"gopkg.in/gomail.v2"

	"github.com/go-playground/log"
)

// NOTES:
// - Run "go test" to run tests
// - Run "gocov test | gocov report" to report on test converage by file
// - Run "gocov test | gocov annotate -" to report on all code and functions, those ,marked with "MISS" were never called

// or

// -- may be a good idea to change to output path to somewherelike /tmp
// go test -coverprofile cover.out && go tool cover -html=cover.out -o cover.html

type Client struct {
	conn    net.Conn
	address string
	time    int64
	bufin   *bufio.Reader
	bufout  *bufio.Writer
}

func (c *Client) w(s string) {
	c.bufout.WriteString(s + "\r\n")
	c.bufout.Flush()
}
func (c *Client) r() string {
	reply, err := c.bufin.ReadString('\n')
	if err != nil {
		fmt.Println("e ", err)
	}
	return reply
}

func handleClient(c *Client, closePrematurly bool) string {

	var msg []byte

	c.w("220 Welcome to the Jungle")
	msg = append(msg, c.r()...)

	c.w("250 No one says helo anymore")
	msg = append(msg, c.r()...)

	c.w("250 Sender")
	msg = append(msg, c.r()...)

	c.w("250 Recipient")
	msg = append(msg, c.r()...)

	c.w("354 Ok Send data ending with <CRLF>.<CRLF>")
	for {
		text := c.r()
		bytes := []byte(text)
		msg = append(msg, bytes...)

		// 46 13 10
		if bytes[0] == 46 && bytes[1] == 13 && bytes[2] == 10 {
			break
		}
	}

	if !closePrematurly {
		c.w("250 server has transmitted the message")
	}

	c.conn.Close()

	return string(msg)
}

func TestEmailHandler(t *testing.T) {

	tests := []struct {
		expected string
	}{
		{
			expected: "from@email.com",
		},
		{
			expected: "to@email.com",
		},
		{
			expected: "Subject: debug",
		},
		{
			expected: "DEBUG",
		},
	}

	email := New("localhost", 3041, "", "", "from@email.com", []string{"to@email.com"})
	email.SetBuffersAndWorkers(1, 0)
	email.SetTimestampFormat("MST")
	email.SetKeepAliveTimout(time.Second * 0)
	email.SetEmailTemplate(defaultTemplate)
	email.SetFilenameDisplay(log.Llongfile)
	// email.SetFormatFunc(testFormatFunc)
	log.RegisterHandler(email, log.InfoLevel, log.DebugLevel)

	var msg string

	server, err := net.Listen("tcp", ":3041")
	if err != nil {
		t.Errorf("Expected <nil> Got '%s'", err)
	}

	defer server.Close()

	proceed := make(chan struct{})
	defer close(proceed)

	go func() {
		for {
			conn, err := server.Accept()
			if err != nil {
				msg = ""
				break
			}

			if conn == nil {
				continue
			}

			c := &Client{
				conn:    conn,
				address: conn.RemoteAddr().String(),
				time:    time.Now().Unix(),
				bufin:   bufio.NewReader(conn),
				bufout:  bufio.NewWriter(conn),
			}

			msg = handleClient(c, false)

			proceed <- struct{}{}
		}
	}()

	log.Debug("debug")

	<-proceed

	for i, tt := range tests {
		if !strings.Contains(msg, tt.expected) {
			t.Errorf("Index %d Expected '%s' Got '%s'", i, tt.expected, msg)
		}
	}

	// this is normally not safe, but in these tests won't cause any issue
	// flipping during runtime.
	email.SetFilenameDisplay(log.Lshortfile)

	log.Debug("debug")

	<-proceed

	for i, tt := range tests {
		if !strings.Contains(msg, tt.expected) {
			t.Errorf("Index %d Expected '%s' Got '%s'", i, tt.expected, msg)
		}
	}
}

func TestBadDial(t *testing.T) {
	email := New("localhost", 3041, "", "", "from@email.com", []string{"to@email.com"})
	email.SetFormatFunc(testFormatFunc)
	log.RegisterHandler(email, log.InfoLevel)

	log.Info("info test")
}

func TestBadEmailTemplate(t *testing.T) {
	badTemplate := `{{ .NonExistentField }}` // referencing non-existent field
	email := New("localhost", 3041, "", "", "from@email.com", []string{"to@email.com"})
	email.SetEmailTemplate(badTemplate)
	log.RegisterHandler(email, log.InfoLevel)

	log.Info("info test")
}

func TestBadSend(t *testing.T) {

	email := New("localhost", 3041, "", "", "from@email.com", []string{"to@email.com"})
	log.RegisterHandler(email, log.InfoLevel)

	server, err := net.Listen("tcp", ":3041")
	if err != nil {
		t.Errorf("Expected <nil> Got '%s'", err)
	}

	defer server.Close()

	go func() {
		for {
			conn, err := server.Accept()
			if err != nil {
				break
			}

			if conn == nil {
				continue
			}

			c := &Client{
				conn:    conn,
				address: conn.RemoteAddr().String(),
				time:    time.Now().Unix(),
				bufin:   bufio.NewReader(conn),
				bufout:  bufio.NewWriter(conn),
			}

			handleClient(c, true)
		}
	}()

	log.Info("info")
}

func testFormatFunc(email *Email) Formatter {
	var err error
	b := new(bytes.Buffer)

	return func(e *log.Entry) *gomail.Message {
		b.Reset()

		message := gomail.NewMessage()
		message.SetHeader("From", email.From())
		message.SetHeader("To", email.To()...)

		if err = email.template.ExecuteTemplate(b, "email", e); err != nil {
			log.WithFields(log.F("error", err)).Error("Error parsing Email handler template")
		}

		message.SetHeader("Subject", e.Message)
		message.SetBody(contentType, b.String())

		return message
	}
}
