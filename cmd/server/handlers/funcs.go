package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
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
