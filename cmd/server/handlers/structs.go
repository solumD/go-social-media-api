package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/solumD/go-social-media-api/cmd/server/database"
	"golang.org/x/crypto/bcrypt"
)

var (
	CurrentUsers = make(map[string]struct{}) // мапа с текущими пользователями
)

type User struct {
	Id       int    `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

// Метод создает пользователя и добавляет его в базу данных
func (u User) CreateUser() (string, error) {
	query := `insert into users(login, password) values (?, ?)`
	data := []any{u.Login, u.Password}
	_, err := database.DBConn.Exec(query, data...)
	if err != nil {
		log.Println(err)
		return "", err
	}
	response := fmt.Sprintf("Successfully created a user: %s", u.Login)
	return response, nil
}

// Метод шифрует пароль
func (u *User) GeneratePassword() error {
	cost := 10
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), cost)
	u.Password = string(hash)
	if err != nil {
		return err
	}
	return nil
}

// Метод проверяет во время входа в аккаунт, существует ли пользователь с введенным логином
func (u User) CheckUserLogin() (string, error) {
	var realUser User
	var message string
	query := `select id, login, password from users where login = ?`
	row := database.DBConn.QueryRow(query, u.Login)
	err := row.Scan(&realUser.Id, &realUser.Login, &realUser.Password)

	if err == sql.ErrNoRows {
		message = "This user doesn't exist!"
		return message, errors.New("invalid login")
	} else if err != nil {
		return "", err
	} else {
		err = bcrypt.CompareHashAndPassword([]byte(realUser.Password), []byte(u.Password))
		if err != nil {
			return "Invalid password", err
		}
		message = fmt.Sprintf("Welcome Back, %s", u.Login)
		return message, nil
	}
}

// Метод проверяет во время регистрации, существует ли пользователь с введенным логином
func (u User) CheckUserRegister() (int, error) {
	var realUser User
	query := `select id, login, password from users where login = ?`
	row := database.DBConn.QueryRow(query, u.Login)
	err := row.Scan(&realUser.Id, &realUser.Login, &realUser.Password)
	if err == sql.ErrNoRows {
		return 1, nil
	} else if err != nil {
		return 2, err
	}
	return 3, nil
}
