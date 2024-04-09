package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/solumD/go-social-media-api/cmd/server/handlers/common"
	"github.com/solumD/go-social-media-api/cmd/server/handlers/jwt"
	"github.com/solumD/go-social-media-api/cmd/server/handlers/person"
	db "github.com/solumD/go-social-media-api/storage"
)

// Главная страница
func Feed(w http.ResponseWriter, r *http.Request) {
	posts, err := db.SelectLatestTenPosts()
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	message := "|  Latest posts  |\n"
	w.Write([]byte(message))

	// выводим последние 10 постов от пользователей
	for _, v := range posts {
		data, err := json.MarshalIndent(v, "", "\t")
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			http.Error(w, "Error in marshalling json", http.StatusBadGateway)
			return
		}
		w.Write(data)
	}

	w.Header().Add("Content-Type", "application/json; charset = UTF-8")
}

// Все посты конкретного пользователя
func GetUserPosts(w http.ResponseWriter, r *http.Request) {
	login := chi.URLParam(r, "user")
	if len(login) == 0 {
		http.Error(w, "expected username after /user/", http.StatusBadRequest)
		return
	}
	posts, err := db.SelectUserPosts(login)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(posts) == 0 {
		message := fmt.Sprintf("%s hasn't post something yet :(", login)
		w.Write([]byte(message))
		return
	} else {
		message := fmt.Sprintf("|  %s's posts |\n\n", login)
		w.Write([]byte(message))
	}

	// выводим все посты пользователя
	for _, v := range posts {
		data, err := json.MarshalIndent(v, "", "\t")
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			http.Error(w, "Error in marshalling json", http.StatusBadGateway)
			return
		}
		w.Write(data)
	}

	w.Header().Add("Content-Type", "application/json; charset = UTF-8")
}

// Создание поста с проверкой jwt токена
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
