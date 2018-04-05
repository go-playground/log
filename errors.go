package log

import (
	"fmt"
	"strings"

	"github.com/go-playground/errors"
)

func errorsWithError(e Entry, err error) Entry {
	ne := newEntry(e)

	switch t := err.(type) {
	case errors.Chain:
		cause := t[0]
		errField := cause.Err.Error()
		types := make([]string, 0, len(t))
		tags := make([]Field, 0, len(t))
		for _, e := range t {
			if e.Prefix != "" {
				errField = fmt.Sprintf("%s: %s", e.Prefix, errField)
			}
			for _, tag := range e.Tags {
				tags = append(tags, Field{Key: tag.Key, Value: tag.Value})

			}
			types = append(types, e.Types...)
		}
		ne.Fields = append(ne.Fields, Field{Key: "error", Value: errField})
		ne.Fields = append(ne.Fields, Field{Key: "source", Value: cause.Source})
		ne.Fields = append(ne.Fields, tags...) // we do it this way to maintain order of error, source as first fields
		if len(types) > 0 {
			ne.Fields = append(ne.Fields, Field{Key: "types", Value: strings.Join(types, ",")})
		}

	default:
		ne.Fields = append(ne.Fields, Field{Key: "error", Value: err.Error()})
		frame := errors.StackLevel(2)
		name := fmt.Sprintf("%n", frame)
		file := fmt.Sprintf("%+s", frame)
		line := fmt.Sprintf("%d", frame)
		parts := strings.Split(file, "\n\t")
		if len(parts) > 1 {
			file = parts[1]
		}
		ne.Fields = append(ne.Fields, Field{Key: "source", Value: fmt.Sprintf("%s: %s:%s", name, file, line)})
	}
	return ne
}
