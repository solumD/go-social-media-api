package storage

import (
	_ "github.com/mattn/go-sqlite3"
)

// Структура Post
type Post struct {
	Id      string `json:"id"`
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
	query := `select posts.id, login, title, content, date_created from posts
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
		if err := rows.Scan(&post.Id, &post.Login, &post.Title, &post.Content, &post.Date); err != nil {
			return nil, err
		}
		post.Date = post.Date[0:10] // обрезаем дату, чтобы отображались только день, месяц и год
		posts = append(posts, post)
	}

	return posts, nil
}

// Функция возвращает последние 10 постов от всех пользователей
func SelectLatestTenPosts() ([]Post, error) {
	query := `select posts.id, login, title, content, date_created from posts
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
		if err := rows.Scan(&post.Id, &post.Login, &post.Title, &post.Content, &post.Date); err != nil {
			return nil, err
		}
		post.Date = post.Date[0:10] // обрезаем дату, чтобы отображались только день, месяц и год
		posts = append(posts, post)
	}
	return posts, nil
}

// Функция возвращает логин пользователя, которому принадлежит пост
func SelectPostLogin(postId string) (string, error) {
	query := `select login from posts inner join users on users.id = posts.user_id where posts.id = ?`
	row := DBConn.QueryRow(query, postId)
	type UserLogin struct {
		login string
	}
	var userLogin UserLogin
	err := row.Scan(&userLogin.login)
	if err != nil {
		return "0", err
	}
	return userLogin.login, nil
}

// Функция удаляет пост пользователя по id
func DeletePost(postId string) error {
	query := `delete from posts where id = ?`
	_, err := DBConn.Exec(query, postId)
	if err != nil {
		return err
	}
	return nil
}
