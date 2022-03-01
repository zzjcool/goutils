# goutils

Golang 常用工具包整理

## reduce

在业务中我们可能经常会遇到某个接口需要高频的调用，例如数据库更新、插入操作，调用某api，在这里每一次的调用都会有额外的成本，所以无论是数据库还是一些api调用接口，都提供了批量操作的方式，通过批量操作来降低开销，提升效率。

位置：[reduce](./reduce/)

## caseWait

Go的groutine使用起来十分的方便，可以帮助开发者快速的实现支持并行的程序，但是并行的程序往往需要根据用户的不同需求进行调度，比较常用的有`sync.WaitGroup`，通常是用于等待多个groutine执行完毕后继续后续的任务。

位置：[caseWait](./caseWait/)

## defaults

一个可以设置struct的默认值的包

位置：[defaults](./defaults/)
