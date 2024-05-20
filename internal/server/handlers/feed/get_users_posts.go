package feed

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	db "github.com/solumD/go-social-media-api/storage"
)

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
