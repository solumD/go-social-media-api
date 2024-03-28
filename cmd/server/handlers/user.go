package handlers

import (
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
	"github.com/solumD/go-social-media-api/cmd/server/database"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Login    string `json:"login"`
	Password string `json:"password,omitempty"`
}

var (
	CurrentUsers = make(map[string]struct{}) // мапа с текущими пользователями
)

// Метод создает пользователя и добавляет его в базу данных
func (u User) CreateUser() error {
	if err := database.InsertUser(u.Login, u.Password); err != nil {
		return err
	}
	return nil
}

// Метод шифрует пароль пользователя
func (u *User) EncryptPassword() error {
	cost := 10
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), cost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return nil
}

// Метод проверяет во время входа в аккаунт, существует ли пользователь с введенным логином
func (u User) CheckUserLogin() error {
	realPass, err := database.SelectUser(u.Login)
	if err == sql.ErrNoRows {
		return errors.New("invalid login")
	} else if err != nil {
		return err
	} else {
		if err = bcrypt.CompareHashAndPassword([]byte(realPass), []byte(u.Password)); err != nil {
			return err
		}
		return nil
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
