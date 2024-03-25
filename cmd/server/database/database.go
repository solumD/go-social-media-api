package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var (
	DBConn *sql.DB
)

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func InsertUser(login, password string) error {
	query := `insert into users(login, password) values (?, ?)`
	data := []any{login, password}
	_, err := DBConn.Exec(query, data...)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func SelectUser(login string) (string, error) {
	var realUser User
	query := `select login, password from users where login = ?`
	row := DBConn.QueryRow(query, login)
	err := row.Scan(&realUser.Login, &realUser.Password)
	if err != nil {
		return "", err
	}
	return realUser.Password, nil
}
