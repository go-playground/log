//go:build go1.21
// +build go1.21

package slog

import (
	"bytes"
	"encoding/json"
	"github.com/go-playground/log/v8"
	"log/slog"
	"strings"
	"testing"
	"testing/slogtest"
)

func TestSlogRedirect(t *testing.T) {
	var buff bytes.Buffer
	log.AddHandler(New(slog.NewJSONHandler(&buff, &slog.HandlerOptions{
		ReplaceAttr: ReplaceAttrFn, // for custom log level output
	})), log.AllLevels...)
	h := slog.Default().Handler()

	results := func() []map[string]any {
		var ms []map[string]any
		for _, line := range bytes.Split(buff.Bytes(), []byte{'\n'}) {
			if len(line) == 0 {
				continue
			}
			var m map[string]any
			if err := json.Unmarshal(line, &m); err != nil {
				panic(err) // In a real test, use t.Fatal.
			}
			ms = append(ms, m)
		}
		return ms
	}
	err := slogtest.TestHandler(h, results)
	if err != nil {
		// if a single error and is time key errors, is ok this logger always sets that.
		// sad this its the only way to hook into these errors because none of concrete and
		// Joined errors has no way to reach into them when not.
		if strings.Count(err.Error(), "\n") != 0 || !strings.Contains(err.Error(), "unexpected key \"time\": a Handler should ignore a zero Record.Time") {
			t.Fatal(err)
		}
	}
}
