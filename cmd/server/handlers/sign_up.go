package handlers

import (
	"fmt"
	"log"
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request) {
	user, err := UnmarshalBody(r)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	answer, err := user.CheckUserRegister() // проверка, есть ли пользователем с введенным логином
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}
	if answer == "exists" { // если answer == exists, то пользователь уже есть, отмена операции
		message := fmt.Sprintf("User %s already exists!", user.Login)
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write([]byte(message))
		return
	}
	if err = user.EncryptPassword(); err != nil { // в случае, если пользователя нет, то пароль шифруется
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = user.CreateUser() // создание нового пользователя
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	userToken, err := GenerateJWTToken(user.Login)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	CurrentUsers[user.Login] = userToken
	message := fmt.Sprintf("Welcome, %s!\nYour jwt-token: %s\nDon't lose it!", user.Login, userToken) // выполнен вход в аккаунт, человек добавляется в список текущик пользователей
	log.Println(CurrentUsers)
	w.Write([]byte(message))
	w.WriteHeader(http.StatusOK)
}
