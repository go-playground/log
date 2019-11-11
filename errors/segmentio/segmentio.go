package segmentio

import (
	"fmt"
	"strings"

	"github.com/go-playground/log/v7"

	"github.com/segmentio/errors-go"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

// ErrorsGoWithError is a custom WithError function that can be used by using log's
// SetWithErrorFn function.
func ErrorsGoWithError(e log.Entry, err error) log.Entry {
	// normally would call newEntry, but instead will shallow copy
	// because it's not exposed.
	ne := new(log.Entry)
	*ne = *(&e)

	flds := make([]log.Field, 0, len(e.Fields))
	flds = append(flds, e.Fields...)
	flds = append(flds, log.Field{Key: "error", Value: err.Error()})
	ne.Fields = flds

	var frame errors.Frame

	_, types, tags, stacks, _ := errors.Inspect(err)

	if len(stacks) > 0 {
		frame = stacks[len(stacks)-1][0]
	} else {
		frame = errors.WithStack(err).(stackTracer).StackTrace()[2:][0]
	}

	name := fmt.Sprintf("%n", frame)
	file := fmt.Sprintf("%+s", frame)
	line := fmt.Sprintf("%d", frame)
	ne.Fields = append(ne.Fields, log.Field{Key: "source", Value: fmt.Sprintf("%s: %s:%s", name, file, line)})

	for _, tag := range tags {
		ne.Fields = append(ne.Fields, log.Field{Key: tag.Name, Value: tag.Value})
	}
	if len(types) > 0 {
		ne.Fields = append(ne.Fields, log.Field{Key: "types", Value: strings.Join(types, ",")})
	}
	return *ne
}
