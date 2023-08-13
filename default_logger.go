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

	c.addFields("", buff, e.Fields)
	buff.B = append(buff.B, newLine)

	c.m.Lock()
	_, _ = c.writer.Write(buff.B)
	c.m.Unlock()

	BytePool().Put(buff)
}

func (c *Logger) addFields(prefix string, buff *Buffer, fields []Field) {
	for _, f := range fields {

		switch t := f.Value.(type) {
		case string:
			printKey(buff, prefix+f.Key)
			buff.B = append(buff.B, t...)
		case int:
			printKey(buff, prefix+f.Key)
			buff.B = strconv.AppendInt(buff.B, int64(t), base10)
		case int8:
			printKey(buff, prefix+f.Key)
			buff.B = strconv.AppendInt(buff.B, int64(t), base10)
		case int16:
			printKey(buff, prefix+f.Key)
			buff.B = strconv.AppendInt(buff.B, int64(t), base10)
		case int32:
			printKey(buff, prefix+f.Key)
			buff.B = strconv.AppendInt(buff.B, int64(t), base10)
		case int64:
			printKey(buff, prefix+f.Key)
			buff.B = strconv.AppendInt(buff.B, t, base10)
		case uint:
			printKey(buff, prefix+f.Key)
			buff.B = strconv.AppendUint(buff.B, uint64(t), base10)
		case uint8:
			printKey(buff, prefix+f.Key)
			buff.B = strconv.AppendUint(buff.B, uint64(t), base10)
		case uint16:
			printKey(buff, prefix+f.Key)
			buff.B = strconv.AppendUint(buff.B, uint64(t), base10)
		case uint32:
			printKey(buff, prefix+f.Key)
			buff.B = strconv.AppendUint(buff.B, uint64(t), base10)
		case uint64:
			printKey(buff, prefix+f.Key)
			buff.B = strconv.AppendUint(buff.B, t, base10)
		case float32:
			printKey(buff, prefix+f.Key)
			buff.B = strconv.AppendFloat(buff.B, float64(t), 'f', -1, 32)
		case float64:
			printKey(buff, prefix+f.Key)
			buff.B = strconv.AppendFloat(buff.B, t, 'f', -1, 64)
		case bool:
			printKey(buff, prefix+f.Key)
			buff.B = strconv.AppendBool(buff.B, t)
		case []Field:
			c.addFields(prefix+f.Key+".", buff, t)
		default:
			printKey(buff, prefix+f.Key)
			buff.B = append(buff.B, fmt.Sprintf(v, f.Value)...)
		}
	}
}

func printKey(buff *Buffer, key string) {
	buff.B = append(buff.B, space)
	buff.B = append(buff.B, key...)
	buff.B = append(buff.B, equals)
}
