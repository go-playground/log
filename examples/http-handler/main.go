package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/log"
	httplogger "github.com/go-playground/log/handlers/http"
	"os"
	"strings"
	"sync"
	"time"
)

func main() {

	hLog, err := httplogger.New(10000, "http://localhost:8888/push-event?key=http-logger-test")

	if err != nil {
		fmt.Println("Could create new http logger: ", err)
		os.Exit(1)
	}

	log.RegisterHandler(hLog, log.AllLevels...)

	/*********************************************************************
	   Set formater for basic text log entry that can be sent to
	*********************************************************************/
	hLog.SetFormatter(func(e *log.Entry) string {
		return fmt.Sprintf("[%s] %s : %s", e.Timestamp.Format(time.RFC3339), strings.ToUpper(e.Level.String()), e.Message)
	})
	hLog.SetNumWorkers(4)

	e := &log.Entry{
		WG:        new(sync.WaitGroup),
		Level:     log.NoticeLevel,
		Message:   "This is a sample message",
		Timestamp: time.Now(),
	}

	for i := 0; i < 30; i++ {
		log.HandleEntry(e)
	}

	/*********************************************************************
	  Set formater for json encoded entry that could be sent to Logstash
	**********************************************************************/
	hLog.SetContentEncoding("application/json")

	hLog.SetFormatter(func(e *log.Entry) string {
		dat := map[string]interface{}{}
		dat["@timestamp"] = e.Timestamp.UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
		dat["log_level"] = e.Level.String()
		dat["message"] = e.Message
		for _, f := range e.Fields {
			dat[f.Key] = f.Value
		}
		msg, _ := json.Marshal(dat)
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

}
