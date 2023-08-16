//go:build !go1.21
// +build !go1.21

package log

import (
	"bufio"
	"io"
	stdlog "log"
	"os"
	"strings"
)

var (
	stdLogWriter     *io.PipeWriter
	redirectComplete chan struct{}
)

// RedirectGoStdLog is used to redirect Go's internal std log output to this logger AND registers a handler for slog
// that redirects slog output to this logger.
func RedirectGoStdLog(redirect bool) {
	if (redirect && stdLogWriter != nil) || (!redirect && stdLogWriter == nil) {
		// already redirected or already not redirected
		return
	}
	if !redirect {
		stdlog.SetOutput(os.Stderr)
		// will stop scanner reading PipeReader
		_ = stdLogWriter.Close()
		stdLogWriter = nil
		<-redirectComplete
		return
	}

	ready := make(chan struct{})
	redirectComplete = make(chan struct{})

	// last option is to redirect
	go func() {
		var r *io.PipeReader
		r, stdLogWriter = io.Pipe()
		defer func() {
			_ = r.Close()
		}()

		stdlog.SetOutput(stdLogWriter)
		defer func() {
			close(redirectComplete)
			redirectComplete = nil
		}()

		scanner := bufio.NewScanner(r)
		close(ready)
		for scanner.Scan() {
			txt := scanner.Text()
			if strings.Contains(txt, "error") {
				WithField("stdlog", true).Error(txt)
			} else if strings.Contains(txt, "warning") {
				WithField("stdlog", true).Warn(txt)
			} else {
				WithField("stdlog", true).Notice(txt)
			}
		}
	}()
	<-ready
}
