package reduce_test

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zzjcool/goutils/reduce"
)

func Example() {

	reduce, error := reduce.Builder[int, int]().
		SetMaxSize(100).
		SetRefreshMillisecond(1000).
		SetHandleFunc(func(datas []int) ([]int, error) {
			// do something
			return datas, nil
		}).
		New()

	if error != nil {
		return
	}

	reduce.Do(1)

}

func TestReduce(t *testing.T) {

	maxSize := 100
	interval := 300
	sum := int64(0)

	reduce, err := reduce.Builder[int, int]().
		SetMaxSize(100).
		SetRefreshMillisecond(interval).
		SetHandleFunc(func(datas []int) ([]int, error) {
			assert.Equal(t, len(datas), maxSize)
			atomic.AddInt64(&sum, int64(len(datas)))
			result := make([]int, len(datas))
			for i := 0; i < len(datas); i++ {
				result[i] = datas[i] + 1
			}
			return result, nil
		}).
		New()

	if err != nil {
		return
	}

	n, m := 100, 1000
	wg := sync.WaitGroup{}

	wg.Add(n)
	forSum := int64(0)
	for i := 0; i < n; i += 1 {
		go func(i int) {
			for j := 0; j < m; j += 1 {
				result, err := reduce.Do(i * j)
				assert.NoError(t, err)
				assert.Equal(t, i*j+1, result)
				atomic.AddInt64(&forSum, 1)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	reduce.Destroy()
	assert.Equal(t, int64(m*n), forSum)
	assert.Equal(t, int64(m*n), sum)
	t.Log(forSum, sum)
}

// TestInterval 测试间隔刷新
func TestInterval(t *testing.T) {

	maxSize := 100
	interval := 300
	sum := int64(0)
	n := 10

	reduce, err := reduce.Builder[int, int]().
		SetMaxSize(100).
		SetRefreshMillisecond(interval).
		SetHandleFunc(func(datas []int) ([]int, error) {
			assert.LessOrEqual(t, len(datas), maxSize)
			atomic.AddInt64(&sum, int64(len(datas)))
			result := make([]int, len(datas))
			for i := 0; i < len(datas); i++ {
				result[i] = datas[i] + 1
			}
			return result, nil
		}).
		New()

	if err != nil {
		return
	}

	for i := 0; i < n; i++ {
		go reduce.Do(1)
	}

	time.Sleep(time.Duration(interval*2) * time.Millisecond)
	assert.Equal(t, int64(n), sum)

}

func BenchmarkReduce2(b *testing.B) {

	reduce, err := reduce.Builder[int, int]().
		SetMaxSize(100).
		SetRefreshMillisecond(500).
		SetHandleFunc(func(datas []int) ([]int, error) {
			return make([]int, len(datas)), nil
		}).
		New()

	if err != nil {
		return
	}

	for n := 0; n < b.N; n++ {
		go reduce.Do(1)
	}
	reduce.Destroy()
}
