package feed

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/solumD/go-social-media-api/internal/server/handlers/common"
	"github.com/solumD/go-social-media-api/internal/server/handlers/jwt"
	"github.com/solumD/go-social-media-api/internal/server/handlers/person"
	db "github.com/solumD/go-social-media-api/storage"
)

// Главная страница
func Feed(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset = UTF-8")
	posts, err := db.SelectLatestTenPosts()
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// выводим последние 10 постов от пользователей
	for _, v := range posts {
		data, err := json.MarshalIndent(v, "", "\t")
		if err != nil {
			log.Println(err)
			resp := fmt.Sprintf(`{"error":"%s"}`, err)
			http.Error(w, resp, http.StatusBadGateway)
			return
		}
		w.Write(data)
	}
}

// Все посты конкретного пользователя
func GetUserPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset = UTF-8")

	login := chi.URLParam(r, "user")
	if len(login) == 0 {
		log.Printf(`error:expected username after /user/`)
		http.Error(w, `{"error":"expected username after /user/"}`, http.StatusBadRequest)
		return
	}
	posts, err := db.SelectUserPosts(login)
	if err != nil {
		log.Println(err)
		resp := fmt.Sprintf(`{"error":"%s"}`, err)
		http.Error(w, resp, http.StatusBadRequest)
		return
	}
	if len(posts) == 0 {
		message := fmt.Sprintf(`{"message":"%s hasn't post something yet :("}`, login)
		log.Println(message)
		w.Write([]byte(message))
		return
	}

	// выводим все посты пользователя
	for _, v := range posts {
		data, err := json.MarshalIndent(v, "", "\t")
		if err != nil {
			log.Println(err)
			resp := fmt.Sprintf(`{"error":"%s"}`, err)
			http.Error(w, resp, http.StatusBadRequest)
			return
		}
		w.Write(data)
	}

}

// Создание поста с проверкой jwt токена
func Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset = UTF-8")
	token := r.Header.Get("Authorization")
	claims, err := jwt.DecodeJWTToken(token)
	if err != nil {
		log.Println(err)
		http.Error(w, `{"error":"invalid jwt token, creation denied"}`, http.StatusUnauthorized)
		return
	}
	login := claims["sub"].(string)
	post, err := common.UnmarshalPost(r)
	if err != nil {
		log.Println(err)
		resp := fmt.Sprintf(`{"error":"%s"}`, err)
		http.Error(w, resp, http.StatusBadRequest)
		return
	}
	post.Login = login
	if len(post.Title) == 0 || len(post.Content) == 0 {
		log.Println("error: post's title and content can't be empty")
		http.Error(w, `{"error":"post's title and content can't be empty"}`, http.StatusBadGateway)
		return
	}
	user_id, err := db.SelectUserId(login)
	if err != nil {
		log.Println(err)
		resp := fmt.Sprintf(`{"error":"%s"}`, err)
		http.Error(w, resp, http.StatusBadRequest)
		return
	}
	if err = person.CreatePost(user_id, post); err != nil {
		log.Println(err)
		resp := fmt.Sprintf(`{"error":"%s"}`, err)
		http.Error(w, resp, http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message":"post created"}`))
}
