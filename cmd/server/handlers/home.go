package handlers

import (
	"log"
	"net/http"

	db "github.com/solumD/go-social-media-api/cmd/server/database"
	"github.com/solumD/go-social-media-api/cmd/server/handlers/common"
	"github.com/solumD/go-social-media-api/cmd/server/handlers/jwt"
	"github.com/solumD/go-social-media-api/cmd/server/handlers/person"
)

// Хендлер домашней страницы
func Home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("All News"))
}

func Create(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	claims, err := jwt.DecodeJWTToken(token)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid jwt token. Creation denied.", http.StatusUnauthorized)
		return
	}
	login := claims["sub"].(string)
	post, err := common.UnmarshalPost(r)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	post.Login = login
	if len(post.Title) == 0 || len(post.Content) == 0 {
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte("Post's title and content can't be empty"))
		return
	}
	user_id, err := db.SelectUserId(login)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = person.CreatePost(user_id, post); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Post created!"))
}
