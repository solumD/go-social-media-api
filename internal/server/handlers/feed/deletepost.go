package feed

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/solumD/go-social-media-api/internal/server/handlers/common"
	"github.com/solumD/go-social-media-api/internal/server/handlers/jwt"
	db "github.com/solumD/go-social-media-api/storage"
)

// Удаление поста по его id
func DeletePost(w http.ResponseWriter, r *http.Request) {
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
