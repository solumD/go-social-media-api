package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/solumD/go-social-media-api/cmd/server/handlers/jwt"
)

func Exit(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	jwtPayload, err := jwt.DecodeJWTToken(token)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}

	login := jwtPayload["sub"].(string)
	delete(CurrentUsers, login)
	goodbye := fmt.Sprintf("See you soon, %s!", login)
	w.Write([]byte(goodbye))
	log.Println(CurrentUsers)
}
