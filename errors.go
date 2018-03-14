package log

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

func pkgErrorsWithError(e Entry, err error) Entry {
	ne := newEntry(e)
	ne.Fields = append(ne.Fields, Field{Key: "error", Value: err.Error()})

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
	ne.Fields = append(ne.Fields, Field{Key: "source", Value: fmt.Sprintf("%s: %s:%s", name, file, line)})
	return ne
}
