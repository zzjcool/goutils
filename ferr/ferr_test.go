package ferr_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/zzjcool/goutils/ferr"
)

func TestNewEqual(t *testing.T) {
	// Different allocations should not be equal.
	if ferr.New("abc") == ferr.New("abc") {
		t.Errorf(`New("abc") == New("abc")`)
	}
	if ferr.New("abc") == ferr.New("xyz") {
		t.Errorf(`New("abc") == New("xyz")`)
	}

	// Same allocation should be equal to itself (not crash).
	err := ferr.New("jkl")
	if err != err {
		t.Errorf(`err != err`)
	}
}

func TestErrorMethod(t *testing.T) {
	err := ferr.New("abc")
	if err.Error() != "abc" {
		t.Errorf(`New("abc").Error() = %q, want %q`, err.Error(), "abc")
	}
}
func TestWrap(t *testing.T) {
	err := ferr.New("abc")
	err2 := ferr.Wrap("def", err)
	if !err2.Contain(err) {
		t.Errorf("err != err2")
	}

	if !err2.Contain(err2) {
		t.Errorf("err != err2")
	}
}

func TestStack(t *testing.T) {
	err := ferr.New("a some err")
	fmt.Println(err.Stack())

	err2 := ferr.Wrap("a other err", err)

	fmt.Println(err2.Stack())

	fmt.Println(err2.TraceStack())

}
func TestTrace(t *testing.T) {
	err := ferr.New("abc")
	err2 := ferr.Wrap("def", err)
	if err2.Trace() != "def: abc" {
		t.Errorf(`err2.Trace() = %q, want %q`, err2.Trace(), "def: abc")
	}
}

func TestTraceStack(t *testing.T) {
	err := ferr.New("abc")
	err2 := ferr.Wrap("def", err)
	fmt.Println(err2.TraceStack())
}

func TestConvert(t *testing.T) {
	err := ferr.New("abc")
	normalErr := error(err)
	err2 := ferr.Convert(normalErr)

	if err != err2 {
		t.Errorf(`New("abc") != New("abc")`)
	}

	if !err2.Contain(err) {
		t.Errorf(`New("abc") != New("abc")`)
	}

	normalErr2 := errors.New("ddd")

	if err2.Contain(normalErr2) {
		t.Errorf(`normalErr2 == normalErr2`)
	}
}
