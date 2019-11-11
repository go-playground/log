## log
<img align="right" src="https://raw.githubusercontent.com/go-playground/log/master/logo.png">![Project status](https://img.shields.io/badge/version-7.0.2-green.svg)
[![Build Status](https://travis-ci.org/go-playground/log.svg?branch=master)](https://travis-ci.org/go-playground/log)
[![Coverage Status](https://coveralls.io/repos/github/go-playground/log/badge.svg?branch=master)](https://coveralls.io/github/go-playground/log?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-playground/log)](https://goreportcard.com/report/github.com/go-playground/log)
[![GoDoc](https://godoc.org/github.com/go-playground/log?status.svg)](https://godoc.org/github.com/go-playground/log)
![License](https://img.shields.io/dub/l/vibe-d.svg)
[![Gitter](https://badges.gitter.im/go-playground/log.svg)](https://gitter.im/go-playground/log?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)

Log is a simple, highly configurable, Structured Logging library

Why another logging library?
----------------------------
There's allot of great stuff out there, but also thought a log library could be made more configurable using per handler log levels.

Features
--------
- [x] Logger is simple, only logic to create the log entry and send it off to the handlers and they take it from there.
- [x] Ability to specify which log levels get sent to each handler
- [x] Built-in console, syslog, http, HipChat, json and email handlers
- [x] Handlers are simple to write + easy to register + easy to remove
- [x] Default logger for quick prototyping and cli applications. It is automatically removed when you register one of your own.
- [x] Logger is a singleton ( one of the few instances a singleton is desired ) so the root package registers which handlers are used and any libraries just follow suit.
- [x] Convenient context helpers `GetContext` & `SetContext`
- [x] Works with go-playground/errors extracting types and tags when used with `WithError`, is the default
- [x] Works with pkg/errors when used with `WithError`, must set using `SetWithErrFn`
- [x] Works with segmentio/errors-go extracting types and tags when used with `WithError`, must set using `SetWithErrFn`

Installation
-----------

Use go get 

```go
go get -u github.com/go-playground/log/v7
``` 

Usage
------
import the log package, it is recommended to set up at least one handler, but there is a default console logger.
```go
package main

import (
	"errors"

	"github.com/go-playground/log/v7"
	"github.com/go-playground/log/v7/handlers/console"
)

func main() {
	// There is a default logger with the same settings
	// once any other logger is registered the default logger is removed.
	cLog := console.New(true)
	log.AddHandler(cLog, log.AllLevels...)

	// Trace
	defer log.WithTrace().Info("time to run")

	log.Debug("debug")
	log.Info("info")
	log.Notice("notice")
	log.Warn("warn")
	log.Error("error")
	// log.Panic("panic") // this will panic
	log.Alert("alert")
	// log.Fatal("fatal") // this will call os.Exit(1)

	err := errors.New("the is an error")
	// logging with fields can be used with any of the above
	log.WithError(err).WithFields(log.F("key", "value")).Info("test info")

	// predefined global fields
	log.WithDefaultFields(log.Fields{
		{"program", "test"},
		{"version", "0.1.3"},
	}...)

	log.WithField("key", "value").Info("testing default fields")

	// or request scoped default fields
	logger := log.WithFields(
		log.F("request", "req"),
		log.F("scoped", "sco"),
	)

	logger.WithField("key", "value").Info("test")
}
```

Adding your own Handler
```go
package main

import (
	"bytes"
	"fmt"

	"github.com/go-playground/log/v7"
)

// CustomHandler is your custom handler
type CustomHandler struct {
	// whatever properties you need
}

// Log accepts log entries to be processed
func (c *CustomHandler) Log(e log.Entry) {

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
	b := new(bytes.Buffer)
	b.Reset()
	b.WriteString(e.Message)

	for _, f := range e.Fields {
		fmt.Fprintf(b, " %s=%v", f.Key, f.Value)
	}
	fmt.Println(b.String())
}

func main() {

	cLog := new(CustomHandler)

	log.AddHandler(cLog, log.AllLevels...)

	// Trace
	defer log.WithTrace().Info("took this long")

	log.Debug("debug")
	log.Info("info")
	log.Notice("notice")
	log.Warn("warn")
	log.Error("error")
	// log.Panic("panic") // this will panic
	log.Alert("alert")
	// log.Fatal("fatal") // this will call os.Exit(1)

	// logging with fields can be used with any of the above
	log.WithField("key", "value").Info("test info")
}
```

Log Level Definitions
---------------------

**DebugLevel** - Info useful to developers for debugging the application, not useful during operations.

**InfoLevel** - Normal operational messages - may be harvested for reporting, measuring throughput, etc. - no action required.

**NoticeLevel** - Normal but significant condition. Events that are unusual but not error conditions - might be summarized in an email to developers or admins to spot potential problems - no immediate action required.

**WarnLevel** - Warning messages, not an error, but indication that an error will occur if action is not taken, e.g. file system 85% full - each item must be resolved within a given time.

**ErrorLevel** - Non-urgent failures, these should be relayed to developers or admins; each item must be resolved within a given time.

**PanicLevel** - A "panic" condition usually affecting multiple apps/servers/sites. At this level it would usually notify all tech staff on call.

**AlertLevel** - Action must be taken immediately. Should be corrected immediately, therefore notify staff who can fix the problem. An example would be the loss of a primary ISP connection.

**FatalLevel** - Should be corrected immediately, but indicates failure in a primary system, an example is a loss of a backup ISP connection. ( same as SYSLOG CRITICAL )

Handlers
-------------
Pull requests for new handlers are welcome, please provide test coverage is all I ask.

| Handler | Description                                                                                                                              | Docs                                                                                                                                                              |
| ------- | ---------------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| console | Allows for log messages to be sent to a any writer, default os.Stderr                                                                    | [![GoDoc](https://godoc.org/github.com/go-playground/log/handlers/console?status.svg)](https://godoc.org/github.com/go-playground/log/handlers/console)           |
| syslog  | Allows for log messages to be sent via syslog, includes TLS support.                                                                     | [![GoDoc](https://godoc.org/github.com/go-playground/log/handlers/syslog?status.svg)](https://godoc.org/github.com/go-playground/log/handlers/syslog)             |
| http    | Allows for log messages to be sent via http. Can use the HTTP handler as a base for creating other handlers requiring http transmission. | [![GoDoc](https://godoc.org/github.com/go-playground/log/handlers/http?status.svg)](https://godoc.org/github.com/go-playground/log/handlers/http)                 |
| email   | Allows for log messages to be sent via email.                                                                                            | [![GoDoc](https://godoc.org/github.com/go-playground/log/handlers/email?status.svg)](https://godoc.org/github.com/go-playground/log/handlers/email)               |
| hipchat | Allows for log messages to be sent to a hipchat room.                                                                                    | [![GoDoc](https://godoc.org/github.com/go-playground/log/handlers/http/hipchat?status.svg)](https://godoc.org/github.com/go-playground/log/handlers/http/hipchat) |
| json    | Allows for log messages to be sent to any wrtier in json format.                                                                         | [![GoDoc](https://godoc.org/github.com/go-playground/log/handlers/json?status.svg)](https://godoc.org/github.com/go-playground/log/handlers/json)                 |

Package Versioning
----------
I'm jumping on the vendoring bandwagon, you should vendor this package as I will not
be creating different version with gopkg.in like allot of my other libraries.

Why? because my time is spread pretty thin maintaining all of the libraries I have + LIFE,
it is so freeing not to worry about it and will help me keep pouring out bigger and better
things for you the community.

Benchmarks
----------
###### Run on Macbook Pro 15-inch 2017 using go version go1.9.4 darwin/amd64
NOTE: only putting benchmarks at others request, by no means does the number of allocations 
make one log library better than another!
```go
go test --bench=. -benchmem=true
goos: darwin
goarch: amd64
pkg: github.com/go-playground/log/benchmarks
BenchmarkLogConsoleTenFieldsParallel-8           2000000               946 ns/op            1376 B/op         16 allocs/op
BenchmarkLogConsoleSimpleParallel-8              5000000               296 ns/op             200 B/op          4 allocs/op
```

Special Thanks
--------------
Special thanks to the following libraries that inspired
* [logrus](https://github.com/Sirupsen/logrus) - Structured, pluggable logging for Go.
* [apex log](https://github.com/apex/log) - Structured logging package for Go.
