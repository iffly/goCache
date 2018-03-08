package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"time"
)

func Init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	Init()

	infoChan := make(chan []string, 1)
	RefreshCase1("http://baidu.com", infoChan, 10*time.Second, HttpGet)

	for {
		select {
		case infos := <-infoChan:
			fmt.Println(infos)
		}
	}
}

func RefreshCase1(in string, out chan []string, invl time.Duration, handle func(string) []string) {
	go func() {
		for {
			now := time.Now()

			ret := handle(in)
			out <- ret

			duration := time.Since(now)
			if duration < invl {
				<-time.NewTimer(invl).C
			} else {
				fmt.Println("timeout")
			}
		}
	}()
}

func HttpGet(url string) (infos []string) {
	cli := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	rsp, err := cli.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rsp.Body.Close()

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	infos = []string{string(body)}
	return
}
