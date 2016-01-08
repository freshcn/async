package async

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"testing"
	"time"
)

func TestAsync(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println("program start")
	useTime := make([]int64, 2)
	startTime := time.Now().UnixNano()
	get("http://www.baidu.com")
	fmt.Println("name: baidu, done")
	get("http://www.sina.com")
	fmt.Println("name: sina, done")
	get("http://www.sohu.com")
	fmt.Println("name: sohu, done")
	get("http://www.163.com")
	fmt.Println("name: 163, done")
	useTime[0] = time.Now().UnixNano() - startTime
	fmt.Printf("single-threaded use time: %d\n", useTime[0])
	fmt.Println("==============\nasync start")

	startTime = time.Now().UnixNano()
	async := NewAsync()
	async.Add("baidu", get, "http://www.baidu.com")
	async.Add("sina", get, "http://www.sina.com")
	async.Add("sohu", get, "http://www.sohu.com")
	async.Add("163", get, "http://www.163.com")

	if chans, ok := async.Run(); ok {
		rs := <-chans

		if len(rs) == 4 {
			for k, _ := range rs {
				fmt.Printf("name: %s, \t done \n", k)
			}
		} else {
			t.Error("async not execution all task")
		}
	}

	useTime[1] = time.Now().UnixNano() - startTime
	fmt.Printf("async use time: %d \n=====\n", useTime[1])

	if useTime[1] < useTime[0] {
		fmt.Printf("async faster than single-threaded: %d \n", useTime[0]-useTime[1])
	} else {
		fmt.Printf("single-threaded faster than async: %d \n", useTime[1]-useTime[0])
	}
}

func get(url string) string {
	if rs, err := http.Get(url); err == nil {
		defer rs.Body.Close()
		if res, err := ioutil.ReadAll(rs.Body); err == nil {
			return string(res)
		}
	}
	return ""
}
