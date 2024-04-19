package authorization

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/solumD/go-social-media-api/internal/server/handlers/common"
	"github.com/solumD/go-social-media-api/internal/server/handlers/jwt"
	"github.com/solumD/go-social-media-api/internal/server/handlers/person"
	db "github.com/solumD/go-social-media-api/storage"
)

type ContextUser string
type UserBody string // тип ключа в контексте

// Вход пользователя в свой аккаунт
func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset = UTF-8")
	user := r.Context().Value(ContextUser("User")).(*person.User) // получаем структуру User из контекста
	userToken, err := jwt.GenerateJWTToken(user.Login)            // генерация jwt токена

	if err != nil {
		log.Println(err)
		resp := fmt.Sprintf(`{"error":"%s"}`, err)
		http.Error(w, resp, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("JWT-Token", userToken)
	resp := fmt.Sprintf(`{"login":"%s"}`, user.Login)
	w.Write([]byte(resp))

	db.CurrentUsers[user.Login] = struct{}{} // выполнен вход в аккаунт, человек добавляется в список текущих пользователей
	log.Println(db.CurrentUsers)
}

// Middleware для декодирования json
func LogUnmarhalMW(next http.HandlerFunc) http.HandlerFunc {
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
		ctx = context.WithValue(ctx, ub, user) // отправляем структуру User в контекст
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// Middleware для проверки существования пользователя
func LogCheckIfExistMW(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset = UTF-8")
		user := r.Context().Value(UserBody("User")).(*person.User)
		if _, loggedIn := db.CurrentUsers[user.Login]; loggedIn { // проверка, выполнен вход или нет
			resp := fmt.Sprintf(`{"error": "user %s already logged in!"}`, user.Login)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(resp))
			return
		}
		err := common.CheckUserLogin(user.Login, user.Password) // проверка на существование пользователя и соответствие введенного пароля
		if err == sql.ErrNoRows {                               // в базе данных не найден пользователь с указанным логином
			log.Println(err)
			resp := fmt.Sprintf(`{"error":"%s"}`, err)
			http.Error(w, resp, http.StatusNotFound)
			return
		} else if err != nil { // ошибка во время исполнения запроса
			log.Println(err)
			resp := fmt.Errorf(`{"error":"%s"}`, err)
			http.Error(w, resp.Error(), http.StatusBadRequest)
			return
		} else {
			ctx := r.Context()
			tp := ContextUser("User")
			ctx = context.WithValue(ctx, tp, user) // отправляем структуру User в контекст
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}
