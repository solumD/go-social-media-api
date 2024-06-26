package common

import (
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
	db "github.com/solumD/go-social-media-api/storage"
	"golang.org/x/crypto/bcrypt"
)

// Фукнция проверяет во время входа в аккаунт, существует ли пользователь с введенным логином. Если существует то проверяет соответствие паролей
func CheckUserLogin(login, password string) error {
	realPass, err := db.SelectUser(login)
	if err == sql.ErrNoRows {
		return errors.New(`user doesn't exist`)
	} else if err != nil {
		return err
	} else {
		if err = bcrypt.CompareHashAndPassword([]byte(realPass), []byte(password)); err != nil {
			return errors.New(`invalid password`)
		}
		return nil
	}
}

// Функция проверяет во время регистрации, существует ли пользователь с введенным логином
func CheckUserRegister(login string) (string, error) {
	_, err := db.SelectUser(login)
	if err == sql.ErrNoRows {
		return "doesn't exist", nil
	} else if err != nil {
		return "error", err
	}
	return "exists", nil
}
