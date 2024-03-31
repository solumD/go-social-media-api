package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	db "github.com/solumD/go-social-media-api/cmd/server/database"
	"github.com/solumD/go-social-media-api/cmd/server/handlers/common"
	"github.com/solumD/go-social-media-api/cmd/server/handlers/jwt"
)

// Вход пользователя в свой аккаунт
func Login(w http.ResponseWriter, r *http.Request) {
	user, err := common.UnmarshalBody(r)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, loggedIn := db.CurrentUsers[user.Login]
	if loggedIn { // проверка, выполнен вход или нет
		message := fmt.Sprintf("%s already logged in!", user.Login)
		w.Write([]byte(message))
	} else {
		var message string
		err := common.CheckUserLogin(user.Login, user.Password) // проверка, есть ли пользователем с введенным логином
		if err == sql.ErrNoRows {
			message := fmt.Sprintf("User %s doesn't exist!", user.Login) // в базе данных не найден пользователь с указанным логином
			log.Println(message)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		} else if err != nil { // ошибка во время исполнения запроса
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		} else { // все ок
			userToken, err := jwt.GenerateJWTToken(user.Login) // генерация jwt токена
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			message = fmt.Sprintf("Welcome Back, %s!\nYour jwt-token: %s\nDon't lose it!", user.Login, userToken)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(message))
			db.CurrentUsers[user.Login] = struct{}{} // выполнен вход в аккаунт, человек добавляется в список текущик пользователе
			log.Println(db.CurrentUsers)
		}
	}
}
