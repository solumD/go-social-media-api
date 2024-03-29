package person

import (
	"time"

	db "github.com/solumD/go-social-media-api/cmd/server/database"
)

type Post struct {
	Login   string `json:"user"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func CreatePost(user_id int, post *Post) error {
	date := time.Now()
	dbDate := date.Format("2006-01-02")
	if err := db.InserPost(user_id, post.Title, post.Content, dbDate); err != nil {
		return err
	}
	return nil
}
