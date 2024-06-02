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

type IO[I any, O any] struct {
	Input  I
	Output O
}

type IOs[I any, O any] []*IO[I, O]

func (i IOs[I, O]) GetInputs(size int64) []I {
	var inputs []I
	if i == nil {
		return inputs
	}
	for idx := int64(0); idx < size; idx++ {
		inputs = append(inputs, i[idx].Input)
	}
	return inputs
}

type reduceOptions[I, O any] struct {
	maxSize            int
	refreshMillisecond int
	handleFunc         ReduceHandle[I, O]
	isEmptyFunc        func(I) bool
	emptyOutput        O
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

func (r *reduceOptions[I, O]) SetEmpty(isEmptyFunc func(I) bool, emptyOutput O) *reduceOptions[I, O] {
	r.isEmptyFunc = isEmptyFunc
	r.emptyOutput = emptyOutput
	return r
}

func (i IOs[I, O]) SetOutputs(outputs []O) {
	for idx := 0; idx < len(outputs); idx++ {
		i[idx].Output = outputs[idx]
	}
}

func (ro *reduceOptions[I, O]) New() (Reduce[I, O], error) {

	if ro.handleFunc == nil {
		return nil, errors.New("handleFunc is nil")
	}
	reduce := &reduce[I, O]{
		refreshDuration: time.Millisecond * time.Duration(ro.refreshMillisecond),
		ticker:          time.NewTicker(time.Millisecond * time.Duration(ro.refreshMillisecond)),
		cleanCh:         make(chan bool),
		maxSize:         int64(ro.maxSize),
		addLock:         sync.Mutex{},
		refreshLock:     sync.RWMutex{},
		cache:           make(IOs[I, O], ro.maxSize),
		do:              ro.handleFunc,
		cw:              castwait.New(),
		cnt:             new(int64),
		isEmptyFunc:     ro.isEmptyFunc,
		emptyOutput:     ro.emptyOutput,
	}
	go reduce.daemon()
	return reduce, nil
}

type reduce[I any, O any] struct {
	ticker          *time.Ticker
	refreshDuration time.Duration
	cleanCh         chan bool
	maxSize         int64
	addLock         sync.Mutex
	refreshLock     sync.RWMutex
	cache           IOs[I, O]
	do              ReduceHandle[I, O]
	cw              castwait.Interface
	cnt             *int64
	isEmptyFunc     func(I) bool
	emptyOutput     O
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

func (r *reduce[I, O]) Do(input I) (O, error) {

	if r.isEmptyFunc != nil && r.isEmptyFunc(input) {
		return r.emptyOutput, nil
	}
	r.addLock.Lock()

	// 读锁保证只上了一把，如果此时正在refresh操作则等待。
	r.refreshLock.RLock()
	// 需要提前获取到cond，避免refresh的时候被刷
	wait := r.cw

	ioData := &IO[I, O]{
		Input: input,
	}
	r.cache[*r.cnt] = ioData
	atomic.AddInt64(r.cnt, 1)

	r.refreshLock.RUnlock()
	if *r.cnt >= r.maxSize {
		r.Refresh()
	}
	r.addLock.Unlock()
	err := wait.Wait()
	return ioData.Output, err
}

func (r *reduce[I, O]) Refresh() {
	r.refreshLock.Lock()
	defer r.refreshLock.Unlock()
	// 如果没有数据不做任何操作
	if *r.cnt == 0 {
		return
	}
	var err error
	output, err := r.do(r.cache.GetInputs(*r.cnt))
	r.cache.SetOutputs(output)

	*r.cnt = 0
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
