package castwait

import (
	"sync"
)

type Interface interface {
	Wait() error
	Done(err error)
}

// New 基于sync.WaitGroup实现
func New() Interface {
	c := &castWait{
		done: false,
		wg:   sync.WaitGroup{},
		err:  nil,
	}
	c.wg.Add(1)
	return c
}

type castWait struct {
	done bool
	wg   sync.WaitGroup
	err  error
}

// Wait 阻塞等待完成
func (c *castWait) Wait() error {
	c.wg.Wait()
	return c.err
}

// Done 完成
func (c *castWait) Done(err error) {
	c.err = err
	c.wg.Done()
}

// NewCond 基于sync.Cond是实现
func NewCond() Interface {
	return &condImpl{
		done: false,
		C:    sync.NewCond(&sync.Mutex{}),
		err:  nil,
	}
}

type condImpl struct {
	done bool
	C    *sync.Cond
	err  error
}

// Wait 阻塞等待完成
func (c *condImpl) Wait() error {
	c.C.L.Lock()
	defer c.C.L.Unlock()

	for !c.done {
		c.C.Wait()
	}
	return c.err
}

// Done 完成
func (c *condImpl) Done(err error) {
	c.err = err
	c.C.L.Lock()
	c.done = true
	c.C.L.Unlock()
	c.C.Broadcast()
}
