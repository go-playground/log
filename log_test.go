package log

import (
	"fmt"
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

func TestColors(t *testing.T) {

	fmt.Printf("%sBlack%s\n", Black, Reset)
	fmt.Printf("%sDarkGray%s\n", DarkGray, Reset)
	fmt.Printf("%sBlue%s\n", Blue, Reset)
	fmt.Printf("%sLightBlue%s\n", LightBlue, Reset)
	fmt.Printf("%sGreen%s\n", Green, Reset)
	fmt.Printf("%sLightGreen%s\n", LightGreen, Reset)
	fmt.Printf("%sCyan%s\n", Cyan, Reset)
	fmt.Printf("%sLightCyan%s\n", LightCyan, Reset)
	fmt.Printf("%sRed%s\n", Red, Reset)
	fmt.Printf("%sLightRed%s\n", LightRed, Reset)
	fmt.Printf("%sMagenta%s\n", Magenta, Reset)
	fmt.Printf("%sLightMagenta%s\n", LightMagenta, Reset)
	fmt.Printf("%sBrown%s\n", Brown, Reset)
	fmt.Printf("%sYellow%s\n", Yellow, Reset)
	fmt.Printf("%sLightGray%s\n", LightGray, Reset)
	fmt.Printf("%sWhite%s\n", White, Reset)

	fmt.Printf("%s%sUnderscoreRed%s\n", Red, Underscore, Reset)
	fmt.Printf("%s%sBlinkRed%s\n", Red, Blink, Reset)
	fmt.Printf("%s%s%sBlinkUnderscoreRed%s\n", Red, Blink, Underscore, Reset)

	fmt.Printf("%s%sRedInverse%s\n", Red, Inverse, Reset)
	fmt.Printf("%sGreenInverse%s\n", Green+Inverse, Reset)
}
