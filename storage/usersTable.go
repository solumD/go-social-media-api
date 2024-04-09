package storage

import (
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// Функция вносит нового пользователя в таблицу users
func InsertUser(login, password string) error {
	query := `insert into users(login, password) values (?, ?)`
	data := []any{login, password}
	if _, err := DBConn.Exec(query, data...); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// Функция ищет пользователся в таблице users, и,
// если находит, возвращает его пароль
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

// Функция возвращает id пользователя из таблицы users
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
