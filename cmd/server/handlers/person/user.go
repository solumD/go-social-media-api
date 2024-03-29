package person

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/solumD/go-social-media-api/cmd/server/database"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Login    string `json:"login"`
	Password string `json:"password,omitempty"`
}

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
