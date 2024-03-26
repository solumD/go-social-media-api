package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var (
	DBConn *sql.DB
)

func InsertUser(login, password string) error {
	query := `insert into users(login, password) values (?, ?)`
	data := []any{login, password}
	if _, err := DBConn.Exec(query, data...); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func SelectUser(login string) (string, error) {
	type TempUser struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	var user TempUser
	query := `select login, password from users where login = ?`
	row := DBConn.QueryRow(query, login)
	err := row.Scan(&user.Login, &user.Password)
	if err != nil {
		return "", err
	}
	return user.Password, nil
}
