package easy

import "time"

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
