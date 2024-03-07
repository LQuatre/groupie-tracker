package user

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type UserStruct struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Mail     string `json:"mail"`
	Starred  string `json:"starred"`
	Grade    string `json:"grade"`
}

type MyDataBase struct {
	Db *sql.DB
}

var myDataBase MyDataBase
var err error

func SetMySQL() *sql.DB {
	myDataBase.Db, err = sql.Open("mysql", "root:@tcp(localhost)/groupie_tracker")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Database successfully connected.")
	}
	return myDataBase.Db
}

func Register(username, password, mail string) (UserStruct, error) {
	stmt, err := myDataBase.Db.Prepare("INSERT INTO user (username, password, mail) VALUES (?, ?, ?)")
	if err != nil {
		return UserStruct{}, fmt.Errorf("failed to prepare the SQL statement: %w", err)
	}

	_, err = stmt.Exec(username, password, mail)
	if err != nil {
		return UserStruct{}, fmt.Errorf("failed to execute the SQL statement: %w", err)
	}

	fmt.Println("User registered successfully.")
	return UserStruct{Username: username, Password: password, Mail: mail}, nil
}

func Login(w http.ResponseWriter, username, password string) (UserStruct, error) {
	stmt, err := myDataBase.Db.Prepare("SELECT username, password, mail, starred, grade FROM user WHERE username = ? AND password = ?")
	if err != nil {
		return UserStruct{}, fmt.Errorf("failed to prepare the SQL statement: %w", err)
	}

	var user UserStruct
	err = stmt.QueryRow(username, password).Scan(&user.Username, &user.Password, &user.Mail, &user.Starred, &user.Grade)
	if err != nil {
		return UserStruct{}, fmt.Errorf("failed to execute the SQL statement: %w", err)
	}

	fmt.Println("User logged in successfully.")

	cookie := http.Cookie{
		Name:  "username",
		Value: user.Username,
	}

	http.SetCookie(w, &cookie)
	fmt.Println("User cookie set successfully.")
	return user, nil
}

func GetUser(username string) (UserStruct, error) {
	stmt, err := myDataBase.Db.Prepare("SELECT username, password, mail, starred, grade FROM user WHERE username = ?")
	if err != nil {
		return UserStruct{}, fmt.Errorf("failed to prepare the SQL statement: %w", err)
	}

	var user UserStruct
	err = stmt.QueryRow(username).Scan(&user.Username, &user.Password, &user.Mail, &user.Starred, &user.Grade)
	if err != nil {
		return UserStruct{}, fmt.Errorf("failed to execute the SQL statement: %w", err)
	}
	return user, nil
}
