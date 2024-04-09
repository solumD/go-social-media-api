package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/solumD/go-social-media-api/cmd/server/handlers/common"
	"github.com/solumD/go-social-media-api/cmd/server/handlers/jwt"
	"github.com/solumD/go-social-media-api/cmd/server/handlers/person"
	db "github.com/solumD/go-social-media-api/storage"
)

type ContextUser string
type UserBody string // тип ключа в контексте

// Вход пользователя в свой аккаунт
func Login(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(ContextUser("User")).(*person.User) // получаем структуру User из контекста
	userToken, err := jwt.GenerateJWTToken(user.Login)            // генерация jwt токена
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	message := fmt.Sprintf("Welcome Back, %s!\nYour jwt-token: %s\nDon't lose it!", user.Login, userToken)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(message))
	db.CurrentUsers[user.Login] = struct{}{} // выполнен вход в аккаунт, человек добавляется в список текущик пользователе
	log.Println(db.CurrentUsers)

}

// Middleware для декодирования json
func LogUnmarhalMW(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := common.UnmarshalBody(r)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ctx := r.Context()
		ub := UserBody("User")
		ctx = context.WithValue(ctx, ub, user) // отправляем структуру User в контекст
		next.ServeHTTP(w, r.WithContext(ctx))

	}
}

// Middleware для проверки существования пользователя
func LogCheckIfExistMW(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserBody("User")).(*person.User)
		message := ""
		if _, loggedIn := db.CurrentUsers[user.Login]; loggedIn { // проверка, выполнен вход или нет
			message = fmt.Sprintf("%s already logged in!", user.Login)
			w.Write([]byte(message))
			return
		}
		err := common.CheckUserLogin(user.Login, user.Password) // проверка, есть ли пользователь с введенным логином
		if err == sql.ErrNoRows {                               // в базе данных не найден пользователь с указанным логином
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		} else if err != nil { // ошибка во время исполнения запроса
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		} else {
			ctx := r.Context()
			tp := ContextUser("User")
			ctx = context.WithValue(ctx, tp, user) // отправляем структуру User в контекст
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}
