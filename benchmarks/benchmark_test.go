package benchmarks

import (
	"errors"
	"io/ioutil"
	"testing"
	"time"

	"github.com/go-playground/log"
	"github.com/go-playground/log/handlers/console"
)

var errExample = errors.New("fail")

type user struct {
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

var _jane = user{
	Name:      "Jane Doe",
	Email:     "jane@test.com",
	CreatedAt: time.Date(1980, 1, 1, 12, 0, 0, 0, time.UTC),
}

const (
	newLine               = byte('\n')
	defaultTS             = "2006-01-02T15:04:05.000000000Z07:00"
	colorFields           = "%s %s%6s%s %-25s"
	colorNoFields         = "%s %s%6s%s %s"
	colorKeyValue         = " %s%s%s=%v"
	colorFieldsCaller     = "%s %s%6s%s %s:%d %-25s"
	colorNoFieldsCaller   = "%s %s%6s%s %s:%d %s"
	noColorFields         = "%s %6s %-25s"
	noColorNoFields       = "%s %6s %s"
	noColorKeyValue       = " %s=%v"
	noColorFieldsCaller   = "%s %6s %s:%d %-25s"
	noColorNoFieldsCaller = "%s %6s %s:%d %s"
	equals                = byte('=')
	v                     = "%v"
	base10                = 10
	space                 = byte(' ')
	colon                 = byte(':')
)

func BenchmarkConsoleParallel(b *testing.B) {

	cLog := console.New()
	cLog.DisplayColor(false)
	cLog.SetWriter(ioutil.Discard)
	cLog.SetBuffersAndWorkers(3, 3)

	log.RegisterHandler(cLog, log.AllLevels...)

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.WithFields(
				log.F("int", 1),
				log.F("int64", int64(1)),
				log.F("float", 3.0),
				log.F("string", "four!"),
				log.F("bool", true),
				log.F("time", time.Unix(0, 0)),
				log.F("error", errExample.Error()),
				log.F("duration", time.Second),
				log.F("user-defined type", _jane),
				log.F("another string", "done!"),
			).Debug("Go fast.")
			// log.Debug("debug")
		}

	})
}

func BenchmarkConsoleSimpleFieldsParallel(b *testing.B) {

	cLog := console.New()
	cLog.DisplayColor(false)
	cLog.SetWriter(ioutil.Discard)
	cLog.SetBuffersAndWorkers(3, 3)

	log.RegisterHandler(cLog, log.AllLevels...)

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Debug("Go fast.")
		}

	})
}

// func BenchmarkConsoleSimple(b *testing.B) {

// 	cLog := New()
// 	cLog.SetWriter(ioutil.Discard)
// 	// cLog.SetFormatFunc(func() Formatter {

// 	// 	b := new(bytes.Buffer)

// 	// 	return func(e *log.Entry) io.WriterTo {
// 	// 		b.WriteString(e.Message)
// 	// 		return b
// 	// 	}
// 	// })

// 	log.RegisterHandler(cLog, log.AllLevels...)

// 	b.ReportAllocs()

// 	for i := 0; i < b.N; i++ {

// 		log.WithFields(
// 			log.F("int", 1),
// 			log.F("int64", int64(1)),
// 			log.F("float", 3.0),
// 			log.F("string", "four!"),
// 			log.F("bool", true),
// 			log.F("time", time.Unix(0, 0)),
// 			log.F("error", errExample.Error()),
// 			log.F("duration", time.Second),
// 			log.F("user-defined type", _jane),
// 			log.F("another string", "done!"),
// 		).Info("Go fast.")
// 		log.Debug("debug")
// 	}
// }
