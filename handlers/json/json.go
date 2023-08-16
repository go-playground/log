// Package json implements a JSON handler.
package json

import (
	stdjson "encoding/json"
	"io"
	"sync"

	log "github.com/go-playground/log/v8"
)

// Handler implementation.
type Handler struct {
	m sync.Mutex
	*stdjson.Encoder
}

// New handler.
func New(w io.Writer) *Handler {
	return &Handler{
		Encoder: stdjson.NewEncoder(w),
	}
}

// Log handles the log entry
func (h *Handler) Log(e log.Entry) {
	h.m.Lock()
	_ = h.Encoder.Encode(e)
	h.m.Unlock()
}
