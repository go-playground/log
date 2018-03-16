package pkg

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/go-playground/log"
	"github.com/pkg/errors"
)

type testHandler struct {
	writer io.Writer
}

func (th *testHandler) Log(e log.Entry) {
	s := e.Level.String() + " "
	s += e.Message

	for _, f := range e.Fields {
		s += fmt.Sprintf(" %s=%v", f.Key, f.Value)
	}
	s += "\n"
	if _, err := th.writer.Write([]byte(s)); err != nil {
		panic(err)
	}
}

func TestWrappedError(t *testing.T) {
	log.SetExitFunc(func(int) {})
	log.SetWithErrorFn(ErrorsWithError)
	buff := new(bytes.Buffer)
	th := &testHandler{
		writer: buff,
	}
	log.AddHandler(th, log.AllLevels...)

	err := fmt.Errorf("this is an %s", "error")
	err = errors.Wrap(err, "prefix")
	log.WithError(err).Error("test")
	expected := "pkg_test.go:41\n"
	if !strings.HasSuffix(buff.String(), expected) {
		t.Errorf("got %s Expected %s", buff.String(), expected)
	}
	buff.Reset()
	expected = "pkg_test.go:50\n"
	err = fmt.Errorf("this is an %s", "error")
	log.WithError(err).Error("test")
	if !strings.HasSuffix(buff.String(), expected) {
		t.Errorf("got %s Expected %s", buff.String(), expected)
	}
}
