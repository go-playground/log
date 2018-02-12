package json

import (
	"bytes"
	"strings"
	"testing"

	"github.com/go-playground/log"
)

func TestJSONLogger(t *testing.T) {
	var buff bytes.Buffer
	l := New(&buff)
	log.AddHandler(l, log.AllLevels...)
	log.WithField("key", "value").Debug("debug")

	expected := `{"message":"debug","timestamp":"","fields":[{"key":"key","value":"value"}],"level":0}`
	if !strings.HasPrefix(buff.String(), `{"message":"debug","timestamp":"`) || !strings.HasSuffix(strings.TrimSpace(buff.String()), `","fields":[{"key":"key","value":"value"}],"level":0}`) {
		t.Errorf("Expected '%s' Got '%s'", expected, buff.String())
	}
}
