package log

import (
	"strconv"
	"strings"

	"github.com/go-playground/errors/v5"
	runtimeext "github.com/go-playground/pkg/v5/runtime"
)

func errorsWithError(e Entry, err error) Entry {
	ne := newEntry(e)

	switch t := err.(type) {
	case errors.Chain:
		types := make([]byte, 0, 64)
		tags := make([]Field, 0, len(t))
		dedupeTags := make(map[Field]bool)
		dedupeType := make(map[string]bool)
		b := BytePool().Get()
		for _, e := range t {
			b = formatLink(e, b)
			b = append(b, ' ')

			for _, tag := range e.Tags {
				field := Field{Key: tag.Key, Value: tag.Value}
				if dedupeTags[field] {
					continue
				}
				dedupeTags[field] = true
				tags = append(tags, field)
			}
			for _, typ := range e.Types {
				if dedupeType[typ] {
					continue
				}
				dedupeType[typ] = true
				types = append(types, typ...)
				types = append(types, ',')
			}
		}
		frame := runtimeext.StackLevel(2)
		b2 := BytePool().Get()
		b2 = extractSource(b2, frame)
		ne.Fields = append(ne.Fields, Field{Key: "source", Value: string(b2[:len(b2)-1])})
		BytePool().Put(b2)
		ne.Fields = append(ne.Fields, Field{Key: "error", Value: string(b[:len(b)-1])})
		BytePool().Put(b)

		ne.Fields = append(ne.Fields, tags...) // we do it this way to maintain order of error, source as first fields
		if len(types) > 0 {
			ne.Fields = append(ne.Fields, Field{Key: "types", Value: string(types[:len(types)-1])})
		}

	default:
		frame := runtimeext.StackLevel(2)
		b := BytePool().Get()
		b = extractSource(b, frame)
		b = append(b, err.Error()...)
		ne.Fields = append(ne.Fields, Field{Key: "error", Value: string(b)})
		BytePool().Put(b)
	}
	return ne
}

func extractSource(b []byte, source runtimeext.Frame) []byte {
	var funcName string

	idx := strings.LastIndexByte(source.Frame.Function, '.')
	if idx == -1 {
		b = append(b, source.File()...)
	} else {
		funcName = source.Frame.Function[idx+1:]
		remaining := source.Frame.Function[:idx]

		idx = strings.LastIndexByte(remaining, '/')
		if idx > -1 {
			b = append(b, source.Frame.Function[:idx+1]...)
			remaining = source.Frame.Function[idx+1:]
		}

		idx = strings.IndexByte(remaining, '.')
		if idx == -1 {
			b = append(b, remaining...)
		} else {
			b = append(b, remaining[:idx]...)
		}
		b = append(b, '/')
		b = append(b, source.File()...)
	}
	b = append(b, ':')
	b = strconv.AppendInt(b, int64(source.Line()), 10)
	if funcName != "" {
		b = append(b, ':')
		b = append(b, funcName...)
	}
	b = append(b, ' ')
	return b
}

func formatLink(l *errors.Link, b []byte) []byte {
	b = extractSource(b, l.Source)
	if l.Prefix != "" {
		b = append(b, l.Prefix...)
	}

	if _, ok := l.Err.(errors.Chain); !ok {
		if l.Prefix != "" {
			b = append(b, ": "...)
		}
		b = append(b, l.Err.Error()...)
	}
	return b
}
