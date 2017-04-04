package main

import (
	"bytes"
	// "crypto/rand"
	"accounts"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"
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

var gIpport string

const (
	index_szie = 10000 * 10
)

func QueryUserRequest() int {
	gLock2.Lock()
	gCounter++
	// fmt.Println("QueryUserRequest, counter=", gCounter)
	gLock2.Unlock()

	idRand := rand.Intn(index_szie)
	iRes := 3
	ipport := gIpport
	url := fmt.Sprintf("http://%s/account/1.0/user?id=%d", ipport, idRand)
	// url := "http://" + ipport + "/account/1.0/user?id=1"
	fmt.Println("URL:>", url)

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

			// fmt.Println("ok")
			iRes = 1
		} else {
			// fmt.Println("failed")
			iRes = 0
		}
	} else {
		//failed
		iRes = 0
	}
	return iRes

}

func SendQuery(ch chan int) {

	//total=10000, query mysql
	//1. time:16.875256 seconds
	//2. time:11.873569 seconds
	//3. time:13.07 seconds
	//4. time:16.187642 seconds

	//total=50000, query mysql
	//1. time:71.957576 seconds
	//2. time:71.764807 seconds
	//3. time:56.857500 seconds
	//4. time:85.119554 seconds

	total := 50000

	startTime := time.Now().UTC()
	var log string
	log = fmt.Sprintf("start to test, total=%d, time=%s", total, startTime)
	fmt.Println(log)

	// for {
	for i := 0; i < total; i++ {
		res := QueryUserRequest()
		ch <- res
	}

	endTime := time.Now().UTC()
	log = fmt.Sprintf("end to test, total=%d, time=%s", total, endTime)
	fmt.Println(log)

	timesexpend := endTime.Sub(startTime)

	log = fmt.Sprintf("==>>time:%f seconds", timesexpend.Seconds())
	fmt.Println(log)

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
		loginfo := fmt.Sprintf("total failed=%d, ok=%d, http error=%d, unknown=%d", gResult[0], gResult[1], gResult[2], gResult[3])
		fmt.Println(loginfo)
		gLock.Unlock()
	}
}

func CreateNewUser(data []byte) int {

	fmt.Println("CreateNewUser...")

	iRes := -1
	ipport := gIpport
	url := "http://" + ipport + "/account/1.0/user"
	// fmt.Println("URL:>", url)

	var postData = bytes.NewBuffer(data)
	// req, err := http.NewRequest("POST", url, body)
	resp, err := http.Post(url, "application/json;charset=utf-8", postData)
	if err != nil {
		fmt.Println(err)
		return iRes
	}

	defer resp.Body.Close()

	// fmt.Println("response Status:", resp.Status)
	// fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

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

func RadmonManyNewUsers() {

	total := 10000 * 10

	for i := 0; i < total; i++ {
		// time.Sleep(1 * time.Second)
		data := RadomNewOneUser()
		fmt.Println(string(data))
		iRes := CreateNewUser(data)
		fmt.Println("1:OK, other Failed. result=", iRes)

	}
}

func RadomNewOneUser() []byte {

	randName, err := GenerateRandomString(12)
	if err != nil {
		return nil
	}
	name := "test_" + string(randName)
	sex := rand.Intn(3)
	age := rand.Intn(100)
	email := "Email_" + string(randName) + "@test.com"
	phone := "12345678901"
	createdate := time.Now().UTC().Unix()
	newUser := accounts.Account{
		Name:       name,
		Sex:        sex,
		Age:        age,
		Email:      email,
		Phone:      phone,
		CreateDate: createdate}

	fmt.Println(name, sex, age, email, phone, createdate)

	data, err := json.Marshal(newUser)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return []byte(data)
}

//Generate random bytes with specific length
func GenerateRandomBytes(length int) ([]byte, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func GenerateRandomString(length int) (string, error) {
	b := make([]byte, length)
	n, err := rand.Read(b)
	if n != len(b) || err != nil {
		return "", fmt.Errorf("Could not successfully read from the system CSPRNG.")
	}
	return hex.EncodeToString(b), nil
}

func getActionFromArgs() (string, string) {

	var host *string = flag.String("h", "localhost:8000", "Usage for example: -h 192.168.1.10:8000")
	var action *string = flag.String("a", "query", "Usage for example: -h 192.168.1.10:8000 -a [query|new]")
	// var action1 string
	// flag.StringVar(&action1, "a", "this is a test", "help msg for dir")
	flag.Parse()
	// if *host == "err" {
	// 	// fmt.Println("Usage:", os.Args[0], "ip:port")
	// 	// fmt.Println("For example:", os.Args[0], "192.168.1.10:8000")
	// 	os.Exit(1)
	// }

	// fmt.Println(gIpport)
	return *host, *action
}

func main() {
	fmt.Println("start to request...")
	if len(os.Args) < 2 {
		fmt.Println("Usage for example", os.Args[0], " -h 192.168.1.10:8000 -a [query|new]")
		os.Exit(1)
	}

	var action string
	gIpport, action = getActionFromArgs()
	fmt.Println(gIpport, action)

	if action == "new" {
		//new user
		fmt.Println("start to create new user...")
		RadmonManyNewUsers()
		return
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
