package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var (
	DBConn *sql.DB
)

var (
	CurrentUsers = make(map[string]struct{}) // мапа с текущими пользователями
)

// Функция вносит нового пользователя в таблицу users
func InsertUser(login, password string) error {
	query := `insert into users(login, password) values (?, ?)`
	data := []any{login, password}
	if _, err := DBConn.Exec(query, data...); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// Функция ищет пользователся в таблице users, и,
// если находит, возвращает его пароль
func SelectUser(login string) (string, error) {
	type TempUser struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	var user TempUser
	query := `select login, password from users where login = ?`
	row := DBConn.QueryRow(query, login)
	err := row.Scan(&user.Login, &user.Password)
	if err != nil {
		return "", err
	}
	return user.Password, nil
}

// Функция возвращает id пользователя из таблицы users
func SelectUserId(login string) (int, error) {
	type ID struct {
		user_id int
	}
	var id ID
	query := `select id from users where login = ?`
	row := DBConn.QueryRow(query, login)
	err := row.Scan(&id.user_id)
	if err != nil {
		return 0, err
	}
	return id.user_id, nil
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

// Структура Post
type Post struct {
	Login   string `json:"author"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Date    string `json:"created on"`
}

// Функция которая возвращает все посты конкретного пользователя
func SelectUserPosts(login string) ([]Post, error) {
	query := `select login, title, content, date_created from posts
	inner join users
	on users.id = posts.user_id
	where login = ?`

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
