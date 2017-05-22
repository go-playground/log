package main

import (
	"bytes"
	"fmt"

	"github.com/go-playground/log"
)

// CustomHandler is your custom handler
type CustomHandler struct {
	// whatever properties you need
	buffer uint // channel buffer
}

// Run starts the logger consuming on the returned channed
func (c *CustomHandler) Run() chan<- *log.Entry {

	// in a big high traffic app, set a higher buffer
	ch := make(chan *log.Entry, c.buffer)

	// can run as many consumers on the channel as you want,
	// depending on the buffer size or your needs
	go func(entries <-chan *log.Entry) {

		// below prints to os.Stderr but could marshal to JSON
		// and send to central logging server
		//																						       ---------
		// 				                                                                 |----------> | console |
		//                                                                               |             ---------
		// i.e. -----------------               -----------------     Unmarshal    -------------       --------
		//     | app log handler | -- json --> | central log app | --    to    -> | log handler | --> | syslog |
		//      -----------------               -----------------       Entry      -------------       --------
		//      																         |             ---------
		//                                  									         |----------> | DataDog |
		//          																	        	   ---------
		var e *log.Entry
		b := new(bytes.Buffer)

		for e = range entries {

			b.Reset()
			b.WriteString(e.Message)

			for _, f := range e.Fields {
				fmt.Fprintf(b, " %s=%v", f.Key, f.Value)
			}

			fmt.Println(b.String())
			e.Consumed() // done writing the entry
		}

	}(ch)

	return ch
}

func main() {

	cLog := &CustomHandler{
		buffer: 0,
	}

	log.RegisterHandler(cLog, log.AllLevels...)

	// Trace
	defer log.Trace("trace").End()

	log.Debug("debug")
	log.Info("info")
	log.Notice("notice")
	log.Warn("warn")
	log.Error("error")
	// log.Panic("panic") // this will panic
	log.Alert("alert")
	// log.Fatal("fatal") // this will call os.Exit(1)

	// logging with fields can be used with any of the above
	log.WithFields(log.F("key", "value")).Info("test info")
}
