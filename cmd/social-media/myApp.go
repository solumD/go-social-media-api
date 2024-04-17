package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"
	"github.com/solumD/go-social-media-api/internal/config"
	h "github.com/solumD/go-social-media-api/internal/server/handlers/authorization"
	f "github.com/solumD/go-social-media-api/internal/server/handlers/feed"
	"github.com/solumD/go-social-media-api/storage"
)

// открытие базы данных и подключение к ней
func initDataBase(cfg *config.Config) {
	var err error
	storage.DBConn, err = sql.Open("sqlite3", cfg.DatabasePath)
	if err != nil {
		log.Println(err)
		return
	}
	err = storage.DBConn.Ping()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("✓ connected to books db")
}

// инициализация хендлеров
func initHandlers(r *chi.Mux) {
	// домашняя страница
	r.Get(`/feed`, f.Feed)

	// регистрация
	r.Post(`/register`, h.RegUnmarhalMW(h.RegCheckIfExistMW(h.Register)))

	// вход в аккаунт
	r.Post(`/login`, h.LogUnmarhalMW(h.LogCheckIfExistMW(h.Login)))

	// выход из аккаунта
	r.Post(`/exit`, h.ExitMiddleware(h.Exit))

	// создание поста пользователем
	r.Post(`/createpost`, f.Create)

	// вывод всех постов конкретного пользователя
	r.Get("/users/{user}", f.GetUserPosts)
}

func main() {
	// добавляем роутер
	r := chi.NewRouter()

	// установка переменной окружения
	os.Setenv("CONFIG_PATH", "./config/local.yaml")

	// инициализация конфига
	cfg := config.MustLoad()
	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      r,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	// инициализируем базу данных
	initDataBase(cfg)
	defer storage.DBConn.Close()

	// инициализируем хендлеры к роутеру
	initHandlers(r)

	// запуск сервера
	if err := srv.ListenAndServe(); err != nil {
		log.Print(err)
	}

	log.Fatal("server closed")
}
