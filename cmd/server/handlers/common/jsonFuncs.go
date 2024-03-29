package common

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/solumD/go-social-media-api/cmd/server/handlers/person"
)

func UnmarshalBody(r *http.Request) (*person.User, error) {
	var user person.User
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
		return nil, errors.New("invalid json User Input")
	}
}

func UnmarshalPost(r *http.Request) (*person.Post, error) {
	var post person.Post
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	if json.Valid(body) {
		if err = json.Unmarshal(body, &post); err != nil {
			return nil, err
		}
		return &post, nil
	} else {
		return nil, errors.New("invalid json Post Input")
	}
}
