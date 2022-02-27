# caseWait

Go的groutine使用起来十分的方便，可以帮助开发者快速的实现支持并行的程序，但是并行的程序往往需要根据用户的不同需求进行调度，比较常用的有`sync.WaitGroup`，通常是用于等待多个groutine执行完毕后继续后续的任务。

但是有时候我们需要多个groutine并行执行后同时等待一个条件满足后再继续执行后续任务，例如之前写的Reduce库提给阻塞的Add方法，当所有的数据flush后所有的Add都停止阻塞，这种行为类似于学校考试，在开考前所有同学提前到考场等待考试开始，考试的开始时间都是统一的，当考试时间到了，广播打铃，所有同学收到信号，开始考试。当然这个时候有同学迟到了，那么因为考试已经开始，他就可以直接进入考场开始考试。

## 需求说明

需要一个`Wait`接口，当调用的时候，如果条件未满足，将阻塞，同时调用`Done`接口的表示条件已经满足，解除所有Wait的阻塞，Wait可以被多个groutine调用，多个Wait接口对应一个Done。

如果发生错误，那么错误将通知给所有的Wait的groutine，

## 接口设计

为了实现上述的需求，我们设计对应的接口：

```go
type Interface interface {
 // Wait 可以阻塞当前Groutine，直到Done被调用，可以获取到Done传入的error
 Wait() error
 // Done 解除所有Wait的阻塞，如果发生错误，将error传入
 Done(err error)
}
```

## 实现细节

需要实现以上的需求，可以想到的方式有两种，一种是使用`sync.WaitGroup`，还有一种是使用`sync.Cond`。同时，根据这个组件的特性，我们把这个库取名为`castwait`

### sync.WaitGroup 实现方式

使用sync.WaitGroup实现起来比较简单，调用Add后，使用Wait后可以产生阻塞。

设计，对应的结构体，为：

```go
type castWait struct {
 wg  sync.WaitGroup
 err error // 保存调用的错误
}
```

对应的接口实现为：

```go
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
```

### sync.Cond 实现方式

sync.Cond包使用的比较少，具体的Cond包的使用可以参考：<https://draveness.me/golang/docs/part3-runtime/ch06-concurrency/golang-sync-primitives/>，下次有时间再补充下sync.Cond的使用。

结构的设计与接口的实现：

```go

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

```

## 使用方式

```go
// import (
//  "github.com/zzjcool/goutils/castWait"
// )
n := 100000
c := New()

wg := sync.WaitGroup{}
wg.Add(n)
for i := 0; i < n; i += 1 {
go func() {
err := c.Wait()
if err!=nil{
    // do ...
}
wg.Done()
}()
}

c.Done(exErr)
wg.Wait()

```

## 仓库地址

<https://github.com/zzjcool/goutils/tree/main/castwait>
