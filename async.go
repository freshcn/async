package async

// 异步执行类，提供异步执行的功能，可快速方便的开启异步执行
// 通过NewAsync() 来创建一个新的异步操作对象
// 通过调用 Add 函数来向异步任务列表中添加新的任务
// 通过调用 Run 函数来获取一个接收返回的channel，当返回结果时将会返回一个map[string][]interface{}
// 的结果集，包括每个异步函数所返回的所有的结果
import (
	"reflect"
)

// 异步执行所需要的数据
type asyncRun struct {
	Handler reflect.Value
	Params  []reflect.Value
}

// 异步执行对象
type Async struct {
	count int
	tasks map[string]asyncRun
}

// 创建一个新的异步执行对象
func NewAsync() Async {
	return Async{tasks: make(map[string]asyncRun)}
}

// 添加异步执行任务
// name 任务名，结果返回时也将放在任务名中
// handler 任务执行函数，将需要被执行的函数导入到程序中
// params 任务执行函数所需要的参数
func (this *Async) Add(name string, handler interface{}, params ...interface{}) bool {
	if _, e := this.tasks[name]; e {
		return false
	}

	handlerValue := reflect.ValueOf(handler)
	if handlerValue.Kind() == reflect.Func {

		paramNum := len(params)

		this.tasks[name] = asyncRun{
			Handler: handlerValue,
			Params:  make([]reflect.Value, paramNum),
		}

		if paramNum > 0 {
			for k, v := range params {
				this.tasks[name].Params[k] = reflect.ValueOf(v)
			}
		}

		this.count++
		return true
	}

	return false
}

// 任务执行函数，成功时将返回一个用于接受结果的channel
// 在所有异步任务都运行完成时，结果channel将会返回一个map[string][]interface{}的结果。
func (this *Async) Run() (chan map[string][]interface{}, bool) {
	if this.count < 1 {
		return nil, false
	}
	result := make(chan map[string][]interface{})
	chans := make(chan map[string]interface{}, this.count)

	go func(result chan map[string][]interface{}, chans chan map[string]interface{}) {
		rs := make(map[string][]interface{})
		defer func(rs map[string][]interface{}) {
			result <- rs
		}(rs)
		for {
			if this.count < 1 {
				break
			}

			select {
			case res := <-chans:
				this.count--
				rs[res["name"].(string)] = res["result"].([]interface{})
			}
		}
	}(result, chans)

	for k, v := range this.tasks {
		go func(name string, chans chan map[string]interface{}, async asyncRun) {
			result := make([]interface{}, 0)
			defer func(name string, chans chan map[string]interface{}) {
				chans <- map[string]interface{}{"name": name, "result": result}
			}(name, chans)

			values := async.Handler.Call(async.Params)

			if valuesNum := len(values); valuesNum > 0 {
				resultItems := make([]interface{}, valuesNum)
				for k, v := range values {
					resultItems[k] = v.Interface()
				}
				result = resultItems
				return
			}
		}(k, chans, v)
	}

	return result, true
}
