package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var (
	DBConn *sql.DB
)

var (
	CurrentUsers = make(map[string]struct{}) // мапа с текущими пользователями
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

func SelectUserId(login string) (int, error) {
	type ID struct {
		user_id int
	}
	var id ID
	query := `select id from users where login = ?`
	row := DBConn.QueryRow(query, login)
	err := row.Scan(&id.user_id)
	if err != nil {
		return 0, err
	}
	return id.user_id, nil
}

type Post struct {
	Login   string `json:"user"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func InserPost(user_id int, title, content, date string) error {
	query := `insert into posts(user_id, title, content, date_created) values (?, ?, ?, ?)`
	data := []any{user_id, title, content, date}
	if _, err := DBConn.Exec(query, data...); err != nil {
		return err
	}
	return nil
}
