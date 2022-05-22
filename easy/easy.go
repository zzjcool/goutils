package easy

import (
	"time"
)

type cancel func()

func Tick(d time.Duration, f func()) cancel {
	t := time.NewTicker(d)

	go func() {
		for range t.C {
			f()
		}
	}()
	return func() {
		t.Stop()
	}
}

// P 获取结构体的指针
func P[T any](val T) *T {
	return &val
}
