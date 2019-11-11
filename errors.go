package log

import (
	"fmt"

	"github.com/go-playground/errors/v5"
	runtimeext "github.com/go-playground/pkg/v4/runtime"
)

func errorsWithError(e Entry, err error) Entry {
	ne := newEntry(e)

	switch t := err.(type) {
	case errors.Chain:
		cause := t[0]
		errField := cause.Err.Error()
		types := make([]byte, 0, 64)
		tags := make([]Field, 0, len(t))
		for i, e := range t {
			if e.Prefix != "" {
				errField = fmt.Sprintf("%s: %s", e.Prefix, errField)
			}
			for _, tag := range e.Tags {
				tags = append(tags, Field{Key: tag.Key, Value: tag.Value})

			}
			for j, typ := range e.Types {
				types = append(types, typ...)
				if i == len(t)-1 && j == len(e.Types)-1 {
					continue
				}
				types = append(types, ',')
			}
		}
		ne.Fields = append(ne.Fields, Field{Key: "error", Value: errField})
		ne.Fields = append(ne.Fields, Field{Key: "source", Value: fmt.Sprintf("%s: %s:%d", cause.Source.Function(), cause.Source.File(), cause.Source.Line())})
		ne.Fields = append(ne.Fields, tags...) // we do it this way to maintain order of error, source as first fields
		if len(types) > 0 {
			ne.Fields = append(ne.Fields, Field{Key: "types", Value: string(types)})
		}

	default:
		ne.Fields = append(ne.Fields, Field{Key: "error", Value: err.Error()})
		frame := runtimeext.StackLevel(2)
		ne.Fields = append(ne.Fields, Field{Key: "source", Value: fmt.Sprintf("%s: %s:%d", frame.Function(), frame.File(), frame.Line())})
	}
	return ne
}
