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
func initDataBase(cfg *config.Config) error {
	var err error
	storage.DBConn, err = sql.Open("sqlite3", cfg.DatabasePath)
	if err != nil {
		return err
	}
	err = storage.DBConn.Ping()
	if err != nil {
		return err
	}
	log.Println("✓ connected to news db")
	return nil
}

// инициализация хендлеров
func initHandlers(r *chi.Mux) {
	// регистрация
	r.Post(`/register`, h.RegUnmarhalMW(h.RegCheckIfExistMW(h.Register)))

	// вход в аккаунт
	r.Post(`/login`, h.LogUnmarhalMW(h.LogCheckIfExistMW(h.Login)))

	// выход из аккаунта
	r.Post(`/exit`, h.ExitMiddleware(h.Exit))

	// домашняя страница
	r.Get(`/feed`, f.Feed)

	// создание поста пользователем
	r.Post(`/createpost`, f.CreatePost)

	// вывод всех постов конкретного пользователя
	r.Get("/users/{user}", f.GetUserPosts)

	// удаление поста
	r.Post("/deletepost", f.DeletePost)
}

func main() {

	// инициализируем роутер
	r := chi.NewRouter()

	// установка пути к конфигу в переменную CONFIG_PATH
	os.Setenv("CONFIG_PATH", "./config/local.yaml")

	// инициализация конфига
	cfg := config.MustLoad()

	// инициализация сервера
	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      r,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	// инициализация базы данных
	if err := initDataBase(cfg); err != nil {
		log.Fatal(err)
	}
	defer storage.DBConn.Close()

	// инициализация хендлеров роутера
	initHandlers(r)

	// запуск сервера
	log.Printf("Starting server at %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
