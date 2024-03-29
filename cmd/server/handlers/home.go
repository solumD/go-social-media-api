package handlers

import (
	"net/http"
)

// Хендлер домашней страницы
func Home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("All News"))
}
