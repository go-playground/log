package segmentio

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/go-playground/log/v8"
	"github.com/segmentio/errors-go"
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
	log.SetWithErrorFn(ErrorsGoWithError)
	buff := new(bytes.Buffer)
	th := &testHandler{
		writer: buff,
	}
	log.AddHandler(th, log.AllLevels...)

	err := fmt.Errorf("this is an %s", "error")
	err = errors.WithTypes(errors.WithTags(errors.Wrap(err, "prefix"),
		errors.T("key", "value"),
	), "Permanent", "Internal")
	log.WithError(err).Error("test")
	expected := "segmentio_test.go:41 key=value types=Internal,Permanent\n"
	if !strings.HasSuffix(buff.String(), expected) {
		t.Errorf("got %s Expected %s", buff.String(), expected)
	}
	buff.Reset()
	expected = "segmentio_test.go:52\n"
	err = fmt.Errorf("this is an %s", "error")
	log.WithError(err).Error("test")
	if !strings.HasSuffix(buff.String(), expected) {
		t.Errorf("got %s Expected %s", buff.String(), expected)
	}
}
