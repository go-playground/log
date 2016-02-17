package log

import "time"

// Traceable interface for a traceable object
type Traceable interface {
	End()
}

// TraceEntry is an object used in creating a trace log entry
type TraceEntry struct {
	start time.Time
	end   time.Time
	entry *Entry
}

// End completes the trace and logs the entry
func (t *TraceEntry) End() {
	t.end = time.Now().UTC()

	if t.entry.Fields == nil {
		t.entry.Fields = make([]Field, 0)
	}

	t.entry.Fields = append(t.entry.Fields,
		F("trace_start", t.start),
		F("trace_end", t.end),
		F("trace_duration_ns", t.end.Sub(t.start).Nanoseconds()),
	)

	Logger.HandleEntry(t.entry)
	Logger.tracePool.Put(t)
}
