package server

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"
	"github.com/solumD/go-social-media-api/cmd/server/database"
	"github.com/solumD/go-social-media-api/cmd/server/handlers"
)

func initDataBase() {
	var err error
	database.DBConn, err = sql.Open("sqlite3", "news.db")
	if err != nil {
		log.Println(err)
	}
	err = database.DBConn.Ping()
	if err != nil {
		log.Println(err)
	}
	log.Println("✓ connected to books db")
}

func initHandlers(r *chi.Mux) {
	// домашняя страница
	r.Get(`/home`, handlers.Home)

	// регистрация
	r.Post(`/register`, handlers.Register)

	// вход в аккаунт
	r.Post(`/login`, handlers.Login)

	// выход из аккаунта
	r.Post(`/exit`, handlers.Exit)

	// создание поста пользователем
	r.Post(`/createpost`, handlers.Create)
}

func Server() {
	// запуск сервера
	r := chi.NewRouter()
	initHandlers(r)
	initDataBase()

	// подключение к базе данных
	defer database.DBConn.Close()
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Println(err)
	}
}
