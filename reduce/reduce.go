package reduce

import (
	"fmt"
	"sync"
	"time"
)

type HandleFunc func(datas []interface{})

type Reduce interface {
	Add(data interface{})
	Destroy()
}

// NewReduce 新建一个Reduce，当间隔时间达到或者缓存达到maxSize的时候触发
// HandleFunc 进行批处理的操作
// refreshMillisecond 刷新缓存处理的间隔毫秒
// maxSize 最大缓存大小
func New(do HandleFunc, refreshMillisecond int64, maxSize int64) Reduce {

	reduce := &ReduceImple{
		ticker:    time.NewTicker(time.Duration(refreshMillisecond)),
		refreshCh: make(chan bool),
		cleanCh:   make(chan bool),
		maxSize:   maxSize,
		mu:        sync.Mutex{},
		cache:     make([]interface{}, 0, maxSize),
		do:        do,
	}
	go reduce.daemon()
	return reduce
}

type ReduceImple struct {
	ticker    *time.Ticker
	refreshCh chan bool
	cleanCh   chan bool
	maxSize   int64
	mu        sync.Mutex
	cache     []interface{}
	do        HandleFunc
}

// daemon 负责处理接收的channel消息
func (r *ReduceImple) daemon() {
	for {
		select {
		// 调用refresh
		case <-r.refreshCh:
		// 定时操作
		case <-r.ticker.C:
			r.refresh()
		case _, ok := <-r.cleanCh:
			if !ok {
				fmt.Println("daemon stop")
			}
			return
		}
	}
}

// refresh 刷新cache中所有的数据，将数据进行批量消费
func (r *ReduceImple) refresh() {
	r.mu.Lock()
	defer r.mu.Unlock()
	// 如果没有数据不做任何操作
	if len(r.cache) == 0 {
		return
	}
	r.do(r.cache)
	r.cache = r.cache[:0]
}

// Destroy 销毁Reduce
func (r *ReduceImple) Destroy() {
	r.ticker.Stop()
	close(r.refreshCh)
	r.refresh()
	close(r.cleanCh)
}

// Add 向缓存中增加数据
func (r *ReduceImple) Add(data interface{}) {
	r.mu.Lock()
	r.cache = append(r.cache, data)
	if len(r.cache) >= int(r.maxSize) {
		r.mu.Unlock()
		r.refresh()
		return
	}
	r.mu.Unlock()
}
