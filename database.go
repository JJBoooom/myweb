//package database
package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

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

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
