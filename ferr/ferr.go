package ferr

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

type Interface interface {
	error
	Cause() Interface
	Stack() string
	Trace() string
	TraceStack() string
	UnWrap() error
	Contain(err error) bool
}

type ferr struct {
	err   error     // error
	pcs   []uintptr // stack
	cause error     // cause
}

func (f *ferr) Error() string {
	return f.err.Error()
}

func (f *ferr) Cause() Interface {
	return Convert(f.cause)
}

func (f *ferr) Stack() string {
	frames := runtime.CallersFrames(f.pcs)
	var b strings.Builder
	b.WriteString(f.Error() + "\n")
	for frame, more := frames.Next(); more; frame, more = frames.Next() {
		b.WriteString(fmt.Sprintf("%s()\n\t%s:%d\n", frame.Function, frame.File, frame.Line))
	}
	return b.String()
}

func (f *ferr) Trace() string {
	var e Interface = f

	var b strings.Builder
	for {
		b.WriteString(e.Error())
		e = e.Cause()
		if e == nil {
			break
		}
		b.WriteString(": ")
	}
	return b.String()
}

func (f *ferr) TraceStack() string {
	var e Interface = f

	var b strings.Builder
	for {
		b.WriteString(e.Stack())
		e = e.Cause()
		if e == nil {
			break
		}
		b.WriteString("cause: ")
	}
	return b.String()
}

func (f *ferr) Contain(err error) bool {
	if f == err {
		return true
	}
	if f.UnWrap() == err {
		return true
	}
	if f.Cause() == nil {
		return false
	}
	return f.Cause().Contain(err)
}
func (f *ferr) UnWrap() error {
	return f.err
}
func New(msg string) Interface {
	return newFerrByMsg(msg, 1)
}

func newFerrByMsg(msg string, skip int) *ferr {
	return newFerrByErr(errors.New(msg), skip+1)
}

func newFerrByErr(err error, skip int) *ferr {
	var pcs [32]uintptr
	n := runtime.Callers(skip+2, pcs[:])
	st := pcs[0:n]
	return &ferr{
		err:   err,
		pcs:   st,
		cause: nil,
	}
}

func Convert(err error) Interface {
	if err == nil {
		return nil
	}
	e, ok := err.(Interface)
	if ok {
		return e
	}
	return newFerrByErr(err, 1)
}

func Wrap(msg string, cause error) Interface {
	f := newFerrByMsg(msg, 1)
	f.cause = Convert(cause)
	return f
}
