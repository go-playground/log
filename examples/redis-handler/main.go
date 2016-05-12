package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/log"
	redislogger "github.com/go-playground/log/handlers/redis"
	msgpack "gopkg.in/vmihailenco/msgpack.v2"
	"os"
	"strings"
	"sync"
	"time"
)

func main() {

	rLog, err := redislogger.New(10000, []string{"127.0.0.1:6379"})
	rLog.SetRedisList("event-logs")

	if err != nil {
		fmt.Println("Could create new redis logger: ", err)
		os.Exit(1)
	}

	log.RegisterHandler(rLog, log.AllLevels...)

	/*************************************************
		Set formater for basic text log entry
	*************************************************/
	rLog.SetFormatter(func(e *log.Entry) string {
		return fmt.Sprintf("[%s] %s : %s", e.Timestamp.Format(time.RFC3339), strings.ToUpper(e.Level.String()), e.Message)
	})

	e := &log.Entry{
		WG:        new(sync.WaitGroup),
		Level:     log.NoticeLevel,
		Message:   "This is a sample message",
		Timestamp: time.Now(),
	}

	log.HandleEntry(e)

	/*************************************************
		Samaple formater that sets formater to create json encoded entry to be sent to Redis
	*************************************************/

	rLog.SetFormatter(func(e *log.Entry) string {
		dat := map[string]interface{}{}
		dat["event_time"] = e.Timestamp.Format(time.RFC3339)
		dat["log_level"] = e.Level.String()
		dat["message"] = e.Message
		for _, f := range e.Fields {
			dat[f.Key] = f.Value
		}
		msg, err := json.Marshal(dat)
		if err != nil {
			fmt.Printf("[ERROR] Could not encoding to JSON: %v\n", err)
		}
		return string(msg)
	})

	e = &log.Entry{
		WG:        new(sync.WaitGroup),
		Level:     log.NoticeLevel,
		Message:   "Sample application error message.",
		Timestamp: time.Now(),
		Fields:    []log.Field{log.Field{Key: "type", Value: "test-log"}, log.Field{Key: "application_id", Value: "abc123"}},
	}

	log.HandleEntry(e)

	/*************************************************
		Samaple formater that sets formater to create msgpack encoded entry to be sent to Redis
	*************************************************/
	rLog.SetRedisList("event-logs-msgpack")
	rLog.SetFormatter(func(e *log.Entry) string {
		dat := map[string]interface{}{}
		dat["event_time"] = e.Timestamp.Format(time.RFC3339)
		dat["log_level"] = e.Level.String()
		dat["message"] = e.Message
		dat["encoding"] = "json"
		for _, f := range e.Fields {
			dat[f.Key] = f.Value
		}
		b, err := msgpack.Marshal(dat)
		if err != nil {
			panic(err)
		}
		return string(b)
	})

	e = &log.Entry{
		WG:        new(sync.WaitGroup),
		Level:     log.NoticeLevel,
		Message:   "Sample application error message encoded in msgpack.",
		Timestamp: time.Now(),
		Fields: []log.Field{
			log.Field{Key: "type", Value: "test-log"},
			log.Field{Key: "application_id", Value: "abc123"},
			log.Field{Key: "encoding", Value: "msgpack"},
		},
	}

	log.HandleEntry(e)

}
