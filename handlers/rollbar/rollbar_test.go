package rollbar

import (
	"testing"

	log "github.com/go-playground/log"
	fake "github.com/go-playground/log/handlers/rollbar/fake"
	rollbar "github.com/rollbar/rollbar-go"
)

func TestRollbarLogger(t *testing.T) {
	client := &fake.Client{}
	log.AddHandler(&Handler{client: client}, log.AllLevels...)

	entry := log.WithField("my-key", "my-value")

	levels := []string{
		rollbar.DEBUG,
		rollbar.INFO,
		rollbar.WARN,
		rollbar.ERR,
		rollbar.CRIT,
	}

	for _, level := range log.AllLevels {
		switch level {
		case log.DebugLevel:
			entry.Debug("oh no!")
		case log.InfoLevel:
			entry.Info("oh no!")
		case log.WarnLevel:
			entry.Warn("oh no!")
		case log.ErrorLevel:
			entry.Error("oh no!")
		case log.AlertLevel:
			entry.Alert("oh no!")
		}
	}

	if count := client.MessageWithExtrasCallCount(); count != 5 {
		t.Errorf("Expected MessageWithExtrasCallCount to be called 5 times but was called %d", count)
	}

	for i := 0; i < 5; i++ {
		level, message, extras := client.MessageWithExtrasArgsForCall(i)

		if text := "oh no!"; message != text {
			t.Errorf("Expect '%v', Got '%v'", message, text)
		}

		if rollbarLevel := levels[i]; rollbarLevel != level {
			t.Errorf("Expect '%v', Got '%v'", level, rollbarLevel)
		}

		if value, ok := extras["my-key"]; !ok {
			t.Errorf("Expect 'my-key' is not found")
		} else if value != "my-value" {
			t.Errorf("Expect 'my-value', Got %v", value)
		}
	}
}
