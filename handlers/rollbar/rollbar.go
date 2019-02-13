// Package rollbar implements a JSON handler.
package rollbar

import (
	log "github.com/go-playground/log"
	rollbar "github.com/rollbar/rollbar-go"
)

// Client that connects to rollbar
type Client interface {
	MessageWithExtras(level string, msg string, extras map[string]interface{})
}

// Config is the configuration of the handler
type Config struct {
	Token       string
	Environment string
	CodeVersion string
	ServerHost  string
	ServerRoot  string
}

// Handler implementation.
type Handler struct {
	client Client
}

// New returns the default implementation of a Client.
func New(config *Config) *Handler {
	return &Handler{
		client: rollbar.NewAsync(
			config.Token,
			config.Environment,
			config.CodeVersion,
			config.ServerHost,
			config.ServerRoot,
		),
	}
}

// Log handles the log entry
func (h *Handler) Log(e log.Entry) {
	var (
		level  string
		extras = make(map[string]interface{}, len(e.Fields))
	)

	switch e.Level {
	case log.DebugLevel:
		level = rollbar.DEBUG
	case log.InfoLevel:
		level = rollbar.INFO
	case log.NoticeLevel, log.WarnLevel:
		level = rollbar.WARN
	case log.ErrorLevel:
		level = rollbar.ERR
	case log.PanicLevel, log.AlertLevel, log.FatalLevel:
		level = rollbar.CRIT
	}

	for _, field := range e.Fields {
		extras[field.Key] = field.Value
	}

	h.client.MessageWithExtras(level, e.Message, extras)
}
