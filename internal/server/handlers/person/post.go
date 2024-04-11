package person

import (
	"time"

	db "github.com/solumD/go-social-media-api/storage"
)

// Структура Post
type Post struct {
	Login   string `json:"user"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// Функция собирает пост с указанием даты создания
// и вносит его в базу данных
func CreatePost(user_id int, post *Post) error {
	date := time.Now()
	dbDate := date.Format("2006-01-02")
	if err := db.InsertPost(user_id, post.Title, post.Content, dbDate); err != nil {
		return err
	}
	return nil
}
