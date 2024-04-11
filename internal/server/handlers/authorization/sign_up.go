package authorization

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/solumD/go-social-media-api/internal/server/handlers/common"
	"github.com/solumD/go-social-media-api/internal/server/handlers/jwt"
	"github.com/solumD/go-social-media-api/internal/server/handlers/person"
	db "github.com/solumD/go-social-media-api/storage"
)

// Регистрация пользователя
func Register(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(ContextUser("User")).(*person.User) // получаем струтуру User из контекста =
	if err := user.EncryptPassword(); err != nil {                // шифрование пароля
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err := user.CreateUser() // создание нового пользователя
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	userToken, err := jwt.GenerateJWTToken(user.Login) // генерация jwt токена
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	db.CurrentUsers[user.Login] = struct{}{}
	message := fmt.Sprintf("Welcome, %s!\nYour jwt-token: %s\nDon't lose it!", user.Login, userToken) // выполнен вход в аккаунт, человек добавляется в список текущик пользователей
	log.Println(db.CurrentUsers)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(message))
}

// Middleware для декодирования json
func RegUnmarhalMW(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := common.UnmarshalBody(r)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ctx := r.Context()
		ub := UserBody("User")
		ctx = context.WithValue(ctx, ub, user) // отправлка структуры User в контекст
		next.ServeHTTP(w, r.WithContext(ctx))
		log.Println("middleware 1")
	}

}

// Middleware для проверки существования пользователя
func RegCheckIfExistMW(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserBody("User")).(*person.User)
		answer, err := common.CheckUserRegister(user.Login) // проверка, есть ли пользователем с введенным логином
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusNotAcceptable)
			return
		}
		if answer == "exists" { // если answer == exists, то пользователь уже есть, отмена операции
			message := fmt.Sprintf("User %s already exists!", user.Login)
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write([]byte(message))
			return
		}
		ctx := r.Context()
		cu := ContextUser("User")
		ctx = context.WithValue(ctx, cu, user) // отправлка структуры User в контекст
		next.ServeHTTP(w, r.WithContext(ctx))
		log.Println("middleware 2")
	}
}
