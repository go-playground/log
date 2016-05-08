package redis

import (
	//"bytes"
	"fmt"
	stdlog "log"
	"math/rand"
	"time"

	redisclient "github.com/garyburd/redigo/redis"
	"github.com/go-playground/log"
)

// FormatFunc is the function that the workers use to create
// a new Formatter per worker allowing reusable go routine safe
// variable to be used within your Formatter function.
type FormatFunc func() Formatter

// Formatter is the function used to format the Redis entry
type Formatter func(e *log.Entry) []byte

// Redis is an instance of the redis logger
type Redis struct {
	buffer     uint // channel buffer
	redisHosts []string
	formatFunc FormatFunc
	encoding   string
	redisList  string
	numWorkers uint
}

// New returns a new instance of the redis logger
func New(redisHosts []string, list string, encoding string) (*Redis, error) {

	r := &Redis{
		buffer:     0,
		redisHosts: redisHosts,
		encoding:   "",
		formatFunc: formatFunc,
		redisList:  "logs",
		numWorkers: 1,
	}

	return r, nil
}

// SetBuffersAndWorkers sets the channels buffer size and number of concurrent workers.
// These settings should be thought about together, hence setting both in the same function.
func (r *Redis) SetBuffersAndWorkers(size uint, workers uint) {
	r.buffer = size

	if workers == 0 {
		// just in case no log registered yet
		stdlog.Println("Invalid number of workers specified, setting to 1")
		log.Warn("Invalid number of workers specified, setting to 1")

		workers = 1
	}

	r.numWorkers = workers
}

// SetFormatFunc sets FormatFunc each worker will call to get
// a Formatter func
func (r *Redis) SetFormatFunc(fn FormatFunc) {
	r.formatFunc = fn
}

// Run starts the logger consuming on the returned channed
func (r *Redis) Run() chan<- *log.Entry {

	ch := make(chan *log.Entry, r.buffer)

	for i := 0; i <= int(r.numWorkers); i++ {
		go r.handleLog(ch)
	}

	return ch
}

func formatFunc() Formatter {

	var b []byte

	return func(e *log.Entry) []byte {
		b = b[0:0]

		b = append(b, "TEST"...)
		return b
	}
}

func (r *Redis) handleLog(entries <-chan *log.Entry) {
	var e *log.Entry
	var b []byte
	formatter := r.formatFunc()

	for e = range entries {

		defer e.Consumed()

		b = formatter(e)

		// TODO: look into a pool of clients, maybe autoreconnect/retry as well.

		// Connect to the redis instance
		c, err := redisclient.DialTimeout("tcp", r.redisHosts[rand.Intn(len(r.redisHosts))], 0, 1*time.Second, 1*time.Second)
		if err != nil {
			log.Info(fmt.Sprintf("[ERROR] Could not connect to Redis: %s\n", err.Error()))
			continue
		}
		defer c.Close()

		// Select the database
		_, err = c.Do("SELECT", "0")
		if err != nil {
			log.Info(fmt.Sprintf("[ERROR] Could not select Redis DB: %s\n", err.Error()))
			continue
		}

		// Issue the command to push the entry onto the designated list
		_, err = c.Do("RPUSH", r.redisList, string(b))
		if err != nil {
			log.Info(fmt.Sprintf("[ERROR] Could not select Redis DB: %s\n", err.Error()))
			continue
		}
	}
}
