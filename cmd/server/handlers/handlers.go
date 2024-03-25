package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// Хендлер домашней страницы
func Home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("All News"))
}

// Хендлер регистрации
func Register(w http.ResponseWriter, r *http.Request) {
	var user User
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	err = json.Unmarshal(body, &user)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	_, loggedIn := CurrentUsers[user.Login]
	if loggedIn { // проверка, выполнен вход или нет
		message := fmt.Sprintf("%s already logged in!", user.Login)
		w.Write([]byte(message))
	} else {
		answer, err := user.CheckUserRegister() // проверка, есть ли пользователем с введенным логином
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		if answer == "exists" { // если answer == exists, то пользователь уже есть, отмена операции
			message := fmt.Sprintf("User %s already exists!", user.Login)
			w.Write([]byte(message))
			return
		}

		err = user.GeneratePassword() // в случае, если пользователя нет, то пароль шифруется
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		response, err := user.CreateUser() // создание нового пользователя
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadGateway)
		}

		CurrentUsers[user.Login] = struct{}{} // выполнен вход в аккаунт, человек добавляется в список текущик пользователей
		log.Println(CurrentUsers)
		w.Write([]byte(response))
	}
}

// Хендлер входа в аккаунт
func Login(w http.ResponseWriter, r *http.Request) {
	var user User
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	err = json.Unmarshal(body, &user)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	_, loggedIn := CurrentUsers[user.Login]
	if loggedIn { // проверка, выполнен вход или нет
		message := fmt.Sprintf("%s already logged in!", user.Login)
		w.Write([]byte(message))
	} else {
		message, err := user.CheckUserLogin() // проверка, есть ли пользователем с введенным логином

		if err == sql.ErrNoRows { // в базе данных не найден пользователь с указанным логином
			log.Println(message)
		} else if err != nil { // ошибка во время исполнения запроса
			log.Println(err)
			w.Write([]byte(message))
		} else { // все ок
			w.Write([]byte(message))
			CurrentUsers[user.Login] = struct{}{} // выполнен вход в аккаунт, человек добавляется в список текущик пользователе
			log.Println(CurrentUsers)
		}
	}
}
