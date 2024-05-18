package feed

import (
	"database/sql"
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
	_, err := db.SelectUserId(login)
	if err == sql.ErrNoRows {
		log.Printf(`error: user doesn't exist`)
		http.Error(w, `{"error":"user doesn't exist"}`, http.StatusBadRequest)
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
		http.Error(w, `{"error":"invalid jwt token, access denied"}`, http.StatusUnauthorized)
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

// Удаление поста по его id
func Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset = UTF-8")
	token := r.Header.Get("Authorization")
	claims, err := jwt.DecodeJWTToken(token)
	if err != nil {
		log.Println(err)
		http.Error(w, `{"error":"invalid jwt token, access denied"}`, http.StatusUnauthorized)
		return
	}
	login := claims["sub"].(string)
	id, err := common.UnmarshalId(r)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	postLogin, err := db.SelectPostLogin(id) // смотрим, кому принадлежит пост
	if err == sql.ErrNoRows {
		log.Printf(`error: post doesn't exist`)
		http.Error(w, `{"error":"post doesn't exist"}`, http.StatusBadRequest)
		return
	} else if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if login != postLogin { // проверяем, совпадает ли логин автора поста с логином отправителя запроса в jwt-токене
		log.Println("error: invalid user, access denied")
		http.Error(w, `{"error":"invalid user, access denied"}`, http.StatusBadRequest)
		return
	}
	err = db.DeletePost(id) // удаляем пост
	if err == sql.ErrNoRows {
		log.Println(err)
		resp := fmt.Sprintf(`{"error":"%s"}`, err)
		http.Error(w, resp, http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`{"message":"post deleted"}`))
}
