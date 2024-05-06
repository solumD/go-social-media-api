package authorization

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/solumD/go-social-media-api/internal/server/handlers/jwt"
	db "github.com/solumD/go-social-media-api/storage"
)

type ContextLogin string

// Выход пользователя из аккаунта по jwt токену
func Exit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset = UTF-8")

	currLogin := r.Context().Value(ContextLogin("Login")).(string) // получаем логин из контекста
	if _, exist := db.CurrentUsers[currLogin]; !exist {
		log.Println("user is not authorized")
		http.Error(w, `{"error":"user is not authorized"}`, http.StatusUnauthorized)
		return
	}

	delete(db.CurrentUsers, currLogin)
	log.Println(db.CurrentUsers)
	message := fmt.Sprintf(`{"message":"see you soon, %s"}`, currLogin)
	w.Write([]byte(message))
}

// Middleware для проверки jwt токена
func ExitMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset = UTF-8")

		token := r.Header.Get("Authorization")
		claims, err := jwt.DecodeJWTToken(token)
		if err != nil {
			log.Println(err)
			http.Error(w, `{"error":"invalid jwt-token, exit denied"}`, http.StatusBadRequest)
			return
		}

		currLogin := claims["sub"].(string)
		ctx := r.Context()
		tp := ContextLogin("Login")
		ctx = context.WithValue(ctx, tp, currLogin) // отправляем логин пользователя в контекст
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
