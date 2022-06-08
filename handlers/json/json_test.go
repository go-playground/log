package json

import (
	"bytes"
	"strings"
	"testing"

	"github.com/go-playground/log/v8"
)

func TestJSONLogger(t *testing.T) {
	var buff bytes.Buffer
	l := New(&buff)
	log.AddHandler(l, log.AllLevels...)
	log.WithField("key", "value").Debug("debug")
	expected := `{"message":"debug","timestamp":"","fields":[{"key":"key","value":"value"}],"level":"DEBUG"}`
	if !strings.HasPrefix(buff.String(), `{"message":"debug","timestamp":"`) || !strings.HasSuffix(strings.TrimSpace(buff.String()), `","fields":[{"key":"key","value":"value"}],"level":"DEBUG"}`) {
		t.Errorf("Expected '%s' Got '%s'", expected, buff.String())
	}
}
