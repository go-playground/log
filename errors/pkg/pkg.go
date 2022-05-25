package pkg

import (
	"fmt"
	"strings"

	"github.com/go-playground/log/v8"
	"github.com/pkg/errors"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

// ErrorsWithError is a custom WithError function that can be used by using log's
// SetWithErrorFn function.
func ErrorsWithError(e log.Entry, err error) log.Entry {
	// normally would call newEntry, but instead will shallow copy
	// because it's not exposed.
	ne := new(log.Entry)
	*ne = *(&e)

	flds := make([]log.Field, 0, len(e.Fields))
	flds = append(flds, e.Fields...)
	flds = append(flds, log.Field{Key: "error", Value: err.Error()})
	ne.Fields = flds

	var frame errors.Frame

	if s, ok := err.(stackTracer); ok {
		frame = s.StackTrace()[0]
	} else {
		frame = errors.WithStack(err).(stackTracer).StackTrace()[2:][0]
	}

	name := fmt.Sprintf("%n", frame)
	file := fmt.Sprintf("%+s", frame)
	line := fmt.Sprintf("%d", frame)
	parts := strings.Split(file, "\n\t")
	if len(parts) > 1 {
		file = parts[1]
	}
	ne.Fields = append(ne.Fields, log.Field{Key: "source", Value: fmt.Sprintf("%s: %s:%s", name, file, line)})
	return *ne
}
