package http

import (
	//"bytes"
	"fmt"
	redisclient "github.com/garyburd/redigo/redis"
	"github.com/go-playground/log"
	"math/rand"
	"time"
)

// Formatter is the function used to format the Redis entry
type Formatter func(e *log.Entry) string

// Redis is an instance of the redis logger
type Redis struct {
	buffer             uint // channel buffer
	redisHosts         []string
	formatter          Formatter
	encoding           string
	hasCustomFormatter bool
	redisList          string
	numWorkers         int
}

// New Redis client
func New(bufferSize uint, redisHosts []string) (*Redis, error) {

	r := &Redis{
		buffer:             0,
		redisHosts:         []string{"127.0.0.1:6379"},
		encoding:           "",
		hasCustomFormatter: false,
		redisList:          "logs",
		numWorkers:         1,
	}

	r.buffer = bufferSize

	r.redisHosts = redisHosts

	return r, nil
}

// SetBuffer sets the buffer for Redis client
func (r *Redis) SetBuffer(buff uint) {
	r.buffer = buff
}

// SetEncoding sets the data encoding type (none, msgpack, etc)
func (r *Redis) SetEncoding(encoding string) {
	r.encoding = encoding
}

// SetRedisHosts sets the list of redis hosts to attempt connecting to
func (r *Redis) SetRedisHosts(hosts []string) {
	r.redisHosts = hosts
}

// SetRedisList sets the list in which to send the events
func (r *Redis) SetRedisList(list string) {
	r.redisList = list
}

// SetFormatter sets the entry formatter
func (r *Redis) SetFormatter(f Formatter) {
	r.formatter = f
	r.hasCustomFormatter = true
}

func (r *Redis) SetNumWorkers(num uint) {
	if num >= 1 {
		r.numWorkers = num
	}
}

// Run starts the logger consuming on the returned channed
func (r *Redis) Run() chan<- *log.Entry {
	ch := make(chan *log.Entry, r.buffer)
	for i := 0; i <= int(r.numWorkers); i++ {
		go r.handleLog(ch)
	}
	return ch
}

func (r *Redis) handleLog(entries <-chan *log.Entry) {
	var e *log.Entry
	var payload string

ItterateOverItems:
	for e = range entries {

		payload = r.formatter(e)

		// Connect to the redis instance
		rand.Seed(int64(time.Now().Nanosecond()))
		c, err := redisclient.DialTimeout("tcp", r.redisHosts[rand.Intn(len(r.redisHosts))], 0, 1*time.Second, 1*time.Second)
		if err != nil {
			log.Info(fmt.Sprintf("[ERROR] Could not connect to Redis: %s\n", err.Error()))
			goto ItterateOverItems
		}
		defer c.Close()
		// Select the database
		_, err = c.Do("SELECT", "0")
		if err != nil {
			c.Close()
			log.Info(fmt.Sprintf("[ERROR] Could not select Redis DB: %s\n", err.Error()))
		}
		// Issue the command to push the entry onto the designated list
		_, err = c.Do("RPUSH", r.redisList, payload)
		if err != nil {
			log.Info(fmt.Sprintf("[ERROR] Could not select Redis DB: %s\n", err.Error()))
		}
		e.WG.Done()
	}
}
