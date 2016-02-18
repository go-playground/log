package log

import (
	"testing"

	. "gopkg.in/go-playground/assert.v1"
)

// NOTES:
// - Run "go test" to run tests
// - Run "gocov test | gocov report" to report on test converage by file
// - Run "gocov test | gocov annotate -" to report on all code and functions, those ,marked with "MISS" were never called
//
// or
//
// -- may be a good idea to change to output path to somewherelike /tmp
// go test -coverprofile cover.out && go tool cover -html=cover.out -o cover.html
//

func TestFatal(t *testing.T) {
	var i int

	exitFunc = func(code int) {
		i = code
	}

	Fatal("fatal")
	Equal(t, i, 1)

	Fatalf("fatalf")
	Equal(t, i, 1)

	Fatalln("fatalln")
	Equal(t, i, 1)

	WithFields(F("key", "value")).Fatal("fatal")
	Equal(t, i, 1)

	WithFields(F("key", "value")).Fatalf("fatalf")
	Equal(t, i, 1)
}
