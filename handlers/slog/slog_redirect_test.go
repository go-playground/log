package slog

import (
	"bytes"
	"github.com/go-playground/log/v8"
	"log/slog"
	"testing"
)

func TestSlogRedirect(t *testing.T) {
	buff := new(bytes.Buffer)
	log.AddHandler(New(slog.NewJSONHandler(buff, &slog.HandlerOptions{
		ReplaceAttr: ReplaceAttrFn,
	})), log.AllLevels...)
	log.WithFields(log.G("grouped", log.F("key", "value"))).Debug("test")

	expected := `,"level":"DEBUG","msg":"test","grouped":{"key":"value"}}`
	if !bytes.Contains(buff.Bytes(), []byte(expected)) {
		t.Errorf("Expected '%s' Got '%s'", expected, buff.String())
	}
}
