## log
<img align="center" src="https://raw.githubusercontent.com/go-playground/log/master/logo.png">![Project status](https://img.shields.io/badge/version-8.1.2-green.svg)
[![Test](https://github.com/go-playground/log/actions/workflows/go.yml/badge.svg)](https://github.com/go-playground/log/actions/workflows/go.yml)
[![Coverage Status](https://coveralls.io/repos/github/go-playground/log/badge.svg?branch=master)](https://coveralls.io/github/go-playground/log?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-playground/log)](https://goreportcard.com/report/github.com/go-playground/log)
[![GoDoc](https://godoc.org/github.com/go-playground/log?status.svg)](https://godoc.org/github.com/go-playground/log)
![License](https://img.shields.io/dub/l/vibe-d.svg)

Log is a simple, highly configurable, Structured Logging library

Why another logging library?
----------------------------
There's a lot of great stuff out there, but this library contains a number of unique features noted below using *.

Features
--------
- [x] Logger is simple, only logic to create the log entry and send it off to the handlers, they take it from there.
- [x] Handlers are simple to write + easy to register + easy to remove
- [x] *Ability to specify which log levels get sent to each handler
- [x] *Handlers & Log Levels are configurable at runtime.
- [x] *`WithError` automatically extracts and adds file, line and package in error output.
- [x] *Convenient context helpers `GetContext` & `SetContext`
- [x] *Works with go-playground/errors extracting wrapped errors, types and tags when used with the `WithError` interface. This is the default but configurable to support more or other error libraries using `SetWithErrorFn`.
- [x] *Default logger for quick prototyping and cli applications. It is automatically removed when you register one of your own.

Installation
-----------

Use go get 

```go
go get github.com/go-playground/log/v8@latest
``` 

Usage
------
import the log package.
```go
package main

import (
	"io"
	stdlog "log"

	"github.com/go-playground/errors/v5"
	"github.com/go-playground/log/v8"
)

func main() {
	log.RedirectGoStdLog(true)

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

	err := errors.New("this is the inner error").AddTags(errors.T("inner", "tag"))
	err = errors.Wrap(err, "this is the wrapping error").AddTags(errors.T("outer", "tag"))

	// logging with fields can be used with any of the above
	log.WithError(err).WithFields(log.F("key", "value")).Info("test info")

	// log unwrapped error
	log.WithError(io.EOF).Error("unwrapped error")

	// predefined global fields
	log.WithDefaultFields(log.Fields{
		log.F("program", "test"),
		log.F("version", "0.1.3"),
	}...)

	log.WithField("key", "value").Info("testing default fields")

	// or request scoped default fields
	logger := log.WithFields(
		log.F("request", "req"),
		log.F("scoped", "sco"),
	)

	logger.WithField("key", "value").Info("test")

	stdlog.Println("This was redirected from Go STD output!")
	log.RedirectGoStdLog(false)
	stdlog.Println("This was NOT redirected from Go STD output!")
}
```

Adding your own Handler
```go
package main

import (
	"bytes"
	"fmt"

	"github.com/go-playground/log/v8"
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
		_, _ = fmt.Fprintf(b, " %s=%v", f.Key, f.Value)
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

#### Go 1.21+ slog compatibility

There is a compatibility layer for slog, which allows redirecting slog to this logger and ability to output to an slog.Handler+.

| type                                        | Definition                                                                                                                                                                          |
|---------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| [Handler](_examples/slog/hanlder/main.go)   | This example demonstrates how to redirect the std log and slog to this logger by using it as an slog.Handler.                                                                       |
| [Redirect](_examples/slog/redirect/main.go) | This example demonstrates how to redirect the std log and slog to this logger and output back out to any slog.Handler, as well as any other handler(s) registered with this logger. |

```go

Log Level Definitions
---------------------

| Level  | Description                                                                                                                                                                                               |
|--------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Debug  | Info useful to developers for debugging the application, not useful during normal operations.                                                                                                             |
| Info   | Normal operational messages which may be harvested for reporting, measuring throughput, etc. no action required.                                                                                          |
| Notice | Normal but significant condition. Events that are unusual but not error conditions eg. might be summarized in an email to developers or admins to spot potential problems - no immediate action required. |
| Warn   | Warning messages, not an error, but indication that an error will occur if action is not taken, e.g. file system 85% full. Each item must be resolved within a given time.                                |
| Error  | Non-urgent failures, these should be relayed to developers or admins; each item must be resolved within a given time.                                                                                     |
| Panic  | A "panic" condition usually affecting multiple apps/servers/sites. At this level it would usually notify all tech staff on call.                                                                          |
| Alert  | Action must be taken immediately. Should be corrected immediately, therefore notify staff who can fix the problem. An example would be the loss of a primary ISP connection.                              |
| Fatal  | Should be corrected immediately, but indicates failure in a primary system, an example is a loss of a backup ISP connection. ( same as SYSLOG CRITICAL )                                                  |

Handlers
-------------
Pull requests for new handlers are welcome when they don't pull in dependencies, it is preferred to have a dedicated package in this case.

| Handler | Description                                                                                                                              | Docs                                                                                                                                                              |
| ------- | ---------------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| json    | Allows for log messages to be sent to any wrtier in json format.                                                                         | [![GoDoc](https://godoc.org/github.com/go-playground/log/handlers/json?status.svg)](https://godoc.org/github.com/go-playground/log/handlers/json)                 |

Package Versioning
----------
This package strictly adheres to semantic versioning guidelines.