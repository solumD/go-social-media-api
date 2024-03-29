package handlers

import (
	"fmt"
	"log"
	"net/http"

	db "github.com/solumD/go-social-media-api/cmd/server/database"
	"github.com/solumD/go-social-media-api/cmd/server/handlers/jwt"
)

func Exit(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	claims, err := jwt.DecodeJWTToken(token)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid jwt token. Exit denied.", http.StatusUnauthorized)
		return
	}

	login := claims["sub"].(string)
	delete(db.CurrentUsers, login)
	goodbye := fmt.Sprintf("See you soon, %s!", login)
	w.Write([]byte(goodbye))
	log.Println(db.CurrentUsers)
}
