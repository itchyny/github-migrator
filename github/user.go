package github

import (
	"encoding/json"
	"fmt"
)

// User represents a user.
type User struct {
	Login   string `json:"login"`
	HTMLURL string `json:"html_url"`
}

type userOrError struct {
	User
	Message string `json:"message"`
}

func (c *client) GetUser() (*User, error) {
	res, err := c.get(c.url("/user"))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var r userOrError
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	if r.Message != "" {
		return nil, fmt.Errorf("GetUser %s: %s", "/user", r.Message)
	}

	return &r.User, nil
}
