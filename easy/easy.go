package easy

import (
	"context"
	"time"
)

type cancel func()

func Tick(d time.Duration, f func()) cancel {
	t := time.NewTicker(d)
	stop := make(chan struct{})

	go func() {
		defer t.Stop()
		for {
			select {
			case <-t.C:
				f()
			case <-stop:
				return
			}
		}
	}()

	return func() {
		close(stop)
	}
}

func Ticker(ctx context.Context, d time.Duration, f func()) {
	t := time.NewTicker(d)

	f()

	go func() {
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				f()
			}
		}
	}()
}

// P 获取结构体的指针
func P[T any](val T) *T {
	return &val
}
