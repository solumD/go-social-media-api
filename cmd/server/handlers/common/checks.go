package common

import (
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
	"github.com/solumD/go-social-media-api/cmd/server/database"
	"golang.org/x/crypto/bcrypt"
)

// Фукнция проверяет во время входа в аккаунт, существует ли пользователь с введенным логином
func CheckUserLogin(login, password string) error {
	realPass, err := database.SelectUser(login)
	if err == sql.ErrNoRows {
		return errors.New("invalid login")
	} else if err != nil {
		return err
	} else {
		if err = bcrypt.CompareHashAndPassword([]byte(realPass), []byte(password)); err != nil {
			return err
		}
		return nil
	}
}

// Функция проверяет во время регистрации, существует ли пользователь с введенным логином
func CheckUserRegister(login string) (string, error) {
	_, err := database.SelectUser(login)
	if err == sql.ErrNoRows {
		return "doesn't exist", nil
	} else if err != nil {
		return "error", err
	}
	return "exists", nil
}
