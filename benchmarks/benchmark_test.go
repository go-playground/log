package benchmarks

import (
	"errors"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"gopkg.in/inconshreveable/log15.v2"

	"github.com/Sirupsen/logrus"
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

// NOTE: log is a singleton, which means handlers need to be
// setup only once otherwise each test just adds another log
// handler and results are cumulative... makes benchmarking
// annoying because you have to manipulate the TestMain before
// running the benchmark you want.
func TestMain(m *testing.M) {

	cLog := console.New()
	cLog.SetDisplayColor(false)
	cLog.SetWriter(ioutil.Discard)
	cLog.SetBuffersAndWorkers(3, 3)

	log.RegisterHandler(cLog, log.AllLevels...)

	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetOutput(ioutil.Discard)
	logrus.SetLevel(logrus.InfoLevel)

	os.Exit(m.Run())
}

func BenchmarkLogConsoleTenFieldsParallel(b *testing.B) {

	b.ResetTimer()
	// log setup in TestMain
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
			).Info("Go fast.")
		}

	})
}

func BenchmarkLogConsoleSimpleParallel(b *testing.B) {

	b.ResetTimer()
	// log setup in TestMain
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Info("Go fast.")
		}

	})
}

func BenchmarkLogrusText10Fields(b *testing.B) {

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logrus.WithFields(logrus.Fields{
				"int":               1,
				"int64":             int64(1),
				"float":             3.0,
				"string":            "four!",
				"bool":              true,
				"time":              time.Unix(0, 0),
				"error":             errExample.Error(),
				"duration":          time.Second,
				"user-defined type": _jane,
				"another string":    "done!",
			}).Info("Go fast.")
		}
	})
}

func BenchmarkLogrusTextSimple(b *testing.B) {

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logrus.Info("Go fast.")
		}
	})
}

func BenchmarkLog1510Fields(b *testing.B) {
	logger := log15.New()
	logger.SetHandler(log15.StreamHandler(ioutil.Discard, log15.TerminalFormat()))
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("Go fast.",
				"int", 1,
				"int64", int64(1),
				"float", 3.0,
				"string", "four!",
				"bool", true,
				"time", time.Unix(0, 0),
				"error", errExample.Error(),
				"duration", time.Second,
				"user-defined type", _jane,
				"another string", "done!",
			)
		}
	})
}

func BenchmarkLog15Simple(b *testing.B) {
	logger := log15.New()
	logger.SetHandler(log15.StreamHandler(ioutil.Discard, log15.TerminalFormat()))
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("Go fast.")
		}
	})
}
