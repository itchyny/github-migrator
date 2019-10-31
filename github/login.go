package github

import (
	"encoding/json"
	"errors"
)

type loginResponse struct {
	Login   string `json:"login"`
	Message string `json:"message"`
}

func (c *client) Login() (string, error) {
	res, err := c.get(c.url("/user"))
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var r loginResponse
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return "", err
	}

	if r.Message != "" {
		return "", errors.New(r.Message)
	}

	return r.Login, nil
}
