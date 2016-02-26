//package webdb
package main

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var (
	db *sql.DB = nil
)

type DbOpts struct {
	user     string
	password string
	ip       string
	port     string
}

type Account struct {
	id     int
	user   string
	passwd string
	date   string
}

/*导入包的时候，建立连接*/
/*
func init() {
	opt := DbOpts{user: "root", password: "123456", ip: "192.168.4.29", port: "3306"}
	err := Open(opts)
	checkErr(err)
}
*/
func main() {
	opt := DbOpts{user: "root", password: "123456", ip: "192.168.2.119", port: "3306"}
	err := Open(opt)
	checkErr(err)

	defer Close()

	Query_data()
	fmt.Println("-----------")

	account := Account{user: "kiongf", passwd: "wangweihong", date: "2015-11-17"}
	insert_data(account)
	Query_data()
	fmt.Println("-----------")
	delete_data(account)
	Query_data()
	fmt.Println("-----------")
	change := Account{user: "123"}
	new := Account{passwd: "654321"}
	update_data(change, new)
	Query_data()
	fmt.Println("-----------")
	Query_password("123")
	Query_password("xxxx")
}

func Open(opt DbOpts) error {
	//("mysql", "root:123456@tcp(127.0.0.1:3306)/test")
	var err error
	connect_addr := opt.user + ":" + opt.password + "@tcp(" + opt.ip + ":" + opt.port + ")/" + "web?charset=utf8,parseTime=true"
	fmt.Println(connect_addr)
	db, err = sql.Open("mysql", connect_addr)

	if err != nil {
		return err
	}
	//Open不一定建立连接,Ping()会
	err = db.Ping()
	if err != nil {
		db.Close()
		return err
	} else {
		return nil
	}
}

func Close() error {
	if db != nil {
		err := db.Close()
		return err
	}
	return errors.New("db is nil")

}

func Query_data() {
	rows, err := db.Query("SELECT * FROM account")
	checkErr(err)

	//account := new(Account)
	var id int
	var user string
	var passwd string
	var date sql.NullString

	for rows.Next() {
		//	err = rows.Scan(&account.id, &account.user, &account.passwd, &account.date)
		err = rows.Scan(&id, &user, &passwd, &date)
		checkErr(err)

		if !date.Valid {
			date.String = ""
		}
		fmt.Println(id, user, passwd, date.String)

	}
}

//获取所有的用户
func Query_password(user string) error {
	if user == "" {
		return errors.New("empty string")
	}
	if db != nil {
		//设置将要进行的sql语句,查询account表中指定username数据，只获取username和password,而不返回所有数据
		stmt, err := db.Prepare("select username, password from account where username = ?")
		checkErr(err)

		//插入需要查询的数据
		rows, err := stmt.Query(user)
		checkErr(err)

		for rows.Next() {
			var user string
			var passwd string
			err = rows.Scan(&user, &passwd)
			checkErr(err)

			fmt.Println(user + ":" + passwd)
		}

	}
	err := errors.New("db is nil")
	return err

}

func insert_data(account Account) error {
	if db != nil {
		//设置将要进行的sql语句
		stmt, err := db.Prepare("INSERT account SET username=?, password=?, created=?")
		checkErr(err)

		//插入需要处理的数据
		_, err = stmt.Exec(account.user, account.passwd, account.date)
		checkErr(err)

	}
	err := errors.New("db is nil")
	return err

}

func delete_data(account Account) error {
	if db != nil {
		//设置将要进行的sql语句
		stmt, err := db.Prepare("delete from account where username=?")
		checkErr(err)

		//插入需要处理的数据
		_, err = stmt.Exec(account.user)
		checkErr(err)

	}
	err := errors.New("db is nil")
	return err

}

func update_data(account Account, new Account) error {
	if db != nil {
		stmt, err := db.Prepare("update account set password=? where username=?")
		checkErr(err)

		_, err = stmt.Exec(new.passwd, account.user)
		checkErr(err)

	}
	err := errors.New("db is nil")
	return err

}

/*
func main() {
	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/test")
	checkErr(err)

	defer db.Close()

	err = db.Ping()
	checkErr(err)

	//插入数据
	stmt, err := db.Prepare("INSERT userinfo SET username=?, departname=?,created=?")
	checkErr(err)

	res, err := stmt.Exec("kiongf", "研发部门", "2016-02-23")
	checkErr(err)

	id, err := res.LastInsertId()
	checkErr(err)

	fmt.Println(id)

	//更新数据
	stmt, err = db.Prepare("update userinfo set username=? where uid=?")
	checkErr(err)

	res, err = stmt.Exec("kiongfupdate", id)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Println(affect)

	//查询数据
	rows, err := db.Query("SELECT * FROM userinfo")
	checkErr(err)

	for rows.Next() {
		var uid int
		var username string
		var department string
		var created string
		err = rows.Scan(&uid, &username, &department, &created)
		checkErr(err)
		fmt.Println(uid)
		fmt.Println(username)
		fmt.Println(department)
		fmt.Println(created)

		//删除数据
		stmt, err = db.Prepare("delete from userinfo where uid=?")
		checkErr(err)

		res, err = stmt.Exec(id)
		checkErr(err)

		fmt.Println(affect)
	}

}
*/
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
