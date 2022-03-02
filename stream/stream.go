package stream

import (
	"io"
	"sync"
)

func Swap(up io.ReadWriter, dn io.ReadWriter) (err error) {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		_, err = io.Copy(up, dn)
	}()
	go func() {
		defer wg.Done()
		_, err = io.Copy(dn, up)
	}()
	wg.Wait()
	return
}
