package log

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"
)

const (
	space   = byte(' ')
	equals  = byte('=')
	newLine = byte('\n')
	base10  = 10
	v       = "%v"
)

// ConsoleBuilder is used to create a new console logger
type ConsoleBuilder struct {
	writer          io.Writer
	timestampFormat string
}

// NewConsoleBuilder creates a new ConsoleBuilder for configuring and creating a new console logger
func NewConsoleBuilder() *ConsoleBuilder {
	return &ConsoleBuilder{
		writer:          os.Stderr,
		timestampFormat: DefaultTimeFormat,
	}
}

func (b *ConsoleBuilder) WithWriter(writer io.Writer) *ConsoleBuilder {
	b.writer = writer
	return b
}

func (b *ConsoleBuilder) WithTimestampFormat(format string) *ConsoleBuilder {
	b.timestampFormat = format
	return b
}

func (b *ConsoleBuilder) Build() *Logger {
	return &Logger{
		writer:          b.writer,
		timestampFormat: b.timestampFormat,
	}
}

// Logger is an instance of the console logger
type Logger struct {
	m               sync.Mutex
	writer          io.Writer
	timestampFormat string
}

// Log handles the log entry
func (c *Logger) Log(e Entry) {
	var lvl string
	var i int
	buff := BytePool().Get()
	buff.B = append(buff.B, e.Timestamp.Format(c.timestampFormat)...)
	buff.B = append(buff.B, space)

	lvl = e.Level.String()

	for i = 0; i < 6-len(lvl); i++ {
		buff.B = append(buff.B, space)
	}

	buff.B = append(buff.B, lvl...)
	buff.B = append(buff.B, space)
	buff.B = append(buff.B, e.Message...)

	for _, f := range e.Fields {
		buff.B = append(buff.B, space)
		buff.B = append(buff.B, f.Key...)
		buff.B = append(buff.B, equals)

		switch t := f.Value.(type) {
		case string:
			buff.B = append(buff.B, t...)
		case int:
			buff.B = strconv.AppendInt(buff.B, int64(t), base10)
		case int8:
			buff.B = strconv.AppendInt(buff.B, int64(t), base10)
		case int16:
			buff.B = strconv.AppendInt(buff.B, int64(t), base10)
		case int32:
			buff.B = strconv.AppendInt(buff.B, int64(t), base10)
		case int64:
			buff.B = strconv.AppendInt(buff.B, t, base10)
		case uint:
			buff.B = strconv.AppendUint(buff.B, uint64(t), base10)
		case uint8:
			buff.B = strconv.AppendUint(buff.B, uint64(t), base10)
		case uint16:
			buff.B = strconv.AppendUint(buff.B, uint64(t), base10)
		case uint32:
			buff.B = strconv.AppendUint(buff.B, uint64(t), base10)
		case uint64:
			buff.B = strconv.AppendUint(buff.B, t, base10)
		case float32:
			buff.B = strconv.AppendFloat(buff.B, float64(t), 'f', -1, 32)
		case float64:
			buff.B = strconv.AppendFloat(buff.B, t, 'f', -1, 64)
		case bool:
			buff.B = strconv.AppendBool(buff.B, t)
		default:
			buff.B = append(buff.B, fmt.Sprintf(v, f.Value)...)
		}
	}
	buff.B = append(buff.B, newLine)

	c.m.Lock()
	_, _ = c.writer.Write(buff.B)
	c.m.Unlock()

	BytePool().Put(buff)
}
