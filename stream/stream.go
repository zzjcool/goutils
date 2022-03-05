package stream

import (
	"context"
	"io"
)

func Swap(up io.ReadWriter, dn io.ReadWriter) (err error) {
	return SwapWithContext(context.Background(), up, dn)
}

func SwapWithContext(ctx context.Context, up io.ReadWriter, dn io.ReadWriter) (err error) {

	done := make(chan bool, 1)
	go func() {
		<-ctx.Done()
		done <- true
	}()
	go func() {

		_, err = io.Copy(up, dn)
		done <- true
	}()
	go func() {
		_, err = io.Copy(dn, up)
		done <- true
	}()
	<-done
	return
}
