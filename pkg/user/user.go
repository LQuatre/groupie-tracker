package userGestion

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	passwordManager "groupietracker.com/m/pkg/password"
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

func ValidateEmail(email string) bool {
	if len(email) < 3 {
		return false
	}
	if len(email) > 254 {
		return false
	}
	at := false
	dot := false
	for i, c := range email {
		if c == '@' {
			if at {
				return false
			}
			at = true
			if i == 0 {
				return false
			}
			if i == len(email)-1 {
				return false
			}
		}
		if c == '.' {
			if !at {
				return false
			}
			if i == 0 {
				return false
			}
			if i == len(email)-1 {
				return false
			}
			dot = true
		}
	}
	return dot
}

func Register(username, password, mail string) (UserStruct, string) {
	if username == "" || password == "" || mail == "" {
		return UserStruct{}, "The username, password, and mail fields are required."
	}
	if len(password) < 8 {
		return UserStruct{}, "password must be at least 8 characters long"
	}
	if len(username) < 4 {
		return UserStruct{}, "username must be at least 4 characters long"
	}
	if !ValidateEmail(mail) {
		return UserStruct{}, "The email address is not valid."
	}

	// check if the user already exists
	_, err := GetUser(username)
	if err == nil {
		return UserStruct{}, "This username is already used."
	}
	// check if the email already exists
	_, err = GetUserByMail(mail)
	if err == nil {
		return UserStruct{}, "This email is already used."
	}

	hashedPassword, err := passwordManager.HashPassword(password)

	stmt, err := myDataBase.Db.Prepare("INSERT INTO user (username, password, mail) VALUES (?, ?, ?)")
	if err != nil {
		return UserStruct{}, ""
	}

	_, err = stmt.Exec(username, hashedPassword, mail)
	if err != nil {
		return UserStruct{}, ""
	}

	// fmt.Println("User registered successfully.")
	return UserStruct{Username: username, Password: hashedPassword, Mail: mail}, ""
}

func Login(w http.ResponseWriter, username, password string) (UserStruct, string) {
	// get the user from the database
	if username == "" || password == "" {
		return UserStruct{}, "The username and password fields are required."
	}
	user, err := GetUser(username)
	if err != nil {
		return UserStruct{}, "The user already exists."
	}

	var PasswordsMatch = passwordManager.DoPasswordsMatch(user.Password, password)

	if !PasswordsMatch {
		return UserStruct{}, "The user does not exist or the password is incorrect."
	}

	cookie := http.Cookie{
		Name:  "username",
		Value: user.Username,
	}

	http.SetCookie(w, &cookie)
	return user, ""
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

func GetUserByMail(mail string) (UserStruct, error) {
	stmt, err := myDataBase.Db.Prepare("SELECT username, password, mail, starred, grade FROM user WHERE mail = ?")
	if err != nil {
		return UserStruct{}, fmt.Errorf("failed to prepare the SQL statement: %w", err)
	}

	var user UserStruct
	err = stmt.QueryRow(mail).Scan(&user.Username, &user.Password, &user.Mail, &user.Starred, &user.Grade)
	if err != nil {
		return UserStruct{}, fmt.Errorf("failed to execute the SQL statement: %w", err)
	}
	return user, nil
}

func GetAllUsers() ([]UserStruct, error) {
	stmt, err := myDataBase.Db.Prepare("SELECT username, password, mail, starred, grade FROM user")
	if err != nil {
		return nil, fmt.Errorf("failed to prepare the SQL statement: %w", err)
	}

	rows, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("failed to execute the SQL statement: %w", err)
	}

	var users []UserStruct
	for rows.Next() {
		var user UserStruct
		err = rows.Scan(&user.Username, &user.Password, &user.Mail, &user.Starred, &user.Grade)
		if err != nil {
			return nil, fmt.Errorf("failed to scan the row: %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}

func DeleteUser(username string) error {
	// Supprimer l'utilisateur de la base de données
	_, err := myDataBase.Db.Exec("DELETE FROM user WHERE username = ?", username)
	if err != nil {
		return fmt.Errorf("erreur lors de la suppression de l'utilisateur: %v", err)
	}
	return nil
}
