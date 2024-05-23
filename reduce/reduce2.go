package reduce

import (
	"sync"
	"time"

	"github.com/zzjcool/goutils/castwait"
)

type Reduce[I any, O any] interface {
	Do(I) (O, error)
	Refresh()
	Destroy()
}

func NewV2[I any, O any](do HandleFunc, refreshMillisecond int, maxSize int) Reduce[I, O] {
	reduce := &reduce[I, O]{
		refreshDuration: time.Millisecond * time.Duration(refreshMillisecond),
		ticker:          time.NewTicker(time.Millisecond * time.Duration(refreshMillisecond)),
		cleanCh:         make(chan bool),
		maxSize:         maxSize,
		addLock:         sync.Mutex{},
		refreshLock:     sync.RWMutex{},
		cache:           []interface{}{},
		do:              do,
		cw:              castwait.New(),
	}
	go reduce.daemon()
	return reduce
}

type reduce[I any, O any] struct {
	ticker          *time.Ticker
	refreshDuration time.Duration
	cleanCh         chan bool
	maxSize         int
	addLock         sync.Mutex
	refreshLock     sync.RWMutex
	cache           []interface{}
	do              HandleFunc
	cw              castwait.Interface
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

// Add 向缓存中增加数据
func (r *reduce[I, O]) Do(data I) (O,error) {
	r.addLock.Lock()
	defer r.addLock.Unlock()
	// 读锁保证只上了一把，如果此时正在refresh操作则等待。
	r.refreshLock.RLock()
	// 需要提前获取到cond，避免refresh的时候被刷
	wait := r.cw
	r.cache = append(r.cache, data)
	if len(r.cache) >= r.maxSize {
		r.refreshLock.RUnlock()
		r.Refresh()
		wait.Wait()
		return *new(O),nil
	}
	r.refreshLock.RUnlock()
	return *new(O),nil
}

func (r *reduce[I, O]) Refresh() {
	r.refreshLock.Lock()
	defer r.refreshLock.Unlock()
	// 如果没有数据不做任何操作
	if len(r.cache) == 0 {
		return
	}
	err := r.do(r.cache)
	r.cache = r.cache[:0]
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
