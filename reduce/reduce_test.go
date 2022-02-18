package reduce

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestReduce 测试cache满时自动触发
func TestReduce(t *testing.T) {
	maxSize := 100
	interval := 300
	sum := int64(0)
	doFunc := func(datas []interface{}) error {
		assert.Equal(t, len(datas), maxSize)
		atomic.AddInt64(&sum, int64(len(datas)))
		return nil
	}
	rdc := New(doFunc, interval, maxSize)
	n, m := 100, 1000
	wg := sync.WaitGroup{}

	wg.Add(n)
	forSum := int64(0)
	for i := 0; i < n; i += 1 {
		go func() {
			for j := 0; j < m; j += 1 {
				rdc.Add(nil)
				atomic.AddInt64(&forSum, 1)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	rdc.Destroy()
	assert.Equal(t, int64(m*n), forSum)
	assert.Equal(t, int64(m*n), sum)
	t.Log(forSum, sum)
}

// TestInterval 测试间隔刷新
func TestInterval(t *testing.T) {
	maxSize := 100
	interval := 300
	n := 10
	sum := int64(0)
	doFunc := func(datas []interface{}) error {
		assert.LessOrEqual(t, len(datas), maxSize)
		atomic.AddInt64(&sum, int64(len(datas)))
		return nil
	}
	rdc := New(doFunc, interval, maxSize)
	for i := 0; i < n; i++ {
		rdc.Add(nil)
	}
	time.Sleep(time.Duration(interval*2) * time.Millisecond)
	assert.Equal(t, int64(n), sum)

}

func BenchmarkReduce(b *testing.B) {
	doFunc := func(datas []interface{}) error { return nil }
	rdc := New(doFunc, 500, 100)

	for n := 0; n < b.N; n++ {
		rdc.Add(nil)
	}
	rdc.Destroy()
}
