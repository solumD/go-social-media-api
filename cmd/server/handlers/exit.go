package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"

	db "github.com/solumD/go-social-media-api/cmd/server/database"
	"github.com/solumD/go-social-media-api/cmd/server/handlers/jwt"
)

// Выход пользователя из аккаунта по jwt токену
func Exit(w http.ResponseWriter, r *http.Request) {
	currLogin := r.Context().Value("Login").(string) // получаем логин из контекста
	if _, exist := db.CurrentUsers[currLogin]; !exist {
		http.Error(w, "User is not authorized", http.StatusUnauthorized)
		return
	}
	delete(db.CurrentUsers, currLogin)
	goodbye := fmt.Sprintf("See you soon, %s!", currLogin)
	w.Write([]byte(goodbye))
	log.Println(db.CurrentUsers)
}

// Middleware для проверки jwt токена
func ExitMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		claims, err := jwt.DecodeJWTToken(token)
		if err != nil {
			log.Println(err)
			http.Error(w, "Invalid jwt token. Exit denied.", http.StatusBadRequest)
			return
		}
		currLogin := claims["sub"].(string)
		ctx := r.Context()
		ctx = context.WithValue(ctx, "Login", currLogin) // отправляем логин пользователя в контекст
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
