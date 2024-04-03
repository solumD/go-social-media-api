package database

import (
	_ "github.com/mattn/go-sqlite3"
)

// Структура Post
type Post struct {
	Login   string `json:"author"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Date    string `json:"created on"`
}

// Функция вносит пост пользователя в таблицу posts
func InsertPost(user_id int, title, content, date string) error {
	query := `insert into posts(user_id, title, content, date_created) values (?, ?, ?, ?)`
	data := []any{user_id, title, content, date}
	if _, err := DBConn.Exec(query, data...); err != nil {
		return err
	}
	return nil
}

// Функция которая возвращает все посты конкретного пользователя
func SelectUserPosts(login string) ([]Post, error) {
	query := `select login, title, content, date_created from posts
	inner join users
	on users.id = posts.user_id
	where login = ?
	order by posts.id desc`

	rows, err := DBConn.Query(query, login)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []Post{}
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.Login, &post.Title, &post.Content, &post.Date); err != nil {
			return nil, err
		}
		post.Date = post.Date[0:10] // обрезаем дату, чтобы отображались только день, месяц и год
		posts = append(posts, post)
	}

	return posts, nil
}

func SelectLatestTenPosts() ([]Post, error) {
	query := `select login, title, content, date_created from posts
	join users on users.id = posts.user_id
	order by posts.id desc
	limit 0, 5;`

	rows, err := DBConn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	posts := []Post{}
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.Login, &post.Title, &post.Content, &post.Date); err != nil {
			return nil, err
		}
		post.Date = post.Date[0:10] // обрезаем дату, чтобы отображались только день, месяц и год
		posts = append(posts, post)
	}
	return posts, nil
}
