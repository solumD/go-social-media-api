package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func UnmarshalBody(r *http.Request) (*User, error) {
	var user User
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	if json.Valid(body) {
		if err = json.Unmarshal(body, &user); err != nil {
			return nil, err
		}
		return &user, nil
	} else {
		return nil, errors.New("invalid json Input")
	}
}

func GenerateJWTToken(login string) (string, error) {
	jwtSecret := []byte("golang-is-very-cool")
	payload := jwt.MapClaims{
		"sub": login,
		"exp": time.Now().Add(time.Hour * 48).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	t, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return t, nil
}
