// Package json implements a JSON handler.
package json

import (
	jsn "encoding/json"
	"io"
	"sync"

	"github.com/go-playground/log"
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
