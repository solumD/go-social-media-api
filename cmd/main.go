package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"
	"github.com/solumD/go-social-media-api/cmd/server/database"
	"github.com/solumD/go-social-media-api/cmd/server/handlers"
	"github.com/solumD/go-social-media-api/internal/config"
)

// открытие базы данных и подключение к ней
func initDataBase() {
	var err error
	database.DBConn, err = sql.Open("sqlite3", "./cmd/server/database/news.db")
	if err != nil {
		log.Println(err)
		return
	}
	err = database.DBConn.Ping()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("✓ connected to books db")
}

// инициализация хендлеров
func initHandlers(r *chi.Mux) {
	// домашняя страница
	r.Get(`/feed`, handlers.Feed)

	// регистрация
	r.Post(`/register`, handlers.RegUnmarhalMW(handlers.RegCheckIfExistMW(handlers.Register)))

	// вход в аккаунт
	r.Post(`/login`, handlers.LogUnmarhalMW(handlers.LogCheckIfExistMW(handlers.Register)))

	// выход из аккаунта
	r.Post(`/exit`, handlers.ExitMiddleware(handlers.Exit))

	// создание поста пользователем
	r.Post(`/createpost`, handlers.Create)

	// вывод всех постов конкретного пользователя
	r.Get("/users/{user}", handlers.GetUserPosts)
}

func main() {
	// добавляем роутер
	r := chi.NewRouter()

	// инициализируем хендлеры к роутеру
	initHandlers(r)

	// инициализируем базу данных
	initDataBase()
	defer database.DBConn.Close()

	//у становка переменной окружения
	os.Setenv("CONFIG_PATH", "./config/local.yaml")

	// создание конфига
	cfg := config.MustLoad()
	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      r,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	// запуск сервера
	if err := srv.ListenAndServe(); err != nil {
		log.Print(err)
	}

	log.Fatal("server closed")
}
