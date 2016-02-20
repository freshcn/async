# golang async

## 简介

通过golang的goruntine来提供一种异步处理程序的能力。

## 应用场景

这在多个耗时长的网络请求（如：调用API接口）时非常有用。其可将顺序执行变为并行计算，可极大提高程序的执行效率。也能更好的发挥出多核CPU的优势。

## 使用

```go
go get github.com/freshcn/async
```

## demo

```go
// 建议程序开启多核支持
runtime.GOMAXPROCS(runtime.NumCPU())

// 耗时操作1
func request1()interface{}{
  //sql request...
}

// 耗时操作2
func request2()interface{}{
  //sql request...
}

// 新建一个async对象
async:=new async.New()

// 添加request1异步请求,第一个参数为该异步请求的唯一logo,第二个参数为异步完成后的回调函数,回调参数类型为func()interface{}
async.Add("request1",request1)
// 添加request2异步请求
async.Add("request2",request2)

// 执行
if chans,ok := async.Run();ok{
    // 将数据从通道中取回,取回的值是一个map[string]interface{}类型,key为async.Add()时添加的logo,interface{}为该logo回调函数返回的结果
    res := <-chans
    // 这里最好判断下是否所有的异步请求都已经执行成功
		if len(res) == 2 {
			for k, v := range res {
				//do something
			}
		} else {
			log.Println("async not execution all task")
		}
}

// 清除掉本次操作的所以数据,方便后续继续使用async对象
async.Clean()
```




