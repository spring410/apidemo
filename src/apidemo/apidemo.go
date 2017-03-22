package main

import (
	"accounts"
	"encoding/json"
	"model"
	"mysqldb"
	_ "utilities"
	"fmt"
	"github.com/alecthomas/log4go"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"strconv"
	"time"
	"path"
)

type ReturnData struct {
	Success string `json:"success"`
	Data    string `json:"data"`
	Message string `message:"message"`
}

var gAccount []accounts.Account
var gCounter int

func HandlerForUndefinitionPath(w http.ResponseWriter, r *http.Request) {
	w.Write(responseNull())
}

/* response data
{
    "success":true,
    "data":{"id":1,"name":"xiaotuan"},
    "message":""
}
{
    "success":false,
    "data":{},
    "message":"no this id or no this email"
}

*/
func HandlerForGetUser(w http.ResponseWriter, r *http.Request) {
	log4go.Debug("entry.", time.Now())

	defer func() {
		gCounter++
		log4go.Debug("exit. counter=%s", gCounter)
	}()

	//read from db

	//query by id first
	var result *accounts.Account
	var err error
	i, err := strconv.Atoi(QueryKeyValue("id", r))
	if err == nil {
		result, err = mysqldb.QueryUserById(i)
	} else {
		//or query by name
		n := QueryKeyValue("name", r)
		if n != "" {
			result, err = mysqldb.QueryUserByName(n)
		}
	}

	// log4go.Info(*result)
	if err == nil {
		w.Write(responseJsonData(*result))
		return
	}
	/* //memory cache
	for _, one := range gAccount {
		if strconv.Itoa(one.ID) == QueryKeyValue("id", r) {
			account = one
			// s := fmt.Sprintf("id=%d, age=%d, email=%s, phone=%s, createdate=%d",
			// 	one.ID, one.Age, one.Email, one.Phone, one.CreateDate)
			// log4go.Info(s)

			w.Write(responseJsonData(account))
			return
		}
	}
	*/

	account := accounts.Account{ID: 0}
	w.Write([]byte(responseJsonData(account)))
}

func responseJsonData(a accounts.Account) []byte {

	if a.ID != 0 {
		aValue, err := json.Marshal(a)
		if err == nil {
			// log4go.Info("aValue=", string(aValue))

			response := ReturnData{
				Success: "true",
				Data:    string(aValue),
				Message: "ok"}

			//注意struct成员要大写，包外可以访问，要不然Json操作不了
			res, _ := json.Marshal(response)
			// log4go.Info(e)
			// log4go.Info(string(res))
			return res
		}
	}

	response := ReturnData{
		Success: "false",

		Message: "not found"}
	res, _ := json.Marshal(response)
	// log4go.Info(string(res))
	return res

}

func nameIsInvalid(msg string) []byte {
	response := ReturnData{
		Success: "false",
		Message: msg}
	res, _ := json.Marshal(response)
	// log4go.Info(string(res))
	return res
}

func responseOk() []byte {
	response := ReturnData{
		Success: "true"}
	res, _ := json.Marshal(response)
	// log4go.Info(string(res))
	return res
}

func responseNull() []byte {
	response := ReturnData{
		Success: "false",
		Message: "Null"}
	res, _ := json.Marshal(response)
	// log4go.Info(string(res))
	return res
}

func HandlerForNewUser(w http.ResponseWriter, r *http.Request) {
	log4go.Debug("entry....")

	defer func() {
		log4go.Debug("exit.")
	}()

	// log4go.Info(*r)
	// keys := r.URL.Query()
	// log4go.Info(keys)

	//check the key if exist
	hasname := QueryKeyValue("name", r)
	log4go.Info("hasname=", hasname)
	if hasname == "" {
		msg := nameIsInvalid("Name is empty.")
		w.Write(msg)
		return
	} else if NameExist(hasname) {
		msg := nameIsInvalid("Name exist.")
		w.Write(msg)
		return
	}

	//add new user to map
	newid := len(gAccount) + 1
	name := QueryKeyValue("name", r)
	sex, _ := strconv.Atoi(QueryKeyValue("sex", r))
	age, _ := strconv.Atoi(QueryKeyValue("age", r))
	email := QueryKeyValue("email", r)
	phone := QueryKeyValue("phone", r)
	createdate := time.Now().UTC().Unix()
	newAccount := accounts.Account{ID: newid,
		Name:       name,
		Sex:        sex,
		Age:        age,
		Email:      email,
		Phone:      phone,
		CreateDate: createdate}
	gAccount = append(gAccount, newAccount)
	// log4go.Info(newAccount)

	err := mysqldb.InsertUser(newid, name, sex, age, email, phone, createdate)
	if err == nil {
		//reponse ok
		w.Write(responseOk())
	} else {
		//failed.
		w.Write(responseNull())
	}

}

func QueryKeyValue(key string, r *http.Request) string {
	vars := r.URL.Query()
	value, ok := vars[key]
	if ok {
		return value[0]
	}
	return ""
}

func NameExist(name string) bool {
	_, err := mysqldb.NameExistInDb(name)
	if err == nil {
		return false
	} else {
		return true
	}
}

func EmailExist(email string) bool {
	for _, one := range gAccount {
		if one.Email == email {
			return true
		}
	}
	return false
}

func PhoneExist(phone string) bool {
	for _, one := range gAccount {
		if one.Phone == phone {
			return true
		}
	}

	return false
}

func initDb() error {
	err := mysqldb.InitDb()
	if err == nil {
		err = mysqldb.PingDb()
		if err != nil {
			log4go.Error("Faild to ping db,", err)
			if mysqldb.DatabaseUnkown(err.Error()) {
				log4go.Error("Please create database, and then try again.")
			}
		} else {
			log4go.Info("Ping db OK.")
			err = mysqldb.CreateUsersTable()
		}
	}

	// log4go.Info(err)
	return err
}

func onlyTest() {

}

func CheckUsage() {
	if len(os.Args) < 2 {
		log4go.Error("Usage:%s ip:port", os.Args[0])
		log4go.Error("For example:%s 192.168.1.10:8000",os.Args[0])

		//waiting log4go to complete
		time.Sleep(1*time.Second)
		os.Exit(1)
	}
}

func initLogger() {
	// execDirAbsPath, _ := os.Getwd()
	// fmt.Println(execDirAbsPath)
	exePath, _ := os.Executable()
	// fmt.Println(exePath)
	exeDir := path.Dir(exePath)
	configPath := exeDir + "/logconfig.xml"	
	fmt.Println("config file:" ,configPath)
	log4go.LoadConfiguration(configPath)
	log4go.Info("Loaded log config file")
}

func main() {
	fmt.Println("Start to run...")
	initLogger()
	log4go.Info("Start to run api demo....pid=%d", os.Getpid())

	//check ip and port
	CheckUsage()

	//init db
	err := initDb()
	if err != nil {
		log4go.Error("Failed to init database.")
		os.Exit(1)
	}

	//only test
	// onlyTest()

	r := mux.NewRouter()

	var strPreFix = model.URL_ROOT + "/" + model.API_VERSOIN
	log4go.Debug("URL path pre is:", strPreFix)
	s := r.PathPrefix(strPreFix).Subrouter()

	// // "/URL_ROOT/"
	// s.HandleFunc("/", ProductsHandler)
	// // "/URL_ROOT/{key}/"
	// s.HandleFunc("/{key}/", ProductHandler)
	// // "/URL_ROOT/{key}/details"
	// s.HandleFunc("/{key}/details", ProductDetailsHandler)

	// Routes consist of a path and a handler function.
	s.HandleFunc("/", HandlerForUndefinitionPath)

	s.HandleFunc("/user", HandlerForGetUser).Methods("GET")

	s.HandleFunc("/user", HandlerForNewUser).Methods("POST")

	// Bind to a port and pass our router in
	ipport := os.Args[1]
	log4go.Info("start listen =%s", ipport)
	log4go.Error(http.ListenAndServe(ipport, r))

	log4go.Info("end for api demo.")
}
