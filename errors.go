package log

import (
	"fmt"
	"github.com/go-playground/errors/v5"
	"strconv"
	"strings"

	runtimeext "github.com/go-playground/pkg/v5/runtime"
)

func errorsWithError(e Entry, err error) Entry {
	frame := runtimeext.StackLevel(2)

	switch t := err.(type) {
	case errors.Chain:
		types := make([]byte, 0, 32)
		tags := make([]Field, 0, len(t))
		dedupeTags := make(map[string]bool)
		dedupeType := make(map[string]bool)
		errorBuff := BytePool().Get()
		for _, e := range t {
			errorBuff.B = formatLink(e, errorBuff.B)
			errorBuff.B = append(errorBuff.B, ' ')

			for _, tag := range e.Tags {
				key := fmt.Sprintf("%s-%v", tag.Key, tag.Value)
				if dedupeTags[key] {
					continue
				}
				dedupeTags[key] = true
				tags = append(tags, Field{Key: tag.Key, Value: tag.Value})
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

		sourceBuff := BytePool().Get()
		sourceBuff.B = extractSource(sourceBuff.B, frame)
		e.Fields = append(e.Fields, Field{Key: "source", Value: string(sourceBuff.B[:len(sourceBuff.B)-1])})
		BytePool().Put(sourceBuff)
		e.Fields = append(e.Fields, Field{Key: "error", Value: string(errorBuff.B[:len(errorBuff.B)-1])})
		BytePool().Put(errorBuff)

		e.Fields = append(e.Fields, tags...) // we do it this way to maintain order of error, source as first fields
		if len(types) > 0 {
			e.Fields = append(e.Fields, Field{Key: "types", Value: string(types[:len(types)-1])})
		}

	default:
		errorBuff := BytePool().Get()
		errorBuff.B = extractSource(errorBuff.B, frame)
		errorBuff.B = append(errorBuff.B, err.Error()...)
		e.Fields = append(e.Fields, Field{Key: "error", Value: string(errorBuff.B)})
		BytePool().Put(errorBuff)
	}
	return e
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
	if l.Err != nil {
		if l.Prefix != "" {
			b = append(b, ": "...)
		}
		b = append(b, l.Err.Error()...)
	}
	return b
}
