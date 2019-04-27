package user

import (
	"database/sql"
	"fmt"
	data "video_project/data"

	_ "github.com/go-sql-driver/mysql"
)

func checkErr(err error) {
	if err != nil {
		//panic(err)
		fmt.Println("register::checkErr::" + err.Error())
	}
}

var (
	dbhostip   = "localhost"
	dbusername = "root"
	dbpassword = "root"
	dbname     = "zither_video"
)

func queryUserByUserName(userName string) (int, data.UserInfo) {
	fmt.Println("user::queryUserByUserName::userName = " + userName)

	db, err := sql.Open("mysql", dbusername+":"+dbpassword+"@tcp("+dbhostip+")/"+dbname+"?charset=utf8")

	checkErr(err)
	stmt, err := db.Prepare("SELECT * FROM user where user_name =?")
	checkErr(err)

	rows, err := stmt.Query(userName)

	checkErr(err)
	for rows.Next() {
		var user_id string
		var user_name string
		var pass_word string
		var user_type int
		err := rows.Scan(&user_id, &user_name, &pass_word, &user_type)
		checkErr(err)
		fmt.Println("user::queryUserByUserName::", user_id, user_name, pass_word, user_type)
		return 0, data.UserInfo{user_id, user_name, pass_word, user_type}
	}
	defer rows.Close()
	defer stmt.Close()
	defer db.Close()
	return -1, data.UserInfo{}
}

func Register(user data.UserInfo) int {
	fmt.Println("user::Register")

	if res, _ := queryUserByUserName(user.UserName); res == 0 {
		return -1
	}

	db, openErr := sql.Open("mysql", dbusername+":"+dbpassword+"@tcp("+dbhostip+")/"+dbname+"?charset=utf8")
	checkErr(openErr)

	stmt, prepareErr := db.Prepare("insert into user (user_name,pass_word,user_type) values (?,?,1)")
	checkErr(prepareErr)

	_, execErr := stmt.Exec(user.UserName, user.PassWord)
	checkErr(execErr)

	defer stmt.Close()
	defer db.Close()

	return 0
}

func Login(user data.UserInfo) (int, data.UserInfo) {
	fmt.Println("user::Login")

	fmt.Println("user::Login::userName = " + user.UserName)
	fmt.Println("user::Login::passWord = " + user.PassWord)

	db, err := sql.Open("mysql", dbusername+":"+dbpassword+"@tcp("+dbhostip+")/"+dbname+"?charset=utf8")

	checkErr(err)
	stmt, err := db.Prepare("SELECT * FROM user where user_name = ? and pass_word = ?")
	checkErr(err)

	rows, err := stmt.Query(user.UserName, user.PassWord)

	checkErr(err)
	if rows.Next() {
		var user_id string
		var user_name string
		var pass_word string
		var user_type int
		err := rows.Scan(&user_id, &user_name, &pass_word, &user_type)
		checkErr(err)
		fmt.Println("user::Login::", user_id, user_name, pass_word, user_type)
		return 0, data.UserInfo{user_id, user_name, pass_word, user_type}
	}
	defer rows.Close()
	defer stmt.Close()
	defer db.Close()
	return -1, data.UserInfo{}
}
