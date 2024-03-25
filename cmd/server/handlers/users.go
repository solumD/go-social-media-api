package handlers

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/solumD/go-social-media-api/cmd/server/database"
	"golang.org/x/crypto/bcrypt"
)

var (
	CurrentUsers = make(map[string]struct{}) // мапа с текущими пользователями
)

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// Метод создает пользователя и добавляет его в базу данных
func (u User) CreateUser() (string, error) {
	err := database.InsertUser(u.Login, u.Password)
	if err != nil {
		return "", err
	}
	response := fmt.Sprintf("Successfully created a user: %s", u.Login)
	return response, nil
}

// Метод шифрует пароль пользователя
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
	var message string
	realPass, err := database.SelectUser(u.Login)

	if err == sql.ErrNoRows {
		message = "This user doesn't exist!"
		return message, errors.New("invalid login")
	} else if err != nil {
		return "", err
	} else {
		err = bcrypt.CompareHashAndPassword([]byte(realPass), []byte(u.Password))
		if err != nil {
			return "Invalid password", err
		}
		message = fmt.Sprintf("Welcome Back, %s", u.Login)
		return message, nil
	}
}

// Метод проверяет во время регистрации, существует ли пользователь с введенным логином
func (u User) CheckUserRegister() (string, error) {
	_, err := database.SelectUser(u.Login)
	if err == sql.ErrNoRows {
		return "doesn't exist", nil
	} else if err != nil {
		return "error", err
	}
	return "exists", nil
}
