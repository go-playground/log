package log

import (
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

func BenchmarkWithErrorExisting(b *testing.B) {
	err := errors.New("new error")
	entry := Entry{}
	for i := 0; i < b.N; i++ {
		_ = errorsWithError(entry, err)
	}
}

//func BenchmarkWithErrorExisting(b *testing.B) {
//
//}
