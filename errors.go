package log

import (
	"strconv"
	"strings"

	"github.com/go-playground/errors/v5"
	runtimeext "github.com/go-playground/pkg/v5/runtime"
	unsafeext "github.com/go-playground/pkg/v5/unsafe"
)

func errorsWithError(e Entry, err error) Entry {
	ne := newEntry(e)
	frame := runtimeext.StackLevel(2)

	switch t := err.(type) {
	case errors.Chain:
		types := make([]byte, 0, 32)
		tags := make([]Field, 0, len(t))
		dedupeTags := make(map[Field]bool)
		dedupeType := make(map[string]bool)
		buff := BytePool().Get()
		for _, e := range t {
			buff.B = formatLink(e, buff.B)
			buff.B = append(buff.B, ' ')

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

		buff2 := BytePool().Get()
		buff2.B = extractSource(buff2.B, frame)
		ne.Fields = append(ne.Fields, Field{Key: "source", Value: unsafeext.BytesToString(buff2.B[:len(buff2.B)-1])})
		BytePool().Put(buff2)
		ne.Fields = append(ne.Fields, Field{Key: "error", Value: unsafeext.BytesToString(buff.B[:len(buff.B)-1])})
		BytePool().Put(buff)

		ne.Fields = append(ne.Fields, tags...) // we do it this way to maintain order of error, source as first fields
		if len(types) > 0 {
			ne.Fields = append(ne.Fields, Field{Key: "types", Value: unsafeext.BytesToString(types[:len(types)-1])})
		}

	default:
		buff := BytePool().Get()
		buff.B = extractSource(buff.B, frame)
		buff.B = append(buff.B, err.Error()...)
		ne.Fields = append(ne.Fields, Field{Key: "error", Value: string(buff.B)})
		BytePool().Put(buff)
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
