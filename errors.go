package log

import (
	"fmt"
	"strings"

	"github.com/go-playground/errors"
)

func errorsWithError(e Entry, err error) Entry {
	ne := newEntry(e)

	if w, ok := err.(*errors.Wrapped); ok {
		cause := errors.Cause(w).(*errors.Wrapped)
		ne.Fields = append(ne.Fields, Field{Key: "error", Value: fmt.Sprintf("%s: %s", cause.Prefix, cause.Err)})
		ne.Fields = append(ne.Fields, Field{Key: "source", Value: cause.Source})
		if len(w.Errors) > 0 {
			// top level error
			types := make([]string, 0, len(w.Errors))
			for _, e := range w.Errors {
				for _, tag := range e.Tags {
					ne.Fields = append(ne.Fields, Field{Key: tag.Key, Value: tag.Value})
				}
				types = append(types, e.Types...)
			}
			if len(types) > 0 {
				ne.Fields = append(ne.Fields, Field{Key: "types", Value: strings.Join(types, ",")})
			}
		} else {
			// not top level, probably cause
			for _, tag := range w.Tags {
				ne.Fields = append(ne.Fields, Field{Key: tag.Key, Value: tag.Value})
			}
			if len(w.Types) > 0 {
				ne.Fields = append(ne.Fields, Field{Key: "types", Value: strings.Join(w.Types, ",")})
			}
		}

	} else {
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
