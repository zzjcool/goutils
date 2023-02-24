package ferr_test

import (
	"fmt"

	"github.com/zzjcool/goutils/ferr"
)

func Example() {
	// New a error
	err := ferr.New("some errors")
	fmt.Print(err)

	// Wrap a error
	errWrap := ferr.Wrap("other errors", err)
	fmt.Print(errWrap)

	// Trace
	fmt.Print(errWrap.TraceStack())

	//TraceStack
	fmt.Print(errWrap.TraceStack())
 
	// UnWrap
	fmt.Print(errWrap.UnWrap())

	// Convert
	normalErr := error(err)
	ferr.Convert(normalErr)

}
