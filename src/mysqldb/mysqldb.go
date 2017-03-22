package mysqldb

import (
	"accounts"
	"database/sql"
	"errors"
	"github.com/alecthomas/log4go"
	_ "github.com/go-sql-driver/mysql"
)

//mysql db
type MysqlDb struct {
	db *sql.DB
}

var MyDbInstance MysqlDb

func InitDb() error {
	log4go.Info("mysqldb Init()...")
	// mydb := new(MysqlDb)
	// db, err := sql.Open("mysql", "test:test@tcp(127.0.0.1:3306)/godb?charset=utf8")
	db, err := sql.Open("mysql", "test:mysql@/Accounts?charset=utf8")
	// db, err := sql.Open("mysql", "test:mysql@/?charset=utf8")
	//args1 ： the db engine
	//args2 :  the db DSN, different engine different DSN
	//https://github.com/go-sql-driver/mysql
	if err != nil {
		log4go.Error("database initialize error : ", err.Error())
		return err
	}
	// log4go.Info(db)
	MyDbInstance.db = db
	// log4go.Info(MyDbInstance.db)
	return nil
}

func CloseDb() {
	if e := CheckIniDb(); e != nil {
		return
	}
	MyDbInstance.db.Close()
}

func CheckIniDb() error {
	if MyDbInstance.db == nil {
		return errors.New("DB did not be inited...")
	}
	return nil
}

func PingDb() error {
	return MyDbInstance.db.Ping()
}

func CreateUsersTable() error {

	createtable := `create table if not exists users (		
		id int unsigned not null AUTO_INCREMENT PRIMARY KEY,
		name char(18) not null,
		sex tinyint unsigned,
		age tinyint unsigned,
        email char(255),
		phone char(18),
        createdate int unsigned,
        ts timestamp);`

	// log4go.Info(createtable)
	stmt, err := MyDbInstance.db.Prepare(createtable)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec()
	if err != nil {
		log4go.Info(result)
		return err
	}

	return nil
}

func InsertUser(id int, name string, sex int, age int, email string, phone string, createdate int64) error {
	// stmt, err := MyDbInstance.db.Prepare(
	// 	"insert into users(id,name,sex,age,email,phone, createdate)values(?,?,?,?,?,?,?)")

	//Don't insert id
	stmt, err := MyDbInstance.db.Prepare(
		"insert into users(name,sex,age,email,phone, createdate)values(?,?,?,?,?,?)")

	if err != nil {
		log4go.Error(err)
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(name, sex, age, email, phone, createdate)
	if err != nil {
		log4go.Error(err)
		return err
	} else {
		log4go.Info(result)
	}

	return nil
}

func feedbackToResult(rows *sql.Rows) (*accounts.Account, error) {
	defer rows.Close()
	if rows == nil {
		return nil, errors.New("can not fild from db.")
	}
	var id, sex, age int
	var name, email, phone string
	var createdate int64

	if rows.Next() {
		err := rows.Scan(&id, &name, &sex, &age, &email, &phone, &createdate)
		if err == nil {
			//ok

		} else {
			log4go.Error(err)
		}
	}

	aa := accounts.Account{ID: id,
		Name:       name,
		Sex:        sex,
		Age:        age,
		Email:      email,
		Phone:      phone,
		CreateDate: createdate}
	return &aa, nil
}

func QueryUserByName(n string) (*accounts.Account, error) {

	stmt, err := MyDbInstance.db.Prepare("select id, name, sex, age, email, phone, createdate from users where name=?")
	if err != nil {
		log4go.Error(err)
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(n)
	if err != nil {
		return nil, err
	}
	// result, err := rows.Columns()
	// if err != nil {
	// 	log4go.Info(err)
	// 	return id, name, sex, age, email, phone, createdate
	// }

	return feedbackToResult(rows)
}

func QueryUserById(i int) (*accounts.Account, error) {
	stmt, err := MyDbInstance.db.Prepare("select id, name, sex, age, email, phone, createdate from users where id=?")
	if err != nil {
		log4go.Error(err)
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(i)
	if err != nil {
		return nil, err
	}

	// result, err := rows.Columns()
	// if err != nil {
	// 	log4go.Info(err)
	// 	return id, name, sex, age, email, phone, createdate
	// }

	return feedbackToResult(rows)
}

func UpdateUser(id int, name string, sex int, age int, email string, phone string) {

}

func DeleteUserByeName(name string) {

}

func NameExistInDb(name string) (string, error) {
	log4go.Info("NameExistInDb, name=", name)
	resName := ""
	if name == "" {
		return resName, errors.New("name is empty.")
	}

	stmt, err := MyDbInstance.db.Prepare("select name from users where name=?")
	if err != nil {
		log4go.Error(err)
		return resName, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(name)
	defer rows.Close()
	if err != nil {
		return resName, err
	}

	has := rows.Next()
	//must use rows.Next()
	if has {
		//name exist
		return name, errors.New("name exist.")
	} else {
		//name did NOT exist
		return resName, nil
	}

}
