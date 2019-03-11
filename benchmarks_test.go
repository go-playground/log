package log

import (
	"bytes"
	stderr "errors"
	"testing"

	"github.com/go-playground/errors"
)

func BenchmarkWithError(b *testing.B) {
	err := stderr.New("new error")
	entry := Entry{}
	for i := 0; i < b.N; i++ {
		_ = errorsWithError(entry, err)
	}
}

func BenchmarkWithErrorParallel(b *testing.B) {
	err := stderr.New("new error")
	entry := Entry{}
	b.RunParallel(func(pb *testing.PB) {
		var buf bytes.Buffer
		for pb.Next() {
			buf.Reset()
			_ = errorsWithError(entry, err)
		}
	})
}

func BenchmarkWithErrorExisting(b *testing.B) {
	err := errors.New("new error")
	entry := Entry{}
	for i := 0; i < b.N; i++ {
		_ = errorsWithError(entry, err)
	}
}

func BenchmarkWithErrorExistingParallel(b *testing.B) {
	err := errors.New("new error")
	entry := Entry{}
	b.RunParallel(func(pb *testing.PB) {
		var buf bytes.Buffer
		for pb.Next() {
			buf.Reset()
			_ = errorsWithError(entry, err)
		}
	})
}
