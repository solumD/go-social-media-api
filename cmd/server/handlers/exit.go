package handlers

import (
	"fmt"
	"log"
	"net/http"
)

func Exit(w http.ResponseWriter, r *http.Request) {
	user, err := UnmarshalBody(r)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	token := r.Header.Get("Token")
	_, loggedIn := CurrentUsers[user.Login]
	if !loggedIn { // проверка, выполнен вход или нет
		message := fmt.Sprintf("User %s is not in system!", user.Login)
		w.Write([]byte(message))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if token != CurrentUsers[user.Login] {
		log.Println("Regected in exit.")
		w.Write([]byte("Invalid jwt-token. Rejected."))
		w.WriteHeader(http.StatusForbidden)
		return
	}
	err = user.CheckUserLogin()
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	delete(CurrentUsers, user.Login)
	goodbye := fmt.Sprintf("See you soon, %s!", user.Login)
	log.Println(CurrentUsers)
	w.Write([]byte(goodbye))
}
