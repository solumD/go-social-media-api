package feed

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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
