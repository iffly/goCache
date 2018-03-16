package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	inChan := make(chan interface{}, 1)
	outChan := make(chan interface{}, 1)

	thrd(inChan, outChan, 1*time.Second, HttpGet)
	inChan <- "http://www.google.com/test.html"

	for {
		select {
		case bodys := <-outChan:
			fmt.Println(bodys.(string))
		}
	}
}

func thrd(in chan interface{}, out chan interface{},
	invl time.Duration, handle func(interface{}) interface{}) {
	var info interface{}
	go func() {
		for {
			select {
			case info = <-in:
				out <- handle(info)
			case <-time.NewTimer(invl).C:
				out <- handle(info)
			}
		}
	}()
}

func HttpGet(url interface{}) (infos interface{}) {
	cli := &http.Client{
		Timeout: 1 * time.Second,
	}

	req, err := http.NewRequest("GET", url.(string), nil)
	if err != nil {
		fmt.Println("Failed to new", err)
		return ""
	}

	rsp, err := cli.Do(req)
	if err != nil {
		fmt.Println("Failed to do", err)
		return ""
	}
	defer rsp.Body.Close()

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		fmt.Println("Failed to read", err)
		return ""
	}

	return string(body)
}
