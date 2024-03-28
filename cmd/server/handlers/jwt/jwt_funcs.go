package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	jwtSecret            = []byte("golang-is-very-cool")
	errWrongJWTTokenType = errors.New("wrong type of JWT token claims")
	errInvalidToken      = errors.New("invalid token")
)

func GenerateJWTToken(login string) (string, error) {
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

func DecodeJWTToken(tokenString string) (jwt.MapClaims, error) {
	secretCode := []byte(jwtSecret)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretCode, nil
	})
	if err != nil {
		return nil, errWrongJWTTokenType
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, errInvalidToken
	}
}
