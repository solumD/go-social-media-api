package handlers

import (
	"net/http"
)

func redirectToLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, `http://localhost:8080/login`, http.StatusSeeOther)
}

func redirectToRegister(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, `http://localhost:8080/register`, http.StatusSeeOther)
}

func redirectToHome(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, `http://localhost:8080/home`, http.StatusSeeOther)
}
