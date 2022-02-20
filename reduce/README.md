# reduce

在业务中我们可能经常会遇到某个接口需要高频的调用，例如数据库更新、插入操作，调用某api，在这里每一次的调用都会有额外的成本，所以无论是数据库还是一些api调用接口，都提供了批量操作的方式，通过批量操作来降低开销，提升效率。

有的旧代码，或者业务本身不支持批量操作，所以我们需要对相关的数据进行聚合，转换为批量操作，个人认为这种操作有点类似于MapReduce操作中的Reduce，所以暂且称之为Reduce吧。

## 需求

提供一个缓存，每次将单次的操作加入到这个缓存中，之后一次性从缓存中取出所有的数据，进行批量操作。

当然何时将缓存全部取出来并进行操作也是个问题，所以我们希望可以传入两个值：

* `maxSize`:表示缓存的最大大小
* `refreshMillisecond`:表示缓存刷新的最短周期

也就是可以有两个条件可以触发批量操作，一个是固定时间，一个是达到缓存上限。

为了模拟原来单次操作产生的阻塞，所以也应当支持阻塞，当该批所有的数据被批量处理完成后，代码可以执行下一步逻辑。

同时应当可以捕获返回的错误，便于后续的处理。

## 接口设计

```go
type ReduceWait interface {
 Wait() error
}
type Interface interface {
 Add(data interface{}) ReduceWait
 Destroy()
}
```

对于Reduce核心操作就是`Add`，其作用就是将数据插入到cache中，之后返回个ReduceWait接口。

`ReduceWait`接口只有一个方法`Wait() error`，当插入数据后，可以调用这个方法产生阻塞，等待到内部的数据被消费后将返回error类型。

## 实现细节

具体实现部分主要是交由`ReduceImple`结构体来实现。除了实现Interface的接口的方法外，还有：

* daemon：负责定时处理cache中积压的数据
* refresh 刷新cache中所有的数据，将数据进行批量消费

在refresh的时候应当对Add操作上锁，保证同一时间只有一个refresh，或者Add在执行。

同时比较麻烦的地方是调用Add的时候可能会因为Cache满了调用refresh，调用refresh的时候要保证没有其他的Add改变cache，所以Add需要一把锁
而refresh需要一把读写锁，其中refresh的读锁给Add。

Add返回的Wait我们使用casewait实现，每次refresh将原来的阻塞接触，开启新的casewait即可。

在初始化的时候需要传入一个`HandleFunc`，这个函数的签名为`type HandleFunc func(datas []interface{}) error`
会在每一次的refresh操作中被调用，这里返回的error会被上抛到wait返回。

## 使用方式

```go
// import (
//  "github.com/zzjcool/goutils/reduce"
// )
// refresh调用的函数
 doFunc := func(datas []interface{}) error {
  fmt.Println(len(datas))
  return nil
 }
// 创建一个reduce，300ms刷新一次，同时最大容量是100
 rdc := reduce.New(doFunc, 300, 100)
// 加入一个数据，同时同步等待完成，如果不调用就是异步操作
 err := rdc.Add("test").Wait()
 if err != nil {
  fmt.Println(err)
 }
```
