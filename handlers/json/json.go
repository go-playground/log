// Package json implements a JSON handler.
package json

import (
	jsn "encoding/json"
	"io"
	"sync"

	"github.com/go-playground/log/v8"
)

// Handler implementation.
type Handler struct {
	*jsn.Encoder
	m sync.Mutex
}

// New handler.
func New(w io.Writer) *Handler {
	return &Handler{
		Encoder: jsn.NewEncoder(w),
	}
}

// Log handles the log entry
func (h *Handler) Log(e log.Entry) {
	h.m.Lock()
	_ = h.Encoder.Encode(e)
	h.m.Unlock()
}

// Close cleans up any resources and de-registers the handler with the logger
func (h *Handler) Close() error {
	log.RemoveHandler(h)
	return nil
}
