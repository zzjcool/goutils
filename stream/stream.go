package stream

import (
	"bufio"
	"context"
	"io"
)

func Swap(up io.ReadWriter, dn io.ReadWriter) (err error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	return SwapWithContext(ctx, up, dn)
}

func SwapWithContext(ctx context.Context, up io.ReadWriter, dn io.ReadWriter) (err error) {
	go func() {
		_, err = io.Copy(up, dn)
	}()
	go func() {
		_, err = io.Copy(dn, up)
	}()
	<-ctx.Done()
	return
}

func BufSwapWithContext(ctx context.Context, up io.ReadWriter, dn io.ReadWriter) (err error) {
	upr := bufio.NewReader(up)
	upw := bufio.NewWriter(up)
	dnr := bufio.NewReader(dn)
	dnw := bufio.NewWriter(dn)
	done := make(chan bool, 1)
	go func() {
		<-ctx.Done()
		done <- true
	}()
	go func() {

		_, err = io.Copy(upw, dnr)
		done <- true
	}()
	go func() {
		_, err = io.Copy(dnw, upr)
		done <- true
	}()
	<-done
	return
}
