package reduce

import (
	"fmt"
	"sync"
	"time"

	"github.com/zzjcool/goutils/castwait"
)

type HandleFunc func(datas []interface{}) error

type ReduceWait interface {
	Wait() error
}
type Interface interface {
	Add(data interface{}) ReduceWait
	Destroy()
}

// NewReduce 新建一个Reduce，当间隔时间达到或者缓存达到maxSize的时候触发
// HandleFunc 进行批处理的操作
// refreshMillisecond 刷新缓存处理的间隔毫秒
// maxSize 最大缓存大小
func New(do HandleFunc, refreshMillisecond int, maxSize int) Interface {

	reduce := &ReduceImple{
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

type ReduceImple struct {
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

// daemon 负责处理接收的channel消息
func (r *ReduceImple) daemon() {
	for {
		select {
		// 定时操作
		case <-r.ticker.C:
			{
				fmt.Println("ticker")
				r.refresh()
			}
		// 关闭清理
		case <-r.cleanCh:
			{
				return
			}
		}
	}
}

// refresh 刷新cache中所有的数据，将数据进行批量消费
func (r *ReduceImple) refresh() {
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

// Add 向缓存中增加数据
func (r *ReduceImple) Add(data interface{}) ReduceWait {
	r.addLock.Lock()
	defer r.addLock.Unlock()
	// 读锁保证只上了一把，如果此时正在refresh操作则等待。
	r.refreshLock.RLock()
	// 需要提前获取到cond，避免refresh的时候被刷
	wait := r.cw
	r.cache = append(r.cache, data)
	if len(r.cache) >= r.maxSize {
		r.refreshLock.RUnlock()
		r.refresh()
		return wait
	}
	r.refreshLock.RUnlock()
	return wait
}

// Destroy 销毁Reduce
func (r *ReduceImple) Destroy() {
	close(r.cleanCh)
	r.ticker.Stop()
	r.refresh()
}
