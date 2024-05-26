package reduce

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/zzjcool/goutils/castwait"
)

// ReduceHandle the order of input and output remains consistent
type ReduceHandle[I any, O any] func(datas []I) ([]O, error)

type Reduce[I any, O any] interface {
	Do(I) (O, error)
	Refresh()
	Destroy()
}

const (
	DefaultMaxSize            = 100
	DefaultRefreshMillisecond = 1000
)

type reduceOptions[I, O any] struct {
	maxSize            int
	refreshMillisecond int
	handleFunc         ReduceHandle[I, O]
}

func Builder[I, O any]() *reduceOptions[I, O] {
	return &reduceOptions[I, O]{
		maxSize:            DefaultMaxSize,
		refreshMillisecond: DefaultRefreshMillisecond,
		handleFunc:         nil,
	}
}

// SetMaxSize Set the max size of cache, default is 100
func (r *reduceOptions[I, O]) SetMaxSize(maxSize int) *reduceOptions[I, O] {
	r.maxSize = maxSize
	return r
}

// SetRefreshMillisecond Set the refresh duration of cache, default is 1000 (1s)
func (r *reduceOptions[I, O]) SetRefreshMillisecond(refreshMillisecond int) *reduceOptions[I, O] {
	r.refreshMillisecond = refreshMillisecond
	return r
}

// SetHandleFunc Set the handle function
func (r *reduceOptions[I, O]) SetHandleFunc(do ReduceHandle[I, O]) *reduceOptions[I, O] {
	r.handleFunc = do
	return r
}

func (r *reduceOptions[I, O]) New() (Reduce[I, O], error) {

	if r.handleFunc == nil {
		return nil, errors.New("handleFunc is nil")
	}
	reduce := &reduce[I, O]{
		refreshDuration: time.Millisecond * time.Duration(r.refreshMillisecond),
		ticker:          time.NewTicker(time.Millisecond * time.Duration(r.refreshMillisecond)),
		cleanCh:         make(chan bool),
		maxSize:         r.maxSize,
		addLock:         sync.Mutex{},
		refreshLock:     sync.RWMutex{},
		cache:           []I{},
		output:          []O{},
		do:              r.handleFunc,
		cw:              castwait.New(),
		cnt:             new(int64),
	}
	go reduce.daemon()
	return reduce, nil
}

type reduce[I any, O any] struct {
	ticker          *time.Ticker
	refreshDuration time.Duration
	cleanCh         chan bool
	maxSize         int
	addLock         sync.Mutex
	refreshLock     sync.RWMutex
	cache           []I
	output          []O
	do              ReduceHandle[I, O]
	cw              castwait.Interface
	cnt             *int64
}

func (r *reduce[I, O]) daemon() {
	for {
		select {
		// 定时操作
		case <-r.ticker.C:
			{
				r.Refresh()
			}
		// 关闭清理
		case <-r.cleanCh:
			{
				return
			}
		}
	}
}

// Do 向缓存中增加数据
func (r *reduce[I, O]) Do(data I) (O, error) {

	r.addLock.Lock()

	// 读锁保证只上了一把，如果此时正在refresh操作则等待。
	r.refreshLock.RLock()
	// 需要提前获取到cond，避免refresh的时候被刷
	wait := r.cw
	i := *r.cnt
	atomic.AddInt64(r.cnt, 1)
	r.cache = append(r.cache, data)

	r.refreshLock.RUnlock()
	if len(r.cache) >= r.maxSize {
		r.Refresh()
	}
	r.addLock.Unlock()
	// FIXME: 这里有锁的问题
	err := wait.Wait()
	return r.output[i], err
}

func (r *reduce[I, O]) Refresh() {
	r.refreshLock.Lock()
	defer r.refreshLock.Unlock()
	// 如果没有数据不做任何操作
	if len(r.cache) == 0 {
		return
	}
	var err error
	output, err := r.do(r.cache)
	r.output = output
	r.cache = r.cache[:0]

	r.cnt = new(int64)
	r.ticker.Reset(r.refreshDuration)
	r.cw.Done(err)
	// 刷新cond
	r.cw = castwait.New()
}

func (r *reduce[I, O]) Destroy() {
	close(r.cleanCh)
	r.ticker.Stop()
	r.Refresh()
}
