package storage

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var (
	DBConn *sql.DB
)

var (
	CurrentUsers = make(map[string]struct{}) // мапа с текущими пользователями
)
