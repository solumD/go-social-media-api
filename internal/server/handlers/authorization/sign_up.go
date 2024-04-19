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
	w.Header().Set("Content-Type", "application/json; charset = UTF-8")
	user := r.Context().Value(ContextUser("User")).(*person.User) // получаем струтуру User из контекста
	if err := user.EncryptPassword(); err != nil {                // шифрование пароля
		log.Println(err)
		resp := fmt.Sprintf(`{"error":"%s"}`, err)
		http.Error(w, resp, http.StatusBadRequest)
		return
	}
	err := user.CreateUser() // создание нового пользователя
	if err != nil {
		log.Println(err)
		resp := fmt.Sprintf(`{"error":"%s"}`, err)
		http.Error(w, resp, http.StatusBadGateway)
		return
	}
	userToken, err := jwt.GenerateJWTToken(user.Login) // генерация jwt токена
	if err != nil {
		log.Println(err)
		resp := fmt.Sprintf(`{"error":"%s"}`, err)
		http.Error(w, resp, http.StatusBadGateway)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("JWT-Token", userToken)
	resp := fmt.Sprintf(`{"login":"%s"}`, user.Login)
	w.Write([]byte(resp))

	db.CurrentUsers[user.Login] = struct{}{} // выполнена регистрация, человек добавляется в список текущих пользователей
	log.Println(db.CurrentUsers)
}

// Middleware для декодирования json
func RegUnmarhalMW(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset = UTF-8")
		user, err := common.UnmarshalBody(r)
		if err != nil {
			log.Println(err)
			resp := fmt.Sprintf(`{"error":"%s"}`, err)
			http.Error(w, resp, http.StatusBadRequest)
			return
		}
		ctx := r.Context()
		ub := UserBody("User")
		ctx = context.WithValue(ctx, ub, user) // отправлка структуры User в контекст
		next.ServeHTTP(w, r.WithContext(ctx))
	}

}

// Middleware для проверки существования пользователя
func RegCheckIfExistMW(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset = UTF-8")
		user := r.Context().Value(UserBody("User")).(*person.User)
		answer, err := common.CheckUserRegister(user.Login) // проверка, есть ли пользователем с введенным логином
		if err != nil {
			log.Println(err)
			resp := fmt.Sprintf(`{"error":"%s"}`, err)
			http.Error(w, resp, http.StatusNotAcceptable)
			return
		}
		if answer == "exists" { // если answer == exists, то пользователь уже есть, отмена операции
			resp := fmt.Sprintf(`{"message": "user %s already exist"}`, user.Login)
			log.Println(resp)
			w.Write([]byte(resp))
			return
		}
		ctx := r.Context()
		cu := ContextUser("User")
		ctx = context.WithValue(ctx, cu, user) // отправлка структуры User в контекст
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
