package feed

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/solumD/go-social-media-api/internal/server/handlers/common"
	"github.com/solumD/go-social-media-api/internal/server/handlers/jwt"
	"github.com/solumD/go-social-media-api/internal/server/handlers/person"
	db "github.com/solumD/go-social-media-api/storage"
)

// Создание поста с проверкой jwt токена
func CreatePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset = UTF-8")

	auth := r.Header.Get("Authorization")
	bearerAndToken := strings.Split(auth, " ")
	claims, err := jwt.DecodeJWTToken(bearerAndToken[1])

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
