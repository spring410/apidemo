package main

import (
	// "bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	// "time"
)

type ReturnData struct {
	Success string `json:"success"`
	Data    string `json:"data"`
	Message string `message:"message"`
}

//the gResult value:
//-1: unknown
//0:false
//1: true
//2: error
var gResult [4]int
var gLock sync.RWMutex
var gLock2 sync.RWMutex
var gCounter int

func QueryUserRequest() int {
	gLock2.Lock()
	gCounter++
	fmt.Println("QueryUserRequest, counter=", gCounter)
	gLock2.Unlock()

	iRes := -1
	ipport := os.Args[1]
	url := "http://" + ipport + "/account/1.0/user?id=1"
	// fmt.Println("URL:>", url)

	// var query = []byte("id=1")
	req, err := http.NewRequest("GET", url, nil)
	// req, err := http.NewRequest("GET", url, bytes.NewBuffer(query))
	// req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
		iRes = 2
		return iRes
	}

	defer resp.Body.Close()

	// fmt.Println("response Status:", resp.Status)
	// fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println("response Body:", string(body))

	var v ReturnData
	e := json.Unmarshal(body, &v)
	// fmt.Println(e)
	if e == nil {
		if v.Success == "true" {
			fmt.Println("ok")
			iRes = 1
		} else {
			fmt.Println("failed")
			iRes = 0
		}
	} else {
		//failed
		iRes = 0
	}
	return iRes

}

func SendQuery(ch chan int) {
	for {
		res := QueryUserRequest()
		ch <- res
	}
}

func CollectResult(ch chan int) {
	for {
		res := <-ch
		gLock.Lock()
		switch res {
		case 0:
			gResult[0]++
		case 1:
			gResult[1]++
		case 2:
			gResult[2]++
		default:
			gResult[3]++
		}
		fmt.Println(gResult[0], gResult[1], gResult[2], gResult[3])
		gLock.Unlock()
	}
}

func main() {
	fmt.Println("start to request...")
	if len(os.Args) < 2 {
		fmt.Println("Usage:", os.Args[0], "ip:port")
		fmt.Println("For example:", os.Args[0], "192.168.1.10:8000")
		os.Exit(1)
	}

	ch := make(chan int, 1000)

	//go counter 2: 10780
	//go counter 3: 10278
	//go counter 4: 9927
	for i := 0; i < 1; i++ {
		go SendQuery(ch)
	}

	// go CollectResult(ch)
	// time.Sleep(10 * time.Second)
	CollectResult(ch)
	fmt.Println("end...")
}
