package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

// Хендлер домашней страницы
func Home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("All News"))
}

// Хендлер регистрации
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if answer == "exists" { // если answer == exists, то пользователь уже есть, отмена операции
		message := fmt.Sprintf("User %s already exists!", user.Login)
		w.Write([]byte(message))
		return
	}
	if err = user.EncryptPassword(); err != nil { // в случае, если пользователя нет, то пароль шифруется
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response, err := user.CreateUser() // создание нового пользователя
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	CurrentUsers[user.Login] = struct{}{} // выполнен вход в аккаунт, человек добавляется в список текущик пользователей
	log.Println(CurrentUsers)
	w.Write([]byte(response))

}

// Хендлер входа в аккаунт
func Login(w http.ResponseWriter, r *http.Request) {
	user, err := UnmarshalBody(r)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, loggedIn := CurrentUsers[user.Login]
	if loggedIn { // проверка, выполнен вход или нет
		message := fmt.Sprintf("%s already logged in!", user.Login)
		w.Write([]byte(message))
	} else {
		var message string
		err := user.CheckUserLogin() // проверка, есть ли пользователем с введенным логином
		if err == sql.ErrNoRows {
			message := fmt.Sprintf("User %s doesn't exist!", user.Login) // в базе данных не найден пользователь с указанным логином
			log.Println(message)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		} else if err != nil { // ошибка во время исполнения запроса
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		} else { // все ок
			message = fmt.Sprintf("Welcome Back, %s!", user.Login)
			w.Write([]byte(message))
			CurrentUsers[user.Login] = struct{}{} // выполнен вход в аккаунт, человек добавляется в список текущик пользователе
			log.Println(CurrentUsers)
		}
	}
}

func Exit(w http.ResponseWriter, r *http.Request) {
	user, err := UnmarshalBody(r)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, loggedIn := CurrentUsers[user.Login]
	if !loggedIn { // проверка, выполнен вход или нет
		message := fmt.Sprintf("User %s is not in system!", user.Login)
		w.Write([]byte(message))
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
